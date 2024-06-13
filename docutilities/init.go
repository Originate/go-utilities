package docutilities

import (
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func redirectTo(path string, permanent bool) func(*gin.Context) {
	if permanent {
		return func(ctx *gin.Context) { ctx.Redirect(http.StatusMovedPermanently, path) }
	}

	return func(ctx *gin.Context) { ctx.Redirect(http.StatusTemporaryRedirect, path) }
}

func InitDocs(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", redirectTo("/swagger/index.html", true))
}
