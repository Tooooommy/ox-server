package logger

import (
	"io"
	"os"
	"sync"
)

// https://github.com/rs/zerolog/blob/d9df1802de8b99f4e97aa88056ea5d2cea663bb4/internal/cbor/string.go#L20
/*
	var bs []byte
	bs = append(bs, '{')
	bs = append(bs, []byte("user")...)
	bs = append(bs, ':')
	bs = append(bs, []byte("tommy")...)
	bs = append(bs, '}')
	fmt.Println(string(bs))
*/

type Level uint8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
	Disabled
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case NoLevel:
		return ""
	}
	return ""
}

type Logger struct {
	mu    *sync.Mutex
	out   io.Writer
	level Level
}

var std = New(os.Stdout)

func New(out io.Writer) *Logger {
	return &Logger{out: out, mu: new(sync.Mutex)}
}

func (l *Logger) SetOutput(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = out
}

func (l *Logger) newEntry() *Entry {
	return newEntry(l.out, l.level)
}

func (l *Logger) DEBUG() *Entry {
	return l.newEntry()

}

func (l *Logger) INFO() *Entry {
	return l.newEntry()

}

func (l *Logger) WARN() *Entry {
	return l.newEntry()
}

func (l *Logger) ERROR() *Entry {
	return l.newEntry()
}

func (l *Logger) FATAL() *Entry {
	return l.newEntry()
}

func (l *Logger) Panic() *Entry {
	return l.newEntry()
}

func DEBUG() *Entry {
	return std.DEBUG()
}

func INFO() *Entry {
	return std.INFO()
}

func WARN() *Entry {
	return std.WARN()

}

func ERROR() *Entry {
	return std.ERROR()
}

func FATAL() *Entry {
	return std.FATAL()
}

func Panic() *Entry {
	return std.Panic()
}
