package file

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"zxi_go/core/logger"
	"zxi_go/core/models"
	"zxi_go/utils/customError"
)

type View struct {
	logger         *logrus.Logger
	handler        *Handler
	errCheck       func(*gin.Context, error, int)
}

func NewView() *View {
	return &View{
		logger:         logger.LogInit(),
		handler:        NewHandler(),
		errCheck:       customError.ErrorCheck,
	}
}

func (v *View) ShowFiles(c *gin.Context) {
	filePath := c.Query("path")
	userInter, _ := c.Get("userInfo")
	userMate := userInter.(models.UserInfo)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"dir_list":  v.handler.GetDirList(filePath, userMate.Id),
		"file_list": v.handler.GetFileList(filePath, userMate.Id),
	})
}
