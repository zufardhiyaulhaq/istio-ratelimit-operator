package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSettings_WithDefaultValues(t *testing.T) {
	// Ensure environment variables are not set
	os.Unsetenv("RATE_LIMIT_SERVICE_IMAGE")
	os.Unsetenv("STATSD_EXPORTER_IMAGE")

	settings, err := NewSettings()

	assert.NoError(t, err)
	assert.Equal(t, "envoyproxy/ratelimit:5e1be594", settings.RateLimitServiceImage)
	assert.Equal(t, "prom/statsd-exporter:v0.23.1", settings.StatsdExporterImage)
}

func TestNewSettings_WithCustomEnvironmentVariables(t *testing.T) {
	testCases := []struct {
		name                          string
		rateLimitServiceImage         string
		statsdExporterImage           string
		expectedRateLimitServiceImage string
		expectedStatsdExporterImage   string
	}{
		{
			name:                          "custom ratelimit image only",
			rateLimitServiceImage:         "custom/ratelimit:v1.0.0",
			statsdExporterImage:           "",
			expectedRateLimitServiceImage: "custom/ratelimit:v1.0.0",
			expectedStatsdExporterImage:   "prom/statsd-exporter:v0.23.1",
		},
		{
			name:                          "custom statsd exporter image only",
			rateLimitServiceImage:         "",
			statsdExporterImage:           "custom/statsd:v2.0.0",
			expectedRateLimitServiceImage: "envoyproxy/ratelimit:5e1be594",
			expectedStatsdExporterImage:   "custom/statsd:v2.0.0",
		},
		{
			name:                          "both custom images",
			rateLimitServiceImage:         "my-registry/ratelimit:latest",
			statsdExporterImage:           "my-registry/statsd:latest",
			expectedRateLimitServiceImage: "my-registry/ratelimit:latest",
			expectedStatsdExporterImage:   "my-registry/statsd:latest",
		},
		{
			name:                          "images with full registry path",
			rateLimitServiceImage:         "gcr.io/my-project/ratelimit:v1.2.3",
			statsdExporterImage:           "docker.io/prom/statsd-exporter:v0.24.0",
			expectedRateLimitServiceImage: "gcr.io/my-project/ratelimit:v1.2.3",
			expectedStatsdExporterImage:   "docker.io/prom/statsd-exporter:v0.24.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up environment variables before each test
			os.Unsetenv("RATE_LIMIT_SERVICE_IMAGE")
			os.Unsetenv("STATSD_EXPORTER_IMAGE")

			// Set environment variables if provided
			if tc.rateLimitServiceImage != "" {
				os.Setenv("RATE_LIMIT_SERVICE_IMAGE", tc.rateLimitServiceImage)
			}
			if tc.statsdExporterImage != "" {
				os.Setenv("STATSD_EXPORTER_IMAGE", tc.statsdExporterImage)
			}

			settings, err := NewSettings()

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRateLimitServiceImage, settings.RateLimitServiceImage)
			assert.Equal(t, tc.expectedStatsdExporterImage, settings.StatsdExporterImage)

			// Clean up after test
			os.Unsetenv("RATE_LIMIT_SERVICE_IMAGE")
			os.Unsetenv("STATSD_EXPORTER_IMAGE")
		})
	}
}

func TestNewSettings_SettingsStructFields(t *testing.T) {
	// Test that Settings struct has expected fields and they are accessible
	settings := Settings{
		RateLimitServiceImage: "test-ratelimit-image",
		StatsdExporterImage:   "test-statsd-image",
	}

	assert.Equal(t, "test-ratelimit-image", settings.RateLimitServiceImage)
	assert.Equal(t, "test-statsd-image", settings.StatsdExporterImage)
}

func TestNewSettings_EmptyStringEnvironmentVariables(t *testing.T) {
	// When env vars are set to empty strings, envconfig should use defaults
	// because the fields have `required:"true"` with defaults
	os.Setenv("RATE_LIMIT_SERVICE_IMAGE", "")
	os.Setenv("STATSD_EXPORTER_IMAGE", "")

	settings, err := NewSettings()

	// envconfig treats empty string as "set" but will use the value as empty
	// Since these are required fields with defaults, empty values are used
	assert.NoError(t, err)
	assert.Equal(t, "", settings.RateLimitServiceImage)
	assert.Equal(t, "", settings.StatsdExporterImage)

	// Clean up
	os.Unsetenv("RATE_LIMIT_SERVICE_IMAGE")
	os.Unsetenv("STATSD_EXPORTER_IMAGE")
}

func TestNewSettings_ReturnsNewSettingsEachTime(t *testing.T) {
	// Ensure that each call to NewSettings returns a fresh Settings struct
	os.Unsetenv("RATE_LIMIT_SERVICE_IMAGE")
	os.Unsetenv("STATSD_EXPORTER_IMAGE")

	settings1, err1 := NewSettings()
	assert.NoError(t, err1)

	// Change environment variable
	os.Setenv("RATE_LIMIT_SERVICE_IMAGE", "changed-image:v2.0.0")

	settings2, err2 := NewSettings()
	assert.NoError(t, err2)

	// The two settings should be different
	assert.NotEqual(t, settings1.RateLimitServiceImage, settings2.RateLimitServiceImage)
	assert.Equal(t, "envoyproxy/ratelimit:5e1be594", settings1.RateLimitServiceImage)
	assert.Equal(t, "changed-image:v2.0.0", settings2.RateLimitServiceImage)

	// Clean up
	os.Unsetenv("RATE_LIMIT_SERVICE_IMAGE")
}

func TestNewSettings_SpecialCharactersInImageName(t *testing.T) {
	testCases := []struct {
		name                  string
		rateLimitServiceImage string
		statsdExporterImage   string
	}{
		{
			name:                  "image with sha256 digest",
			rateLimitServiceImage: "envoyproxy/ratelimit@sha256:abcdef1234567890",
			statsdExporterImage:   "prom/statsd-exporter@sha256:1234567890abcdef",
		},
		{
			name:                  "image with port in registry",
			rateLimitServiceImage: "localhost:5000/ratelimit:v1.0.0",
			statsdExporterImage:   "registry.example.com:8080/statsd:latest",
		},
		{
			name:                  "image with underscores and hyphens",
			rateLimitServiceImage: "my_registry/rate-limit_service:v1.0.0-beta",
			statsdExporterImage:   "my-registry/statsd_exporter:v0.23.1-rc1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("RATE_LIMIT_SERVICE_IMAGE", tc.rateLimitServiceImage)
			os.Setenv("STATSD_EXPORTER_IMAGE", tc.statsdExporterImage)

			settings, err := NewSettings()

			assert.NoError(t, err)
			assert.Equal(t, tc.rateLimitServiceImage, settings.RateLimitServiceImage)
			assert.Equal(t, tc.statsdExporterImage, settings.StatsdExporterImage)

			// Clean up
			os.Unsetenv("RATE_LIMIT_SERVICE_IMAGE")
			os.Unsetenv("STATSD_EXPORTER_IMAGE")
		})
	}
}

func TestSettings_ZeroValue(t *testing.T) {
	// Test zero value of Settings struct
	var settings Settings

	assert.Equal(t, "", settings.RateLimitServiceImage)
	assert.Equal(t, "", settings.StatsdExporterImage)
}
