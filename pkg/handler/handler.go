package handler

import "github.com/gin-gonic/gin"

// Handler .
type Handler struct {
	Name        string
	Middlewares gin.HandlersChain
	Actions     map[string]*Action
	SubHandlers []Handler
}

// Mount handler
func (h *Handler) Mount(r *gin.Engine) {
	h.mount(&r.RouterGroup)
}

func (h *Handler) mount(g *gin.RouterGroup) {
	g = g.Group(h.Name)

	g.Use(func() []gin.HandlerFunc {
		middlewares := make([]gin.HandlerFunc, len(h.Middlewares)+1)
		middlewares[0] = func(c *gin.Context) {
			c.Next()
		}
		copy(middlewares[1:], h.Middlewares)
		return middlewares
	}()...)

	for name, action := range h.Actions {
		if action.Method&GET != 0 {
			g.GET(name, action.Handler())
		}
		if action.Method&POST != 0 {
			g.POST(name, action.Handler())
		}
	}
	for _, sub := range h.SubHandlers {
		sub.mount(g)
	}
}
