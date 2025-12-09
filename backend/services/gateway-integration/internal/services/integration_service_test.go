package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
)

// Mock repository for testing
type MockIntegrationRepository struct {
	mock.Mock
}

func (m *MockIntegrationRepository) CreateIntegration(ctx context.Context, integration *models.Integration) error {
	args := m.Called(ctx, integration)
	return args.Error(0)
}

func (m *MockIntegrationRepository) GetIntegration(ctx context.Context, id string) (*models.Integration, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Integration), args.Error(1)
}

func (m *MockIntegrationRepository) GetIntegrations(ctx context.Context, filters map[string]interface{}) ([]*models.Integration, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Integration), args.Error(1)
}

func (m *MockIntegrationRepository) UpdateIntegration(ctx context.Context, integration *models.Integration) error {
	args := m.Called(ctx, integration)
	return args.Error(0)
}

func (m *MockIntegrationRepository) DeleteIntegration(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockIntegrationRepository) UpdateIntegrationHealth(ctx context.Context, id string, health *models.HealthStatus) error {
	args := m.Called(ctx, id, health)
	return args.Error(0)
}

func (m *MockIntegrationRepository) UpdateIntegrationLastSync(ctx context.Context, id string, lastSync time.Time) error {
	args := m.Called(ctx, id, lastSync)
	return args.Error(0)
}

func (m *MockIntegrationRepository) GetIntegrationStats(ctx context.Context) (*models.IntegrationStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IntegrationStats), args.Error(1)
}

// Mock Kafka producer for testing
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) Produce(topic string, key []byte, value []byte) error {
	args := m.Called(topic, key, value)
	return args.Error(0)
}

func (m *MockKafkaProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Mock logger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Fatal(msg string, args ...interface{}) {
	m.Called(msg, args)
}

// Test data helpers
func createTestIntegration() *models.Integration {
	return &models.Integration{
		ID:     "test-integration-id",
		Name:   "Test Kong Gateway",
		Type:   models.GatewayTypeKong,
		Status: models.IntegrationStatusActive,
		Config: map[string]interface{}{
			"admin_url": "http://kong-admin:8001",
			"proxy_url": "http://kong-proxy:8000",
		},
		Endpoints: []models.Endpoint{
			{
				ID:       "1",
				Name:     "Admin API",
				URL:      "http://kong-admin:8001",
				Protocol: "http",
				Port:     8001,
				Timeout:  30000,
			},
		},
		Health: &models.HealthStatus{
			Status:    "healthy",
			Message:   "Kong is running",
			LastCheck: time.Now(),
			Latency:   45,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Test cases
func TestNewIntegrationService(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.Equal(t, mockProducer, service.kafkaProducer)
	assert.Equal(t, mockLogger, service.logger)
	assert.NotNil(t, service.gatewayClients)
}

func TestIntegrationService_CreateIntegration_Success(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	integration := createTestIntegration()

	mockRepo.On("CreateIntegration", ctx, integration).Return(nil)
	mockProducer.On("Produce", "integration_events", mock.Anything, mock.Anything).Return(nil)
	mockLogger.On("Info", "Integration created", mock.Anything).Return()

	result, err := service.CreateIntegration(ctx, integration)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, integration.ID, result.ID)
	assert.Equal(t, integration.Name, result.Name)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestIntegrationService_CreateIntegration_RepositoryError(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	integration := createTestIntegration()

	expectedError := assert.AnError
	mockRepo.On("CreateIntegration", ctx, integration).Return(expectedError)

	result, err := service.CreateIntegration(ctx, integration)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestIntegrationService_GetIntegration_Success(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	integrationID := "test-integration-id"
	expectedIntegration := createTestIntegration()

	mockRepo.On("GetIntegration", ctx, integrationID).Return(expectedIntegration, nil)

	result, err := service.GetIntegration(ctx, integrationID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedIntegration.ID, result.ID)
	assert.Equal(t, expectedIntegration.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestIntegrationService_GetIntegration_NotFound(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	integrationID := "non-existent-id"

	mockRepo.On("GetIntegration", ctx, integrationID).Return(nil, assert.AnError)

	result, err := service.GetIntegration(ctx, integrationID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestIntegrationService_GetIntegrations_Success(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	filters := map[string]interface{}{
		"type":   "kong",
		"status": "active",
	}
	expectedIntegrations := []*models.Integration{createTestIntegration()}

	mockRepo.On("GetIntegrations", ctx, filters).Return(expectedIntegrations, nil)

	result, err := service.GetIntegrations(ctx, filters)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedIntegrations[0].ID, result[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestIntegrationService_UpdateIntegration_Success(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	integration := createTestIntegration()
	integration.Name = "Updated Kong Gateway"

	mockRepo.On("UpdateIntegration", ctx, integration).Return(nil)
	mockProducer.On("Produce", "integration_events", mock.Anything, mock.Anything).Return(nil)
	mockLogger.On("Info", "Integration updated", mock.Anything).Return()

	result, err := service.UpdateIntegration(ctx, integration)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, integration.Name, result.Name)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestIntegrationService_DeleteIntegration_Success(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	integrationID := "test-integration-id"

	mockRepo.On("DeleteIntegration", ctx, integrationID).Return(nil)
	mockProducer.On("Produce", "integration_events", mock.Anything, mock.Anything).Return(nil)
	mockLogger.On("Info", "Integration deleted", mock.Anything).Return()

	err := service.DeleteIntegration(ctx, integrationID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestIntegrationService_TestIntegration_Success(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	integrationID := "test-integration-id"
	integration := createTestIntegration()
	expectedHealth := &models.HealthStatus{
		Status:    "healthy",
		Message:   "Connection successful",
		LastCheck: time.Now(),
		Latency:   25,
	}

	mockRepo.On("GetIntegration", ctx, integrationID).Return(integration, nil)
	mockRepo.On("UpdateIntegrationHealth", ctx, integrationID, expectedHealth).Return(nil)

	// Mock the gateway client test
	service.gatewayClients[models.GatewayTypeKong] = func(logger interface{}) GatewayClient {
		return &MockGatewayClient{
			testConnectionFunc: func(ctx context.Context, integration *models.Integration) error {
				return nil
			},
			getStatusFunc: func(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error) {
				return expectedHealth, nil
			},
		}
	}

	health, err := service.TestIntegration(ctx, integrationID)

	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Equal(t, expectedHealth.Status, health.Status)
	mockRepo.AssertExpectations(t)
}

func TestIntegrationService_SyncIntegration_Success(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	integrationID := "test-integration-id"
	integration := createTestIntegration()
	expectedResult := &models.SyncResult{
		Success:   true,
		Message:   "Configuration synchronized successfully",
		Changes:   []models.Change{},
		Timestamp: time.Now(),
		Duration:  time.Duration(100 * time.Millisecond),
	}

	mockRepo.On("GetIntegration", ctx, integrationID).Return(integration, nil)
	mockRepo.On("UpdateIntegrationLastSync", ctx, integrationID, mock.Anything).Return(nil)

	// Mock the gateway client sync
	service.gatewayClients[models.GatewayTypeKong] = func(logger interface{}) GatewayClient {
		return &MockGatewayClient{
			syncConfigurationFunc: func(ctx context.Context, integration *models.Integration) (*models.SyncResult, error) {
				return expectedResult, nil
			},
		}
	}

	result, err := service.SyncIntegration(ctx, integrationID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResult.Success, result.Success)
	assert.Equal(t, expectedResult.Message, result.Message)
	mockRepo.AssertExpectations(t)
}

func TestIntegrationService_GetIntegrationStats_Success(t *testing.T) {
	mockRepo := &MockIntegrationRepository{}
	mockProducer := &MockKafkaProducer{}
	mockLogger := &MockLogger{}

	service := NewIntegrationService(mockRepo, mockProducer, mockLogger)
	ctx := context.Background()
	expectedStats := &models.IntegrationStats{
		Total:        5,
		Active:       3,
		Error:        1,
		Pending:      1,
		KongCount:    2,
		NginxCount:   1,
		TraefikCount: 1,
		EnvoyCount:   0,
		HAProxyCount: 1,
	}

	mockRepo.On("GetIntegrationStats", ctx).Return(expectedStats, nil)

	result, err := service.GetIntegrationStats(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStats.Total, result.Total)
	assert.Equal(t, expectedStats.Active, result.Active)
	assert.Equal(t, expectedStats.KongCount, result.KongCount)
	mockRepo.AssertExpectations(t)
}

// Mock GatewayClient for testing
type MockGatewayClient struct {
	testConnectionFunc    func(ctx context.Context, integration *models.Integration) error
	getStatusFunc        func(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error)
	syncConfigurationFunc func(ctx context.Context, integration *models.Integration) (*models.SyncResult, error)
}

func (m *MockGatewayClient) TestConnection(ctx context.Context, integration *models.Integration) error {
	if m.testConnectionFunc != nil {
		return m.testConnectionFunc(ctx, integration)
	}
	return nil
}

func (m *MockGatewayClient) GetStatus(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error) {
	if m.getStatusFunc != nil {
		return m.getStatusFunc(ctx, integration)
	}
	return nil, nil
}

func (m *MockGatewayClient) SyncConfiguration(ctx context.Context, integration *models.Integration) (*models.SyncResult, error) {
	if m.syncConfigurationFunc != nil {
		return m.syncConfigurationFunc(ctx, integration)
	}
	return nil, nil
} 