package logger

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var entryPool = &sync.Pool{
	New: func() interface{} {
		return &Entry{
			buf: make([]byte, 0, 500),
		}
	},
}

type Entry struct {
	buf   []byte
	w     io.Writer
	level Level
	done  func(msg string)
}

func putEntry(e *Entry) {
	const maxSize = 1 << 16
	if cap(e.buf) > maxSize {
		return
	}
	entryPool.Put(e)
}

func newEntry(w io.Writer, level Level) *Entry {
	e := entryPool.Get().(*Entry)
	e.buf = e.buf[:0]
	e.buf = e.AppendBeginMarker(e.buf)
	e.Str("level", level.String())
	e.Int("time", time.Now().Unix())
	_, file, line, _ := runtime.Caller(0)
	e.Str("file", file)
	e.Int("line", int64(line))
	e.level = level
	e.w = w
	return e
}

func (e *Entry) AppendBeginMarker(dst []byte) []byte {
	return append(dst, '{')
}

func (e *Entry) AppendEndMarker(dst []byte) []byte {
	return append(dst, '}')
}

func (e *Entry) AppendLineBreak(dst []byte) []byte {
	return append(dst, '\n')
}

func (e *Entry) appendKey(dst []byte, key string) []byte {
	if len(dst) > 1 && dst[len(dst)-1] != '{' {
		dst = append(dst, ',')
	}
	dst = e.appendString(dst, key)
	return append(dst, ':')
}

func (e *Entry) Msg(msg string) error {
	if e == nil {
		return nil
	}
	e.buf = e.appendString(e.appendKey(e.buf, "message"), msg)
	return e.Send()
}

// 输出数据
func (e *Entry) Send() error {
	if e == nil {
		return nil
	}
	defer putEntry(e)
	e.buf = e.AppendEndMarker(e.buf)
	e.buf = e.AppendLineBreak(e.buf)
	if e.w != nil {
		if _, err := e.w.Write(e.buf); err != nil {
			return err
		}
	}
	return nil
}

func (e *Entry) Str(key, val string) *Entry {
	if e == nil {
		fmt.Println("aa")
		return e
	}
	e.buf = e.appendString(e.appendKey(e.buf, key), val)
	return e
}

func (e *Entry) appendString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	dst = append(dst, s...)
	return append(dst, '"')
}

func (e *Entry) Int(key string, i int64) *Entry {
	if e == nil {
		return e
	}
	e.buf = e.appendInt(e.appendKey(e.buf, key), i)
	return e
}

func (e *Entry) appendInt(dst []byte, val int64) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func (e *Entry) Float(key string, f float64) *Entry {
	if e == nil {
		return e
	}
	e.buf = e.appendFloat(e.appendKey(e.buf, key), f)
	return e
}

func (e *Entry) appendFloat(dst []byte, val float64) []byte {
	return strconv.AppendFloat(dst, val, 'f', -1, 64)
}

func (e *Entry) Bool(key string, b bool) *Entry {
	if e == nil {
		return e
	}
	e.buf = e.appendBool(e.appendKey(e.buf, key), b)
	return e
}

func (e *Entry) appendBool(dst []byte, val bool) []byte {
	return strconv.AppendBool(dst, val)
}

func (e *Entry) Strs(key string, vals []string) *Entry {
	if e == nil {
		return e
	}
	e.buf = e.appendStrings(e.appendKey(e.buf, key), vals)
	return e
}

func (e *Entry) appendStrings(dst []byte, vals []string) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = e.appendString(dst, vals[0])
	if len(vals) > 1 {
		for _, val := range vals {
			dst = e.appendString(append(dst, ','), val)
		}
	}
	dst = append(dst, ']')
	return dst
}

func (e *Entry) Bytes(key string, val []byte) *Entry {
	if e == nil {
		return e
	}
	e.buf = e.appendBytes(e.appendKey(e.buf, key), val)
	return e
}

func (e *Entry) appendBytes(dst []byte, val []byte) []byte {
	dst = append(dst, '"')
	dst = append(dst, val...)
	return append(dst, '"')
}

func (e *Entry) Err(key string, val error) *Entry {
	if e == nil {
		return e
	}
	e.buf = e.appendString(e.appendKey(e.buf, key), val.Error())
	return e
}

func (e *Entry) Errs(key string, vals []error) *Entry {
	if e == nil {
		return e
	}
	var es []string
	for _, e := range vals {
		es = append(es, e.Error())
	}
	e.buf = e.appendStrings(e.appendKey(e.buf, key), es)
	return e
}
func (e *Entry) Bools(key string, bs []bool) *Entry {
	if e == nil {
		return e
	}
	e.buf = e.appendBools(e.appendKey(e.buf, key), bs)
	return e
}

func (e *Entry) appendBools(dst []byte, bs []bool) []byte {
	if len(bs) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendBool(dst, bs[0])
	if len(bs) > 1 {
		for _, b := range bs[1:] {
			dst = strconv.AppendBool(append(dst, ','), b)
		}
	}
	dst = append(dst, ']')
	return dst
}

func (e *Entry) Ints(key string, is []int64) *Entry {
	if e == nil {
		return e
	}
	e.buf = e.appendInts(e.appendKey(e.buf, key), is)
	return e
}

func (e *Entry) appendInts(dst []byte, is []int64) []byte {
	if len(is) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(is[0]), 10)
	if len(is) > 1 {
		for _, i := range is[1:] {
			dst = strconv.AppendInt(append(dst, ','), i, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}
