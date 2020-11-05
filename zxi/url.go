package zxi

import (
	"github.com/gin-gonic/gin"
	"zxi_go/zxi/api"
)

func UrlMap(router *gin.Engine) {
	//router.POST("/register", api.Register)
	//router.GET("/logout", api.Logout)
	//
	authorized := router.Group("/zxi/auth", LoginRequired)
	authorized.POST("/file/uploadfiles", api.UploadFiles)
	authorized.POST("/file/uploadfile", api.UploadFile)
	authorized.GET("/file/show", api.ShowFiles)
	//authorized.GET("/user/show", api.ShowUser)
	//authorized.POST("/user/change/password", api.PasswordChange)
	//authorized.POST("/user/change/username", api.UsernameChange)
	//authorized.POST("/user/change/phone", api.PhoneChange)
	//authorized.POST("/user/change/email", api.EmailChange)
	//
	//authorized.POST("/file/init", api.InitFile)
	//authorized.POST("/file/upload", api.Upload)
	//authorized.POST("/file/rapid_upload", api.RapidUpload)
	//authorized.GET("/file/download", api.Download)
	//authorized.GET("/file/public_download", api.PublicDownload)
	//authorized.POST("/file/update", api.UpdateFileName)
	//authorized.POST("/file/delete", api.Delete)
	//authorized.GET("/file/upload_show", api.UploadShow)
	//
	//authorized.POST("/block/init", api.InitBlockUpload)
	//authorized.GET("/block/resume", api.ResumeFromBreakPoint)
	//authorized.POST("/block/upload", api.BlockUpload)
	//authorized.GET("/block/progress", api.UploadProgress)
	//authorized.POST("/block/merge", api.BlockMerge)
	//authorized.POST("/block/remove", api.RemoveBlock)
	//
	//authorized.GET("/dir/show", api.ShowDir)
	//authorized.POST("/dir/save", api.SaveDir)
	//authorized.POST("/dir/change", api.ChangeDir)
	//authorized.POST("/dir/remove", api.RemoveDir)
}
