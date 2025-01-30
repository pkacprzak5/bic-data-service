package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {
	originalEnv := map[string]string{
		"PORT":        os.Getenv("PORT"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_NAME":     os.Getenv("DB_NAME"),
	}
	t.Cleanup(func() {
		for k, v := range originalEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	})

	t.Run("fallback values when no env vars", func(t *testing.T) {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_NAME")

		config := initConfig()

		assert.Equal(t, "8080", config.Port)
		assert.Equal(t, "5432", config.DB_Port)
		assert.Equal(t, "example_user", config.User)
		assert.Equal(t, "Passwd@1234", config.Password)
		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, "bicdatabase", config.Database)
	})

	t.Run("environment variables override fallbacks", func(t *testing.T) {
		t.Setenv("PORT", "9090")
		t.Setenv("DB_PORT", "6432")
		t.Setenv("DB_USER", "test_user")
		t.Setenv("DB_PASSWORD", "Test@1234")
		t.Setenv("DB_HOST", "test_host")
		t.Setenv("DB_NAME", "test_db")

		config := initConfig()

		assert.Equal(t, "9090", config.Port)
		assert.Equal(t, "6432", config.DB_Port)
		assert.Equal(t, "test_user", config.User)
		assert.Equal(t, "Test@1234", config.Password)
		assert.Equal(t, "test_host", config.Host)
		assert.Equal(t, "test_db", config.Database)
	})

	t.Run("valid .env file overrides fallbacks", func(t *testing.T) {
		tempDir := t.TempDir()
		oldWD, _ := os.Getwd()
		os.Chdir(tempDir)
		t.Cleanup(func() { os.Chdir(oldWD) })

		envContent := `PORT=7070
					DB_PORT=5433
					DB_USER=env_user
					DB_PASSWORD=Env@1234
					DB_HOST=env_host
					DB_NAME=env_db`
		require.NoError(t, os.WriteFile(".env", []byte(envContent), 0644))

		config := initConfig()

		assert.Equal(t, "7070", config.Port)
		assert.Equal(t, "5433", config.DB_Port)
		assert.Equal(t, "env_user", config.User)
		assert.Equal(t, "Env@1234", config.Password)
		assert.Equal(t, "env_host", config.Host)
		assert.Equal(t, "env_db", config.Database)
	})

	t.Run("invalid .env file uses fallbacks", func(t *testing.T) {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_NAME")

		tempDir := t.TempDir()
		oldWD, _ := os.Getwd()
		os.Chdir(tempDir)
		t.Cleanup(func() { os.Chdir(oldWD) })

		envContent := "invalid_ENV_file\n"
		require.NoError(t, os.WriteFile(".env", []byte(envContent), 0644))

		config := initConfig()

		assert.Equal(t, "8080", config.Port)
		assert.Equal(t, "5432", config.DB_Port)
		assert.Equal(t, "example_user", config.User)
		assert.Equal(t, "Passwd@1234", config.Password)
		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, "bicdatabase", config.Database)
	})

	t.Run(".env file overrides system environment", func(t *testing.T) {
		t.Setenv("PORT", "9090")

		tempDir := t.TempDir()
		oldWD, _ := os.Getwd()
		os.Chdir(tempDir)
		t.Cleanup(func() { os.Chdir(oldWD) })
		require.NoError(t, os.WriteFile(".env", []byte("PORT=7070"), 0644))

		config := initConfig()

		assert.Equal(t, "7070", config.Port, ".env value should override system environment")
	})
}
