package config

import (
	"os"
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
)

func TestUseNonePrefixEnvAsDefaultValue(t *testing.T) {
	type demo struct {
		EndpointURL string `env:"ENDPOINT_URL"`
	}

	t.Run("fallback to ENDPOINT_URL when TEST_ENDPOINT_URL is empty or not set", func(t *testing.T) {
		os.Setenv("ENDPOINT_URL", "hello:27017/someservice")
		os.Setenv("TEST_ENDPOINT_URL", "")
		defer os.Unsetenv("ENDPOINT_URL")
		defer os.Unsetenv("TEST_ENDPOINT_URL")

		opts := env.Options{
			Prefix: prefix("TEST"),
		}
		de, err := parseEnv[demo](opts)

		assert.Nil(t, err)
		assert.Equal(t, "hello:27017/someservice", de.EndpointURL)
	})

	t.Run("override ENDPOINT_URL when TEST_ENDPOINT_URL have value", func(t *testing.T) {
		os.Setenv("ENDPOINT_URL", "hello:27017/someservice")
		os.Setenv("TEST_ENDPOINT_URL", "hey:27017/demoservice")
		defer os.Unsetenv("ENDPOINT_URL")
		defer os.Unsetenv("TEST_ENDPOINT_URL")

		opts := env.Options{
			Prefix: prefix("TEST"),
		}
		de, err := parseEnv[demo](opts)

		assert.Nil(t, err)
		assert.Equal(t, "hey:27017/demoservice", de.EndpointURL)
	})

	type defaultDemo struct {
		EndpointURL string `env:"ENDPOINT_URL" envDefault:"localhost"`
	}

	t.Run("it will use envDefault when env with Prefix is not set", func(t *testing.T) {
		os.Setenv("ENDPOINT_URL", "hello:27017/someservice")
		os.Setenv("TEST_ENDPOINT_URL", "")
		defer os.Unsetenv("ENDPOINT_URL")
		defer os.Unsetenv("TEST_ENDPOINT_URL")

		opts := env.Options{
			Prefix: prefix("TEST"),
		}
		de, err := parseEnv[defaultDemo](opts)

		assert.Nil(t, err)
		assert.Equal(t, "localhost", de.EndpointURL)
	})
}
