package file

import (
	"github.com/gin-gonic/gin"
	"zxi_go/core/middleware"
)

func UrlMap(router *gin.Engine) {
	view := NewView()

	authorized := router.Group("/zxi/auth", middleware.LoginRequired)
	authorized.GET("/file/show", view.ShowFiles)
}
