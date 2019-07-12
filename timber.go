package timber

import (
	"fmt"
	"strings"
)

const (
	defaultStackDepth = 3
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

type LoggerBase interface {
	// Log will write a raw entry to the log, it accepts an array of interfaces which will
	// be converted to strings if they are not already.
	Log(lvl Level, v ...interface{})

	// With will create a new Logger interface that will prefix all log entries written
	// from the new interface with the keys specified here. It will also include any
	// keys that are specified in the current Logger instance.
	// This means that you can chain multiple of these together to add/remove keys that
	// are written with every message.
	With(keys Keys) Logger
}

func New() Logger {
	return &logger{
		stackDepth: defaultStackDepth,
		keys:       make(Keys),
	}
}

type logger struct {
	stackDepth int
	keys       Keys
}

func keys(keys ...Keys) string {
	msg := make([]string, 0)
	for _, keySet := range keys {
		for k, v := range keySet {
			// Exclude items where the value is null.
			if v == nil {
				continue
			}
			msg = append(msg, fmt.Sprintf(`%s: %v`, k, v))
		}
	}
	if len(msg) == 0 {
		return ""
	}
	return fmt.Sprintf("[ %s ]", strings.Join(msg, ", "))
}

func (l *logger) log(stack int, lvl Level, m Keys, v ...interface{}) {
	k := keys(l.keys, m)
	switch len(k) {
	case 0:
		fmt.Println(fmt.Sprintf("[%s]", strings.ToUpper(string(lvl))), CallerInfo(stack), fmt.Sprint(v...))
	default:
		fmt.Println(fmt.Sprintf("[%s]", strings.ToUpper(string(lvl))), CallerInfo(stack), k, fmt.Sprint(v...))
	}
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
	lg := *l
	for k, v := range keys {
		lg.keys[k] = v
	}
	return &lg
}

func With(keys Keys) Logger {
	return defaultLogger.With(keys)
}
