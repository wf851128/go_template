package routes

import (
	"github.com/gin-gonic/gin"
	"go_template/logger"
	"net/http"
)

func Setup() (g *gin.Engine) {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})
	return r
}
