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
	authorized.POST("/file/uploadfiles", api.UploadFiles)
	authorized.POST("/file/uploadfile", api.UploadFile)
	authorized.GET("/file/show", api.ShowFiles)
}
