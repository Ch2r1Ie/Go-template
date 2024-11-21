package config

import (
	"testing"
)

func TestPrefixEnv(t *testing.T) {
	t.Run("Prefix is empty", func(t *testing.T) {
		prefix := prefix("")

		if prefix != "" {
			t.Error("Expected empty string but got", prefix)
		}
	})

	t.Run("Prefix is not empty", func(t *testing.T) {
		prefix := prefix("LOCAL")

		if prefix != "LOCAL_" {
			t.Error("Expected LOCAL_ but got", prefix)
		}
	})
}
