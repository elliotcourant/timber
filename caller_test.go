package timber

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallerInfo(t *testing.T) {
	t.Run("bad stack index", func(t *testing.T) {
		info := CallerInfo(19999)
		assert.Equal(t, "unknown:0", info)
	})
}
