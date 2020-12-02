package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrBadReq     = 0001
	ErrDBConn   = 1001
	ErrDBUpdate = 1002
	ErrDBDelete = 1003
	ErrDBCreate = 1004
	ErrDBSelect = 1005
	ErrSaveFile = 2001
	ErrCreatePath = 2002
)

func CustomError (g *gin.Context, err error, errCode int){
	if err != nil{
		Logging.Error(err)
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": errCode,
			"err_msg":  err.Error(),
		})
		panic(err)
	}
}