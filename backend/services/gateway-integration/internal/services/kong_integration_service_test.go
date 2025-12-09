package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"scopeapi.local/backend/services/gateway-integration/internal/models"
)

// MockKongClient is a mock implementation of KongClient
type MockKongClient struct {
	mock.Mock
}

func (m *MockKongClient) GetStatus() (*models.HealthStatus, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.HealthStatus), args.Error(1)
}

func (m *MockKongClient) GetServices() ([]map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockKongClient) GetRoutes() ([]map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockKongClient) GetPlugins() ([]map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockKongClient) GetConsumers() ([]map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockKongClient) GetUpstreams() ([]map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func TestNewKongIntegrationService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := NewKongIntegrationService()
		
		assert.NotNil(t, service)
		assert.NotNil(t, service.clientFactory)
	})
}

func TestKongIntegrationService_GetStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := &MockKongClient{}
		service := &KongIntegrationService{
			clientFactory: func(config map[string]interface{}) (KongClient, error) {
				return mockClient, nil
			},
		}

		expectedStatus := &models.HealthStatus{
			Status:    "healthy",
			Message:   "Kong is running",
			Timestamp: "2024-01-01T00:00:00Z",
		}

		integration := &models.Integration{
			Configuration: map[string]interface{}{
				"admin_url": "http://localhost:8001",
			},
		}

		mockClient.On("GetStatus").Return(expectedStatus, nil)

		status, err := service.GetStatus(context.Background(), integration)

		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, status)
		mockClient.AssertExpectations(t)
	})

	t.Run("client factory error", func(t *testing.T) {
		service := &KongIntegrationService{
			clientFactory: func(config map[string]interface{}) (KongClient, error) {
				return nil, assert.AnError
			},
		}

		integration := &models.Integration{
			Configuration: map[string]interface{}{
				"admin_url": "http://localhost:8001",
			},
		}

		status, err := service.GetStatus(context.Background(), integration)

		assert.Error(t, err)
		assert.Nil(t, status)
	})

	t.Run("client error", func(t *testing.T) {
		mockClient := &MockKongClient{}
		service := &KongIntegrationService{
			clientFactory: func(config map[string]interface{}) (KongClient, error) {
				return mockClient, nil
			},
		}

		integration := &models.Integration{
			Configuration: map[string]interface{}{
				"admin_url": "http://localhost:8001",
			},
		}

		mockClient.On("GetStatus").Return(nil, assert.AnError)

		status, err := service.GetStatus(context.Background(), integration)

		assert.Error(t, err)
		assert.Nil(t, status)
		mockClient.AssertExpectations(t)
	})
}

func TestKongIntegrationService_GetConfiguration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := &MockKongClient{}
		service := &KongIntegrationService{
			clientFactory: func(config map[string]interface{}) (KongClient, error) {
				return mockClient, nil
			},
		}

		expectedServices := []map[string]interface{}{
			{"id": "service1", "name": "api-service", "protocol": "http"},
		}
		expectedRoutes := []map[string]interface{}{
			{"id": "route1", "name": "api-route", "protocols": []string{"http"}},
		}
		expectedPlugins := []map[string]interface{}{
			{"id": "plugin1", "name": "rate-limiting", "enabled": true},
		}
		expectedConsumers := []map[string]interface{}{
			{"id": "consumer1", "username": "test-user"},
		}
		expectedUpstreams := []map[string]interface{}{
			{"id": "upstream1", "name": "api-upstream"},
		}

		integration := &models.Integration{
			Configuration: map[string]interface{}{
				"admin_url": "http://localhost:8001",
			},
		}

		mockClient.On("GetServices").Return(expectedServices, nil)
		mockClient.On("GetRoutes").Return(expectedRoutes, nil)
		mockClient.On("GetPlugins").Return(expectedPlugins, nil)
		mockClient.On("GetConsumers").Return(expectedConsumers, nil)
		mockClient.On("GetUpstreams").Return(expectedUpstreams, nil)

		config, err := service.GetConfiguration(context.Background(), integration)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, expectedServices, config["services"])
		assert.Equal(t, expectedRoutes, config["routes"])
		assert.Equal(t, expectedPlugins, config["plugins"])
		assert.Equal(t, expectedConsumers, config["consumers"])
		assert.Equal(t, expectedUpstreams, config["upstreams"])
		mockClient.AssertExpectations(t)
	})

	t.Run("services error", func(t *testing.T) {
		mockClient := &MockKongClient{}
		service := &KongIntegrationService{
			clientFactory: func(config map[string]interface{}) (KongClient, error) {
				return mockClient, nil
			},
		}

		integration := &models.Integration{
			Configuration: map[string]interface{}{
				"admin_url": "http://localhost:8001",
			},
		}

		mockClient.On("GetServices").Return(nil, assert.AnError)

		config, err := service.GetConfiguration(context.Background(), integration)

		assert.Error(t, err)
		assert.Nil(t, config)
		mockClient.AssertExpectations(t)
	})
}

func TestKongIntegrationService_SyncConfiguration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := &MockKongClient{}
		service := &KongIntegrationService{
			clientFactory: func(config map[string]interface{}) (KongClient, error) {
				return mockClient, nil
			},
		}

		expectedServices := []map[string]interface{}{
			{"id": "service1", "name": "api-service", "protocol": "http"},
		}
		expectedRoutes := []map[string]interface{}{
			{"id": "route1", "name": "api-route", "protocols": []string{"http"}},
		}

		integration := &models.Integration{
			Configuration: map[string]interface{}{
				"admin_url": "http://localhost:8001",
			},
		}

		mockClient.On("GetServices").Return(expectedServices, nil)
		mockClient.On("GetRoutes").Return(expectedRoutes, nil)
		mockClient.On("GetPlugins").Return([]map[string]interface{}{}, nil)
		mockClient.On("GetConsumers").Return([]map[string]interface{}{}, nil)
		mockClient.On("GetUpstreams").Return([]map[string]interface{}{}, nil)

		result, err := service.SyncConfiguration(context.Background(), integration)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "completed", result.Status)
		assert.Contains(t, result.Message, "Sync completed")
		assert.NotEmpty(t, result.Changes)
		mockClient.AssertExpectations(t)
	})

	t.Run("client factory error", func(t *testing.T) {
		service := &KongIntegrationService{
			clientFactory: func(config map[string]interface{}) (KongClient, error) {
				return nil, assert.AnError
			},
		}

		integration := &models.Integration{
			Configuration: map[string]interface{}{
				"admin_url": "http://localhost:8001",
			},
		}

		result, err := service.SyncConfiguration(context.Background(), integration)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
} 