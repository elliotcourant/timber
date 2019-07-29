package timber

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"strings"
	"sync"
)

const (
	defaultStackDepth = 3
)

var (
	level     = Level(0)
	levelSync sync.RWMutex
)

var (
	defaultLogger *logger
)

func init() {
	defaultLogger = &logger{
		stackDepth: defaultStackDepth,
		keys:       make(Keys),
	}
}

func New() Logger {
	return &logger{
		stackDepth: defaultStackDepth,
		keys:       make(Keys),
	}
}

func shouldLog(lvl Level) bool {
	levelSync.RLock()
	defer levelSync.RUnlock()
	return lvl >= level
}

type logger struct {
	stackDepth int
	keys       Keys
	keysLock   sync.RWMutex

	prefix     string
	prefixLock sync.RWMutex
}

func (l *logger) getKeysString(keys Keys) string {
	l.keysLock.RLock()
	defer l.keysLock.RUnlock()
	msg := make([]string, 0)
	for _, keySet := range append([]Keys{}, keys, l.keys) {
		for k, v := range keySet {
			// Exclude items where the value is null.
			if v == nil {
				continue
			}
			msg = append(msg, fmt.Sprintf(`%s: %v`, k, aurora.White(v)))
		}
	}
	if len(msg) == 0 {
		return ""
	}
	return fmt.Sprint(aurora.BrightBlack("{ "), strings.Join(msg, ", "), aurora.BrightBlack(" }"))
}

func (l *logger) getPrefixString() string {
	l.prefixLock.RLock()
	defer l.prefixLock.RUnlock()
	return l.prefix
}

func (l *logger) log(stack int, lvl Level, m Keys, v ...interface{}) {
	// If the message is below our level threshold then do not write it to
	// stdout.
	if !shouldLog(lvl) {
		return
	}
	k := l.getKeysString(m)
	p := l.getPrefixString()
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
	items := []interface{}{
		prefix,
	}
	if len(p) > 0 {
		items = append(items, aurora.White(fmt.Sprintf("[%s]", p)))
	}
	items = append(items, CallerInfo(stack))
	if len(k) > 0 {
		items = append(items, k, aurora.BrightBlack("|"))
	}
	items = append(items, fmt.Sprint(v...))
	fmt.Println(items...)
}

// SetDepth will change the number of stacks that will be skipped to find
// the filepath and line number of the executed code.
func (l *logger) SetDepth(depth int) Logger {
	l.stackDepth = defaultStackDepth + depth
	return l
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
	l.keysLock.Lock()
	defer l.keysLock.Unlock()
	lg := &logger{
		stackDepth: l.stackDepth,
		keys:       l.keys,
		prefix:     l.prefix,
	}
	for k, v := range keys {
		lg.keys[k] = v
	}
	return lg
}

// Prefix will add a small string before the file path.
func (l *logger) Prefix(prefix string) Logger {
	l.prefixLock.Lock()
	defer l.prefixLock.Unlock()
	l.prefix = prefix
	return l
}

func With(keys Keys) Logger {
	return defaultLogger.With(keys)
}

// SetLevel will set the minimum message level that will be output to stdout.
// This level is inherited by new logging instances created via With. But does
// not affect completely new logging instances.
func SetLevel(lvl Level) {
	levelSync.Lock()
	defer levelSync.Unlock()
	level = lvl
}

// GetLevel will return the current minimum logging level for the global
// logger.
func GetLevel() Level {
	levelSync.RLock()
	defer levelSync.RUnlock()
	return level
}

// Log will write a raw entry to the log, it accepts an array of interfaces which will
// be converted to strings if they are not already.
func Log(lvl Level, v ...interface{}) {
	defaultLogger.Log(lvl, v...)
}
