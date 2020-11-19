package file

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	var file FileInfo
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	if err := g.Bind(&file); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}

	handlers.CreateOrIgnoreFile(file.Hash, file.Size)
	lastId := uploadHandlers.CreateOrIgnoreUpload(file.Hash, file.Path, userMate.Id)
	uploadHandlers.CreateOrIgnoreUploadBlock(lastId, file.Size)
	handlers.CreateUserFile(file.Hash, file.Name, userMate.Id, "")

	g.JSON(http.StatusOK, gin.H{"success": true})
}

func (v *View) SaveFilesInfo(g *gin.Context) {
	var postParser DirInfo
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	if err := g.Bind(&postParser); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"err_code": core.ERR_BAD_REQ,
			"err_msg":  err.Error(),
		})
		return
	}

	for _, file := range postParser.Files {
		fileRelativePath := handlers.AbsolutePathToRelativePath(file.Path, postParser.Root)
		dirPath, _ := handlers.PathSplit(fileRelativePath)

		handlers.CreatePath(dirPath, userMate.Id)
		handlers.CreateOrIgnoreFile(file.Hash, file.Size)
		lastId := uploadHandlers.CreateOrIgnoreUpload(file.Hash, file.Path, userMate.Id)
		uploadHandlers.CreateOrIgnoreUploadBlock(lastId, file.Size)
		handlers.CreateUserFile(file.Hash, file.Name, userMate.Id, dirPath)
	}

	g.JSON(http.StatusOK, gin.H{"success": true})
}
