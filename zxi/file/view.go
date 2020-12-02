package file

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zxi_go/core"
	"zxi_go/zxi/models"
	"zxi_go/zxi/upload"
)

type View struct {
}

var view = View{}
var uploadHandlers = upload.Handlers{}

func (v *View) ShowFiles(g *gin.Context) {
	filePath := g.Query("path")
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)

	g.JSON(http.StatusOK, gin.H{
		"success":   true,
		"dir_list":  handlers.GetDirList(filePath, userMate.Id),
		"file_list": handlers.GetFileList(filePath, userMate.Id),
	})
}

func (v *View) SaveFileInfo(g *gin.Context) {
	sizeStr := g.PostForm("size")
	hash := g.PostForm("hash")
	name := g.PostForm("name")
	path := g.PostForm("path")
	root := g.PostForm("root")
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	size, err := strconv.Atoi(sizeStr)
	core.CustomError(g, err, core.ErrBadReq)

	var dirPath string
	if root != ""{
		fileRelativePath := handlers.AbsolutePathToRelativePath(path, root)
		dirPath, _ = handlers.PathSplit(fileRelativePath)
		handlers.CreatePath(dirPath, userMate.Id)
	}
	handlers.CreateOrIgnoreFile(hash, size)
	lastId := uploadHandlers.CreateOrIgnoreUpload(hash, path, userMate.Id)
	uploadHandlers.CreateOrIgnoreUploadBlock(lastId, size)
	handlers.CreateUserFile(hash, name, userMate.Id, dirPath)

	g.JSON(http.StatusOK, gin.H{"success": true})
}
