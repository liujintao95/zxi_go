package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zxi_go/utils"
	"zxi_go/zxi/file/handlers"
	"zxi_go/zxi/file/parser"
	"zxi_go/zxi/models"
)

var errCheck = utils.ErrCheck

func UploadFile(g *gin.Context) {
	fileJson := g.PostForm("file_json")
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)

	postFile := new(parser.PostFileInfo)
	err := json.Unmarshal([]byte(fileJson), postFile)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	_, err = handlers.CreateOrIgnoreFile(postFile.Hash, postFile.Size)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	lastId, err := handlers.CreateOrIgnoreUpload(
		postFile.Hash, postFile.Path, userMate.Id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	if lastId != 0 {
		handlers.CreateUploadBlock(lastId, postFile.Size)
	}

	_, err = handlers.CreateUserFile(
		postFile.Hash, postFile.Name, userMate.Id, 0)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	g.JSON(http.StatusOK, gin.H{"err": nil})
}

func UploadFiles(g *gin.Context) {
	filesJsonList := g.PostFormArray("files_json")
	rootPath := g.PostForm("root")
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)

	for _, fileJson := range filesJsonList {
		postFile := new(parser.PostFileInfo)
		err := json.Unmarshal([]byte(fileJson), postFile)
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"err": err})
			return
		}
		relativePath := handlers.AbsolutePathToRelativePath(postFile.Path, rootPath)
		dirId, err := handlers.CreatePath(relativePath, userMate.Id)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		_, err = handlers.CreateOrIgnoreFile(postFile.Hash, postFile.Size)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		_, err = handlers.CreateOrIgnoreUpload(
			postFile.Hash, postFile.Path, userMate.Id)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		_, err = handlers.CreateUserFile(
			postFile.Hash, postFile.Name, userMate.Id, dirId)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}
	}

	g.JSON(http.StatusOK, gin.H{
		"err": "ok",
	})
}

func ShowFiles(g *gin.Context) {
	path := g.Query("path")
	sDirId := g.Query("dir_id")
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)

	dirId, err := strconv.ParseInt(sDirId, 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	dirList, err := handlers.GetDirList(path, userMate.Id)
	fileList, err := handlers.GetFileList(dirId, userMate.Id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"err":       "ok",
		"dir_list":  dirList,
		"file_list": fileList,
	})
}
