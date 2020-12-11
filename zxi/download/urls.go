package download

import (
	"github.com/gin-gonic/gin"
	"zxi_go/core/middleware"
)

func UrlMap(router *gin.Engine) {
	view := NewView()

	authorized := router.Group("/zxi/auth", middleware.LoginRequired)
	authorized.GET("/download/show", view.ShowDownload)
    authorized.GET("/download/info", view.ShowDownloadInfo)
	authorized.GET("/download/progress", view.ShowProgress)
	authorized.POST("/download/create", view.CreateDownload)
	authorized.POST("/download/buffer", view.DownloadBuffer)
	authorized.POST("/download/pause", view.PauseDownload)
	authorized.POST("/download/start", view.StartDownload)
	authorized.POST("/download/cancel", view.CancelDownload)
}
