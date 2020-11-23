package upload

import (
	"github.com/gin-gonic/gin"
	"zxi_go/zxi/middleware"
)

func UrlMap(router *gin.Engine) {
	authorized := router.Group("/zxi/auth", middleware.LoginRequired)
	authorized.GET("/upload/show", view.ShowUploads)
	authorized.GET("/upload/info", view.ShowUploadInfo)
	authorized.POST("/upload/file", view.UploadFile)
}
