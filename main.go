package main

import (
	"github.com/gin-gonic/gin"
	"zxi_go/core"
	"zxi_go/zxi/file"
	"zxi_go/zxi/upload"
)

func main() {
	router := gin.Default()
	router.Use(core.Cors())
	file.UrlMap(router)
	upload.UrlMap(router)
	_ = router.Run(":5000")
}
