package customError

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorCheck(c *gin.Context, err error, errCode int) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": errCode,
			"err_msg":  err.Error(),
		})
		panic(err)
	}
}
