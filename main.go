package main

import (
	"github.com/gin-gonic/gin"
	"zxi_go/core"
	network "zxi_go/zxi"
)


func main() {
	core.LogInit()
	core.MySqlInit()

	// Disable Console Color
	// gin.DisableConsoleColor()

	// 使用默认中间件创建一个gin路由器
	// logger and recovery (crash-free) 中间件
	router := gin.Default()
	router.Use(core.Cors())
	network.UrlMap(router)
	_ = router.Run(":5000")
}
