package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"scopeapi.local/backend/services/gateway-integration/internal/models"
)

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *IntegrationRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewIntegrationRepository(db)
	return db, mock, repo
}

func TestCreateIntegration(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		integration := &models.Integration{
			Name:        "Test Kong",
			Type:        models.GatewayTypeKong,
			Description: "Test integration",
			Status:      models.IntegrationStatusActive,
			Endpoints: []models.Endpoint{
				{
					URL:     "http://localhost:8001",
					Type:    "admin",
					Timeout: 30,
				},
			},
			Credentials: models.Credentials{
				Type:     models.CredentialTypeAPIKey,
				Username: "admin",
				Password: "password",
			},
			Configuration: map[string]interface{}{
				"version": "2.8.0",
			},
		}

		mock.ExpectQuery("INSERT INTO integrations").
			WithArgs(
				integration.Name,
				integration.Type,
				integration.Description,
				integration.Status,
				sqlmock.AnyArg(), // endpoints JSON
				sqlmock.AnyArg(), // credentials JSON
				sqlmock.AnyArg(), // configuration JSON
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("test-id"))

		err := repo.CreateIntegration(context.Background(), integration)

		assert.NoError(t, err)
		assert.Equal(t, "test-id", integration.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		integration := &models.Integration{
			Name: "Test Kong",
			Type: models.GatewayTypeKong,
		}

		mock.ExpectQuery("INSERT INTO integrations").
			WillReturnError(sql.ErrConnDone)

		err := repo.CreateIntegration(context.Background(), integration)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetIntegration(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		expectedID := "test-id"
		expectedIntegration := &models.Integration{
			ID:          expectedID,
			Name:        "Test Kong",
			Type:        models.GatewayTypeKong,
			Description: "Test integration",
			Status:      models.IntegrationStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "type", "description", "status", "endpoints", "credentials", 
			"configuration", "health_status", "last_sync_at", "created_at", "updated_at",
		}).AddRow(
			expectedIntegration.ID,
			expectedIntegration.Name,
			expectedIntegration.Type,
			expectedIntegration.Description,
			expectedIntegration.Status,
			`[{"url":"http://localhost:8001","type":"admin","timeout":30}]`,
			`{"type":"api_key","username":"admin","password":"password"}`,
			`{"version":"2.8.0"}`,
			`{"status":"healthy","message":"OK","timestamp":"2024-01-01T00:00:00Z"}`,
			time.Now(),
			expectedIntegration.CreatedAt,
			expectedIntegration.UpdatedAt,
		)

		mock.ExpectQuery("SELECT (.+) FROM integrations WHERE id = \\$1").
			WithArgs(expectedID).
			WillReturnRows(rows)

		integration, err := repo.GetIntegration(context.Background(), expectedID)

		assert.NoError(t, err)
		assert.NotNil(t, integration)
		assert.Equal(t, expectedID, integration.ID)
		assert.Equal(t, "Test Kong", integration.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		expectedID := "non-existent"

		mock.ExpectQuery("SELECT (.+) FROM integrations WHERE id = \\$1").
			WithArgs(expectedID).
			WillReturnError(sql.ErrNoRows)

		integration, err := repo.GetIntegration(context.Background(), expectedID)

		assert.Error(t, err)
		assert.Nil(t, integration)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetIntegrations(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	t.Run("success with filters", func(t *testing.T) {
		filters := map[string]interface{}{
			"type":   models.GatewayTypeKong,
			"status": models.IntegrationStatusActive,
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "type", "description", "status", "endpoints", "credentials", 
			"configuration", "health_status", "last_sync_at", "created_at", "updated_at",
		}).AddRow(
			"1", "Kong 1", models.GatewayTypeKong, "First Kong", models.IntegrationStatusActive,
			`[{"url":"http://localhost:8001"}]`, `{"type":"api_key"}`,
			`{"version":"2.8.0"}`, `{"status":"healthy"}`, time.Now(), time.Now(), time.Now(),
		).AddRow(
			"2", "Kong 2", models.GatewayTypeKong, "Second Kong", models.IntegrationStatusActive,
			`[{"url":"http://localhost:8002"}]`, `{"type":"api_key"}`,
			`{"version":"2.8.0"}`, `{"status":"healthy"}`, time.Now(), time.Now(), time.Now(),
		)

		mock.ExpectQuery("SELECT (.+) FROM integrations WHERE type = \\$1 AND status = \\$2").
			WithArgs(models.GatewayTypeKong, models.IntegrationStatusActive).
			WillReturnRows(rows)

		integrations, err := repo.GetIntegrations(context.Background(), filters)

		assert.NoError(t, err)
		assert.Len(t, integrations, 2)
		assert.Equal(t, "Kong 1", integrations[0].Name)
		assert.Equal(t, "Kong 2", integrations[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without filters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "name", "type", "description", "status", "endpoints", "credentials", 
			"configuration", "health_status", "last_sync_at", "created_at", "updated_at",
		}).AddRow(
			"1", "Kong", models.GatewayTypeKong, "Kong", models.IntegrationStatusActive,
			`[]`, `{}`, `{}`, `{}`, time.Now(), time.Now(), time.Now(),
		)

		mock.ExpectQuery("SELECT (.+) FROM integrations").
			WillReturnRows(rows)

		integrations, err := repo.GetIntegrations(context.Background(), nil)

		assert.NoError(t, err)
		assert.Len(t, integrations, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM integrations").
			WillReturnError(sql.ErrConnDone)

		integrations, err := repo.GetIntegrations(context.Background(), nil)

		assert.Error(t, err)
		assert.Nil(t, integrations)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateIntegration(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		integration := &models.Integration{
			ID:          "test-id",
			Name:        "Updated Kong",
			Type:        models.GatewayTypeKong,
			Description: "Updated description",
			Status:      models.IntegrationStatusActive,
		}

		mock.ExpectExec("UPDATE integrations SET").
			WithArgs(
				integration.Name,
				integration.Type,
				integration.Description,
				integration.Status,
				sqlmock.AnyArg(), // endpoints JSON
				sqlmock.AnyArg(), // credentials JSON
				sqlmock.AnyArg(), // configuration JSON
				sqlmock.AnyArg(), // updated_at
				integration.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateIntegration(context.Background(), integration)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		integration := &models.Integration{
			ID:   "non-existent",
			Name: "Updated Kong",
			Type: models.GatewayTypeKong,
		}

		mock.ExpectExec("UPDATE integrations SET").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateIntegration(context.Background(), integration)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteIntegration(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		expectedID := "test-id"

		mock.ExpectExec("DELETE FROM integrations WHERE id = \\$1").
			WithArgs(expectedID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteIntegration(context.Background(), expectedID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		expectedID := "non-existent"

		mock.ExpectExec("DELETE FROM integrations WHERE id = \\$1").
			WithArgs(expectedID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteIntegration(context.Background(), expectedID)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateHealthStatus(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		expectedID := "test-id"
		healthStatus := &models.HealthStatus{
			Status:    "healthy",
			Message:   "Connection successful",
			Timestamp: time.Now().Format(time.RFC3339),
		}

		mock.ExpectExec("UPDATE integrations SET health_status = \\$1, updated_at = \\$2 WHERE id = \\$3").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), expectedID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateHealthStatus(context.Background(), expectedID, healthStatus)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateLastSyncTime(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		expectedID := "test-id"
		lastSyncAt := time.Now()

		mock.ExpectExec("UPDATE integrations SET last_sync_at = \\$1, updated_at = \\$2 WHERE id = \\$3").
			WithArgs(lastSyncAt, sqlmock.AnyArg(), expectedID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateLastSyncTime(context.Background(), expectedID, lastSyncAt)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetIntegrationStats(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		// Mock the stats query
		rows := sqlmock.NewRows([]string{
			"total_count", "healthy_count", "unhealthy_count", "unknown_count",
			"gateway_type", "type_count",
		}).AddRow(
			5, 3, 1, 1, // Total stats
			"kong", 2, // By type
		).AddRow(
			5, 3, 1, 1,
			"nginx", 1,
		).AddRow(
			5, 3, 1, 1,
			"traefik", 1,
		).AddRow(
			5, 3, 1, 1,
			"envoy", 1,
		)

		mock.ExpectQuery("SELECT (.+) FROM integrations").
			WillReturnRows(rows)

		stats, err := repo.GetIntegrationStats(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 5, stats.TotalIntegrations)
		assert.Equal(t, 3, stats.HealthyCount)
		assert.Equal(t, 1, stats.UnhealthyCount)
		assert.Equal(t, 1, stats.UnknownCount)
		assert.Equal(t, 2, stats.ByType["kong"])
		assert.Equal(t, 1, stats.ByType["nginx"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM integrations").
			WillReturnError(sql.ErrConnDone)

		stats, err := repo.GetIntegrationStats(context.Background())

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
} 