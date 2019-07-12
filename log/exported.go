package log

var (
	defLog = NewLogger()
)

func Panic(msg ...interface{}) {
	defLog.Panic(msg...)
}

func Fatal(msg ...interface{}) {
	defLog.Fatal(msg...)
}

func Error(msg ...interface{}) {
	defLog.Fatal(msg...)
}

func Warning(msg ...interface{}) {
	defLog.Warning(msg...)
}

func Info(msg ...interface{}) {
	defLog.Info(msg...)
}

func Debug(msg ...interface{}) {
	defLog.Debug(msg...)
}

func Trace(msg ...interface{}) {
	defLog.Trace(msg...)
}

func Default() *Logger {
	return defLog
}
