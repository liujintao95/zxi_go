package file

import (
	"github.com/gin-gonic/gin"
	"zxi_go/zxi/middleware"
)

func UrlMap(router *gin.Engine) {
	authorized := router.Group("/zxi/auth", middleware.LoginRequired)
	authorized.GET("/file/showfiles", view.ShowFiles)
	authorized.POST("/file/savefilesinfo", view.SaveFilesInfo)
	authorized.POST("/file/savefileinfo", view.SaveFileInfo)
}
