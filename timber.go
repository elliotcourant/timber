package timber

import (
	"fmt"
	"strings"
	"sync"
)

const (
	defaultStackDepth = 3
)

var (
	defaultLevel = Level(0)
	globalSync   sync.Mutex
)

var (
	defaultLogger *logger
)

func init() {
	defaultLogger = &logger{
		stackDepth: defaultStackDepth,
		keys:       make(Keys),
		level:      defaultLevel,
	}
}

func New() Logger {
	return &logger{
		stackDepth: defaultStackDepth,
		keys:       make(Keys),
		level:      defaultLevel,
	}
}

type logger struct {
	stackDepth int
	keys       Keys
	level      Level
	withLock   sync.RWMutex
}

func keys(keys ...Keys) string {
	msg := make([]string, 0)
	for _, keySet := range keys {
		for k, v := range keySet {
			// Exclude items where the value is null.
			if v == nil {
				continue
			}
			msg = append(msg, fmt.Sprintf(`[%s]: %v`, k, v))
		}
	}
	if len(msg) == 0 {
		return ""
	}
	return fmt.Sprintf("{ %s }", strings.Join(msg, ", "))
}

func (l *logger) log(stack int, lvl Level, m Keys, v ...interface{}) {
	// If the message is below our level threshold then do not write it to
	// stdout.
	if lvl < l.GetLevel() {
		return
	}
	k := keys(l.keys, m)
	foregroundColor, ok := foregroundColors[lvl]
	var prefix interface{}
	s := fmt.Sprintf("[%s]", shortLevelNames[lvl])
	if ok {
		prefix = foregroundColor(s)
	} else {
		prefix = s
	}
	backgroundColor, ok := backgroundColors[lvl]
	if ok {
		prefix = backgroundColor(s)
	}
	switch len(k) {
	case 0:
		fmt.Println(prefix, CallerInfo(stack), fmt.Sprint(v...))
	default:
		fmt.Println(prefix, CallerInfo(stack), k, fmt.Sprint(v...))
	}
}

func (l *logger) clone() *logger {
	l.withLock.Lock()
	defer l.withLock.Unlock()
	return &logger{
		stackDepth: l.stackDepth,
		level:      l.level,
		keys:       l.keys,
	}
}

// SetDepth will change the number of stacks that will be skipped to find
// the filepath and line number of the executed code.
func (l *logger) SetDepth(depth int) {
	l.stackDepth = defaultStackDepth + depth
}

// Log will write a raw entry to the log, it accepts an array of interfaces which will
// be converted to strings if they are not already.
func (l *logger) Log(lvl Level, v ...interface{}) {
	l.log(l.stackDepth, lvl, nil, v...)
}

// With will create a new Logger interface that will prefix all log entries written
// from the new interface with the keys specified here. It will also include any
// keys that are specified in the current Logger instance.
// This means that you can chain multiple of these together to add/remove keys that
// are written with every message.
func (l *logger) With(keys Keys) Logger {
	globalSync.Lock()
	defer globalSync.Unlock()
	lg := l.clone()
	lg.withLock.Lock()
	defer lg.withLock.Unlock()
	for k, v := range keys {
		lg.keys[k] = v
	}
	return lg
}

// SetLevel will set the minimum message level that will be output to stdout.
// This level is inherited by new logging instances created via With. But does
// not affect completely new logging instances.
func (l *logger) SetLevel(lvl Level) {
	l.withLock.Lock()
	defer l.withLock.Unlock()
	l.level = lvl
}

// GetLevel will return the current minimum logging level for this instance of
// the logger object.
func (l *logger) GetLevel() Level {
	l.withLock.RLock()
	defer l.withLock.RUnlock()
	return l.level
}

func With(keys Keys) Logger {
	return defaultLogger.With(keys)
}

// SetLevel will set the minimum message level that will be output to stdout.
// This level is inherited by new logging instances created via With. But does
// not affect completely new logging instances.
func SetLevel(lvl Level) {
	defaultLogger.SetLevel(lvl)
}

// GetLevel will return the current minimum logging level for the global
// logger.
func GetLevel() Level {
	return defaultLogger.GetLevel()
}

// SetDefaultLevel will define the level that is used for new loggers created.
// It will not change how the global logger or loggers that already exist behave.
// Existing loggers should be defined explicitly. The global logger should also
// be changed directly..
func SetDefaultLevel(lvl Level) {
	defaultLevel = lvl
}

// GetDefaultLevel will return the current level that is used when new loggers
// are created.
func GetDefaultLevel() Level {
	return defaultLevel
}

// Log will write a raw entry to the log, it accepts an array of interfaces which will
// be converted to strings if they are not already.
func Log(lvl Level, v ...interface{}) {
	defaultLogger.Log(lvl, v...)
}
