package file

import (
	"github.com/jinzhu/gorm"
	"strings"
	"zxi_go/core/database"
	"zxi_go/core/models"
)

type Handler struct {
	localDB *gorm.DB
}

func NewHandler() *Handler {
	return &Handler{
		localDB: database.MySqlInit(),
	}
}

func (h *Handler) GetDirList(path string, userId int) []models.Directory {
	var dirList []models.Directory
	if path == "" {
		path = "/"
	}
	h.localDB.Where(&models.Directory{
		Recycled:   "N",
		Path:       path,
		UserInfoId: userId,
	}).Find(&dirList)
	return dirList
}

func (h *Handler) GetFileList(path string, userId int) []models.UserFile {
	var dirPath, dirName string
	var dirMate models.Directory
	var fileList []models.UserFile
	if path == "" {
		dirName, dirPath = "/", "-"
	} else {
		dirPath, dirName = h.PathSplit(path)
	}
	h.localDB.Where(&models.Directory{
		Recycled:   "N",
		Path:       dirPath,
		Name:       dirName,
		UserInfoId: userId,
	}).First(&dirMate)
	h.localDB.Where(&models.UserFile{
		Recycled: "N",
		File: models.File{
			IsComplete: 1,
			Recycled:   "N",
		},
		DirectoryId: dirMate.Id,
		UserInfoId:  userId,
	}).Find(&fileList)
	return fileList
}

func (h *Handler) GetFileInfo(fileId int) models.File {
	var fileMate models.File
	h.localDB.Where(&models.Upload{
		Recycled: "N",
		Id:       fileId,
	}).First(&fileMate)
	return fileMate
}

func (h *Handler) AbsolutePathToRelativePath(absolutePath string, rootPath string) string {
	prefix, _ := h.PathSplit(rootPath)
	absolutePath = strings.Replace(absolutePath, `\`, "/", -1)
	return strings.Replace(absolutePath, prefix, ``, 1)
}

func (h *Handler) PathSplit(path string) (string, string) {
	path = strings.Replace(path, `\`, "/", -1)
	i := strings.LastIndex(path, "/")
	if path[:i] != "" {
		return path[:i], path[i+1:]
	} else {
		return "/", path[i+1:]
	}
}
