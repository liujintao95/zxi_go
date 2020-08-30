package zxi

import (
	"github.com/gin-gonic/gin"
	"zxi_network_disk_go/zxi/models"
)

func LoginRequired(g *gin.Context) {
	userMate := models.UserInfo{}
	userMate.Id = 1
	userMate.Name = "北风忆夕"
	userMate.User = "ljt"
	userMate.Pwd = "dgewgdsfw^&(^r1426"

	g.Set("userInfo", userMate)
	g.Next()
}
