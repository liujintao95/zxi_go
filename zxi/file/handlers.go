package file

import (
	"path/filepath"
	"strings"
	. "zxi_go/core"
	"zxi_go/zxi/models"
)

type Handlers struct {
}

var handlers = Handlers{}

func (h *Handlers) GetDirList(path string, userId int) []models.Directory {
	var dirList []models.Directory
	if path == "" {
		path = "/"
	}
	LocalDB.Where(&models.Directory{
		Recycled:   "N",
		Path:       path,
		UserInfoId: userId,
	}).Find(&dirList)
	return dirList
}

func (h *Handlers) GetFileList(path string, userId int) []models.UserFile {
	var dirPath, dirName string
	var dirMate models.Directory
	var fileList []models.UserFile
	if path == "" {
		dirName, dirPath = "/", "-"
	} else {
		dirPath, dirName = h.PathSplit(path)
	}
	LocalDB.Where(&models.Directory{
		Recycled:   "N",
		Path:       dirPath,
		Name:       dirName,
		UserInfoId: userId,
	}).First(&dirMate)
	LocalDB.Where(&models.UserFile{
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

func (h *Handlers) GetFileInfo(fileId int) models.File {
	var fileMate models.File
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       fileId,
	}).First(&fileMate)
	return fileMate
}

func (h *Handlers) CreateUserFile(hash string, fileName string, userId int, path string) int {
	var dirPath, dirName string
	var fileMate models.File
	var dirMate models.Directory
	if path == "" {
		dirName, dirPath = "/", "-"
	} else {
		dirPath, dirName = h.PathSplit(path)
	}
	LocalDB.Where(&models.File{
		Recycled: "N",
		Hash:     hash,
	}).First(&fileMate)
	LocalDB.Where(&models.Directory{
		Recycled:   "N",
		UserInfoId: userId,
		Name:       dirName,
		Path:       dirPath,
	}).First(&dirMate)
	userFileMate := &models.UserFile{
		Name:        fileName,
		UserInfoId:  userId,
		FileId:      fileMate.Id,
		DirectoryId: dirMate.Id,
	}
	LocalDB.Create(&userFileMate)
	return userFileMate.Id
}

func (h *Handlers) CreateOrIgnoreFile(hash string, size int) int {
	var fileMate models.File
	LocalDB.Where(&models.File{
		Recycled: "N",
		Hash:     hash,
	}).First(&fileMate)
	if fileMate.Id == 0 {
		fileMate.Hash = hash
		fileMate.Path = filepath.Join(SAVE_PATH, hash)
		fileMate.Size = size
		LocalDB.Create(&fileMate)
	}
	return fileMate.Id
}

func (h *Handlers) CreateOrIgnoreDir(dirName string, dirPath string, userId int) int {
	var dirMate models.Directory
	LocalDB.Where(&models.Directory{
		Recycled:   "N",
		Name:       dirName,
		Path:       dirPath,
		UserInfoId: userId,
	}).First(&dirMate)
	if dirMate.Id == 0 {
		dirMate.Name = dirName
		dirMate.Path = dirPath
		dirMate.UserInfoId = userId
		dirMate.IsKey = 0
		LocalDB.Create(&dirMate)
	}
	return dirMate.Id
}

func (h *Handlers) CreatePath(path string, userId int) {
	var dirPath, dirName string
	dirPath = path
	for {
		dirPath, dirName = h.PathSplit(dirPath)
		h.CreateOrIgnoreDir(dirName, dirPath, userId)
		if dirPath == "/" {
			h.CreateOrIgnoreDir("/", "-", userId)
			break
		}
	}
}

func (h *Handlers) AbsolutePathToRelativePath(absolutePath string, rootPath string) string {
	prefix, _ := h.PathSplit(rootPath)
	absolutePath = strings.Replace(absolutePath, `\`, "/", -1)
	return strings.Replace(absolutePath, prefix, ``, 1)
}

func (h *Handlers) PathSplit(path string) (string, string) {
	path = strings.Replace(path, `\`, "/", -1)
	i := strings.LastIndex(path, "/")
	if path[:i] != "" {
		return path[:i], path[i+1:]
	} else {
		return "/", path[i+1:]
	}
}
