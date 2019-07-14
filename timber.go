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
	return fmt.Sprintf("{ %s }", strings.Join(msg, ", "))
}

func (l *logger) log(stack int, lvl Level, m Keys, v ...interface{}) {
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
	lg := *l
	for k, v := range keys {
		lg.keys[k] = v
	}
	return &lg
}

func With(keys Keys) Logger {
	return defaultLogger.With(keys)
}
