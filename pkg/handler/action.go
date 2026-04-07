package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// ActionFunc handle the requests
	ActionFunc func(c *Context) (ActionResponse, error)

	// ActionResponse for response data
	ActionResponse any
)

// BasicResponse response struct
type BasicResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ActionResponseWithMessage struct {
	Data    any
	Message string
}

// ActionResponseRaw 原生响应，由 Render 函数完全控制响应输出
type ActionResponseRaw struct {
	Render func(c *gin.Context)
}

// ActionOption action option
type ActionOption struct {
	LoginRequired bool
	// ...
}

// Action request action
type Action struct {
	Method        Method
	Action        ActionFunc
	LoginRequired bool

	handler gin.HandlerFunc
}

// Handler for request handler
func (a *Action) Handler() gin.HandlerFunc {
	return a.handler
}

// NewAction creates new action
func NewAction(method Method, handler ActionFunc) *Action {
	a := Action{Method: method, Action: handler}
	a.handler = func(c *gin.Context) {
		ctx := NewContext(c)
		data, err := a.Action(ctx)
		if err != nil {
			e := DecodeError(err)
			c.AbortWithStatusJSON(e.Status, BasicResponse{Code: e.Code, Message: e.Message, Data: nil})
			return
		}
		if v, ok := data.(*ActionResponseRaw); ok {
			if v.Render != nil {
				v.Render(c)
			} else {
				c.Status(http.StatusNoContent)
			}
			return
		}
		if v, ok := data.(*ActionResponseWithMessage); ok {
			c.JSON(http.StatusOK, BasicResponse{Code: 0, Message: v.Message, Data: v.Data})
			return
		}
		c.JSON(http.StatusOK, BasicResponse{Code: 0, Message: "", Data: data})
	}
	return &a
}

// NewActionResponseWithMessage creates new action response with message
func NewActionResponseWithMessage(data any, message string) *ActionResponseWithMessage {
	return &ActionResponseWithMessage{Data: data, Message: message}
}

// NewActionResponseRaw 创建建新的原生响应，Render 函数完全控制响应输出（render 为 nil 时，默认返回 204 No Content）
func NewActionResponseRaw(render func(c *gin.Context)) *ActionResponseRaw {
	return &ActionResponseRaw{Render: render}
}
