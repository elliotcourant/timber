package timber

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
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

func TestPrefix(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		With(map[string]interface{}{
			"things": "stuff",
		}).Prefix("12.0.0.1:54313").Debug("test")
	})

	t.Run("with null value", func(t *testing.T) {
		With(map[string]interface{}{
			"things":      "stuff",
			"otherThings": nil,
		}).Prefix("12.0.0.1:54313").Debug("test")
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

func TestLogNewline(t *testing.T) {
	Log(Level_Debug, `test
`)
}

func TestSetLevel(t *testing.T) {
	SetLevel(Level_Info)
	Log(Level_Debug, "test")   // Will not be written.
	Log(Level_Warning, "test") // Will be written.
}

func TestGetLevel(t *testing.T) {
	SetLevel(0)
	newLevel := Level_Warning
	firstLevel := GetLevel()
	SetLevel(newLevel)
	finalLevel := GetLevel()
	assert.Equal(t, newLevel, finalLevel)
	assert.NotEqual(t, firstLevel, finalLevel)
}

func TestConcurrent(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		log := New()
		numberOfRoutines := 100
		wg := new(sync.WaitGroup)
		wg.Add(numberOfRoutines)
		for i := 0; i < numberOfRoutines; i++ {
			go func(i int, lg Logger) {
				defer wg.Done()
				lg.With(Keys{
					"thread": i,
				}).Infof("this is a message")
			}(i, log)
		}
		wg.Wait()
	})

	t.Run("global logger", func(t *testing.T) {
		numberOfRoutines := 100
		wg := new(sync.WaitGroup)
		wg.Add(numberOfRoutines)
		for i := 0; i < numberOfRoutines; i++ {
			go func(i int) {
				defer wg.Done()
				With(Keys{
					"thread": i,
				}).Infof("this is a message")
			}(i)
		}
		wg.Wait()
	})

	t.Run("global logger levels", func(t *testing.T) {
		numberOfRoutines := 100
		wg := new(sync.WaitGroup)
		wg.Add(numberOfRoutines)
		for i := 0; i < numberOfRoutines; i++ {
			go func(i int) {
				defer wg.Done()
				SetLevel(Level_Error)
				With(Keys{
					"thread": i,
				}).Infof("this is a message")
			}(i)
		}
		wg.Wait()
	})

	t.Run("global logger levels changing concurrently", func(t *testing.T) {
		numberOfRoutines := 1000
		wg := new(sync.WaitGroup)
		wg.Add(numberOfRoutines)
		SetLevel(Level_Error)
		var done struct {
			done     bool
			doneSync sync.RWMutex
		}
		go func() {
			for {
				if GetLevel() == Level_Critical {
					SetLevel(Level_Error)
				} else {
					SetLevel(Level_Critical)
				}

				if func() bool {
					done.doneSync.RLock()
					defer done.doneSync.RUnlock()
					return done.done
				}() {
					return
				}
			}
		}()
		l := New()
		for i := 0; i < numberOfRoutines; i++ {
			go func(i int) {
				defer wg.Done()
				l.With(Keys{
					"thread": i,
				}).Infof("this is a message")
			}(i)
		}
		wg.Wait()
		done.doneSync.Lock()
		done.done = true
		done.doneSync.Unlock()
	})

	t.Run("logger levels changing concurrently", func(t *testing.T) {
		numberOfRoutines := 1000
		wg := new(sync.WaitGroup)
		wg.Add(numberOfRoutines)
		SetLevel(Level_Verbose)
		var done struct {
			done     bool
			doneSync sync.RWMutex
		}
		go func() {
			for {
				if GetLevel() == Level_Debug {
					SetLevel(Level_Verbose)
				} else {
					SetLevel(Level_Debug)
				}

				if func() bool {
					done.doneSync.RLock()
					defer done.doneSync.RUnlock()
					return done.done
				}() {
					return
				}
			}
		}()
		l := New()
		for i := 0; i < numberOfRoutines; i++ {
			go func(i int) {
				defer wg.Done()
				f := l.With(Keys{
					fmt.Sprintf("thread-%d", i): i,
				})
				f.Infof("test")
				f.Debugf("test")
			}(i)
		}
		wg.Wait()
		done.doneSync.Lock()
		done.done = true
		done.doneSync.Unlock()
	})
}
