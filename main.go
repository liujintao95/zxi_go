package main

import (
	"github.com/gin-gonic/gin"
	"zxi_go/core/middleware"
	"zxi_go/zxi/download"
	"zxi_go/zxi/file"
	"zxi_go/zxi/upload"
)

func main() {
	router := gin.Default()
	router.Use(middleware.Cors())
	file.UrlMap(router)
	upload.UrlMap(router)
	download.UrlMap(router)
	_ = router.Run(":5000")
}
