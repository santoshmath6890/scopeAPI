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

type AgentManagementService struct {
	agentRepo       repository.AgentRepository
	logger          logging.Logger
	kafkaProducer   kafka.Producer
	agents          map[string]*models.Agent
	agentSessions   map[string]*models.AgentSession
	agentMetrics    map[string]*models.AgentMetrics
	mutex           sync.RWMutex
	sessionTimeout  time.Duration
	heartbeatTicker *time.Ticker
	ctx             context.Context
	cancel          context.CancelFunc
}

type AgentManagementServiceConfig struct {
	SessionTimeout    time.Duration
	HeartbeatInterval time.Duration
	MaxInactiveTime   time.Duration
}

func NewAgentManagementService(
	agentRepo repository.AgentRepository,
	logger logging.Logger,
	kafkaProducer kafka.Producer,
	config *AgentManagementServiceConfig,
) *AgentManagementService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &AgentManagementService{
		agentRepo:       agentRepo,
		logger:          logger,
		kafkaProducer:   kafkaProducer,
		agents:          make(map[string]*models.Agent),
		agentSessions:   make(map[string]*models.AgentSession),
		agentMetrics:    make(map[string]*models.AgentMetrics),
		sessionTimeout:  config.SessionTimeout,
		ctx:             ctx,
		cancel:          cancel,
	}

	// Start background tasks
	service.startHeartbeatMonitoring(config.HeartbeatInterval)
	service.loadAgents()

	return service
}

// Agent Registration and Management

func (s *AgentManagementService) RegisterAgent(ctx context.Context, request *models.AgentRegistrationRequest) (*models.AgentRegistrationResponse, error) {
	s.logger.Info("Registering new agent",
		"agent_id", request.AgentID,
		"hostname", request.Hostname,
		"version", request.Version)

	// Validate registration request
	if err := s.validateRegistrationRequest(request); err != nil {
		return nil, fmt.Errorf("registration validation failed: %w", err)
	}

	// Check if agent already exists
	existingAgent, err := s.agentRepo.GetAgent(ctx, request.AgentID)
	if err == nil && existingAgent != nil {
		// Update existing agent
		return s.updateExistingAgent(ctx, existingAgent, request)
	}

	// Create new agent
	agent := &models.Agent{
		ID:           request.AgentID,
		Name:         request.Name,
		Hostname:     request.Hostname,
		IPAddress:    request.IPAddress,
		Version:      request.Version,
		Type:         request.Type,
		Capabilities: request.Capabilities,
		Configuration: models.AgentConfiguration{
			LogLevel:        "info",
			ReportInterval:  30 * time.Second,
			MaxConnections:  1000,
			BufferSize:      10000,
			EnableMetrics:   true,
			EnableTracing:   false,
		},
		Status:       models.AgentStatusPending,
		RegisteredAt: time.Now(),
		UpdatedAt:    time.Now(),
		Metadata:     request.Metadata,
	}

	// Save agent to repository
	if err := s.agentRepo.CreateAgent(ctx, agent); err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	// Add to in-memory cache
	s.mutex.Lock()
	s.agents[agent.ID] = agent
	s.mutex.Unlock()

	// Create agent session
	session := &models.AgentSession{
		ID:           uuid.New().String(),
		AgentID:      agent.ID,
		SessionToken: s.generateSessionToken(),
		StartedAt:    time.Now(),
		LastSeen:     time.Now(),
		Status:       models.SessionStatusActive,
		IPAddress:    request.IPAddress,
		UserAgent:    request.UserAgent,
	}

	s.mutex.Lock()
	s.agentSessions[agent.ID] = session
	s.mutex.Unlock()

	// Initialize agent metrics
	s.initializeAgentMetrics(agent.ID)

	// Publish agent registration event
	s.publishAgentEvent(ctx, "agent_registered", agent)

	response := &models.AgentRegistrationResponse{
		AgentID:      agent.ID,
		SessionToken: session.SessionToken,
		Status:       agent.Status,
		Configuration: agent.Configuration,
		Policies:     s.getAgentPolicies(ctx, agent.ID),
		Message:      "Agent registered successfully",
	}

	s.logger.Info("Agent registered successfully", "agent_id", agent.ID)
	return response, nil
}

func (s *AgentManagementService) UnregisterAgent(ctx context.Context, agentID string) error {
	s.logger.Info("Unregistering agent", "agent_id", agentID)

	// Get agent
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Update agent status
	agent.Status = models.AgentStatusDecommissioned
	agent.UpdatedAt = time.Now()

	if err := s.agentRepo.UpdateAgent(ctx, agent); err != nil {
		return fmt.Errorf("failed to update agent status: %w", err)
	}

	// Remove from in-memory caches
	s.mutex.Lock()
	delete(s.agents, agentID)
	delete(s.agentSessions, agentID)
	delete(s.agentMetrics, agentID)
	s.mutex.Unlock()

	// Publish agent unregistration event
	s.publishAgentEvent(ctx, "agent_unregistered", agent)

	s.logger.Info("Agent unregistered successfully", "agent_id", agentID)
	return nil
}

func (s *AgentManagementService) GetAgent(ctx context.Context, agentID string) (*models.Agent, error) {
	// Try in-memory cache first
	s.mutex.RLock()
	if agent, exists := s.agents[agentID]; exists {
		s.mutex.RUnlock()
		return agent, nil
	}
	s.mutex.RUnlock()

	// Fallback to repository
	return s.agentRepo.GetAgent(ctx, agentID)
}

func (s *AgentManagementService) GetAgents(ctx context.Context, filter *models.AgentFilter) ([]*models.Agent, error) {
	return s.agentRepo.GetAgents(ctx, filter)
}

func (s *AgentManagementService) UpdateAgent(ctx context.Context, agent *models.Agent) error {
	agent.UpdatedAt = time.Now()

	if err := s.agentRepo.UpdateAgent(ctx, agent); err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	// Update in-memory cache
	s.mutex.Lock()
	s.agents[agent.ID] = agent
	s.mutex.Unlock()

	// Publish agent update event
	s.publishAgentEvent(ctx, "agent_updated", agent)

	s.logger.Info("Agent updated successfully", "agent_id", agent.ID)
	return nil
}

// Agent Session Management

func (s *AgentManagementService) AuthenticateAgent(ctx context.Context, agentID, sessionToken string) (*models.AgentSession, error) {
	s.mutex.RLock()
	session, exists := s.agentSessions[agentID]
	s.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("agent session not found")
	}

	if session.SessionToken != sessionToken {
		return nil, fmt.Errorf("invalid session token")
	}

	if session.Status != models.SessionStatusActive {
		return nil, fmt.Errorf("session is not active")
	}

	// Check session timeout
	if time.Since(session.LastSeen) > s.sessionTimeout {
		session.Status = models.SessionStatusExpired
		return nil, fmt.Errorf("session expired")
	}

	// Update last seen
	session.LastSeen = time.Now()

	return session, nil
}

func (s *AgentManagementService) RefreshSession(ctx context.Context, agentID string) (*models.AgentSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.agentSessions[agentID]
	if !exists {
		return nil, fmt.Errorf("agent session not found")
	}

	// Generate new session token
	session.SessionToken = s.generateSessionToken()
	session.LastSeen = time.Now()
	session.Status = models.SessionStatusActive

	s.logger.Info("Agent session refreshed", "agent_id", agentID)
	return session, nil
}

func (s *AgentManagementService) TerminateSession(ctx context.Context, agentID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.agentSessions[agentID]
	if !exists {
		return fmt.Errorf("agent session not found")
	}

	session.Status = models.SessionStatusTerminated
	session.EndedAt = &time.Time{}
	*session.EndedAt = time.Now()

	s.logger.Info("Agent session terminated", "agent_id", agentID)
	return nil
}

// Agent Heartbeat and Health Monitoring

func (s *AgentManagementService) ProcessHeartbeat(ctx context.Context, heartbeat *models.AgentHeartbeat) error {
	s.logger.Debug("Processing agent heartbeat", "agent_id", heartbeat.AgentID)

	// Validate heartbeat
	if err := s.validateHeartbeat(heartbeat); err != nil {
		return fmt.Errorf("heartbeat validation failed: %w", err)
	}

	// Update agent session
	s.mutex.Lock()
	if session, exists := s.agentSessions[heartbeat.AgentID]; exists {
		session.LastSeen = time.Now()
		session.Status = models.SessionStatusActive
	}
	s.mutex.Unlock()

	// Update agent status
	if agent, exists := s.agents[heartbeat.AgentID]; exists {
		agent.Status = models.AgentStatusActive
		agent.LastHeartbeat = time.Now()
		agent.UpdatedAt = time.Now()
	}

	// Update agent metrics
	s.updateAgentMetrics(heartbeat.AgentID, heartbeat.Metrics)

	// Process health status
	if heartbeat.HealthStatus != nil {
		s.processAgentHealthStatus(heartbeat.AgentID, heartbeat.HealthStatus)
	}

	// Publish heartbeat event
	s.publishHeartbeatEvent(ctx, heartbeat)

	return nil
}

func (s *AgentManagementService) GetAgentHealth(ctx context.Context, agentID string) (*models.AgentHealthStatus, error) {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	s.mutex.RLock()
	session := s.agentSessions[agentID]
	metrics := s.agentMetrics[agentID]
	s.mutex.RUnlock()

	health := &models.AgentHealthStatus{
		AgentID:        agentID,
		Status:         agent.Status,
		LastHeartbeat:  agent.LastHeartbeat,
		Uptime:         time.Since(agent.RegisteredAt),
		Version:        agent.Version,
		CPUUsage:       0,
		MemoryUsage:    0,
		DiskUsage:      0,
		NetworkLatency: 0,
		ErrorRate:      0,
		CheckedAt:      time.Now(),
	}

	if session != nil {
		health.SessionActive = session.Status == models.SessionStatusActive
		health.LastSeen = session.LastSeen
	}

	if metrics != nil {
		health.CPUUsage = metrics.CPUUsage
		health.MemoryUsage = metrics.MemoryUsage
		health.DiskUsage = metrics.DiskUsage
		health.NetworkLatency = metrics.NetworkLatency
		health.ErrorRate = metrics.ErrorRate
		health.RequestsProcessed = metrics.RequestsProcessed
		health.BlockedRequests = metrics.BlockedRequests
	}

	// Determine overall health
	health.OverallHealth = s.calculateOverallHealth(health)

	return health, nil
}

// Agent Configuration Management

func (s *AgentManagementService) UpdateAgentConfiguration(ctx context.Context, agentID string, config *models.AgentConfiguration) error {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Validate configuration
	if err := s.validateAgentConfiguration(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Update agent configuration
	agent.Configuration = *config
	agent.UpdatedAt = time.Now()

	if err := s.UpdateAgent(ctx, agent); err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	// Send configuration update to agent
	configUpdate := &models.AgentConfigurationUpdate{
		AgentID:       agentID,
		Configuration: *config,
		UpdatedAt:     time.Now(),
	}

	if err := s.sendConfigurationUpdate(ctx, configUpdate); err != nil {
		s.logger.Error("Failed to send configuration update to agent",
			"agent_id", agentID,
			"error", err)
	}

	s.logger.Info("Agent configuration updated", "agent_id", agentID)
	return nil
}

func (s *AgentManagementService) GetAgentConfiguration(ctx context.Context, agentID string) (*models.AgentConfiguration, error) {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	return &agent.Configuration, nil
}

// Agent Policy Management

func (s *AgentManagementService) AssignPolicyToAgent(ctx context.Context, agentID, policyID string) error {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Check if policy is already assigned
	for _, assignedPolicyID := range agent.AssignedPolicies {
		if assignedPolicyID == policyID {
			return fmt.Errorf("policy already assigned to agent")
		}
	}

	// Add policy to agent
	agent.AssignedPolicies = append(agent.AssignedPolicies, policyID)
	agent.UpdatedAt = time.Now()

	if err := s.UpdateAgent(ctx, agent); err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	// Send policy update to agent
	policyUpdate := &models.AgentPolicyUpdate{
		AgentID:   agentID,
		PolicyID:  policyID,
		Action:    "assign",
		UpdatedAt: time.Now(),
	}

	if err := s.sendPolicyUpdate(ctx, policyUpdate); err != nil {
		s.logger.Error("Failed to send policy update to agent",
			"agent_id", agentID,
			"policy_id", policyID,
			"error", err)
	}

	s.logger.Info("Policy assigned to agent", "agent_id", agentID, "policy_id", policyID)
	return nil
}

func (s *AgentManagementService) UnassignPolicyFromAgent(ctx context.Context, agentID, policyID string) error {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	// Remove policy from agent
	var updatedPolicies []string
	found := false
	for _, assignedPolicyID := range agent.AssignedPolicies {
		if assignedPolicyID != policyID {
			updatedPolicies = append(updatedPolicies, assignedPolicyID)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("policy not assigned to agent")
	}

	agent.AssignedPolicies = updatedPolicies
	agent.UpdatedAt = time.Now()

	if err := s.UpdateAgent(ctx, agent); err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	// Send policy update to agent
	policyUpdate := &models.AgentPolicyUpdate{
		AgentID:   agentID,
		PolicyID:  policyID,
		Action:    "unassign",
		UpdatedAt: time.Now(),
	}

	if err := s.sendPolicyUpdate(ctx, policyUpdate); err != nil {
		s.logger.Error("Failed to send policy update to agent",
			"agent_id", agentID,
			"policy_id", policyID,
			"error", err)
	}

	s.logger.Info("Policy unassigned from agent", "agent_id", agentID, "policy_id", policyID)
	return nil
}

func (s *AgentManagementService) GetAgentPolicies(ctx context.Context, agentID string) ([]string, error) {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	return agent.AssignedPolicies, nil
}

// Agent Commands and Control

func (s *AgentManagementService) SendCommandToAgent(ctx context.Context, agentID string, command *models.AgentCommand) error {
	// Validate agent exists and is active
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	if agent.Status != models.AgentStatusActive {
		return fmt.Errorf("agent is not active")
	}

	// Validate command
	if err := s.validateAgentCommand(command); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	// Set command metadata
	command.ID = uuid.New().String()
	command.AgentID = agentID
	command.SentAt = time.Now()
	command.Status = models.CommandStatusPending

	// Store command for tracking
	if err := s.agentRepo.CreateAgentCommand(ctx, command); err != nil {
		return fmt.Errorf("failed to store command: %w", err)
	}

	// Send command via Kafka
	commandData, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	topic := fmt.Sprintf("agent-commands-%s", agentID)
	if err := s.kafkaProducer.Produce(ctx, topic, commandData); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	s.logger.Info("Command sent to agent",
		"agent_id", agentID,
		"command_id", command.ID,
		"command_type", command.Type)

	return nil
}

func (s *AgentManagementService) GetAgentCommands(ctx context.Context, agentID string, filter *models.AgentCommandFilter) ([]*models.AgentCommand, error) {
	return s.agentRepo.GetAgentCommands(ctx, agentID, filter)
}

func (s *AgentManagementService) UpdateCommandStatus(ctx context.Context, commandID string, status models.CommandStatus, result *models.CommandResult) error {
	command, err := s.agentRepo.GetAgentCommand(ctx, commandID)
	if err != nil {
		return fmt.Errorf("command not found: %w", err)
	}

	command.Status = status
	command.UpdatedAt = time.Now()

	if result != nil {
		command.Result = result
		command.CompletedAt = &time.Time{}
		*command.CompletedAt = time.Now()
	}

	if err := s.agentRepo.UpdateAgentCommand(ctx, command); err != nil {
		return fmt.Errorf("failed to update command: %w", err)
	}

	s.logger.Info("Command status updated",
		"command_id", commandID,
		"status", status)

	return nil
}

// Agent Metrics and Analytics

func (s *AgentManagementService) GetAgentMetrics(ctx context.Context, agentID string, timeRange *models.TimeRange) (*models.AgentMetrics, error) {
	// Get current metrics from memory
	s.mutex.RLock()
	currentMetrics := s.agentMetrics[agentID]
	s.mutex.RUnlock()

	if currentMetrics == nil {
		return nil, fmt.Errorf("agent metrics not found")
	}

	// If no time range specified, return current metrics
	if timeRange == nil {
		return currentMetrics, nil
	}

	// Get historical metrics from repository
	historicalMetrics, err := s.agentRepo.GetAgentMetricsHistory(ctx, agentID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical metrics: %w", err)
	}

	// Combine current and historical metrics
	combinedMetrics := s.combineMetrics(currentMetrics, historicalMetrics)
	return combinedMetrics, nil
}

func (s *AgentManagementService) GetAgentStatistics(ctx context.Context, filter *models.AgentStatsFilter) (*models.AgentStatistics, error) {
	stats, err := s.agentRepo.GetAgentStatistics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent statistics: %w", err)
	}

	// Add real-time statistics
	s.mutex.RLock()
	stats.ActiveAgents = 0
	stats.InactiveAgents = 0
	stats.TotalAgents = len(s.agents)

	for _, agent := range s.agents {
		if agent.Status == models.AgentStatusActive {
			stats.ActiveAgents++
		} else {
			stats.InactiveAgents++
		}
	}
	s.mutex.RUnlock()

	return stats, nil
}

// Agent Deployment and Updates

func (s *AgentManagementService) DeployAgentUpdate(ctx context.Context, agentID string, update *models.AgentUpdate) error {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Validate update
	if err := s.validateAgentUpdate(update); err != nil {
		return fmt.Errorf("update validation failed: %w", err)
	}

	// Create update command
	command := &models.AgentCommand{
		Type: models.CommandTypeUpdate,
		Payload: map[string]interface{}{
			"version":     update.Version,
			"download_url": update.DownloadURL,
			"checksum":    update.Checksum,
			"restart_required": update.RestartRequired,
		},
		Priority: models.CommandPriorityHigh,
		Timeout:  30 * time.Minute,
	}

	if err := s.SendCommandToAgent(ctx, agentID, command); err != nil {
		return fmt.Errorf("failed to send update command: %w", err)
	}

	// Update agent status
	agent.Status = models.AgentStatusUpdating
	agent.UpdatedAt = time.Now()

	if err := s.UpdateAgent(ctx, agent); err != nil {
		return fmt.Errorf("failed to update agent status: %w", err)
	}

	s.logger.Info("Agent update deployed",
		"agent_id", agentID,
		"version", update.Version)

	return nil
}

func (s *AgentManagementService) GetAgentLogs(ctx context.Context, agentID string, filter *models.AgentLogFilter) ([]*models.AgentLog, error) {
	return s.agentRepo.GetAgentLogs(ctx, agentID, filter)
}

// Private helper methods

func (s *AgentManagementService) validateRegistrationRequest(request *models.AgentRegistrationRequest) error {
	if request.AgentID == "" {
		return fmt.Errorf("agent ID is required")
	}

	if request.Name == "" {
		return fmt.Errorf("agent name is required")
	}

	if request.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	if request.Version == "" {
		return fmt.Errorf("version is required")
	}

	validTypes := []models.AgentType{
		models.AgentTypeBlocking,
		models.AgentTypeDiscovery,
		models.AgentTypeMonitoring,
	}

	found := false
	for _, validType := range validTypes {
		if request.Type == validType {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid agent type: %s", request.Type)
	}

	return nil
}

func (s *AgentManagementService) updateExistingAgent(ctx context.Context, agent *models.Agent, request *models.AgentRegistrationRequest) (*models.AgentRegistrationResponse, error) {
	// Update agent information
	agent.Name = request.Name
	agent.Hostname = request.Hostname
	agent.IPAddress = request.IPAddress
	agent.Version = request.Version
	agent.Capabilities = request.Capabilities
	agent.Status = models.AgentStatusActive
	agent.UpdatedAt = time.Now()
	agent.Metadata = request.Metadata

	if err := s.UpdateAgent(ctx, agent); err != nil {
		return nil, fmt.Errorf("failed to update existing agent: %w", err)
	}

	// Create new session
	session := &models.AgentSession{
		ID:           uuid.New().String(),
		AgentID:      agent.ID,
		SessionToken: s.generateSessionToken(),
		StartedAt:    time.Now(),
		LastSeen:     time.Now(),
		Status:       models.SessionStatusActive,
		IPAddress:    request.IPAddress,
		UserAgent:    request.UserAgent,
	}

	s.mutex.Lock()
	s.agentSessions[agent.ID] = session
	s.mutex.Unlock()

	response := &models.AgentRegistrationResponse{
		AgentID:      agent.ID,
		SessionToken: session.SessionToken,
		Status:       agent.Status,
		Configuration: agent.Configuration,
		Policies:     s.getAgentPolicies(ctx, agent.ID),
		Message:      "Agent re-registered successfully",
	}

	return response, nil
}

func (s *AgentManagementService) generateSessionToken() string {
	return uuid.New().String()
}

func (s *AgentManagementService) getAgentPolicies(ctx context.Context, agentID string) []string {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return []string{}
	}
	return agent.AssignedPolicies
}

func (s *AgentManagementService) initializeAgentMetrics(agentID string) {
	s.mutex.Lock()
	s.agentMetrics[agentID] = &models.AgentMetrics{
		AgentID:           agentID,
		CPUUsage:          0,
		MemoryUsage:       0,
		DiskUsage:         0,
		NetworkLatency:    0,
		RequestsProcessed: 0,
		BlockedRequests:   0,
		ErrorRate:         0,
		LastUpdated:       time.Now(),
	}
	s.mutex.Unlock()
}

func (s *AgentManagementService) validateHeartbeat(heartbeat *models.AgentHeartbeat) error {
	if heartbeat.AgentID == "" {
		return fmt.Errorf("agent ID is required")
	}

	if heartbeat.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	// Check if heartbeat is not too old
	if time.Since(heartbeat.Timestamp) > 5*time.Minute {
		return fmt.Errorf("heartbeat is too old")
	}

	return nil
}

func (s *AgentManagementService) updateAgentMetrics(agentID string, metrics *models.AgentMetrics) {
	if metrics == nil {
		return
	}

	s.mutex.Lock()
	if existingMetrics, exists := s.agentMetrics[agentID]; exists {
		existingMetrics.CPUUsage = metrics.CPUUsage
		existingMetrics.MemoryUsage = metrics.MemoryUsage
		existingMetrics.DiskUsage = metrics.DiskUsage
		existingMetrics.NetworkLatency = metrics.NetworkLatency
		existingMetrics.RequestsProcessed = metrics.RequestsProcessed
		existingMetrics.BlockedRequests = metrics.BlockedRequests
		existingMetrics.ErrorRate = metrics.ErrorRate
		existingMetrics.LastUpdated = time.Now()
	}
	s.mutex.Unlock()
}

func (s *AgentManagementService) processAgentHealthStatus(agentID string, healthStatus *models.AgentHealthStatus) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	agent, exists := s.agents[agentID]
	if !exists {
		return
	}

	// Update agent status based on health
	previousStatus := agent.Status
	if healthStatus.OverallHealth == "healthy" {
		agent.Status = models.AgentStatusActive
	} else if healthStatus.OverallHealth == "degraded" {
		agent.Status = models.AgentStatusDegraded
	} else {
		agent.Status = models.AgentStatusUnhealthy
	}

	// Log status changes
	if previousStatus != agent.Status {
		s.logger.Info("Agent status changed",
			"agent_id", agentID,
			"previous_status", previousStatus,
			"new_status", agent.Status,
			"health", healthStatus.OverallHealth)
	}

	agent.UpdatedAt = time.Now()
}

func (s *AgentManagementService) calculateOverallHealth(health *models.AgentHealthStatus) string {
	// Define thresholds
	const (
		cpuThreshold     = 80.0
		memoryThreshold  = 85.0
		diskThreshold    = 90.0
		errorRateThreshold = 5.0
	)

	unhealthyCount := 0
	degradedCount := 0

	// Check CPU usage
	if health.CPUUsage > cpuThreshold {
		if health.CPUUsage > 95.0 {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}

	// Check memory usage
	if health.MemoryUsage > memoryThreshold {
		if health.MemoryUsage > 95.0 {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}

	// Check disk usage
	if health.DiskUsage > diskThreshold {
		if health.DiskUsage > 98.0 {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}

	// Check error rate
	if health.ErrorRate > errorRateThreshold {
		if health.ErrorRate > 15.0 {
			unhealthyCount++
		} else {
			degradedCount++
		}
	}

	// Check last heartbeat
	if time.Since(health.LastHeartbeat) > 2*time.Minute {
		unhealthyCount++
	} else if time.Since(health.LastHeartbeat) > time.Minute {
		degradedCount++
	}

	// Determine overall health
	if unhealthyCount > 0 {
		return "unhealthy"
	} else if degradedCount > 0 {
		return "degraded"
	}

	return "healthy"
}

func (s *AgentManagementService) validateAgentConfiguration(config *models.AgentConfiguration) error {
	if config.ReportInterval < time.Second {
		return fmt.Errorf("report interval must be at least 1 second")
	}

	if config.MaxConnections < 1 {
		return fmt.Errorf("max connections must be at least 1")
	}

	if config.BufferSize < 100 {
		return fmt.Errorf("buffer size must be at least 100")
	}

	validLogLevels := []string{"debug", "info", "warn", "error"}
	found := false
	for _, level := range validLogLevels {
		if config.LogLevel == level {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid log level: %s", config.LogLevel)
	}

	return nil
}

func (s *AgentManagementService) validateAgentCommand(command *models.AgentCommand) error {
	if command.Type == "" {
		return fmt.Errorf("command type is required")
	}

	validTypes := []models.CommandType{
		models.CommandTypeRestart,
		models.CommandTypeUpdate,
		models.CommandTypeConfigUpdate,
		models.CommandTypePolicyUpdate,
		models.CommandTypeHealthCheck,
		models.CommandTypeLogLevel,
		models.CommandTypeShutdown,
	}

	found := false
	for _, validType := range validTypes {
		if command.Type == validType {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid command type: %s", command.Type)
	}

	if command.Timeout < time.Second {
		return fmt.Errorf("command timeout must be at least 1 second")
	}

	return nil
}

func (s *AgentManagementService) validateAgentUpdate(update *models.AgentUpdate) error {
	if update.Version == "" {
		return fmt.Errorf("update version is required")
	}

	if update.DownloadURL == "" {
		return fmt.Errorf("download URL is required")
	}

	if update.Checksum == "" {
		return fmt.Errorf("checksum is required")
	}

	return nil
}

func (s *AgentManagementService) sendConfigurationUpdate(ctx context.Context, configUpdate *models.AgentConfigurationUpdate) error {
	updateData, err := json.Marshal(configUpdate)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration update: %w", err)
	}

	topic := fmt.Sprintf("agent-config-updates-%s", configUpdate.AgentID)
	return s.kafkaProducer.Produce(ctx, topic, updateData)
}

func (s *AgentManagementService) sendPolicyUpdate(ctx context.Context, policyUpdate *models.AgentPolicyUpdate) error {
	updateData, err := json.Marshal(policyUpdate)
	if err != nil {
		return fmt.Errorf("failed to marshal policy update: %w", err)
	}

	topic := fmt.Sprintf("agent-policy-updates-%s", policyUpdate.AgentID)
	return s.kafkaProducer.Produce(ctx, topic, updateData)
}

func (s *AgentManagementService) combineMetrics(current *models.AgentMetrics, historical []*models.AgentMetrics) *models.AgentMetrics {
	combined := &models.AgentMetrics{
		AgentID:           current.AgentID,
		CPUUsage:          current.CPUUsage,
		MemoryUsage:       current.MemoryUsage,
		DiskUsage:         current.DiskUsage,
		NetworkLatency:    current.NetworkLatency,
		RequestsProcessed: current.RequestsProcessed,
		BlockedRequests:   current.BlockedRequests,
		ErrorRate:         current.ErrorRate,
		LastUpdated:       current.LastUpdated,
		Historical:        historical,
	}

	return combined
}

func (s *AgentManagementService) publishAgentEvent(ctx context.Context, eventType string, agent *models.Agent) {
	event := &models.AgentEvent{
		ID:        uuid.New().String(),
		Type:      eventType,
		AgentID:   agent.ID,
		AgentName: agent.Name,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"hostname":    agent.Hostname,
			"ip_address":  agent.IPAddress,
			"version":     agent.Version,
			"status":      agent.Status,
			"agent_type":  agent.Type,
		},
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		s.logger.Error("Failed to marshal agent event", "error", err)
		return
	}

	if err := s.kafkaProducer.Produce(ctx, "agent-events", eventData); err != nil {
		s.logger.Error("Failed to publish agent event", "error", err)
	}
}

func (s *AgentManagementService) publishHeartbeatEvent(ctx context.Context, heartbeat *models.AgentHeartbeat) {
	event := &models.AgentEvent{
		ID:        uuid.New().String(),
		Type:      "agent_heartbeat",
		AgentID:   heartbeat.AgentID,
		Timestamp: heartbeat.Timestamp,
		Data: map[string]interface{}{
			"status":         heartbeat.Status,
			"health_status":  heartbeat.HealthStatus,
			"metrics":        heartbeat.Metrics,
		},
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		s.logger.Error("Failed to marshal heartbeat event", "error", err)
		return
	}

	if err := s.kafkaProducer.Produce(ctx, "agent-heartbeats", eventData); err != nil {
		s.logger.Error("Failed to publish heartbeat event", "error", err)
	}
}

func (s *AgentManagementService) loadAgents() {
	ctx := context.Background()
	agents, err := s.agentRepo.GetAgents(ctx, &models.AgentFilter{})
	if err != nil {
		s.logger.Error("Failed to load agents from repository", "error", err)
		return
	}

	s.mutex.Lock()
	for _, agent := range agents {
		s.agents[agent.ID] = agent
		s.initializeAgentMetrics(agent.ID)
	}
	s.mutex.Unlock()

	s.logger.Info("Agents loaded from repository", "count", len(agents))
}

func (s *AgentManagementService) startHeartbeatMonitoring(interval time.Duration) {
	s.heartbeatTicker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				s.heartbeatTicker.Stop()
				return
			case <-s.heartbeatTicker.C:
				s.checkAgentHeartbeats()
			}
		}
	}()
}

func (s *AgentManagementService) checkAgentHeartbeats() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for agentID, agent := range s.agents {
		// Check if agent hasn't sent heartbeat recently
		if agent.Status == models.AgentStatusActive && 
		   !agent.LastHeartbeat.IsZero() && 
		   now.Sub(agent.LastHeartbeat) > 2*time.Minute {
			
			s.logger.Warn("Agent heartbeat timeout", 
				"agent_id", agentID,
				"last_heartbeat", agent.LastHeartbeat)
			
			agent.Status = models.AgentStatusInactive
			agent.UpdatedAt = now
			
			// Update in repository
			go func(a *models.Agent) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				
				if err := s.agentRepo.UpdateAgent(ctx, a); err != nil {
					s.logger.Error("Failed to update agent status", 
						"agent_id", a.ID, 
						"error", err)
				}
			}(agent)
			
			// Publish agent inactive event
			go func(a *models.Agent) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				s.publishAgentEvent(ctx, "agent_inactive", a)
			}(agent)
		}
	}
}

// Bulk Operations

func (s *AgentManagementService) BulkUpdateAgentConfiguration(ctx context.Context, agentIDs []string, config *models.AgentConfiguration) (*models.BulkOperationResult, error) {
	result := &models.BulkOperationResult{
		TotalItems:    len(agentIDs),
		SuccessCount:  0,
		FailureCount:  0,
		Errors:        make([]string, 0),
		StartedAt:     time.Now(),
	}

	for _, agentID := range agentIDs {
		if err := s.UpdateAgentConfiguration(ctx, agentID, config); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Agent %s: %v", agentID, err))
		} else {
			result.SuccessCount++
		}
	}

	result.CompletedAt = time.Now()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)

	s.logger.Info("Bulk agent configuration update completed",
		"total", result.TotalItems,
		"success", result.SuccessCount,
		"failures", result.FailureCount)

	return result, nil
}

func (s *AgentManagementService) BulkAssignPolicy(ctx context.Context, agentIDs []string, policyID string) (*models.BulkOperationResult, error) {
	result := &models.BulkOperationResult{
		TotalItems:    len(agentIDs),
		SuccessCount:  0,
		FailureCount:  0,
		Errors:        make([]string, 0),
		StartedAt:     time.Now(),
	}

	for _, agentID := range agentIDs {
		if err := s.AssignPolicyToAgent(ctx, agentID, policyID); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Agent %s: %v", agentID, err))
		} else {
			result.SuccessCount++
		}
	}

	result.CompletedAt = time.Now()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)

	s.logger.Info("Bulk policy assignment completed",
		"policy_id", policyID,
		"total", result.TotalItems,
		"success", result.SuccessCount,
		"failures", result.FailureCount)

	return result, nil
}

func (s *AgentManagementService) BulkSendCommand(ctx context.Context, agentIDs []string, command *models.AgentCommand) (*models.BulkOperationResult, error) {
	result := &models.BulkOperationResult{
		TotalItems:    len(agentIDs),
		SuccessCount:  0,
		FailureCount:  0,
		Errors:        make([]string, 0),
		StartedAt:     time.Now(),
	}

	for _, agentID := range agentIDs {
		// Create a copy of the command for each agent
		agentCommand := *command
		if err := s.SendCommandToAgent(ctx, agentID, &agentCommand); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Agent %s: %v", agentID, err))
		} else {
			result.SuccessCount++
		}
	}

	result.CompletedAt = time.Now()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)

	s.logger.Info("Bulk command send completed",
		"command_type", command.Type,
		"total", result.TotalItems,
		"success", result.SuccessCount,
		"failures", result.FailureCount)

	return result, nil
}

// Agent Groups Management

func (s *AgentManagementService) CreateAgentGroup(ctx context.Context, group *models.AgentGroup) error {
	// Validate group
	if group.Name == "" {
		return fmt.Errorf("group name is required")
	}

	group.ID = uuid.New().String()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	if err := s.agentRepo.CreateAgentGroup(ctx, group); err != nil {
		return fmt.Errorf("failed to create agent group: %w", err)
	}

	s.logger.Info("Agent group created", "group_id", group.ID, "name", group.Name)
	return nil
}

func (s *AgentManagementService) AddAgentToGroup(ctx context.Context, agentID, groupID string) error {
	// Validate agent exists
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Validate group exists
	group, err := s.agentRepo.GetAgentGroup(ctx, groupID)
	if err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	// Check if agent is already in group
	for _, memberID := range group.AgentIDs {
		if memberID == agentID {
			return fmt.Errorf("agent already in group")
		}
	}

	// Add agent to group
	group.AgentIDs = append(group.AgentIDs, agentID)
	group.UpdatedAt = time.Now()

	if err := s.agentRepo.UpdateAgentGroup(ctx, group); err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	// Update agent's group membership
	agent.GroupIDs = append(agent.GroupIDs, groupID)
	if err := s.UpdateAgent(ctx, agent); err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	s.logger.Info("Agent added to group", "agent_id", agentID, "group_id", groupID)
	return nil
}

func (s *AgentManagementService) RemoveAgentFromGroup(ctx context.Context, agentID, groupID string) error {
	// Get group
	group, err := s.agentRepo.GetAgentGroup(ctx, groupID)
	if err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	// Remove agent from group
	var updatedAgentIDs []string
	found := false
	for _, memberID := range group.AgentIDs {
		if memberID != agentID {
			updatedAgentIDs = append(updatedAgentIDs, memberID)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("agent not in group")
	}

	group.AgentIDs = updatedAgentIDs
	group.UpdatedAt = time.Now()

	if err := s.agentRepo.UpdateAgentGroup(ctx, group); err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	// Update agent's group membership
	agent, err := s.GetAgent(ctx, agentID)
	if err == nil {
		var updatedGroupIDs []string
		for _, gID := range agent.GroupIDs {
			if gID != groupID {
				updatedGroupIDs = append(updatedGroupIDs, gID)
			}
		}
		agent.GroupIDs = updatedGroupIDs
		s.UpdateAgent(ctx, agent)
	}

	s.logger.Info("Agent removed from group", "agent_id", agentID, "group_id", groupID)
	return nil
}

func (s *AgentManagementService) GetAgentGroups(ctx context.Context, filter *models.AgentGroupFilter) ([]*models.AgentGroup, error) {
	return s.agentRepo.GetAgentGroups(ctx, filter)
}

func (s *AgentManagementService) GetAgentGroup(ctx context.Context, groupID string) (*models.AgentGroup, error) {
	return s.agentRepo.GetAgentGroup(ctx, groupID)
}

// Agent Backup and Recovery

func (s *AgentManagementService) BackupAgentConfiguration(ctx context.Context, agentID string) (*models.AgentBackup, error) {
	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	backup := &models.AgentBackup{
		ID:            uuid.New().String(),
		AgentID:       agentID,
		Configuration: agent.Configuration,
		Policies:      agent.AssignedPolicies,
		Metadata:      agent.Metadata,
		BackupTime:    time.Now(),
		Version:       agent.Version,
	}

	if err := s.agentRepo.CreateAgentBackup(ctx, backup); err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	s.logger.Info("Agent configuration backed up", "agent_id", agentID, "backup_id", backup.ID)
	return backup, nil
}

func (s *AgentManagementService) RestoreAgentConfiguration(ctx context.Context, agentID, backupID string) error {
	backup, err := s.agentRepo.GetAgentBackup(ctx, backupID)
	if err != nil {
		return fmt.Errorf("backup not found: %w", err)
	}

	if backup.AgentID != agentID {
		return fmt.Errorf("backup does not belong to agent")
	}

	agent, err := s.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Restore configuration
	agent.Configuration = backup.Configuration
	agent.AssignedPolicies = backup.Policies
	agent.Metadata = backup.Metadata
	agent.UpdatedAt = time.Now()

	if err := s.UpdateAgent(ctx, agent); err != nil {
		return fmt.Errorf("failed to restore agent: %w", err)
	}

	// Send configuration update to agent
	configUpdate := &models.AgentConfigurationUpdate{
		AgentID:       agentID,
		Configuration: backup.Configuration,
		UpdatedAt:     time.Now(),
	}

	if err := s.sendConfigurationUpdate(ctx, configUpdate); err != nil {
		s.logger.Error("Failed to send restored configuration to agent", "error", err)
	}

	s.logger.Info("Agent configuration restored", "agent_id", agentID, "backup_id", backupID)
	return nil
}

// Agent Maintenance and Cleanup

func (s *AgentManagementService) PerformAgentMaintenance(ctx context.Context) error {
	s.logger.Info("Starting agent maintenance")

	// Clean up expired sessions
	if err := s.cleanupExpiredSessions(ctx); err != nil {
		s.logger.Error("Failed to cleanup expired sessions", "error", err)
	}

	// Clean up old metrics
	if err := s.cleanupOldMetrics(ctx); err != nil {
		s.logger.Error("Failed to cleanup old metrics", "error", err)
	}

	// Clean up old commands
	if err := s.cleanupOldCommands(ctx); err != nil {
		s.logger.Error("Failed to cleanup old commands", "error", err)
	}

	// Clean up old logs
	if err := s.cleanupOldLogs(ctx); err != nil {
		s.logger.Error("Failed to cleanup old logs", "error", err)
	}

	s.logger.Info("Agent maintenance completed")
	return nil
}

func (s *AgentManagementService) cleanupExpiredSessions(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	expiredSessions := make([]string, 0)
	for agentID, session := range s.agentSessions {
		if time.Since(session.LastSeen) > s.sessionTimeout {
			session.Status = models.SessionStatusExpired
			expiredSessions = append(expiredSessions, agentID)
		}
	}

	for _, agentID := range expiredSessions {
		delete(s.agentSessions, agentID)
	}

	if len(expiredSessions) > 0 {
		s.logger.Info("Cleaned up expired sessions", "count", len(expiredSessions))
	}

	return nil
}

func (s *AgentManagementService) cleanupOldMetrics(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -30) // 30 days ago
	return s.agentRepo.DeleteOldMetrics(ctx, cutoffTime)
}

func (s *AgentManagementService) cleanupOldCommands(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7) // 7 days ago
	return s.agentRepo.DeleteOldCommands(ctx, cutoffTime)
}

func (s *AgentManagementService) cleanupOldLogs(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -14) // 14 days ago
	return s.agentRepo.DeleteOldLogs(ctx, cutoffTime)
}

// Service lifecycle management

func (s *AgentManagementService) Start(ctx context.Context) error {
	s.logger.Info("Starting Agent Management Service")

	// Start maintenance routine
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Daily maintenance
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.PerformAgentMaintenance(ctx); err != nil {
					s.logger.Error("Maintenance failed", "error", err)
				}
			}
		}
	}()

	return nil
}

func (s *AgentManagementService) Stop(ctx context.Context) error {
	s.logger.Info("Stopping Agent Management Service")

	// Cancel background tasks
	s.cancel()

	// Stop heartbeat monitoring
	if s.heartbeatTicker != nil {
		s.heartbeatTicker.Stop()
	}

	// Gracefully disconnect all agents
	s.mutex.RLock()
	agentIDs := make([]string, 0, len(s.agents))
	for agentID := range s.agents {
		agentIDs = append(agentIDs, agentID)
	}
	s.mutex.RUnlock()

	// Send shutdown commands to active agents
	shutdownCommand := &models.AgentCommand{
		Type:     models.CommandTypeShutdown,
		Priority: models.CommandPriorityHigh,
		Timeout:  30 * time.Second,
		Payload: map[string]interface{}{
			"reason": "service_shutdown",
		},
	}

	for _, agentID := range agentIDs {
		if err := s.SendCommandToAgent(ctx, agentID, shutdownCommand); err != nil {
			s.logger.Error("Failed to send shutdown command to agent",
				"agent_id", agentID,
				"error", err)
		}
	}

	// Wait for graceful shutdown
	time.Sleep(5 * time.Second)

	s.logger.Info("Agent Management Service stopped")
	return nil
}

// Health check for the service itself
func (s *AgentManagementService) HealthCheck(ctx context.Context) error {
	// Check if we can access the repository
	if _, err := s.agentRepo.GetAgents(ctx, &models.AgentFilter{Limit: 1}); err != nil {
		return fmt.Errorf("repository health check failed: %w", err)
	}

	// Check Kafka producer
	if err := s.kafkaProducer.HealthCheck(ctx); err != nil {
		return fmt.Errorf("kafka producer health check failed: %w", err)
	}

	// Check in-memory state
	s.mutex.RLock()
	agentCount := len(s.agents)
	sessionCount := len(s.agentSessions)
	s.mutex.RUnlock()

	s.logger.Debug("Agent Management Service health check",
		"agents", agentCount,
		"sessions", sessionCount)

	return nil
}

// GetServiceMetrics returns metrics about the service itself
func (s *AgentManagementService) GetServiceMetrics(ctx context.Context) *models.ServiceMetrics {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	activeAgents := 0
	inactiveAgents := 0
	activeSessions := 0

	for _, agent := range s.agents {
		if agent.Status == models.AgentStatusActive {
			activeAgents++
		} else {
			inactiveAgents++
		}
	}

	for _, session := range s.agentSessions {
		if session.Status == models.SessionStatusActive {
			activeSessions++
		}
	}

	return &models.ServiceMetrics{
		ServiceName:     "agent-management",
		TotalAgents:     len(s.agents),
		ActiveAgents:    activeAgents,
		InactiveAgents:  inactiveAgents,
		ActiveSessions:  activeSessions,
		TotalSessions:   len(s.agentSessions),
		LastUpdated:     time.Now(),
	}
}
