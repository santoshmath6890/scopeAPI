package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
)

// MockIntegrationService is a mock implementation of IntegrationService
type MockIntegrationService struct {
	mock.Mock
}

func (m *MockIntegrationService) CreateIntegration(ctx context.Context, integration *models.Integration) error {
	args := m.Called(ctx, integration)
	return args.Error(0)
}

func (m *MockIntegrationService) GetIntegration(ctx context.Context, id string) (*models.Integration, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Integration), args.Error(1)
}

func (m *MockIntegrationService) GetIntegrations(ctx context.Context, filters map[string]interface{}) ([]*models.Integration, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Integration), args.Error(1)
}

func (m *MockIntegrationService) UpdateIntegration(ctx context.Context, id string, integration *models.Integration) error {
	args := m.Called(ctx, id, integration)
	return args.Error(0)
}

func (m *MockIntegrationService) DeleteIntegration(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockIntegrationService) TestIntegration(ctx context.Context, id string) (*models.HealthStatus, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.HealthStatus), args.Error(1)
}

func (m *MockIntegrationService) SyncIntegration(ctx context.Context, id string) (*models.SyncResult, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SyncResult), args.Error(1)
}

func (m *MockIntegrationService) GetIntegrationStats(ctx context.Context) (*models.IntegrationStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IntegrationStats), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockIntegrationService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	mockService := &MockIntegrationService{}
	handler := NewIntegrationHandler(mockService)
	
	// Setup routes
	api := router.Group("/api/v1")
	{
		integrations := api.Group("/integrations")
		{
			integrations.GET("", handler.GetIntegrations)
			integrations.GET("/:id", handler.GetIntegration)
			integrations.POST("", handler.CreateIntegration)
			integrations.PUT("/:id", handler.UpdateIntegration)
			integrations.DELETE("/:id", handler.DeleteIntegration)
			integrations.POST("/:id/test", handler.TestIntegration)
			integrations.POST("/:id/sync", handler.SyncIntegration)
			integrations.GET("/stats", handler.GetIntegrationStats)
		}
	}
	
	return router, mockService
}

func TestGetIntegrations(t *testing.T) {
	router, mockService := setupTestRouter()

	t.Run("success", func(t *testing.T) {
		expectedIntegrations := []*models.Integration{
			{
				ID:   "1",
				Name: "Test Kong",
				Type: models.GatewayTypeKong,
			},
		}

		mockService.On("GetIntegrations", mock.Anything, mock.Anything).Return(expectedIntegrations, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/integrations", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Integration
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, "Test Kong", response[0].Name)

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService.On("GetIntegrations", mock.Anything, mock.Anything).Return(nil, assert.AnError)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/integrations", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetIntegration(t *testing.T) {
	router, mockService := setupTestRouter()

	t.Run("success", func(t *testing.T) {
		expectedIntegration := &models.Integration{
			ID:   "1",
			Name: "Test Kong",
			Type: models.GatewayTypeKong,
		}

		mockService.On("GetIntegration", mock.Anything, "1").Return(expectedIntegration, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/integrations/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Integration
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Kong", response.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetIntegration", mock.Anything, "999").Return(nil, services.ErrIntegrationNotFound)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/integrations/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCreateIntegration(t *testing.T) {
	router, mockService := setupTestRouter()

	t.Run("success", func(t *testing.T) {
		integration := &models.Integration{
			Name: "New Kong",
			Type: models.GatewayTypeKong,
		}

		mockService.On("CreateIntegration", mock.Anything, mock.AnythingOfType("*models.Integration")).Return(nil)

		body, _ := json.Marshal(integration)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/integrations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request", func(t *testing.T) {
		invalidBody := `{"name": "", "type": "invalid"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/integrations", bytes.NewBufferString(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		integration := &models.Integration{
			Name: "New Kong",
			Type: models.GatewayTypeKong,
		}

		mockService.On("CreateIntegration", mock.Anything, mock.AnythingOfType("*models.Integration")).Return(assert.AnError)

		body, _ := json.Marshal(integration)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/integrations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateIntegration(t *testing.T) {
	router, mockService := setupTestRouter()

	t.Run("success", func(t *testing.T) {
		integration := &models.Integration{
			Name: "Updated Kong",
			Type: models.GatewayTypeKong,
		}

		mockService.On("UpdateIntegration", mock.Anything, "1", mock.AnythingOfType("*models.Integration")).Return(nil)

		body, _ := json.Marshal(integration)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/integrations/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		integration := &models.Integration{
			Name: "Updated Kong",
			Type: models.GatewayTypeKong,
		}

		mockService.On("UpdateIntegration", mock.Anything, "999", mock.AnythingOfType("*models.Integration")).Return(services.ErrIntegrationNotFound)

		body, _ := json.Marshal(integration)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/integrations/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteIntegration(t *testing.T) {
	router, mockService := setupTestRouter()

	t.Run("success", func(t *testing.T) {
		mockService.On("DeleteIntegration", mock.Anything, "1").Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/integrations/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("DeleteIntegration", mock.Anything, "999").Return(services.ErrIntegrationNotFound)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/integrations/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTestIntegration(t *testing.T) {
	router, mockService := setupTestRouter()

	t.Run("success", func(t *testing.T) {
		expectedHealth := &models.HealthStatus{
			Status:    "healthy",
			Message:   "Connection successful",
			Timestamp: "2024-01-01T00:00:00Z",
		}

		mockService.On("TestIntegration", mock.Anything, "1").Return(expectedHealth, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/integrations/1/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.HealthStatus
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("TestIntegration", mock.Anything, "999").Return(nil, services.ErrIntegrationNotFound)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/integrations/999/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSyncIntegration(t *testing.T) {
	router, mockService := setupTestRouter()

	t.Run("success", func(t *testing.T) {
		expectedSync := &models.SyncResult{
			Status:      "completed",
			Message:     "Sync completed successfully",
			Changes:     []models.Change{},
			LastSyncAt:  "2024-01-01T00:00:00Z",
		}

		mockService.On("SyncIntegration", mock.Anything, "1").Return(expectedSync, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/integrations/1/sync", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SyncResult
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "completed", response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("SyncIntegration", mock.Anything, "999").Return(nil, services.ErrIntegrationNotFound)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/integrations/999/sync", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetIntegrationStats(t *testing.T) {
	router, mockService := setupTestRouter()

	t.Run("success", func(t *testing.T) {
		expectedStats := &models.IntegrationStats{
			TotalIntegrations: 5,
			HealthyCount:      3,
			UnhealthyCount:    1,
			UnknownCount:      1,
			ByType: map[string]int{
				"kong":    2,
				"nginx":   1,
				"traefik": 1,
				"envoy":   1,
			},
		}

		mockService.On("GetIntegrationStats", mock.Anything).Return(expectedStats, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/integrations/stats", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.IntegrationStats
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 5, response.TotalIntegrations)
		assert.Equal(t, 3, response.HealthyCount)

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService.On("GetIntegrationStats", mock.Anything).Return(nil, assert.AnError)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/integrations/stats", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
} 