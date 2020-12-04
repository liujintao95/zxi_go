package upload

import (
	"github.com/gin-gonic/gin"
	"zxi_go/core/middleware"
)

func UrlMap(router *gin.Engine) {
	view := NewView()

	authorized := router.Group("/zxi/auth", middleware.LoginRequired)
	authorized.GET("/upload/show", view.ShowUploads)
	authorized.GET("/upload/info", view.ShowUploadInfo)
	authorized.GET("/upload/progress", view.ShowProgress)
	authorized.POST("/upload/buffer", view.UploadBuffer)
	authorized.POST("/upload/pause", view.PauseUpload)
	authorized.POST("/upload/start", view.StartUpload)
	authorized.POST("/upload/cancel", view.CancelUpload)
}
