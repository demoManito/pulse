package controller

import (
	"github.com/demoManito/pulse/internal/app/controller/api"
	"github.com/demoManito/pulse/pkg/handler"
)

// Handler return all controllers
func Handler() *handler.Handler {
	return &handler.Handler{
		SubHandlers: []handler.Handler{
			api.Handler(),
		},
	}
}
