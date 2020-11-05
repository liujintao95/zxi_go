package utils

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "os"
    "zxi_go/core"
)

func ErrCheck(g *gin.Context, err error, msg string, httpCode int) {
	if err != nil && err != sql.ErrNoRows {
		core.Logging.Error(msg + ":" + err.Error())
		if httpCode != 0 {
			g.JSON(httpCode, gin.H{
				"errmsg": msg,
			})
			panic(msg + ":" + err.Error())
		}
	}
}
