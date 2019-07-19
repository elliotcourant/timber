package timber

import (
	"github.com/stretchr/testify/assert"
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

func TestLogger_SetDepth(t *testing.T) {
	logger := New()
	logger.SetDepth(1)
	logger.Log(Level_Debug, "test")
}

func TestLog(t *testing.T) {
	Log(Level_Debug, "test")
}

func TestLogger_SetLevel(t *testing.T) {
	logger := New()
	logger.SetLevel(Level_Info)
	logger.Log(Level_Debug, "test")   // Will not be written.
	logger.Log(Level_Warning, "test") // Will be written.
}

func TestLogger_GetLevel(t *testing.T) {
	logger := New()
	newLevel := Level_Warning
	firstLevel := logger.GetLevel()
	logger.SetLevel(newLevel)
	finalLevel := logger.GetLevel()
	assert.Equal(t, newLevel, finalLevel)
	assert.NotEqual(t, firstLevel, finalLevel)
}

func TestSetLevel(t *testing.T) {
	SetLevel(Level_Info)
	Log(Level_Debug, "test")   // Will not be written.
	Log(Level_Warning, "test") // Will be written.
}

func TestGetLevel(t *testing.T) {
	newLevel := Level_Warning
	firstLevel := GetLevel()
	SetLevel(newLevel)
	finalLevel := GetLevel()
	assert.Equal(t, newLevel, finalLevel)
	assert.NotEqual(t, firstLevel, finalLevel)
}

func TestSetDefaultLevel(t *testing.T) {
	newDefault := Level_Warning
	SetDefaultLevel(newDefault)
	logger := New()
	level := logger.GetLevel()
	assert.Equal(t, newDefault, level)
}

func TestGetDefaultLevel(t *testing.T) {
	originalDefault := GetDefaultLevel()
	newDefault := Level_Warning
	SetDefaultLevel(newDefault)
	finalDefault := GetDefaultLevel()
	assert.Equal(t, newDefault, finalDefault)
	assert.NotEqual(t, originalDefault, finalDefault)
}
