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
	authorized.GET("/download/buffer", view.DownloadBuffer)
	authorized.GET("/download/pause", view.PauseDownload)
	authorized.GET("/download/start", view.StartDownload)
	authorized.GET("/download/cancel", view.CancelDownload)
}
