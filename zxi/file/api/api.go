package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strconv"
	"zxi_go/zxi/file/handlers"
	"zxi_go/zxi/file/parser"
	"zxi_go/zxi/models"
)

func SaveFileInfo(g *gin.Context) {
	var postParser parser.FileInfo
	if err := g.Bind(&postParser); err != nil{
		g.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)

	_, err := handlers.CreateOrIgnoreFile(postParser.Hash, postParser.Size)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	lastId, err := handlers.CreateOrIgnoreUpload(
		postParser.Hash, postParser.Path, userMate.Id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	} else if lastId != 0 {
		err = handlers.CreateUploadBlock(lastId, postParser.Size)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}
	}

	_, err = handlers.CreateUserFile(
		postParser.Hash, postParser.Name, userMate.Id, 0)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	g.JSON(http.StatusOK, gin.H{"err": nil})
}

func SaveFilesInfo(g *gin.Context) {
	var postParser parser.DirInfo
	if err := g.Bind(&postParser); err != nil{
		g.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)

	for _, file := range postParser.Files {
		relativePath := handlers.AbsolutePathToRelativePath(file.Path, postParser.Root)
		dirId, err := handlers.CreatePath(relativePath, userMate.Id)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		_, err = handlers.CreateOrIgnoreFile(file.Hash, file.Size)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		_, err = handlers.CreateOrIgnoreUpload(
			file.Hash, file.Path, userMate.Id)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		_, err = handlers.CreateUserFile(
			file.Hash, file.Name, userMate.Id, dirId)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}
	}

	g.JSON(http.StatusOK, gin.H{"err": nil})
}

func ShowFiles(g *gin.Context) {
	file_path := g.Query("path")
	sDirId := g.Query("dir_id")
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)

	dirId, err := strconv.ParseInt(sDirId, 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	dirList, err := handlers.GetDirList(file_path, userMate.Id)
	fileList, err := handlers.GetFileList(dirId, userMate.Id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"err":       nil,
		"dir_list":  dirList,
		"file_list": fileList,
	})
}

func ShowUploads(g *gin.Context) {
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	uploadList, err := handlers.GetUploadList(userMate.Id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"err":         nil,
		"upload_list": uploadList,
	})
}

func ShowUploadInfo(g *gin.Context) {
	uploadIdStr := g.Query("id")
	uploadId, err := strconv.ParseInt(uploadIdStr, 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	uploadInfo, err := handlers.GetUploadInfo(uploadId)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"err":         nil,
		"upload_info": uploadInfo,
	})
}

func UploadFile(g *gin.Context) {
	var postParser parser.UploadFile
	err := g.Bind(&postParser)
	f, err := g.FormFile("file")
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	err = g.SaveUploadedFile(f, path.Join("file", f.Filename))
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	err = handlers.SetIsComplete(postParser.Id, 1)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	g.JSON(http.StatusOK, gin.H{"err": nil})
}

func UploadBlock(g *gin.Context) {}

func ShowProgress(g *gin.Context) {}
