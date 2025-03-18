package gee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")

	// 先序列化到缓冲区
	var buff bytes.Buffer
	encoder := json.NewEncoder(&buff)
	if err := encoder.Encode(obj); err != nil {
		// 确保状态码为 500，且未发送过头部
		c.Status(http.StatusInternalServerError)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 确认无错误后，设置状态码并写入数据
	c.Status(code)
	if _, err := c.Writer.Write(buff.Bytes()); err != nil {
		// TODO: 处理写入错误（如记录日志）
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
