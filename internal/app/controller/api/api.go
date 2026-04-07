package api

import (
	"github.com/demoManito/pulse/internal/app/controller/api/health"
	"github.com/demoManito/pulse/pkg/handler"
)

// Handler return all controllers
func Handler() handler.Handler {
	return handler.Handler{
		Name: "api",
		SubHandlers: []handler.Handler{
			health.Handler(),
		},
	}
}
