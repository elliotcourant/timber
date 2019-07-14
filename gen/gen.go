package main

import (
	"flag"
	"os"
	"strings"
	"text/template"
)

type data struct {
	Levels []string
}

func main() {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	var d data
	var items string

	var testFlag bool

	flag.StringVar(&items, "levels", "", "List of levels")
	flag.BoolVar(&testFlag, "test", false, "Generate tests")
	flag.Parse()

	if items != "" {
		d.Levels = strings.Split(items, ",")
	}

	a := template.Must(template.New("1").Funcs(funcMap).Parse(levelsTemplate))
	b := template.Must(template.New("2").Funcs(funcMap).Parse(testTemplate))

	switch testFlag {
	case true:
		b.Execute(os.Stdout, d)
	case false:
		a.Execute(os.Stdout, d)
	}
}

var levelsTemplate = `// Code generated by gen/gen.go - DO NOT EDIT.
// This code can be regenerated by running the go generate below.
//go:generate make generated

package timber

import (
	"fmt"
)

type Keys map[string]interface{}

type Level string

const ({{range $index, $level := .Levels}}
	Level_{{$level}} {{if (eq $index 0)}} Level{{end}} = "{{ $level | ToLower }}"{{end}}
)

type Logger interface {
{{range .Levels}}
	// {{.}} writes the provided string to the log.
	{{.}}(msg string)

	// {{.}}f writes a formatted string using the arguments provided to the log.
	{{.}}f(msg string, args ...interface{})

	// {{.}}Ex writes a formatted string using the arguments provided to the log
	// but also will prefix the log message with they keys provided to help print
	// runtime variables.
	{{.}}Ex(keys Keys, msg string, args ...interface{})
{{end}}
	// Log will write a raw entry to the log, it accepts an array of interfaces which will
	// be converted to strings if they are not already.
	Log(lvl Level, v ...interface{})

	// With will create a new Logger interface that will prefix all log entries written
	// from the new interface with the keys specified here. It will also include any
	// keys that are specified in the current Logger instance.
	// This means that you can chain multiple of these together to add/remove keys that
	// are written with every message.
	With(keys Keys) Logger
}{{range .Levels}}

// {{.}} writes the provided string to the log.
func (l *logger) {{.}}(msg string) {
	l.log(l.stackDepth, Level_{{.}}, nil, msg)
}

// {{.}}f writes a formatted string using the arguments provided to the log.
func (l *logger) {{.}}f(msg string, args ...interface{}) {
	l.log(l.stackDepth, Level_{{.}}, nil, fmt.Sprintf(msg, args...))
}

// {{.}}Ex writes a formatted string using the arguments provided to the log
// but also will prefix the log message with they keys provided to help print
// runtime variables.
func (l *logger) {{.}}Ex(keys Keys, msg string, args ...interface{}) {
	l.log(l.stackDepth, Level_{{.}}, keys, fmt.Sprintf(msg, args...))
}{{else}}
// No levels
{{end}}


{{range .Levels}}

// {{.}} writes the provided string to the log.
func {{.}}(msg string) {
	defaultLogger.log(defaultLogger.stackDepth, Level_{{.}}, nil, msg)
}

// {{.}}f writes a formatted string using the arguments provided to the log.
func {{.}}f(msg string, args ...interface{}) {
	defaultLogger.log(defaultLogger.stackDepth, Level_{{.}}, nil, fmt.Sprintf(msg, args...))
}

// {{.}}Ex writes a formatted string using the arguments provided to the log
// but also will prefix the log message with they keys provided to help print
// runtime variables.
func {{.}}Ex(keys Keys, msg string, args ...interface{}) {
	defaultLogger.log(defaultLogger.stackDepth, Level_{{.}}, keys, fmt.Sprintf(msg, args...))
}{{else}}
// No levels
{{end}}`

var testTemplate = `// Code generated by gen/gen.go - DO NOT EDIT.
// This code can be regenerated by running the go generate below.
//go:generate make generated

package timber

import (
	"testing"
)
{{range .Levels}}

func Test{{.}}(t *testing.T) {
	{{.}}("test")
}

func Test{{.}}f(t *testing.T) {
	{{.}}f("test %s", "format")
}

func Test{{.}}Ex(t *testing.T) {
	{{.}}Ex(map[string]interface{}{
		"thing": "stuff",
	}, "test")
}

func TestLogger_{{.}}(t *testing.T) {
	New().{{.}}("test")
}

func TestLogger_{{.}}f(t *testing.T) {
	New().{{.}}f("test %s", "format")
}

func TestLogger_{{.}}Ex(t *testing.T) {
	New().{{.}}Ex(map[string]interface{}{
		"thing": "stuff",
	}, "test")
}
{{else}}
// No levels
{{end}}`
