package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/licaonfee/ivent"
	"github.com/licaonfee/ivent/stream"
)

//Level default event class, just like log levels
type Level int64

const (
	//PanicLevel max severity level this send event and panic
	PanicLevel = Level(iota)
	//FatalLevel equivalent to Print and os.Exit(1) this terminate application
	FatalLevel
	//ErrorLevel used to send errors that require action
	ErrorLevel
	//WarningLevel used to send informational non-critical messages
	WarningLevel
	//InfoLevel normal function
	InfoLevel
	//DebugLevel very verbose
	DebugLevel
	//TraceLevel more verbose tan Debug
	TraceLevel
)

//ensure interface
var (
	_ ivent.Class = new(Level)
)

var levels = []string{"Panic", "Fatal", "Error", "Warning", "Info", "Debug", "Trace"}

//Value return Level value as integer type
func (l Level) Value() int64 {
	return int64(l)
}

func (l Level) String() string {
	if int(l) >= len(levels) {
		return "Unknow"
	}
	return levels[l]
}

//ParseLevel return a Type from string or error
func ParseLevel(t string) (Level, error) {
	t1 := strings.ToLower(t)
	for i := 0; i < len(levels); i++ {
		if t1 == strings.ToLower(levels[i]) {
			return Level(i), nil
		}
	}
	return 0, fmt.Errorf("unknow level %s", t)
}

//Logger implementation of ivent.Client
type Logger struct {
	tags   map[string]string
	stream ivent.Stream
	mtx    sync.Mutex
	lv     *Level
}

func (l *Logger) logEvent(lv Level, tags map[string]string, data interface{}) ivent.Event {
	if lv < *l.lv {

	}
	nt := make(map[string]string)
	for k, v := range tags {
		nt[k] = v
	}
	return ivent.NewEvent(lv, nt, data)
}

func (l *Logger) send(e ivent.Event) {
	if e.T.Value() >= int64(*l.lv) {
		l.stream.Send(e)
	}
}

func (l *Logger) Trace(msg ...interface{}) {
	l.send(l.logEvent(TraceLevel, l.tags, msg))
}

func (l *Logger) Tracef(format string, msg ...interface{}) {
	l.Trace(fmt.Sprintf(format, msg...))
}

func (l *Logger) Debug(msg ...interface{}) {
	l.send(l.logEvent(DebugLevel, l.tags, msg))
}

func (l *Logger) Debugf(format string, msg ...interface{}) {
	l.Debug(fmt.Sprintf(format, msg...))
}

func (l *Logger) Info(msg ...interface{}) {
	l.send(l.logEvent(InfoLevel, l.tags, msg))
}

func (l *Logger) Infof(format string, msg ...interface{}) {
	l.Info(fmt.Sprintf(format, msg...))
}

func (l *Logger) Warning(msg ...interface{}) {
	l.send(l.logEvent(WarningLevel, l.tags, msg))
}

func (l *Logger) Warningf(format string, msg ...interface{}) {
	l.Warning(fmt.Sprintf(format, msg...))
}

func (l *Logger) Error(msg ...interface{}) {
	l.send(l.logEvent(ErrorLevel, l.tags, msg))
}

func (l *Logger) Errorf(format string, msg ...interface{}) {
	l.Error(fmt.Sprintf(format, msg...))
}

func (l *Logger) Fatal(msg ...interface{}) {
	l.send(l.logEvent(FatalLevel, l.tags, msg))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, msg ...interface{}) {
	l.Fatal(fmt.Sprintf(format, msg...))
}

func (l *Logger) Panic(msg ...interface{}) {
	l.send(l.logEvent(PanicLevel, l.tags, msg))
	panic(msg)
}

func (l *Logger) Panicf(format string, msg ...interface{}) {
	l.Panic(fmt.Sprintf(format, msg...))
}

func NewLogger() *Logger {
	l := &Logger{}
	l.tags = make(map[string]string)
	l.stream = stream.NewNoop()
	l.lv = new(Level)
	*l.lv = InfoLevel
	return l
}

func (l *Logger) WithStream(stream ivent.Stream) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.stream = stream

}

func (l *Logger) WithTags(tags map[string]string) *Logger {
	nl := l.Copy()
	for k, v := range tags {
		nl.tags[k] = v
	}
	return nl
}

func (l *Logger) WithTag(key string, value string) *Logger {
	nl := l.Copy()
	nl.tags[key] = value
	return nl
}

//Copy return a new Logger that share same Stream
func (l *Logger) Copy() *Logger {
	n := NewLogger()
	n.WithStream(l.stream)
	for k, v := range l.tags {
		n.tags[k] = v
	}
	return n
}

func (l *Logger) SetLogLevel(lv Level) {
	atomic.StoreInt64((*int64)(l.lv), int64(lv))
}
