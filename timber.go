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

type loggerBase interface {
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

func (l *logger) log(stack int, lvl Level, m map[string]interface{}, v ...interface{}) {
	switch len(m) {
	case 0:
		fmt.Println(fmt.Sprintf("[%s]", strings.ToUpper(string(lvl))), CallerInfo(stack), fmt.Sprint(v...))
	default:
		fmt.Println(fmt.Sprintf("[%s]", strings.ToUpper(string(lvl))), CallerInfo(stack), m, fmt.Sprint(v...))
	}
}

func (l *logger) Log(lvl Level, v ...interface{}) {
	l.log(l.stackDepth, lvl, nil, v...)
}

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
