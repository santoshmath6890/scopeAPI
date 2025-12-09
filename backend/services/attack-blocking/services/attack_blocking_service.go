package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/attack-blocking/internal/models"
	"scopeapi.local/backend/services/attack-blocking/internal/repository"
	"scopeapi.local/backend/shared/messaging/kafka"
)

type AttackBlockingService struct {
	blockingRepo         repository.BlockingRepository
	policyRepo           repository.PolicyRepository
	kafkaProducer        kafka.Producer
	logger               *slog.Logger
	blockingRules        map[string]*models.BlockingRule
	blockingPolicies     map[string]*models.BlockingPolicy
	activeBlocks         map[string]*models.ActiveBlock
	rateLimiters         map[string]*models.RateLimiter
	ipWhitelist          map[string]bool
	ipBlacklist          map[string]bool
	geoBlocking          map[string]bool
	signatureDetectors   map[string]*models.SignatureDetector
	anomalyDetectors     map[string]*models.AnomalyDetector
	cloudIntelligence    *models.CloudIntelligence
	mutex                sync.RWMutex
	config               *AttackBlockingConfig
}

type AttackBlockingConfig struct {
	EnableRealTimeBlocking    bool          `json:"enable_real_time_blocking"`
	EnableCloudIntelligence   bool          `json:"enable_cloud_intelligence"`
	EnableGeoBlocking         bool          `json:"enable_geo_blocking"`
	EnableRateLimiting        bool          `json:"enable_rate_limiting"`
	EnableSignatureDetection  bool          `json:"enable_signature_detection"`
	EnableAnomalyDetection    bool          `json:"enable_anomaly_detection"`
	BlockingTimeout           time.Duration `json:"blocking_timeout"`
	MaxConcurrentBlocks       int           `json:"max_concurrent_blocks"`
	DefaultBlockDuration      time.Duration `json:"default_block_duration"`
	WhitelistEnabled          bool          `json:"whitelist_enabled"`
	BlacklistEnabled          bool          `json:"blacklist_enabled"`
	LogBlockedRequests        bool          `json:"log_blocked_requests"`
	NotifyOnBlock             bool          `json:"notify_on_block"`
	AutoUnblockEnabled        bool          `json:"auto_unblock_enabled"`
	AutoUnblockThreshold      int           `json:"auto_unblock_threshold"`
}

func NewAttackBlockingService(
	blockingRepo repository.BlockingRepository,
	policyRepo repository.PolicyRepository,
	kafkaProducer kafka.Producer,
	logger *slog.Logger,
	config *AttackBlockingConfig,
) *AttackBlockingService {
	service := &AttackBlockingService{
		blockingRepo:         blockingRepo,
		policyRepo:           policyRepo,
		kafkaProducer:        kafkaProducer,
		logger:               logger,
		blockingRules:        make(map[string]*models.BlockingRule),
		blockingPolicies:     make(map[string]*models.BlockingPolicy),
		activeBlocks:         make(map[string]*models.ActiveBlock),
		rateLimiters:         make(map[string]*models.RateLimiter),
		ipWhitelist:          make(map[string]bool),
		ipBlacklist:          make(map[string]bool),
		geoBlocking:          make(map[string]bool),
		signatureDetectors:   make(map[string]*models.SignatureDetector),
		anomalyDetectors:     make(map[string]*models.AnomalyDetector),
		config:               config,
	}

	// Initialize cloud intelligence if enabled
	if config.EnableCloudIntelligence {
		service.cloudIntelligence = &models.CloudIntelligence{
			ThreatFeeds:      make(map[string]*models.ThreatFeed),
			ReputationScores: make(map[string]float64),
			LastUpdated:      time.Now(),
		}
	}

	// Load initial configuration
	service.loadConfiguration()

	return service
}

func (s *AttackBlockingService) ProcessRequest(ctx context.Context, request *models.AttackBlockingRequest) (*models.AttackBlockingResult, error) {
	startTime := time.Now()

	// Generate request ID if not provided
	if request.RequestID == "" {
		request.RequestID = uuid.New().String()
	}

	s.logger.Info("Processing attack blocking request",
		"request_id", request.RequestID,
		"ip_address", request.IPAddress,
		"endpoint", request.Endpoint)

	// Check if IP is whitelisted
	if s.isWhitelisted(request.IPAddress) {
		return &models.AttackBlockingResult{
			RequestID:      request.RequestID,
			Action:         models.ActionAllow,
			Reason:         "IP address is whitelisted",
			ProcessingTime: time.Since(startTime),
			ProcessedAt:    time.Now(),
		}, nil
	}

	// Check if IP is blacklisted
	if s.isBlacklisted(request.IPAddress) {
		return s.blockRequest(ctx, request, "IP address is blacklisted", models.BlockReasonBlacklist, startTime)
	}

	// Check for active blocks
	if activeBlock := s.getActiveBlock(request.IPAddress); activeBlock != nil {
		return &models.AttackBlockingResult{
			RequestID:      request.RequestID,
			Action:         models.ActionBlock,
			Reason:         fmt.Sprintf("IP is currently blocked: %s", activeBlock.Reason),
			BlockID:        activeBlock.ID,
			BlockedUntil:   &activeBlock.ExpiresAt,
			ProcessingTime: time.Since(startTime),
			ProcessedAt:    time.Now(),
		}, nil
	}

	// Apply rate limiting
	if s.config.EnableRateLimiting {
		if rateLimitExceeded, limiter := s.checkRateLimit(request); rateLimitExceeded {
			return s.blockRequest(ctx, request, 
				fmt.Sprintf("Rate limit exceeded: %d requests in %v", limiter.RequestCount, limiter.Window),
				models.BlockReasonRateLimit, startTime)
		}
	}

	// Check geo-blocking
	if s.config.EnableGeoBlocking {
		if blocked, reason := s.checkGeoBlocking(request); blocked {
			return s.blockRequest(ctx, request, reason, models.BlockReasonGeoBlocking, startTime)
		}
	}

	// Apply signature detection
	if s.config.EnableSignatureDetection {
		if detected, signature := s.detectSignatures(request); detected {
			return s.blockRequest(ctx, request,
				fmt.Sprintf("Malicious signature detected: %s", signature.Name),
				models.BlockReasonSignature, startTime)
		}
	}

	// Apply anomaly detection
	if s.config.EnableAnomalyDetection {
		if anomalous, anomaly := s.detectAnomalies(request); anomalous {
			return s.blockRequest(ctx, request,
				fmt.Sprintf("Anomalous behavior detected: %s", anomaly.Description),
				models.BlockReasonAnomaly, startTime)
		}
	}

	// Check cloud intelligence
	if s.config.EnableCloudIntelligence {
		if threat, score := s.checkCloudIntelligence(request); threat {
			return s.blockRequest(ctx, request,
				fmt.Sprintf("Threat intelligence match (score: %.2f)", score),
				models.BlockReasonThreatIntelligence, startTime)
		}
	}

	// Apply custom blocking rules
	if blocked, rule := s.applyBlockingRules(request); blocked {
		return s.blockRequest(ctx, request,
			fmt.Sprintf("Custom rule triggered: %s", rule.Name),
			models.BlockReasonCustomRule, startTime)
	}

	// Request is allowed
	result := &models.AttackBlockingResult{
		RequestID:      request.RequestID,
		Action:         models.ActionAllow,
		Reason:         "Request passed all security checks",
		ProcessingTime: time.Since(startTime),
		ProcessedAt:    time.Now(),
	}

	// Log allowed request if configured
	if s.config.LogBlockedRequests {
		s.logRequest(ctx, request, result)
	}

	return result, nil
}

func (s *AttackBlockingService) blockRequest(ctx context.Context, request *models.AttackBlockingRequest, reason string, blockReason models.BlockReason, startTime time.Time) (*models.AttackBlockingResult, error) {
	blockID := uuid.New().String()
	blockDuration := s.config.DefaultBlockDuration
	expiresAt := time.Now().Add(blockDuration)

	// Create active block
	activeBlock := &models.ActiveBlock{
		ID:          blockID,
		IPAddress:   request.IPAddress,
		Reason:      reason,
		BlockReason: blockReason,
		RequestID:   request.RequestID,
		APIID:       request.APIID,
		EndpointID:  request.EndpointID,
		UserAgent:   request.UserAgent,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		Active:      true,
	}

	// Store active block
	s.mutex.Lock()
	s.activeBlocks[request.IPAddress] = activeBlock
	s.mutex.Unlock()

	// Persist block to repository
	if err := s.blockingRepo.CreateActiveBlock(ctx, activeBlock); err != nil {
		s.logger.Error("Failed to persist active block", "error", err, "block_id", blockID)
	}

	// Create blocking result
	result := &models.AttackBlockingResult{
		RequestID:      request.RequestID,
		Action:         models.ActionBlock,
		Reason:         reason,
		BlockID:        blockID,
		BlockedUntil:   &expiresAt,
		ProcessingTime: time.Since(startTime),
		ProcessedAt:    time.Now(),
	}

	// Publish blocking event
	if err := s.publishBlockingEvent(ctx, request, result, activeBlock); err != nil {
		s.logger.Error("Failed to publish blocking event", "error", err, "block_id", blockID)
	}

	// Send notification if configured
	if s.config.NotifyOnBlock {
		s.sendBlockingNotification(ctx, request, result, activeBlock)
	}

	// Log blocked request
	s.logRequest(ctx, request, result)

	s.logger.Warn("Request blocked",
		"request_id", request.RequestID,
		"block_id", blockID,
		"ip_address", request.IPAddress,
		"reason", reason,
		"expires_at", expiresAt)

	return result, nil
}

func (s *AttackBlockingService) isWhitelisted(ipAddress string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.config.WhitelistEnabled && s.ipWhitelist[ipAddress]
}

func (s *AttackBlockingService) isBlacklisted(ipAddress string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.config.BlacklistEnabled && s.ipBlacklist[ipAddress]
}

func (s *AttackBlockingService) getActiveBlock(ipAddress string) *models.ActiveBlock {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	if block, exists := s.activeBlocks[ipAddress]; exists {
		if block.Active && time.Now().Before(block.ExpiresAt) {
			return block
		}
		// Block has expired, remove it
		delete(s.activeBlocks, ipAddress)
	}
	return nil
}

func (s *AttackBlockingService) checkRateLimit(request *models.AttackBlockingRequest) (bool, *models.RateLimiter) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := fmt.Sprintf("%s:%s", request.IPAddress, request.EndpointID)
	limiter, exists := s.rateLimiters[key]
	
	if !exists {
		limiter = &models.RateLimiter{
			Key:          key,
			RequestCount: 0,
			Window:       time.Minute,
			Limit:        100, // Default limit
			WindowStart:  time.Now(),
		}
		s.rateLimiters[key] = limiter
	}

	// Reset window if expired
	if time.Since(limiter.WindowStart) > limiter.Window {
		limiter.RequestCount = 0
		limiter.WindowStart = time.Now()
	}

	limiter.RequestCount++
	limiter.LastRequest = time.Now()

	return limiter.RequestCount > limiter.Limit, limiter
}

func (s *AttackBlockingService) checkGeoBlocking(request *models.AttackBlockingRequest) (bool, string) {
	// This would integrate with a GeoIP service
	// For now, return false (not blocked)
	return false, ""
}

func (s *AttackBlockingService) detectSignatures(request *models.AttackBlockingRequest) (bool, *models.AttackSignature) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, detector := range s.signatureDetectors {
		if !detector.Enabled {
			continue
		}

		for _, signature := range detector.Signatures {
			if s.matchSignature(request, signature) {
				return true, signature
			}
		}
	}

	return false, nil
}

func (s *AttackBlockingService) matchSignature(request *models.AttackBlockingRequest, signature *models.AttackSignature) bool {
	content := strings.ToLower(request.RequestBody + request.QueryString + request.Headers["User-Agent"])
	
	for _, pattern := range signature.Patterns {
		if strings.Contains(content, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

func (s *AttackBlockingService) detectAnomalies(request *models.AttackBlockingRequest) (bool, *models.AnomalyPattern) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, detector := range s.anomalyDetectors {
		if !detector.Enabled {
			continue
		}

		for _, pattern := range detector.Patterns {
			if s.matchAnomalyPattern(request, pattern) {
				return true, pattern
			}
		}
	}

	return false, nil
}

func (s *AttackBlockingService) matchAnomalyPattern(request *models.AttackBlockingRequest, pattern *models.AnomalyPattern) bool {
	// Implement anomaly detection logic based on pattern type
	switch pattern.Type {
	case "request_size":
		if len(request.RequestBody) > pattern.Threshold {
			return true
		}
	case "request_frequency":
		// Check request frequency for this IP
		return s.checkRequestFrequency(request.IPAddress, pattern.Threshold)
	case "unusual_headers":
		return s.checkUnusualHeaders(request.Headers, pattern)
	}

	return false
}

func (s *AttackBlockingService) checkRequestFrequency(ipAddress string, threshold int) bool {
	// Implementation would check request frequency from logs/metrics
	return false
}

func (s *AttackBlockingService) checkUnusualHeaders(headers map[string]string, pattern *models.AnomalyPattern) bool {
	// Implementation would check for unusual header patterns
	return false
}

func (s *AttackBlockingService) checkCloudIntelligence(request *models.AttackBlockingRequest) (bool, float64) {
	if s.cloudIntelligence == nil {
		return false, 0.0
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Check reputation score
	if score, exists := s.cloudIntelligence.ReputationScores[request.IPAddress]; exists {
		if score < 0.3 { // Threshold for blocking
			return true, score
		}
	}

	// Check threat feeds
	for _, feed := range s.cloudIntelligence.ThreatFeeds {
		if feed.ContainsIP(request.IPAddress) {
			return true, 0.0
		}
	}

	return false, 0.0
}

func (s *AttackBlockingService) applyBlockingRules(request *models.AttackBlockingRequest) (bool, *models.BlockingRule) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Sort rules by priority
	var rules []*models.BlockingRule
	for _, rule := range s.blockingRules {
		if rule.Enabled {
			rules = append(rules, rule)
		}
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	// Apply rules in priority order
	for _, rule := range rules {
		if s.matchBlockingRule(request, rule) {
			return true, rule
		}
	}

	return false, nil
}

func (s *AttackBlockingService) matchBlockingRule(request *models.AttackBlockingRequest, rule *models.BlockingRule) bool {
	// Check all conditions in the rule
	for _, condition := range rule.Conditions {
		if !s.evaluateCondition(request, condition) {
			return false
		}
	}
	return true
}

func (s *AttackBlockingService) evaluateCondition(request *models.AttackBlockingRequest, condition models.RuleCondition) bool {
	var value string

	// Extract value based on field
	switch condition.Field {
	case "ip_address":
		value = request.IPAddress
	case "user_agent":
		value = request.UserAgent
	case "endpoint":
		value = request.Endpoint
	case "method":
		value = request.Method
	case "query_string":
		value = request.QueryString
	case "request_body":
		value = request.RequestBody
	case "header":
		if headerName, exists := condition.Context["header_name"]; exists {
			value = request.Headers[headerName.(string)]
		}
	default:
		return false
	}

	// Apply operator
	switch condition.Operator {
	case "equals":
		return value == condition.Value
	case "contains":
		return strings.Contains(value, condition.Value)
	case "starts_with":
		return strings.HasPrefix(value, condition.Value)
	case "ends_with":
		return strings.HasSuffix(value, condition.Value)
	case "regex":
		// Implementation would use regex matching
		return false
	case "greater_than":
		// For numeric comparisons
		return len(value) > len(condition.Value)
	case "less_than":
		return len(value) < len(condition.Value)
	default:
		return false
	}
}

func (s *AttackBlockingService) publishBlockingEvent(ctx context.Context, request *models.AttackBlockingRequest, result *models.AttackBlockingResult, block *models.ActiveBlock) error {
	event := map[string]interface{}{
		"event_type":   "attack_blocked",
		"request_id":   request.RequestID,
		"block_id":     result.BlockID,
		"ip_address":   request.IPAddress,
		"endpoint":     request.Endpoint,
		"reason":       result.Reason,
		"blocked_at":   result.ProcessedAt,
		"blocked_until": result.BlockedUntil,
		"api_id":       request.APIID,
		"endpoint_id":  request.EndpointID,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal blocking event: %w", err)
	}

	return s.kafkaProducer.Produce(ctx, "attack-blocking-events", eventData)
}

func (s *AttackBlockingService) sendBlockingNotification(ctx context.Context, request *models.AttackBlockingRequest, result *models.AttackBlockingResult, block *models.ActiveBlock) {
	// Implementation would send notifications via email, Slack, etc.
	s.logger.Info("Sending blocking notification",
		"block_id", result.BlockID,
		"ip_address", request.IPAddress,
		"reason", result.Reason)
}

func (s *AttackBlockingService) logRequest(ctx context.Context, request *models.AttackBlockingRequest, result *models.AttackBlockingResult) {
	logEntry := map[string]interface{}{
		"request_id":     request.RequestID,
		"ip_address":     request.IPAddress,
		"endpoint":       request.Endpoint,
		"method":         request.Method,
		"user_agent":     request.UserAgent,
		"action":         result.Action,
		"reason":         result.Reason,
		"processing_time": result.ProcessingTime.Milliseconds(),
		"processed_at":   result.ProcessedAt,
	}

	if result.BlockID != "" {
		logEntry["block_id"] = result.BlockID
		logEntry["blocked_until"] = result.BlockedUntil
	}

	logData, _ := json.Marshal(logEntry)
	s.logger.Info("Attack blocking request processed", "data", string(logData))
}

func (s *AttackBlockingService) loadConfiguration() {
	// Load blocking rules
	if rules, err := s.blockingRepo.GetAllBlockingRules(context.Background()); err == nil {
		s.mutex.Lock()
		for _, rule := range rules {
			s.blockingRules[rule.ID] = rule
		}
		s.mutex.Unlock()
	}

	// Load policies
	if policies, err := s.policyRepo.GetAllPolicies(context.Background()); err == nil {
		s.mutex.Lock()
		for _, policy := range policies {
			s.blockingPolicies[policy.ID] = policy
		}
		s.mutex.Unlock()
	}

	// Load IP lists
	s.loadIPLists()

	// Load signature detectors
	s.loadSignatureDetectors()

	// Load anomaly detectors
	s.loadAnomalyDetectors()
}

func (s *AttackBlockingService) loadIPLists() {
	// Load whitelist
	if whitelist, err := s.blockingRepo.GetIPWhitelist(context.Background()); err == nil {
		s.mutex.Lock()
		s.ipWhitelist = make(map[string]bool)
		for _, ip := range whitelist {
			s.ipWhitelist[ip] = true
		}
		s.mutex.Unlock()
	}

	// Load blacklist
	if blacklist, err := s.blockingRepo.GetIPBlacklist(context.Background()); err == nil {
		s.mutex.Lock()
		s.ipBlacklist = make(map[string]bool)
		for _, ip := range blacklist {
			s.ipBlacklist[ip] = true
		}
		s.mutex.Unlock()
	}
}

func (s *AttackBlockingService) loadSignatureDetectors() {
	// Implementation would load signature detectors from repository
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Example signature detector
	s.signatureDetectors["sql_injection"] = &models.SignatureDetector{
		ID:      "sql_injection",
		Name:    "SQL Injection Detector",
		Type:    "signature",
		Enabled: true,
		Signatures: []*models.AttackSignature{
			{
				ID:       "sql_injection_1",
				Name:     "SQL Injection Pattern 1",
				Type:     "sql_injection",
				Patterns: []string{"union select", "drop table", "'; --", "' or 1=1"},
				Severity: "high",
			},
		},
	}
}

func (s *AttackBlockingService) loadAnomalyDetectors() {
	// Implementation would load anomaly detectors from repository
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Example anomaly detector
	s.anomalyDetectors["request_size"] = &models.AnomalyDetector{
		ID:      "request_size",
		Name:    "Request Size Anomaly Detector",
		Type:    "anomaly",
		Enabled: true,
		Patterns: []*models.AnomalyPattern{
			{
				ID:          "large_request",
				Name:        "Large Request Body",
				Type:        "request_size",
				Threshold:   1024 * 1024, // 1MB
				Description: "Request body exceeds normal size",
			},
		},
	}
}

// Additional service methods

func (s *AttackBlockingService) GetActiveBlocks(ctx context.Context, filter *models.ActiveBlockFilter) ([]*models.ActiveBlock, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var blocks []*models.ActiveBlock
	for _, block := range s.activeBlocks {
		if s.matchesFilter(block, filter) {
			blocks = append(blocks, block)
		}
	}

	// Sort by creation time (newest first)
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].CreatedAt.After(blocks[j].CreatedAt)
	})

	// Apply pagination
	start := filter.Offset
	end := start + filter.Limit
	if start > len(blocks) {
		return []*models.ActiveBlock{}, nil
	}
	if end > len(blocks) {
		end = len(blocks)
	}

	return blocks[start:end], nil
}

func (s *AttackBlockingService) matchesFilter(block *models.ActiveBlock, filter *models.ActiveBlockFilter) bool {
	if filter.IPAddress != "" && block.IPAddress != filter.IPAddress {
		return false
	}
	if filter.APIID != "" && block.APIID != filter.APIID {
		return false
	}
	if filter.EndpointID != "" && block.EndpointID != filter.EndpointID {
		return false
	}
	if filter.BlockReason != "" && string(block.BlockReason) != filter.BlockReason {
		return false
	}
	if filter.StartDate != nil && block.CreatedAt.Before(*filter.StartDate) {
		return false
	}
	if filter.EndDate != nil && block.CreatedAt.After(*filter.EndDate) {
		return false
	}
	return true
}

func (s *AttackBlockingService) UnblockIP(ctx context.Context, ipAddress string, reason string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if block, exists := s.activeBlocks[ipAddress]; exists {
		block.Active = false
		block.UnblockedAt = &time.Time{}
		*block.UnblockedAt = time.Now()
		block.UnblockReason = reason

		// Remove from active blocks
		delete(s.activeBlocks, ipAddress)

		// Update in repository
		if err := s.blockingRepo.UpdateActiveBlock(ctx, block); err != nil {
			return fmt.Errorf("failed to update block in repository: %w", err)
		}

		s.logger.Info("IP unblocked", "ip_address", ipAddress, "reason", reason)
		return nil
	}

	return fmt.Errorf("no active block found for IP: %s", ipAddress)
}

func (s *AttackBlockingService) GetBlockingStatistics(ctx context.Context, filter *models.BlockingStatsFilter) (*models.BlockingStatistics, error) {
	stats, err := s.blockingRepo.GetBlockingStatistics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocking statistics: %w", err)
	}

	// Add real-time statistics
	s.mutex.RLock()
	stats.ActiveBlocks = len(s.activeBlocks)
	s.mutex.RUnlock()

	return stats, nil
}

func (s *AttackBlockingService) CreateBlockingRule(ctx context.Context, rule *models.BlockingRule) error {
	rule.ID = uuid.New().String()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	if err := s.blockingRepo.CreateBlockingRule(ctx, rule); err != nil {
		return fmt.Errorf("failed to create blocking rule: %w", err)
	}

	// Add to in-memory cache
	s.mutex.Lock()
	s.blockingRules[rule.ID] = rule
	s.mutex.Unlock()

	s.logger.Info("Blocking rule created", "rule_id", rule.ID, "rule_name", rule.Name)
	return nil
}

func (s *AttackBlockingService) UpdateBlockingRule(ctx context.Context, rule *models.BlockingRule) error {
	rule.UpdatedAt = time.Now()

	if err := s.blockingRepo.UpdateBlockingRule(ctx, rule); err != nil {
		return fmt.Errorf("failed to update blocking rule: %w", err)
	}

	// Update in-memory cache
	s.mutex.Lock()
	s.blockingRules[rule.ID] = rule
	s.mutex.Unlock()

	s.logger.Info("Blocking rule updated", "rule_id", rule.ID, "rule_name", rule.Name)
	return nil
}

func (s *AttackBlockingService) DeleteBlockingRule(ctx context.Context, ruleID string) error {
	if err := s.blockingRepo.DeleteBlockingRule(ctx, ruleID); err != nil {
		return fmt.Errorf("failed to delete blocking rule: %w", err)
	}

	// Remove from in-memory cache
	s.mutex.Lock()
	delete(s.blockingRules, ruleID)
	s.mutex.Unlock()

	s.logger.Info("Blocking rule deleted", "rule_id", ruleID)
	return nil
}

func (s *AttackBlockingService) GetBlockingRules(ctx context.Context, filter *models.BlockingRuleFilter) ([]*models.BlockingRule, error) {
	return s.blockingRepo.GetBlockingRules(ctx, filter)
}

func (s *AttackBlockingService) GetBlockingRule(ctx context.Context, ruleID string) (*models.BlockingRule, error) {
	s.mutex.RLock()
	if rule, exists := s.blockingRules[ruleID]; exists {
		s.mutex.RUnlock()
		return rule, nil
	}
	s.mutex.RUnlock()

	return s.blockingRepo.GetBlockingRule(ctx, ruleID)
}

func (s *AttackBlockingService) AddToWhitelist(ctx context.Context, ipAddress string, reason string) error {
	if err := s.blockingRepo.AddToWhitelist(ctx, ipAddress, reason); err != nil {
		return fmt.Errorf("failed to add IP to whitelist: %w", err)
	}

	// Update in-memory cache
	s.mutex.Lock()
	s.ipWhitelist[ipAddress] = true
	s.mutex.Unlock()

	s.logger.Info("IP added to whitelist", "ip_address", ipAddress, "reason", reason)
	return nil
}

func (s *AttackBlockingService) RemoveFromWhitelist(ctx context.Context, ipAddress string) error {
	if err := s.blockingRepo.RemoveFromWhitelist(ctx, ipAddress); err != nil {
		return fmt.Errorf("failed to remove IP from whitelist: %w", err)
	}

	// Update in-memory cache
	s.mutex.Lock()
	delete(s.ipWhitelist, ipAddress)
	s.mutex.Unlock()

	s.logger.Info("IP removed from whitelist", "ip_address", ipAddress)
	return nil
}

func (s *AttackBlockingService) AddToBlacklist(ctx context.Context, ipAddress string, reason string) error {
	if err := s.blockingRepo.AddToBlacklist(ctx, ipAddress, reason); err != nil {
		return fmt.Errorf("failed to add IP to blacklist: %w", err)
	}

	// Update in-memory cache
	s.mutex.Lock()
	s.ipBlacklist[ipAddress] = true
	s.mutex.Unlock()

	s.logger.Info("IP added to blacklist", "ip_address", ipAddress, "reason", reason)
	return nil
}

func (s *AttackBlockingService) RemoveFromBlacklist(ctx context.Context, ipAddress string) error {
	if err := s.blockingRepo.RemoveFromBlacklist(ctx, ipAddress); err != nil {
		return fmt.Errorf("failed to remove IP from blacklist: %w", err)
	}
		// Update in-memory cache
	s.mutex.Lock()
	delete(s.ipBlacklist, ipAddress)
	s.mutex.Unlock()

	s.logger.Info("IP removed from blacklist", "ip_address", ipAddress)
	return nil
@}

func (s *AttackBlockingService) UpdateCloudIntelligence(ctx context.Context) error {
	if !s.config.EnableCloudIntelligence || s.cloudIntelligence == nil {
		return nil
	}

	s.logger.Info("Updating cloud intelligence data")

	// Update threat feeds (implementation would fetch from external sources)
	// Update reputation scores
	// This is a placeholder implementation

	s.mutex.Lock()
	s.cloudIntelligence.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.logger.Info("Cloud intelligence data updated")
	return nil
}

func (s *AttackBlockingService) TestBlockingRule(ctx context.Context, ruleID string, testRequest *models.AttackBlockingRequest) (*models.RuleTestResult, error) {
	s.mutex.RLock()
	rule, exists := s.blockingRules[ruleID]
	s.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("blocking rule not found: %s", ruleID)
	}

	result := &models.RuleTestResult{
		RuleID:    ruleID,
		RuleName:  rule.Name,
		TestTime:  time.Now(),
		Matched:   s.matchBlockingRule(testRequest, rule),
		Details:   make(map[string]interface{}),
	}

	// Test each condition
	for i, condition := range rule.Conditions {
		conditionResult := s.evaluateCondition(testRequest, condition)
		result.Details[fmt.Sprintf("condition_%d", i)] = map[string]interface{}{
			"field":    condition.Field,
			"operator": condition.Operator,
			"value":    condition.Value,
			"matched":  conditionResult,
		}
	}

	return result, nil
}

func (s *AttackBlockingService) GetBlockingHealth(ctx context.Context) (*models.BlockingHealthStatus, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	health := &models.BlockingHealthStatus{
		Status:           "healthy",
		ActiveBlocks:     len(s.activeBlocks),
		BlockingRules:    len(s.blockingRules),
		WhitelistEntries: len(s.ipWhitelist),
		BlacklistEntries: len(s.ipBlacklist),
		LastUpdated:      time.Now(),
		Components:       make(map[string]string),
	}

	// Check component health
	health.Components["rate_limiting"] = s.getRateLimitingHealth()
	health.Components["signature_detection"] = s.getSignatureDetectionHealth()
	health.Components["anomaly_detection"] = s.getAnomalyDetectionHealth()
	health.Components["cloud_intelligence"] = s.getCloudIntelligenceHealth()

	// Determine overall health
	for _, componentHealth := range health.Components {
		if componentHealth != "healthy" {
			health.Status = "degraded"
			break
		}
	}

	return health, nil
}

func (s *AttackBlockingService) getRateLimitingHealth() string {
	if !s.config.EnableRateLimiting {
		return "disabled"
	}
	return "healthy"
}

func (s *AttackBlockingService) getSignatureDetectionHealth() string {
	if !s.config.EnableSignatureDetection {
		return "disabled"
	}
	return "healthy"
}

func (s *AttackBlockingService) getAnomalyDetectionHealth() string {
	if !s.config.EnableAnomalyDetection {
		return "disabled"
	}
	return "healthy"
}

func (s *AttackBlockingService) getCloudIntelligenceHealth() string {
	if !s.config.EnableCloudIntelligence {
		return "disabled"
	}
	if s.cloudIntelligence != nil && time.Since(s.cloudIntelligence.LastUpdated) < time.Hour {
		return "healthy"
	}
	return "stale"
}

// Cleanup expired blocks periodically
func (s *AttackBlockingService) StartCleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 5) // Cleanup every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.cleanupExpiredBlocks(ctx)
		}
	}
}

func (s *AttackBlockingService) cleanupExpiredBlocks(ctx context.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	var expiredIPs []string

	for ip, block := range s.activeBlocks {
		if now.After(block.ExpiresAt) {
			expiredIPs = append(expiredIPs, ip)
		}
	}

	for _, ip := range expiredIPs {
		block := s.activeBlocks[ip]
		block.Active = false
		
		// Update in repository
		if err := s.blockingRepo.UpdateActiveBlock(ctx, block); err != nil {
			s.logger.Error("Failed to update expired block", "error", err, "ip", ip)
		}

		delete(s.activeBlocks, ip)
		s.logger.Info("Expired block cleaned up", "ip_address", ip, "block_id", block.ID)
	}

	if len(expiredIPs) > 0 {
		s.logger.Info("Cleaned up expired blocks", "count", len(expiredIPs))
	}
}
