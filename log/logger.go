package log

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

type Entry struct {
}
type Logger struct {
	mu     *sync.Mutex
	out    io.Writer
	prefix string
}

var std = New(os.Stdin, "")

func New(out io.Writer, prefix string) *Logger {
	return &Logger{out: out, prefix: prefix, mu: new(sync.Mutex)}
}

func (l *Logger) SetOutput(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = out
}

func (l *Logger) output(format string, obj ...interface{}) {

}

func (l *Logger) DEBUG() {

}

func (l *Logger) INFO() {

}

func (l *Logger) WARNING() {

}

func (l *Logger) ERROR() {

}

func (l *Logger) FATAL() {

}

func DEBUG() {
	std.DEBUG()
}

func INFO() {
	std.INFO()
}

func WARNING() {
	std.WARNING()

}

func ERROR() {
	std.ERROR()
}

func FATAL() {
	std.FATAL()
}
