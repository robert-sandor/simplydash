package internal

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func Test_boolFromLabel(t *testing.T) {
	container := types.Container{Labels: map[string]string{}}

	t.Run("returns default value when label is missing", func(t *testing.T) {
		assert.Equal(t, false, boolFromLabel(container, "nonexistent_label", false))
		assert.Equal(t, true, boolFromLabel(container, "nonexistent_label", true))
	})

	t.Run("returns default value when label is not a bool", func(t *testing.T) {
		container.Labels["invalid_label"] = "not_a_bool"
		assert.Equal(t, false, boolFromLabel(container, "invalid_label", false))
		assert.Equal(t, true, boolFromLabel(container, "invalid_label", true))
	})

	t.Run("returns the bool value when label is a bool", func(t *testing.T) {
		container.Labels["valid_label"] = "false"
		assert.Equal(t, false, boolFromLabel(container, "valid_label", false))
		container.Labels["valid_label"] = "true"
		assert.Equal(t, true, boolFromLabel(container, "valie_label", true))
	})
}

func Test_durationFromLabel(t *testing.T) {
	// A mock container to use in testing.
	container := types.Container{
		Labels: map[string]string{},
	}

	t.Run("returns default value when label is not present", func(t *testing.T) {
		defaultValue := 10 * time.Second
		assert.Equal(t, defaultValue, durationFromLabel(container, "nonexistent_label", defaultValue))
	})

	t.Run("returns default value when label's value is not correctly formatted", func(t *testing.T) {
		defaultValue := 10 * time.Second
		container.Labels["invalid_label"] = "this_is_not_a_duration"
		assert.Equal(t, defaultValue, durationFromLabel(container, "invalid_label", defaultValue))
	})

	t.Run("returns default value when label's value is non-positive", func(t *testing.T) {
		defaultValue := 10 * time.Second
		container.Labels["nonpositive_label"] = "-5s"
		assert.Equal(t, defaultValue, durationFromLabel(container, "nonpositive_label", defaultValue))
	})

	t.Run("returns parsed value when label's value is correctly formatted and positive", func(t *testing.T) {
		container.Labels["valid_label"] = "5s"
		assert.Equal(t, 5*time.Second, durationFromLabel(container, "valid_label", 10*time.Second))
	})
}
