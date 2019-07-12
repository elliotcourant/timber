package timber

import (
	"testing"
)

func TestNew(t *testing.T) {
	logger := New()
	logger.Debug("test")
}

func TestWith(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		With(map[string]interface{}{
			"things": "stuff",
		}).Debug("test")
	})

	t.Run("with null value", func(t *testing.T) {
		With(map[string]interface{}{
			"things":      "stuff",
			"otherThings": nil,
		}).Debug("test")
	})
}

func TestLogger_Log(t *testing.T) {
	logger := New()
	logger.Log(Level_Debug, "test")
}
