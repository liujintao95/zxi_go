package upload

import (
	"github.com/gin-gonic/gin"
	"zxi_go/zxi/middleware"
)

func UrlMap(router *gin.Engine) {
	authorized := router.Group("/zxi/auth", middleware.LoginRequired)
	authorized.GET("/upload/show", view.ShowUploads)
	authorized.GET("/upload/info", view.ShowUploadInfo)
	authorized.GET("/upload/progress", view.ShowProgress)
	authorized.POST("/upload/file", view.UploadFile)
	authorized.POST("/upload/block", view.UploadBlock)
	authorized.POST("/upload/pause", view.PauseUpload)
	authorized.POST("/upload/start", view.StartUpload)
	authorized.POST("/upload/cancel", view.CancelUpload)
}
