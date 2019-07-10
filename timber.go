package timber

type logger struct {
	stackDepth int
	keys       map[string]interface{}
}

func (l *logger) log(stack int, lvl Level, m map[string]interface{}, v ...interface{}) {

}

func (l *logger) Log(lvl Level, v ...interface{}) {
	l.log(l.stackDepth, lvl, nil, v...)
}
