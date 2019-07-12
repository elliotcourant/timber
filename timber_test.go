package timber

import (
	"testing"
)

func TestNew(t *testing.T) {
	logger := New()
	logger.Debug("test")
}

func TestWith(t *testing.T) {
	With(map[string]interface{}{
		"things": "stuff",
	}).Debug("test")
}

func TestLogger_Log(t *testing.T) {
	logger := New()
	logger.Log(Level_Debug, "test")
}
