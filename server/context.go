package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type M map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	Params     map[string]string
	handlers   []HandlerFunc
	index      int
	app        *Application
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// middleware function Next
func (c *Context) Next() {
	c.index++
	for s := len(c.handlers); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// get value from post form
func (c *Context) PostForm(key string) string {
	return c.Req.PostForm.Get(key)
}

// get from query
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// get from params
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// get from body
func (c *Context) Body() ([]byte, error) {
	return ioutil.ReadAll(c.Req.Body)
}

// get from body binding struct or map
func (c *Context) Bind(bind interface{}) error {
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(bind)
}

// set header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// set status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// fail to json
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, M{"message": err})
}

// string
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	_, _ = fmt.Fprintf(c.Writer, format, values)
}

// byte
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)

}

// write json data
func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// write html
func (c *Context) HTML(code int, name string, data interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	if err := c.app.templates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
