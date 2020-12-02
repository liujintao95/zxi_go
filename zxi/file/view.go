package file

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"zxi_go/core"
	"zxi_go/zxi/models"
	"zxi_go/zxi/upload"
)

type View struct {
	logging *logrus.Logger
	handler *Handler
	uploadHandlers *upload.Handler
}

func NewView() *View {
	return &View{
		logging: core.LogInit(),
		handler: NewHandler(),
		uploadHandlers: upload.NewHandler(),
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

func (v *View) SaveFileInfo(c *gin.Context) {
	sizeStr := c.PostForm("size")
	hash := c.PostForm("hash")
	name := c.PostForm("name")
	path := c.PostForm("path")
	root := c.PostForm("root")
	userInter, _ := c.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	size, err := strconv.Atoi(sizeStr)
	core.CustomError(c, err, core.ErrBadReq)

	var dirPath string
	if root != ""{
		fileRelativePath := v.handler.AbsolutePathToRelativePath(path, root)
		dirPath, _ = v.handler.PathSplit(fileRelativePath)
		v.handler.CreatePath(dirPath, userMate.Id)
	}
	v.handler.CreateOrIgnoreFile(hash, size)
	lastId := v.uploadHandlers.CreateOrIgnoreUpload(hash, path, userMate.Id)
	v.uploadHandlers.CreateOrIgnoreUploadBlock(lastId, size)
	v.handler.CreateUserFile(hash, name, userMate.Id, dirPath)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
