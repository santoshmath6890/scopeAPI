package services
import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/attack-blocking/internal/models"
	"scopeapi.local/backend/services/attack-blocking/internal/repository"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

// CloudIntelligenceService provides cloud-based threat intelligence and adaptive security
type CloudIntelligenceService struct {
	repo                    repository.ThreatIntelligenceRepository
	logger                  logging.Logger
	kafkaProducer          kafka.Producer
	kafkaConsumer          kafka.Consumer
	threatFeeds            map[string]*models.ThreatFeed
	intelligenceCache      map[string]*models.ThreatIntelligence
	adaptiveRules          map[string]*models.AdaptiveRule
	cloudProviders         map[string]CloudProvider
	mlModelClient          MLModelClient
	mutex                  sync.RWMutex
	updateInterval         time.Duration
	cacheExpiration        time.Duration
	ctx                    context.Context
	cancel                 context.CancelFunc
	feedUpdateTicker       *time.Ticker
	ruleEvaluationTicker   *time.Ticker
}

// CloudProvider interface for different cloud service integrations
type CloudProvider interface {
	GetThreatIntelligence(ctx context.Context, query *models.ThreatQuery) (*models.CloudThreatData, error)
	ReportThreat(ctx context.Context, threat *models.ThreatReport) error
	GetReputation(ctx context.Context, indicator *models.ThreatIndicator) (*models.ReputationScore, error)
	SubmitSample(ctx context.Context, sample *models.ThreatSample) (*models.AnalysisResult, error)
}

// MLModelClient interface for machine learning model interactions
type MLModelClient interface {
	PredictThreatScore(ctx context.Context, features *models.ThreatFeatures) (*models.ThreatPrediction, error)
	UpdateModel(ctx context.Context, trainingData *models.TrainingData) error
	GetModelMetrics(ctx context.Context, modelID string) (*models.ModelMetrics, error)
}

type CloudIntelligenceServiceConfig struct {
	UpdateInterval    time.Duration
	CacheExpiration   time.Duration
	MaxCacheSize      int
	ThreatFeedURLs    []string
	CloudProviders    map[string]CloudProviderConfig
	MLModelEndpoint   string
	EnableAdaptiveRules bool
}

type CloudProviderConfig struct {
	Name        string
	APIKey      string
	Endpoint    string
	RateLimit   int
	Timeout     time.Duration
	Enabled     bool
}

func NewCloudIntelligenceService(
	repo repository.ThreatIntelligenceRepository,
	logger logging.Logger,
	kafkaProducer kafka.Producer,
	kafkaConsumer kafka.Consumer,
	mlModelClient MLModelClient,
	config *CloudIntelligenceServiceConfig,
) *CloudIntelligenceService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &CloudIntelligenceService{
		repo:              repo,
		logger:            logger,
		kafkaProducer:     kafkaProducer,
		kafkaConsumer:     kafkaConsumer,
		mlModelClient:     mlModelClient,
		threatFeeds:       make(map[string]*models.ThreatFeed),
		intelligenceCache: make(map[string]*models.ThreatIntelligence),
		adaptiveRules:     make(map[string]*models.AdaptiveRule),
		cloudProviders:    make(map[string]CloudProvider),
		updateInterval:    config.UpdateInterval,
		cacheExpiration:   config.CacheExpiration,
		ctx:               ctx,
		cancel:            cancel,
	}

	// Initialize cloud providers
	service.initializeCloudProviders(config.CloudProviders)

	// Initialize threat feeds
	service.initializeThreatFeeds(config.ThreatFeedURLs)

	// Start background tasks
	service.startFeedUpdates()
	service.startRuleEvaluation()
	service.startEventConsumer()

	return service
}

// Threat Intelligence Operations

func (s *CloudIntelligenceService) GetThreatIntelligence(ctx context.Context, indicator *models.ThreatIndicator) (*models.ThreatIntelligence, error) {
	s.logger.Debug("Getting threat intelligence", "indicator", indicator.Value, "type", indicator.Type)

	// Check cache first
	cacheKey := s.generateCacheKey(indicator)
	s.mutex.RLock()
	if cached, exists := s.intelligenceCache[cacheKey]; exists {
		if time.Since(cached.LastUpdated) < s.cacheExpiration {
			s.mutex.RUnlock()
			return cached, nil
		}
	}
	s.mutex.RUnlock()

	// Query multiple cloud providers
	intelligence := &models.ThreatIntelligence{
		ID:           uuid.New().String(),
		Indicator:    *indicator,
		Sources:      make([]*models.IntelligenceSource, 0),
		Reputation:   &models.ReputationScore{},
		LastUpdated:  time.Now(),
		Confidence:   0,
		ThreatTypes:  make([]string, 0),
		Attributes:   make(map[string]interface{}),
	}

	// Collect intelligence from all providers
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for providerName, provider := range s.cloudProviders {
		wg.Add(1)
		go func(name string, p CloudProvider) {
			defer wg.Done()
			
			query := &models.ThreatQuery{
				Indicator: *indicator,
				MaxResults: 10,
				IncludeContext: true,
			}
			
			cloudData, err := p.GetThreatIntelligence(ctx, query)
			if err != nil {
				s.logger.Error("Failed to get intelligence from provider", 
					"provider", name, 
					"error", err)
				return
			}
			
			mu.Lock()
			source := &models.IntelligenceSource{
				Provider:    name,
				Data:        cloudData,
				Confidence:  cloudData.Confidence,
				LastUpdated: time.Now(),
			}
			intelligence.Sources = append(intelligence.Sources, source)
			
			// Update reputation score
			if cloudData.Reputation != nil {
				intelligence.Reputation.Score += cloudData.Reputation.Score
				intelligence.Reputation.Sources = append(intelligence.Reputation.Sources, name)
			}
			
			// Merge threat types
			for _, threatType := range cloudData.ThreatTypes {
				found := false
				for _, existing := range intelligence.ThreatTypes {
					if existing == threatType {
						found = true
						break
					}
				}
				if !found {
					intelligence.ThreatTypes = append(intelligence.ThreatTypes, threatType)
				}
			}
			
			// Merge attributes
			for key, value := range cloudData.Attributes {
				intelligence.Attributes[key] = value
			}
			mu.Unlock()
		}(providerName, provider)
	}
	
	wg.Wait()

	// Calculate overall confidence and reputation
	if len(intelligence.Sources) > 0 {
		totalConfidence := 0.0
		for _, source := range intelligence.Sources {
			totalConfidence += source.Confidence
		}
		intelligence.Confidence = totalConfidence / float64(len(intelligence.Sources))
		
		if len(intelligence.Reputation.Sources) > 0 {
			intelligence.Reputation.Score = intelligence.Reputation.Score / float64(len(intelligence.Reputation.Sources))
		}
	}

	// Enhance with ML predictions
	if err := s.enhanceWithMLPredictions(ctx, intelligence); err != nil {
		s.logger.Error("Failed to enhance with ML predictions", "error", err)
	}

	// Store in repository
	if err := s.repo.CreateThreatIntelligence(ctx, intelligence); err != nil {
		s.logger.Error("Failed to store threat intelligence", "error", err)
	}

	// Cache the result
	s.mutex.Lock()
	s.intelligenceCache[cacheKey] = intelligence
	s.mutex.Unlock()

	s.logger.Info("Threat intelligence retrieved", 
		"indicator", indicator.Value,
		"confidence", intelligence.Confidence,
		"sources", len(intelligence.Sources))

	return intelligence, nil
}

func (s *CloudIntelligenceService) ReportThreat(ctx context.Context, threat *models.ThreatReport) error {
	s.logger.Info("Reporting threat to cloud providers", "threat_id", threat.ID)

	// Validate threat report
	if err := s.validateThreatReport(threat); err != nil {
		return fmt.Errorf("threat report validation failed: %w", err)
	}

	// Report to all enabled cloud providers
	var wg sync.WaitGroup
	errors := make([]error, 0)
	var mu sync.Mutex

	for providerName, provider := range s.cloudProviders {
		wg.Add(1)
		go func(name string, p CloudProvider) {
			defer wg.Done()
			
			if err := p.ReportThreat(ctx, threat); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("provider %s: %w", name, err))
				mu.Unlock()
				s.logger.Error("Failed to report threat to provider", 
					"provider", name, 
					"error", err)
			} else {
				s.logger.Debug("Threat reported successfully", "provider", name)
			}
		}(providerName, provider)
	}

	wg.Wait()

	// Store threat report
	threat.ReportedAt = time.Now()
	if err := s.repo.CreateThreatReport(ctx, threat); err != nil {
		return fmt.Errorf("failed to store threat report: %w", err)
	}

	// Publish threat report event
	if err := s.publishThreatReportEvent(ctx, threat); err != nil {
		s.logger.Error("Failed to publish threat report event", "error", err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("some providers failed: %v", errors)
	}

	return nil
}

func (s *CloudIntelligenceService) GetReputationScore(ctx context.Context, indicator *models.ThreatIndicator) (*models.ReputationScore, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("reputation:%s:%s", indicator.Type, indicator.Value)
	s.mutex.RLock()
	if cached, exists := s.intelligenceCache[cacheKey]; exists {
		if time.Since(cached.LastUpdated) < s.cacheExpiration {
			s.mutex.RUnlock()
			return cached.Reputation, nil
		}
	}
	s.mutex.RUnlock()

	// Query providers for reputation
	reputation := &models.ReputationScore{
		Indicator:   *indicator,
		Score:       0,
		Sources:     make([]string, 0),
		LastUpdated: time.Now(),
		Categories:  make(map[string]float64),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	totalScore := 0.0
	sourceCount := 0

	for providerName, provider := range s.cloudProviders {
		wg.Add(1)
		go func(name string, p CloudProvider) {
			defer wg.Done()
			
			providerReputation, err := p.GetReputation(ctx, indicator)
			if err != nil {
				s.logger.Error("Failed to get reputation from provider", 
					"provider", name, 
					"error", err)
				return
			}
			
			mu.Lock()
			totalScore += providerReputation.Score
			sourceCount++
			reputation.Sources = append(reputation.Sources, name)
			
			// Merge categories
			for category, score := range providerReputation.Categories {
				if existing, exists := reputation.Categories[category]; exists {
					reputation.Categories[category] = (existing + score) / 2
				} else {
					reputation.Categories[category] = score
				}
			}
			mu.Unlock()
		}(providerName, provider)
	}

	wg.Wait()

	if sourceCount > 0 {
		reputation.Score = totalScore / float64(sourceCount)
	}

	// Store in cache
	intelligence := &models.ThreatIntelligence{
		Indicator:   *indicator,
		Reputation:  reputation,
		LastUpdated: time.Now(),
	}

	s.mutex.Lock()
	s.intelligenceCache[cacheKey] = intelligence
	s.mutex.Unlock()

	return reputation, nil
}

// Adaptive Rules Management

func (s *CloudIntelligenceService) CreateAdaptiveRule(ctx context.Context, rule *models.AdaptiveRule) error {
	// Validate rule
	if err := s.validateAdaptiveRule(rule); err != nil {
		return fmt.Errorf("rule validation failed: %w", err)
	}

	rule.ID = uuid.New().String()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	rule.Status = models.RuleStatusActive

	// Store in repository
	if err := s.repo.CreateAdaptiveRule(ctx, rule); err != nil {
		return fmt.Errorf("failed to create adaptive rule: %w", err)
	}

	// Add to memory
	s.mutex.Lock()
	s.adaptiveRules[rule.ID] = rule
	s.mutex.Unlock()

	// Publish rule creation event
	if err := s.publishRuleEvent(ctx, "rule_created", rule); err != nil {
		s.logger.Error("Failed to publish rule creation event", "error", err)
	}

	s.logger.Info("Adaptive rule created", "rule_id", rule.ID, "name", rule.Name)
	return nil
}

func (s *CloudIntelligenceService) UpdateAdaptiveRule(ctx context.Context, rule *models.AdaptiveRule) error {
	// Validate rule
	if err := s.validateAdaptiveRule(rule); err != nil {
		return fmt.Errorf("rule validation failed: %w", err)
	}

	rule.UpdatedAt = time.Now()

	// Update in repository
	if err := s.repo.UpdateAdaptiveRule(ctx, rule); err != nil {
		return fmt.Errorf("failed to update adaptive rule: %w", err)
	}

	// Update in memory
	s.mutex.Lock()
	s.adaptiveRules[rule.ID] = rule
	s.mutex.Unlock()

	// Publish rule update event
	if err := s.publishRuleEvent(ctx, "rule_updated", rule); err != nil {
		s.logger.Error("Failed to publish rule update event", "error", err)
	}

	s.logger.Info("Adaptive rule updated", "rule_id", rule.ID)
	return nil
}

func (s *CloudIntelligenceService) EvaluateAdaptiveRules(ctx context.Context, event *models.SecurityEvent) ([]*models.RuleMatch, error) {
	matches := make([]*models.RuleMatch, 0)

	s.mutex.RLock()
	rules := make([]*models.AdaptiveRule, 0, len(s.adaptiveRules))
	for _, rule := range s.adaptiveRules {
		if rule.Status == models.RuleStatusActive {
			rules = append(rules, rule)
		}
	}
	s.mutex.RUnlock()

	// Evaluate each rule
	for _, rule := range rules {
		match, err := s.evaluateRule(ctx, rule, event)
		if err != nil {
			s.logger.Error("Failed to evaluate rule", 
				"rule_id", rule.ID, 
				"error", err)
			continue
		}

		if match != nil {
			matches = append(matches, match)
			
			// Update rule statistics
			rule.MatchCount++
			rule.LastMatched = time.Now()
			
			// Trigger rule actions
			if err := s.executeRuleActions(ctx, rule, match, event); err != nil {
				s.logger.Error("Failed to execute rule actions", 
					"rule_id", rule.ID, 
					"error", err)
			}
		}
	}

	return matches, nil
}

func (s *CloudIntelligenceService) GetAdaptiveRules(ctx context.Context, filter *models.AdaptiveRuleFilter) ([]*models.AdaptiveRule, error) {
	return s.repo.GetAdaptiveRules(ctx, filter)
}

// Threat Feed Management

func (s *CloudIntelligenceService) UpdateThreatFeeds(ctx context.Context) error {
	s.logger.Info("Updating threat feeds")

	var wg sync.WaitGroup
	errors := make([]error, 0)
	var mu sync.Mutex

	s.mutex.RLock()
	feeds := make([]*models.ThreatFeed, 0, len(s.threatFeeds))
	for _, feed := range s.threatFeeds {
		feeds = append(feeds, feed)
	}
	s.mutex.RUnlock()

	for _, feed := range feeds {
		wg.Add(1)
		go func(f *models.ThreatFeed) {
			defer wg.Done()
			
			if err := s.updateThreatFeed(ctx, f); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("feed %s: %w", f.Name, err))
				mu.Unlock()
				s.logger.Error("Failed to update threat feed", 
					"feed", f.Name, 
					"error", err)
			}
		}(feed)
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("some feeds failed to update: %v", errors)
	}

	s.logger.Info("Threat feeds updated successfully")
	return nil
}

func (s *CloudIntelligenceService) AddThreatFeed(ctx context.Context, feed *models.ThreatFeed) error {
	// Validate feed
	if err := s.validateThreatFeed(feed); err != nil {
		return fmt.Errorf("feed validation failed: %w", err)
	}

	feed.ID = uuid.New().String()
	feed.CreatedAt = time.Now()
	feed.UpdatedAt = time.Now()
	feed.Status = models.FeedStatusActive

	// Store in repository
	if err := s.repo.CreateThreatFeed(ctx, feed); err != nil {
		return fmt.Errorf("failed to create threat feed: %w", err)
	}

	// Add to memory
	s.mutex.Lock()
	s.threatFeeds[feed.ID] = feed
	s.mutex.Unlock()

	s.logger.Info("Threat feed added", "feed_id", feed.ID, "name", feed.Name)
	return nil
}

// ML Model Integration

func (s *CloudIntelligenceService) PredictThreatScore(ctx context.Context, features *models.ThreatFeatures) (*models.ThreatPrediction, error) {
	prediction, err := s.mlModelClient.PredictThreatScore(ctx, features)
	if err != nil {
		return nil, fmt.Errorf("ML prediction failed: %w", err)
	}

	// Store prediction for model improvement
	predictionRecord := &models.ThreatPredictionRecord{
		ID:          uuid.New().String(),
		Features:    *features,
		Prediction:  *prediction,
		Timestamp:   time.Now(),
		ModelVersion: prediction.ModelVersion,
	}

	if err := s.repo.CreateThreatPrediction(ctx, predictionRecord); err != nil {
		s.logger.Error("Failed to store prediction record", "error", err)
	}

	return prediction, nil
}

func (s *CloudIntelligenceService) UpdateMLModel(ctx context.Context, trainingData *models.TrainingData) error {
	s.logger.Info("Updating ML model with new training data")

	if err := s.mlModelClient.UpdateModel(ctx, trainingData); err != nil {
		return fmt.Errorf("failed to update ML model: %w", err)
	}

	// Record model update
	updateRecord := &models.ModelUpdateRecord{
		ID:            uuid.New().String(),
		ModelType:     trainingData.ModelType,
		TrainingSize:  trainingData.SampleCount,
		UpdatedAt:     time.Now(),
		Version:       trainingData.Version,
		Metrics:       trainingData.ValidationMetrics,
	}

	if err := s.repo.CreateModelUpdate(ctx, updateRecord); err != nil {
		s.logger.Error("Failed to record model update", "error", err)
	}

	s.logger.Info("ML model updated successfully")
	return nil
}

// Analytics and Reporting

func (s *CloudIntelligenceService) GetThreatIntelligenceAnalytics(ctx context.Context, filter *models.AnalyticsFilter) (*models.ThreatIntelligenceAnalytics, error) {
	analytics, err := s.repo.GetThreatIntelligenceAnalytics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}

	// Enhance with real-time data
	s.mutex.RLock
	analytics.CachedIntelligenceCount = len(s.intelligenceCache)
	analytics.ActiveRulesCount = 0
	for _, rule := range s.adaptiveRules {
		if rule.Status == models.RuleStatusActive {
			analytics.ActiveRulesCount++
		}
	}
	analytics.ActiveFeedsCount = 0
	for _, feed := range s.threatFeeds {
		if feed.Status == models.FeedStatusActive {
			analytics.ActiveFeedsCount++
		}
	}
	s.mutex.RUnlock()

	analytics.LastUpdated = time.Now()
	return analytics, nil
}

func (s *CloudIntelligenceService) GenerateIntelligenceReport(ctx context.Context, request *models.ReportRequest) (*models.IntelligenceReport, error) {
	report := &models.IntelligenceReport{
		ID:          uuid.New().String(),
		Type:        request.Type,
		Period:      request.Period,
		GeneratedAt: time.Now(),
		Sections:    make([]*models.ReportSection, 0),
	}

	// Generate threat intelligence summary
	if request.IncludeThreatSummary {
		summary, err := s.generateThreatSummary(ctx, request.Period)
		if err != nil {
			return nil, fmt.Errorf("failed to generate threat summary: %w", err)
		}
		report.Sections = append(report.Sections, summary)
	}

	// Generate adaptive rules performance
	if request.IncludeRulesPerformance {
		rulesSection, err := s.generateRulesPerformance(ctx, request.Period)
		if err != nil {
			return nil, fmt.Errorf("failed to generate rules performance: %w", err)
		}
		report.Sections = append(report.Sections, rulesSection)
	}

	// Generate feed statistics
	if request.IncludeFeedStats {
		feedStats, err := s.generateFeedStatistics(ctx, request.Period)
		if err != nil {
			return nil, fmt.Errorf("failed to generate feed statistics: %w", err)
		}
		report.Sections = append(report.Sections, feedStats)
	}

	// Generate ML model performance
	if request.IncludeMLPerformance {
		mlSection, err := s.generateMLPerformance(ctx, request.Period)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ML performance: %w", err)
		}
		report.Sections = append(report.Sections, mlSection)
	}

	// Store report
	if err := s.repo.CreateIntelligenceReport(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to store report: %w", err)
	}

	return report, nil
}

// Private helper methods

func (s *CloudIntelligenceService) initializeCloudProviders(configs map[string]CloudProviderConfig) {
	for name, config := range configs {
		if !config.Enabled {
			continue
		}

		var provider CloudProvider
		switch name {
		case "virustotal":
			provider = NewVirusTotalProvider(config)
		case "threatcrowd":
			provider = NewThreatCrowdProvider(config)
		case "alienvault":
			provider = NewAlienVaultProvider(config)
		case "malwarebazaar":
			provider = NewMalwareBazaarProvider(config)
		default:
			s.logger.Warn("Unknown cloud provider", "provider", name)
			continue
		}

		s.cloudProviders[name] = provider
		s.logger.Info("Cloud provider initialized", "provider", name)
	}
}

func (s *CloudIntelligenceService) initializeThreatFeeds(feedURLs []string) {
	for i, url := range feedURLs {
		feed := &models.ThreatFeed{
			ID:          uuid.New().String(),
			Name:        fmt.Sprintf("Feed-%d", i+1),
			URL:         url,
			Type:        models.FeedTypeIOC,
			Format:      models.FeedFormatJSON,
			Status:      models.FeedStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		s.threatFeeds[feed.ID] = feed
	}
}

func (s *CloudIntelligenceService) startFeedUpdates() {
	s.feedUpdateTicker = time.NewTicker(s.updateInterval)
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				s.feedUpdateTicker.Stop()
				return
			case <-s.feedUpdateTicker.C:
				if err := s.UpdateThreatFeeds(s.ctx); err != nil {
					s.logger.Error("Failed to update threat feeds", "error", err)
				}
			}
		}
	}()
}

func (s *CloudIntelligenceService) startRuleEvaluation() {
	s.ruleEvaluationTicker = time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				s.ruleEvaluationTicker.Stop()
				return
			case <-s.ruleEvaluationTicker.C:
				s.performRuleMaintenance()
			}
		}
	}()
}

func (s *CloudIntelligenceService) startEventConsumer() {
	go func() {
		topics := []string{"security-events", "threat-detections", "attack-attempts"}
		
		for _, topic := range topics {
			go func(t string) {
				if err := s.kafkaConsumer.Subscribe(t, s.handleSecurityEvent); err != nil {
					s.logger.Error("Failed to subscribe to topic", "topic", t, "error", err)
				}
			}(topic)
		}
	}()
}

func (s *CloudIntelligenceService) handleSecurityEvent(message []byte) error {
	var event models.SecurityEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal security event: %w", err)
	}

	// Evaluate adaptive rules
	matches, err := s.EvaluateAdaptiveRules(context.Background(), &event)
	if err != nil {
		s.logger.Error("Failed to evaluate adaptive rules", "error", err)
		return err
	}

	// Process matches
	for _, match := range matches {
		s.logger.Info("Adaptive rule matched", 
			"rule_id", match.RuleID,
			"event_id", event.ID,
			"confidence", match.Confidence)
	}

	return nil
}

func (s *CloudIntelligenceService) generateCacheKey(indicator *models.ThreatIndicator) string {
	return fmt.Sprintf("%s:%s", indicator.Type, indicator.Value)
}

func (s *CloudIntelligenceService) enhanceWithMLPredictions(ctx context.Context, intelligence *models.ThreatIntelligence) error {
	features := &models.ThreatFeatures{
		IndicatorType:    intelligence.Indicator.Type,
		IndicatorValue:   intelligence.Indicator.Value,
		SourceCount:      len(intelligence.Sources),
		ReputationScore:  intelligence.Reputation.Score,
		ThreatTypes:      intelligence.ThreatTypes,
		Attributes:       intelligence.Attributes,
	}

	prediction, err := s.PredictThreatScore(ctx, features)
	if err != nil {
		return err
	}

	intelligence.MLPrediction = prediction
	intelligence.Confidence = (intelligence.Confidence + prediction.Confidence) / 2

	return nil
}

func (s *CloudIntelligenceService) validateThreatReport(report *models.ThreatReport) error {
	if report.ThreatType == "" {
		return fmt.Errorf("threat type is required")
	}

	if len(report.Indicators) == 0 {
		return fmt.Errorf("at least one indicator is required")
	}

	if report.Severity == "" {
		return fmt.Errorf("severity is required")
	}

	return nil
}

func (s *CloudIntelligenceService) validateAdaptiveRule(rule *models.AdaptiveRule) error {
	if rule.Name == "" {
		return fmt.Errorf("rule name is required")
	}

	if len(rule.Conditions) == 0 {
		return fmt.Errorf("at least one condition is required")
	}

	if len(rule.Actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}

	// Validate conditions
	for _, condition := range rule.Conditions {
		if condition.Field == "" {
			return fmt.Errorf("condition field is required")
		}
		if condition.Operator == "" {
			return fmt.Errorf("condition operator is required")
		}
	}

	// Validate actions
	for _, action := range rule.Actions {
		if action.Type == "" {
			return fmt.Errorf("action type is required")
		}
	}

	return nil
}

func (s *CloudIntelligenceService) validateThreatFeed(feed *models.ThreatFeed) error {
	if feed.Name == "" {
		return fmt.Errorf("feed name is required")
	}

	if feed.URL == "" {
		return fmt.Errorf("feed URL is required")
	}

	if feed.Type == "" {
		return fmt.Errorf("feed type is required")
	}

	return nil
}

func (s *CloudIntelligenceService) evaluateRule(ctx context.Context, rule *models.AdaptiveRule, event *models.SecurityEvent) (*models.RuleMatch, error) {
	// Check if all conditions are met
	conditionsMet := true
	matchedConditions := make([]*models.ConditionMatch, 0)

	for _, condition := range rule.Conditions {
		match, err := s.evaluateCondition(condition, event)
		if err != nil {
			return nil, err
		}

		if match == nil {
			conditionsMet = false
			break
		}

		matchedConditions = append(matchedConditions, match)
	}

	if !conditionsMet {
		return nil, nil
	}

	// Calculate confidence based on condition matches
	totalConfidence := 0.0
	for _, match := range matchedConditions {
		totalConfidence += match.Confidence
	}
	confidence := totalConfidence / float64(len(matchedConditions))

	return &models.RuleMatch{
		ID:                uuid.New().String(),
		RuleID:            rule.ID,
		EventID:           event.ID,
		Confidence:        confidence,
		MatchedConditions: matchedConditions,
		Timestamp:         time.Now(),
	}, nil
}

func (s *CloudIntelligenceService) evaluateCondition(condition *models.RuleCondition, event *models.SecurityEvent) (*models.ConditionMatch, error) {
	// Extract field value from event
	fieldValue, err := s.extractFieldValue(condition.Field, event)
	if err != nil {
		return nil, err
	}

	// Evaluate condition based on operator
	matched := false
	confidence := 0.0

	switch condition.Operator {
	case "equals":
		matched = fieldValue == condition.Value
		confidence = 1.0
	case "contains":
		if str, ok := fieldValue.(string); ok {
			matched = strings.Contains(str, condition.Value.(string))
			confidence = 0.8
		}
	case "regex":
		if str, ok := fieldValue.(string); ok {
			regex, err := regexp.Compile(condition.Value.(string))
			if err != nil {
				return nil, err
			}
			matched = regex.MatchString(str)
			confidence = 0.9
		}
	case "greater_than":
		if num, ok := fieldValue.(float64); ok {
			if conditionNum, ok := condition.Value.(float64); ok {
				matched = num > conditionNum
				confidence = 0.7
			}
		}
	case "less_than":
		if num, ok := fieldValue.(float64); ok {
			if conditionNum, ok := condition.Value.(float64); ok {
				matched = num < conditionNum
				confidence = 0.7
			}
		}
	default:
		return nil, fmt.Errorf("unsupported operator: %s", condition.Operator)
	}

	if !matched {
		return nil, nil
	}

	return &models.ConditionMatch{
		Condition:  *condition,
		Value:      fieldValue,
		Confidence: confidence,
		Timestamp:  time.Now(),
	}, nil
}

func (s *CloudIntelligenceService) extractFieldValue(field string, event *models.SecurityEvent) (interface{}, error) {
	switch field {
	case "source_ip":
		return event.SourceIP, nil
	case "destination_ip":
		return event.DestinationIP, nil
	case "threat_type":
		return event.ThreatType, nil
	case "severity":
		return event.Severity, nil
	case "confidence":
		return event.Confidence, nil
	case "user_agent":
		return event.UserAgent, nil
	case "request_method":
		return event.RequestMethod, nil
	case "request_path":
		return event.RequestPath, nil
	default:
		// Check custom attributes
		if value, exists := event.Attributes[field]; exists {
			return value, nil
		}
		return nil, fmt.Errorf("unknown field: %s", field)
	}
}

func (s *CloudIntelligenceService) executeRuleActions(ctx context.Context, rule *models.AdaptiveRule, match *models.RuleMatch, event *models.SecurityEvent) error {
	for _, action := range rule.Actions {
		switch action.Type {
		case models.ActionTypeBlock:
			if err := s.executeBlockAction(ctx, action, match, event); err != nil {
				return err
			}
		case models.ActionTypeAlert:
			if err := s.executeAlertAction(ctx, action, match, event); err != nil {
				return err
			}
		case models.ActionTypeLog:
			if err := s.executeLogAction(ctx, action, match, event); err != nil {
				return err
			}
		case models.ActionTypeUpdateRule:
			if err := s.executeUpdateRuleAction(ctx, action, match, event); err != nil {
				return err
			}
		default:
			s.logger.Warn("Unknown action type", "type", action.Type)
		}
	}

	return nil
}

func (s *CloudIntelligenceService) executeBlockAction(ctx context.Context, action *models.RuleAction, match *models.RuleMatch, event *models.SecurityEvent) error {
	blockRequest := &models.BlockRequest{
		ID:        uuid.New().String(),
		SourceIP:  event.SourceIP,
		Reason:    fmt.Sprintf("Adaptive rule match: %s", match.ID),
		Duration:  time.Duration(action.Parameters["duration"].(float64)) * time.Second,
		RuleID:    match.RuleID,
		Timestamp: time.Now(),
	}

	// Publish block request event
	eventData, _ := json.Marshal(blockRequest)
	if err := s.kafkaProducer.Publish("block-requests", eventData); err != nil {
		return fmt.Errorf("failed to publish block request: %w", err)
	}

	s.logger.Info("Block action executed", "rule_id", match.RuleID, "source_ip", event.SourceIP)
	return nil
}

func (s *CloudIntelligenceService) executeAlertAction(ctx context.Context, action *models.RuleAction, match *models.RuleMatch, event *models.SecurityEvent) error {
	alert := &models.SecurityAlert{
		ID:          uuid.New().String(),
		Type:        "adaptive_rule_match",
		Severity:    action.Parameters["severity"].(string),
		Title:       fmt.Sprintf("Adaptive Rule Triggered: %s", match.RuleID),
		Description: fmt.Sprintf("Rule %s matched event %s with confidence %.2f", match.RuleID, event.ID, match.Confidence),
		SourceIP:    event.SourceIP,
		RuleID:      match.RuleID,
		EventID:     event.ID,
		Timestamp:   time.Now(),
		Metadata: map[string]interface{}{
			"rule_match": match,
			"event":      event,
		},
	}

	// Publish alert event
	eventData, _ := json.Marshal(alert)
	if err := s.kafkaProducer.Publish("security-alerts", eventData); err != nil {
		return fmt.Errorf("failed to publish alert: %w", err)
	}

	s.logger.Info("Alert action executed", "rule_id", match.RuleID, "alert_id", alert.ID)
	return nil
}

func (s *CloudIntelligenceService) executeLogAction(ctx context.Context, action *models.RuleAction, match *models.RuleMatch, event *models.SecurityEvent) error {
	logEntry := &models.SecurityLogEntry{
		ID:        uuid.New().String(),
		Level:     action.Parameters["level"].(string),
		Message:   fmt.Sprintf("Adaptive rule %s matched event %s", match.RuleID, event.ID),
		RuleID:    match.RuleID,
		EventID:   event.ID,
		Timestamp: time.Now(),
		Context: map[string]interface{}{
			"match":      match,
			"event":      event,
			"confidence": match.Confidence,
		},
	}

	// Store log entry
	if err := s.repo.CreateSecurityLog(ctx, logEntry); err != nil {
		return fmt.Errorf("failed to create log entry: %w", err)
	}

	s.logger.Info("Log action executed", "rule_id", match.RuleID, "log_id", logEntry.ID)
	return nil
}

func (s *CloudIntelligenceService) executeUpdateRuleAction(ctx context.Context, action *models.RuleAction, match *models.RuleMatch, event *models.SecurityEvent) error {
	// Get the rule to update
	s.mutex.RLock()
	rule, exists := s.adaptiveRules[match.RuleID]
	s.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("rule not found: %s", match.RuleID)
	}

	// Update rule based on action parameters
	if threshold, ok := action.Parameters["confidence_threshold"]; ok {
		rule.ConfidenceThreshold = threshold.(float64)
	}

	if priority, ok := action.Parameters["priority"]; ok {
		rule.Priority = int(priority.(float64))
	}

	// Update rule
	if err := s.UpdateAdaptiveRule(ctx, rule); err != nil {
		return fmt.Errorf("failed to update rule: %w", err)
	}

	s.logger.Info("Update rule action executed", "rule_id", match.RuleID)
	return nil
}

func (s *CloudIntelligenceService) updateThreatFeed(ctx context.Context, feed *models.ThreatFeed) error {
	s.logger.Debug("Updating threat feed", "feed", feed.Name, "url", feed.URL)

	// Fetch feed data
	feedData, err := s.fetchFeedData(ctx, feed)
	if err != nil {
		return fmt.Errorf("failed to fetch feed data: %w", err)
	}

	// Parse feed data
	indicators, err := s.parseFeedData(feedData, feed.Format)
	if err != nil {
		return fmt.Errorf("failed to parse feed data: %w", err)
	}

	// Update feed statistics
	feed.LastUpdated = time.Now()
	feed.IndicatorCount = len(indicators)
	feed.Status = models.FeedStatusActive

	// Store indicators
	for _, indicator := range indicators {
		indicator.FeedID = feed.ID
		indicator.Source = feed.Name
		indicator.LastSeen = time.Now()

		if err := s.repo.CreateThreatIndicator(ctx, indicator); err != nil {
			s.logger.Error("Failed to store threat indicator", 
				"indicator", indicator.Value, 
				"error", err)
		}
	}

	// Update feed in repository
	if err := s.repo.UpdateThreatFeed(ctx, feed); err != nil {
		return fmt.Errorf("failed to update feed: %w", err)
	}

	// Update in memory
	s.mutex.Lock()
	s.threatFeeds[feed.ID] = feed
	s.mutex.Unlock()

	s.logger.Info("Threat feed updated", 
		"feed", feed.Name, 
		"indicators", len(indicators))

	return nil
}

func (s *CloudIntelligenceService) fetchFeedData(ctx context.Context, feed *models.ThreatFeed) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feed.URL, nil)
	if err != nil {
		return nil, err
	}

	// Add authentication if configured
	if feed.AuthType == models.AuthTypeAPIKey && feed.APIKey != "" {
		req.Header.Set("X-API-Key", feed.APIKey)
	} else if feed.AuthType == models.AuthTypeBearer && feed.Token != "" {
		req.Header.Set("Authorization", "Bearer "+feed.Token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (s *CloudIntelligenceService) parseFeedData(data []byte, format models.FeedFormat) ([]*models.ThreatIndicator, error) {
	indicators := make([]*models.ThreatIndicator, 0)

	switch format {
	case models.FeedFormatJSON:
		var feedData models.JSONFeedData
		if err := json.Unmarshal(data, &feedData); err != nil {
			return nil, err
		}

		for _, item := range feedData.Indicators {
			indicator := &models.ThreatIndicator{
				ID:          uuid.New().String(),
				Type:        item.Type,
				Value:       item.Value,
				Confidence:  item.Confidence,
				ThreatTypes: item.ThreatTypes,
				Tags:        item.Tags,
				FirstSeen:   item.FirstSeen,
				LastSeen:    time.Now(),
				Attributes:  item.Attributes,
			}
			indicators = append(indicators, indicator)
		}

	case models.FeedFormatCSV:
		reader := csv.NewReader(strings.NewReader(string(data)))
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}

		// Skip header row
		for i, record := range records {
			if i == 0 {
				continue
			}

			if len(record) < 2 {
				continue
			}

			indicator := &models.ThreatIndicator{
				ID:       uuid.New().String(),
				Type:     record[0],
				Value:    record[1],
				LastSeen: time.Now(),
			}

			if len(record) > 2 {
				if confidence, err := strconv.ParseFloat(record[2], 64); err == nil {
					indicator.Confidence = confidence
				}
			}

			indicators = append(indicators, indicator)
		}

	case models.FeedFormatSTIX:
		// Parse STIX format (simplified)
		var stixData models.STIXBundle
		if err := json.Unmarshal(data, &stixData); err != nil {
			return nil, err
		}

		for _, object := range stixData.Objects {
			if object.Type == "indicator" {
				indicator := &models.ThreatIndicator{
					ID:          uuid.New().String(),
					Type:        object.Pattern.Type,
					Value:       object.Pattern.Value,
					Confidence:  object.Confidence,
					ThreatTypes: object.Labels,
					FirstSeen:   object.Created,
					LastSeen:    time.Now(),
				}
				indicators = append(indicators, indicator)
			}
		}

	default:
		return nil, fmt.Errorf("unsupported feed format: %s", format)
	}

	return indicators, nil
}

func (s *CloudIntelligenceService) performRuleMaintenance() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for _, rule := range s.adaptiveRules {
		// Disable rules that haven't matched in a long time
		if rule.Status == models.RuleStatusActive && 
		   !rule.LastMatched.IsZero() && 
		   now.Sub(rule.LastMatched) > 30*24*time.Hour {
			rule.Status = models.RuleStatusInactive
			rule.UpdatedAt = now
			
			// Update in repository
			if err := s.repo.UpdateAdaptiveRule(context.Background(), rule); err != nil {
				s.logger.Error("Failed to update rule status", 
					"rule_id", rule.ID, 
					"error", err)
			}
		}

		// Update rule effectiveness score
		if rule.MatchCount > 0 {
			rule.EffectivenessScore = float64(rule.TruePositives) / float64(rule.MatchCount)
		}
	}
}

func (s *CloudIntelligenceService) publishThreatReportEvent(ctx context.Context, report *models.ThreatReport) error {
	event := &models.ThreatReportEvent{
		ID:        uuid.New().String(),
		ReportID:  report.ID,
		Type:      "threat_reported",
		Timestamp: time.Now(),
		Data:      report,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.kafkaProducer.Publish("threat-reports", eventData)
}

func (s *CloudIntelligenceService) publishRuleEvent(ctx context.Context, eventType string, rule *models.AdaptiveRule) error {
	event := &models.RuleEvent{
		ID:        uuid.New().String(),
		RuleID:    rule.ID,
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      rule,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.kafkaProducer.Publish("adaptive-rules", eventData)
}

func (s *CloudIntelligenceService) generateThreatSummary(ctx context.Context, period *models.TimePeriod) (*models.ReportSection, error) {
	summary, err := s.repo.GetThreatSummary(ctx, period)
	if err != nil {
		return nil, err
	}

	return &models.ReportSection{
		Title: "Threat Intelligence Summary",
		Type:  "threat_summary",
		Data:  summary,
	}, nil
}

func (s *CloudIntelligenceService) generateRulesPerformance(ctx context.Context, period *models.TimePeriod) (*models.ReportSection, error) {
	performance, err := s.repo.GetRulesPerformance(ctx, period)
	if err != nil {
		return nil, err
	}

	return &models.ReportSection{
		Title: "Adaptive Rules Performance",
		Type:  "rules_performance",
		Data:  performance,
	}, nil
}

func (s *CloudIntelligenceService) generateFeedStatistics(ctx context.Context, period *models.TimePeriod) (*models.ReportSection, error) {
	stats, err := s.repo.GetFeedStatistics(ctx, period)
	if err != nil {
		return nil, err
	}

	return &models.ReportSection{
		Title: "Threat Feed Statistics",
		Type:  "feed_statistics",
		Data:  stats,
	}, nil
}

func (s *CloudIntelligenceService) generateMLPerformance(ctx context.Context, period *models.TimePeriod) (*models.ReportSection, error) {
	performance, err := s.repo.GetMLPerformance(ctx, period)
	if err != nil {
		return nil, err
	}

	return &models.ReportSection{
		Title: "ML Model Performance",
		Type:  "ml_performance",
		Data:  performance,
	}, nil
}

// Service lifecycle management

func (s *CloudIntelligenceService) Start(ctx context.Context) error {
	s.logger.Info("Starting Cloud Intelligence Service")

	// Load existing rules and feeds from repository
	if err := s.loadExistingData(ctx); err != nil {
		return fmt.Errorf("failed to load existing data: %w", err)
	}

	// Start background tasks
	s.startFeedUpdates()
	s.startRuleEvaluation()
	s.startEventConsumer()

	return nil
}

func (s *CloudIntelligenceService) Stop(ctx context.Context) error {
	s.logger.Info("Stopping Cloud Intelligence Service")

	// Cancel background tasks
	s.cancel()

	// Stop tickers
	if s.feedUpdateTicker != nil {
		s.feedUpdateTicker.Stop()
	}
	if s.ruleEvaluationTicker != nil {
		s.ruleEvaluationTicker.Stop()
	}
	s.logger.Info("Cloud Intelligence Service stopped")
	return nil
}

func (s *CloudIntelligenceService) loadExistingData(ctx context.Context) error {
	// Load adaptive rules
	rules, err := s.repo.GetAdaptiveRules(ctx, &models.AdaptiveRuleFilter{
		Status: models.RuleStatusActive,
	})
	if err != nil {
		return fmt.Errorf("failed to load adaptive rules: %w", err)
	}

	s.mutex.Lock()
	for _, rule := range rules {
		s.adaptiveRules[rule.ID] = rule
	}
	s.mutex.Unlock()

	s.logger.Info("Loaded adaptive rules", "count", len(rules))

	// Load threat feeds
	feeds, err := s.repo.GetThreatFeeds(ctx, &models.ThreatFeedFilter{
		Status: models.FeedStatusActive,
	})
	if err != nil {
		return fmt.Errorf("failed to load threat feeds: %w", err)
	}

	s.mutex.Lock()
	for _, feed := range feeds {
		s.threatFeeds[feed.ID] = feed
	}
	s.mutex.Unlock()

	s.logger.Info("Loaded threat feeds", "count", len(feeds))

	return nil
}

// Health check
func (s *CloudIntelligenceService) HealthCheck(ctx context.Context) error {
	// Check repository connectivity
	if err := s.repo.HealthCheck(ctx); err != nil {
		return fmt.Errorf("repository health check failed: %w", err)
	}

	// Check Kafka connectivity
	if err := s.kafkaProducer.HealthCheck(); err != nil {
		return fmt.Errorf("kafka producer health check failed: %w", err)
	}

	// Check ML model client
	if s.mlModelClient != nil {
		if _, err := s.mlModelClient.GetModelMetrics(ctx, "default"); err != nil {
			s.logger.Warn("ML model client health check failed", "error", err)
		}
	}

	// Check cloud providers
	healthyProviders := 0
	for name, provider := range s.cloudProviders {
		// Simple health check by getting reputation for a known good indicator
		testIndicator := &models.ThreatIndicator{
			Type:  "ip",
			Value: "8.8.8.8", // Google DNS - should be clean
		}
		
		_, err := provider.GetReputation(ctx, testIndicator)
		if err != nil {
			s.logger.Warn("Cloud provider health check failed", 
				"provider", name, 
				"error", err)
		} else {
			healthyProviders++
		}
	}

	if healthyProviders == 0 && len(s.cloudProviders) > 0 {
		return fmt.Errorf("no cloud providers are healthy")
	}

	return nil
}

// Metrics collection
func (s *CloudIntelligenceService) GetMetrics(ctx context.Context) (*models.CloudIntelligenceMetrics, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	metrics := &models.CloudIntelligenceMetrics{
		CachedIntelligenceCount: len(s.intelligenceCache),
		ActiveRulesCount:        0,
		InactiveRulesCount:      0,
		ActiveFeedsCount:        0,
		InactiveFeedsCount:      0,
		CloudProvidersCount:     len(s.cloudProviders),
		LastUpdated:            time.Now(),
	}

	// Count rule statuses
	for _, rule := range s.adaptiveRules {
		if rule.Status == models.RuleStatusActive {
			metrics.ActiveRulesCount++
		} else {
			metrics.InactiveRulesCount++
		}
	}

	// Count feed statuses
	for _, feed := range s.threatFeeds {
		if feed.Status == models.FeedStatusActive {
			metrics.ActiveFeedsCount++
		} else {
			metrics.InactiveFeedsCount++
		}
	}

	// Get repository metrics
	repoMetrics, err := s.repo.GetMetrics(ctx)
	if err != nil {
		s.logger.Error("Failed to get repository metrics", "error", err)
	} else {
		metrics.TotalThreatIntelligence = repoMetrics.TotalThreatIntelligence
		metrics.TotalThreatReports = repoMetrics.TotalThreatReports
		metrics.TotalRuleMatches = repoMetrics.TotalRuleMatches
	}

	return metrics, nil
}

// Cache management
func (s *CloudIntelligenceService) ClearCache(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.intelligenceCache = make(map[string]*models.ThreatIntelligence)
	s.logger.Info("Intelligence cache cleared")
	return nil
}

func (s *CloudIntelligenceService) GetCacheStats(ctx context.Context) (*models.CacheStats, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := &models.CacheStats{
		TotalEntries:    len(s.intelligenceCache),
		ExpiredEntries:  0,
		HitRate:         0.0, // Would need to track hits/misses
		LastUpdated:     time.Now(),
	}

	// Count expired entries
	now := time.Now()
	for _, intelligence := range s.intelligenceCache {
		if now.Sub(intelligence.LastUpdated) > s.cacheExpiration {
			stats.ExpiredEntries++
		}
	}

	return stats, nil
}

// Cleanup expired cache entries
func (s *CloudIntelligenceService) cleanupExpiredCache() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for key, intelligence := range s.intelligenceCache {
		if now.Sub(intelligence.LastUpdated) > s.cacheExpiration {
			delete(s.intelligenceCache, key)
		}
	}
}

// Batch operations
func (s *CloudIntelligenceService) BatchCreateAdaptiveRules(ctx context.Context, rules []*models.AdaptiveRule) ([]*models.AdaptiveRule, error) {
	createdRules := make([]*models.AdaptiveRule, 0, len(rules))
	
	for _, rule := range rules {
		if err := s.CreateAdaptiveRule(ctx, rule); err != nil {
			s.logger.Error("Failed to create rule in batch", 
				"rule_name", rule.Name, 
				"error", err)
			continue
		}
		createdRules = append(createdRules, rule)
	}

	s.logger.Info("Batch rule creation completed", 
		"requested", len(rules), 
		"created", len(createdRules))

	return createdRules, nil
}

func (s *CloudIntelligenceService) BatchUpdateAdaptiveRules(ctx context.Context, rules []*models.AdaptiveRule) ([]*models.AdaptiveRule, error) {
	updatedRules := make([]*models.AdaptiveRule, 0, len(rules))
	
	for _, rule := range rules {
		if err := s.UpdateAdaptiveRule(ctx, rule); err != nil {
			s.logger.Error("Failed to update rule in batch", 
				"rule_id", rule.ID, 
				"error", err)
			continue
		}
		updatedRules = append(updatedRules, rule)
	}

	s.logger.Info("Batch rule update completed", 
		"requested", len(rules), 
		"updated", len(updatedRules))

	return updatedRules, nil
}

// Export/Import functionality
func (s *CloudIntelligenceService) ExportAdaptiveRules(ctx context.Context, filter *models.AdaptiveRuleFilter) ([]byte, error) {
	rules, err := s.repo.GetAdaptiveRules(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules for export: %w", err)
	}

	exportData := &models.RuleExport{
		Version:     "1.0",
		ExportedAt:  time.Now(),
		Rules:       rules,
		TotalCount:  len(rules),
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal export data: %w", err)
	}

	s.logger.Info("Rules exported", "count", len(rules))
	return data, nil
}

func (s *CloudIntelligenceService) ImportAdaptiveRules(ctx context.Context, data []byte, overwrite bool) (*models.ImportResult, error) {
	var importData models.RuleExport
	if err := json.Unmarshal(data, &importData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	result := &models.ImportResult{
		TotalRules:    len(importData.Rules),
		ImportedRules: 0,
		SkippedRules:  0,
		Errors:        make([]string, 0),
	}

	for _, rule := range importData.Rules {
		// Check if rule already exists
		existing, err := s.repo.GetAdaptiveRule(ctx, rule.ID)
		if err == nil && existing != nil && !overwrite {
			result.SkippedRules++
			result.Errors = append(result.Errors, 
				fmt.Sprintf("Rule %s already exists (skipped)", rule.ID))
			continue
		}

		// Reset timestamps for import
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()

		if err := s.CreateAdaptiveRule(ctx, rule); err != nil {
			result.Errors = append(result.Errors, 
				fmt.Sprintf("Failed to import rule %s: %v", rule.ID, err))
			continue
		}

		result.ImportedRules++
	}

	s.logger.Info("Rules import completed", 
		"total", result.TotalRules,
		"imported", result.ImportedRules,
		"skipped", result.SkippedRules,
		"errors", len(result.Errors))

	return result, nil
}
