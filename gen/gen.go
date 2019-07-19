package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type LevelItem struct {
	Order           int     `json:"order"`
	Name            string  `json:"name"`
	ShortName       string  `json:"shortName"`
	ForegroundColor *string `json:"foregroundColor"`
	BackgroundColor *string `json:"backgroundColor"`
}

type Data struct {
	Levels []LevelItem
}

func main() {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	var testFlag bool
	var levelsPath string

	flag.StringVar(&levelsPath, "path", "levels.json", "Path to levels.json file.")
	flag.BoolVar(&testFlag, "test", false, "Generate tests")

	flag.Parse()

	a := template.Must(template.New("1").Funcs(funcMap).Parse(levelsTemplate))
	b := template.Must(template.New("2").Funcs(funcMap).Parse(testTemplate))

	if levelsPath == "" {
		panic("path cannot be blank")
	}

	var data Data

	if j, err := ioutil.ReadFile(levelsPath); err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal(j, &data); err != nil {
			panic(err)
		}
	}

	switch testFlag {
	case true:
		b.Execute(os.Stdout, data)
	case false:
		a.Execute(os.Stdout, data)
	}
}

var levelsTemplate = `// Code generated by gen/gen.go - DO NOT EDIT.
// This code can be regenerated by running the go generate below.
//go:generate make generated

package timber

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

type colorFunc func(arg interface{}) aurora.Value

type Keys map[string]interface{}

type Level int

const ({{range $index, $level := .Levels}}
	Level_{{$level.Name}} Level = {{$level.Order}}{{end}}
)

var (
	foregroundColors = map[Level]colorFunc{ {{range .Levels}}{{ if .ForegroundColor }}
		Level_{{.Name}}: aurora.{{.ForegroundColor}},{{end}}{{end}} 
	}

	backgroundColors = map[Level]colorFunc{ {{range .Levels}}{{ if .BackgroundColor }}
		Level_{{.Name}}: aurora.Bg{{.BackgroundColor}},{{end}}{{end}} 
	}

	levelNames = map[Level]string{ {{range .Levels}}
		Level_{{.Name}}: "{{.Name}}",{{end}}
	}

	shortLevelNames = map[Level]string{ {{range .Levels}}
		Level_{{.Name}}: "{{.ShortName}}",{{end}}
	}
)

type Logger interface { {{range .Levels}}
	// {{.Name}} writes the provided string to the log.
	{{.Name}}(msg interface{})

	// {{.Name}}f writes a formatted string using the arguments provided to the log.
	{{.Name}}f(msg string, args ...interface{})

	// {{.Name}}Ex writes a formatted string using the arguments provided to the log
	// but also will prefix the log message with they keys provided to help print
	// runtime variables.
	{{.Name}}Ex(keys Keys, msg string, args ...interface{})
{{end}}
	// SetDepth will change the number of stacks that will be skipped to find
	// the filepath and line number of the executed code.
	SetDepth(depth int)

	// Log will write a raw entry to the log, it accepts an array of interfaces which will
	// be converted to strings if they are not already.
	Log(lvl Level, v ...interface{})

	// With will create a new Logger interface that will prefix all log entries written
	// from the new interface with the keys specified here. It will also include any
	// keys that are specified in the current Logger instance.
	// This means that you can chain multiple of these together to add/remove keys that
	// are written with every message.
	With(keys Keys) Logger

	// SetLevel will set the minimum message level that will be output to stdout.
	// This level is inherited by new logging instances created via With. But does
	// not affect completely new logging instances.
	SetLevel(lvl Level)

	// GetLevel will return the current minimum logging level for this instance of
	// the logger object.
	GetLevel() Level
}{{range .Levels}}

// {{.Name}} writes the provided string to the log.
func (l *logger) {{.Name}}(msg interface{}) {
	l.log(l.stackDepth, Level_{{.Name}}, nil, msg)
}

// {{.Name}}f writes a formatted string using the arguments provided to the log.
func (l *logger) {{.Name}}f(msg string, args ...interface{}) {
	l.log(l.stackDepth, Level_{{.Name}}, nil, fmt.Sprintf(msg, args...))
}

// {{.Name}}Ex writes a formatted string using the arguments provided to the log
// but also will prefix the log message with they keys provided to help print
// runtime variables.
func (l *logger) {{.Name}}Ex(keys Keys, msg string, args ...interface{}) {
	l.log(l.stackDepth, Level_{{.Name}}, keys, fmt.Sprintf(msg, args...))
}{{else}}
// No levels
{{end}}


{{range .Levels}}

// {{.Name}} writes the provided string to the log.
func {{.Name}}(msg interface{}) {
	defaultLogger.log(defaultLogger.stackDepth, Level_{{.Name}}, nil, msg)
}

// {{.Name}}f writes a formatted string using the arguments provided to the log.
func {{.Name}}f(msg string, args ...interface{}) {
	defaultLogger.log(defaultLogger.stackDepth, Level_{{.Name}}, nil, fmt.Sprintf(msg, args...))
}

// {{.Name}}Ex writes a formatted string using the arguments provided to the log
// but also will prefix the log message with they keys provided to help print
// runtime variables.
func {{.Name}}Ex(keys Keys, msg string, args ...interface{}) {
	defaultLogger.log(defaultLogger.stackDepth, Level_{{.Name}}, keys, fmt.Sprintf(msg, args...))
}{{else}}
// No levels
{{end}}
`

var testTemplate = `// Code generated by gen/gen.go - DO NOT EDIT.
// This code can be regenerated by running the go generate below.
//go:generate make generated

package timber

import (
	"testing"
)
{{range .Levels}}

func Test{{.Name}}(t *testing.T) {
	{{.Name}}("test")
}

func Test{{.Name}}f(t *testing.T) {
	{{.Name}}f("test %s", "format")
}

func Test{{.Name}}Ex(t *testing.T) {
	{{.Name}}Ex(map[string]interface{}{
		"thing": "stuff",
	}, "test")
}

func TestLogger_{{.Name}}(t *testing.T) {
	New().{{.Name}}("test")
}

func TestLogger_{{.Name}}f(t *testing.T) {
	New().{{.Name}}f("test %s", "format")
}

func TestLogger_{{.Name}}Ex(t *testing.T) {
	New().{{.Name}}Ex(map[string]interface{}{
		"thing": "stuff",
	}, "test")
}
{{else}}
// No levels
{{end}}`
