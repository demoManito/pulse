package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Context request context
type Context struct {
	C *gin.Context
}

func NewContext(c *gin.Context) *Context {
	return &Context{
		C: c,
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.C.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.C.Done()
}

func (c *Context) Err() error {
	return c.C.Err()
}

func (c *Context) Value(key any) any {
	return c.C.Value(key)
}
