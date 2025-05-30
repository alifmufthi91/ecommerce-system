package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/cmd"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
)

type TestEnvironment struct {
	PostgresContainer testcontainers.Container
	APIBaseURL        string
	DatabaseURL       string
	Cleanup           func()
}

func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	ctx := context.Background()

	// 1. Start PostgreSQL container
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	databaseURL := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	// 2. Override environment variables for testing
	originalEnv := map[string]string{}
	testEnvVars := map[string]string{
		"DATABASE_URL":   databaseURL,
		"APP_PORT":       "0", // Let OS assign a free port
		"APP_ENV":        "test",
		"LOG_LEVEL":      "info",
		"JWT_SECRET_KEY": "test-secret-key",
		"DB_HOST":        host,
		"DB_PORT":        port.Port(),
		"DB_NAME":        "testdb",
		"DB_USER":        "testuser",
		"DB_PASSWORD":    "testpass",
		"DB_SSL_MODE":    "disable",
	}

	// Backup original env vars and set test ones
	for key, value := range testEnvVars {
		if original, exists := os.LookupEnv(key); exists {
			originalEnv[key] = original
		}
		os.Setenv(key, value)
	}

	// 3. Start API server in a goroutine
	serverDone := make(chan bool)
	var apiBaseURL string

	go func() {
		defer func() {
			serverDone <- true
		}()

		// Load test config
		testConfig := config.Config{
			App: config.App{
				Name: "warehouse-test",
				Port: "8080",
			},
			DB: config.DB{
				DSN: databaseURL,
			},
			Token: config.Token{
				JWTStatic: "test-static-token",
			},
		}

		// Start server (this will block)
		cmd.StartServerForTesting(&testConfig)
	}()

	// 4. Wait for API to be ready
	apiPort := "8080" // Default port, or detect the assigned port
	apiBaseURL = fmt.Sprintf("http://localhost:%s", apiPort)

	require.NoError(t, waitForAPI(apiBaseURL, 30*time.Second))

	// 5. Return test environment
	cleanup := func() {
		// Restore original environment variables
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
		for key := range testEnvVars {
			if _, exists := originalEnv[key]; !exists {
				os.Unsetenv(key)
			}
		}

		// Stop containers
		postgresContainer.Terminate(ctx)

		// Wait for server to stop (you might need to implement graceful shutdown)
		select {
		case <-serverDone:
		case <-time.After(10 * time.Second):
			t.Log("Server didn't stop gracefully")
		}
	}

	return &TestEnvironment{
		PostgresContainer: postgresContainer,
		APIBaseURL:        apiBaseURL,
		DatabaseURL:       databaseURL,
		Cleanup:           cleanup,
	}
}

func waitForAPI(baseURL string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("API did not become ready within %v", timeout)
		case <-ticker.C:
			resp, err := http.Get(fmt.Sprintf("%s/health", baseURL))
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				return nil
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}
