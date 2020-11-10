package url

import (
	"github.com/gin-gonic/gin"
	"zxi_go/zxi/file/api"
	"zxi_go/zxi/middleware"
)

func UrlMap(router *gin.Engine) {
	//router.POST("/register", api.Register)
	//router.GET("/logout", api.Logout)
	//
	authorized := router.Group("/zxi/auth", middleware.LoginRequired)
	authorized.POST("/file/savefilesinfo", api.SaveFilesInfo)
	authorized.POST("/file/savefileinfo", api.SaveFileInfo)
	authorized.GET("/file/showfiles", api.ShowFiles)
	authorized.GET("/file/showuploads", api.ShowUploads)
	authorized.GET("/file/showupload", api.ShowUploadInfo)
	authorized.POST("/file/uploadfile", api.UploadFile)
}
