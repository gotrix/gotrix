package gotrix

import (
	"io"
	"net/http"
)

// Context interface.
type Context interface {
	Writer() io.Writer
	Request() *http.Request
	AddAsyncCSS(...string)
	AddCSS(...string)
	AddAsyncJS(...string)
	AddDeferJS(...string)
	AddJS(...string)
	Component(string, ...interface{}) string
	SetData(string, interface{})
	Data(string) interface{}
}
