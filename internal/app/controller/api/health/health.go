package health

import (
	"github.com/demoManito/pulse/pkg/handler"
)

// Handler return middlewares
func Handler() handler.Handler {
	return handler.Handler{
		Name: "health",
		Actions: map[string]*handler.Action{
			"": handler.NewAction(handler.GET, ActionHealth),
		},
	}
}

// ActionHealth health check
// @Summary      健康检查
// @Description  健康检查
// @Tags         健康检查
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]any "{"code":0,"message":"success","data":null}"
// @Router       /api/akashic/health [get]
func ActionHealth(*handler.Context) (handler.ActionResponse, error) {
	return handler.NewActionResponseWithMessage(nil, "success"), nil
}
