package timber

type Keys map[string]interface{}

type Level string

const (
	Trace    Level = "trace"
	Verbose        = "verbose"
	Debug          = "debug"
	Info           = "info"
	Warning        = "warning"
	Error          = "error"
	Critical       = "critical"
	Fatal          = "fatal"
)

type Logger interface {
	Trace(msg string)
	Tracef(msg string, args ...interface{})
	TraceEx(keys Keys, msg string, args ...interface{})

	Verbose(msg string)
	Verbosef(msg string, args ...interface{})
	VerboseEx(keys Keys, msg string, args ...interface{})

	Debug(msg string)
	Debugf(msg string, args ...interface{})
	DebugEx(keys Keys, msg string, args ...interface{})

	Info(msg string)
	Infof(msg string, args ...interface{})
	InfoEx(keys Keys, msg string, args ...interface{})

	Warning(msg string)
	Warningf(msg string, args ...interface{})
	WarningEx(keys Keys, msg string, args ...interface{})

	Error(msg string)
	Errorf(msg string, args ...interface{})
	ErrorEx(keys Keys, msg string, args ...interface{})

	Critical(msg string)
	Criticalf(msg string, args ...interface{})
	CriticalEx(keys Keys, msg string, args ...interface{})

	Fatal(msg string)
	Fatalf(msg string, args ...interface{})
	FatalEx(keys Keys, msg string, args ...interface{})

	Log(lvl Level, v ...interface{})

	With(keys Keys) Logger
}

type logger struct {
	stackDepth int
	keys       map[string]interface{}
}

func (l *logger) log(stack int, lvl Level, m map[string]interface{}, v ...interface{}) {

}

func (l *logger) Log(lvl Level, v ...interface{}) {
	l.log(l.stackDepth, lvl, nil, v...)
}
