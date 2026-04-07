package httpsvc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/demoManito/pulse/config"
	"github.com/demoManito/pulse/internal/app/controller"
)

// NewServer new http server
func NewServer(cfg *config.Config) http.Server {
	gin.ForceConsoleColor()

	router := gin.Default()
	router.Use(gin.Logger())
	controller.Handler().Mount(router)

	return http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port),
		Handler: router,
	}
}
