package timber

import (
	"testing"
)

func TestDebug(t *testing.T) {
	Debug("test")
}

func TestDebugf(t *testing.T) {
	Debugf("test")
}

func TestDebugEx(t *testing.T) {
	DebugEx(map[string]interface{}{
		"thing": "stuff",
	}, "test")
}
