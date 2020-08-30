package utils

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "os"
)

func ErrCheck(g *gin.Context, err error, msg string, httpCode int) {
	if err != nil && err != sql.ErrNoRows {
		Logging.Error(msg + ":" + err.Error())
		if httpCode != 0 {
			g.JSON(httpCode, gin.H{
				"errmsg": msg,
				"data":   nil,
			})
			panic(msg + ":" + err.Error())
		}
	}
}
