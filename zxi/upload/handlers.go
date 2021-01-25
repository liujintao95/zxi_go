package upload

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"path/filepath"
	"strconv"
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

func (h *Handler) CreateOrIgnoreUploadBlock(uploadId int, size int) {
	var uploadMate models.Upload
	h.localDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	if uploadMate.IsComplete == 0 {
		blockNum := (size / BlockSize) + 1
		for i := 0; i < blockNum; i++ {
			var blockSize int
			var uploadBlockMate models.UploadBlock
			if i+1 == blockNum {
				blockSize = size % BlockSize
			} else {
				blockSize = BlockSize
			}
			h.localDB.Where(&models.UploadBlock{
				Recycled: "N",
				UploadId: uploadId,
				Offset:   i,
				Size:     blockSize,
			}).First(&uploadBlockMate)
			if uploadBlockMate.Id == 0 {
				uploadBlockMate.UploadId = uploadId
				uploadBlockMate.Offset = i
				uploadBlockMate.Size = blockSize
				h.localDB.Create(&uploadBlockMate)
			}
		}
	}
}

func (h *Handler) CreateOrIgnoreUpload(hash string, path string, userId int) int {
	var fileMate models.File
	var uploadMate models.Upload
	h.localDB.Where(&models.File{
		Recycled: "N",
		Hash:     hash,
	}).First(&fileMate)
	h.localDB.Where(&models.Upload{
		Recycled:   "N",
		FileId:     fileMate.Id,
		UserInfoId: userId,
	}).First(&uploadMate)
	if uploadMate.Id == 0 {
		uploadMate.FileId = fileMate.Id
		uploadMate.UserInfoId = userId
		uploadMate.BlockSize = BlockSize
		uploadMate.LocalPath = path
		if fileMate.IsComplete == 1 {
			uploadMate.Uploading = 0
			uploadMate.IsComplete = 1
		} else {
			uploadMate.Uploading = 1
			uploadMate.IsComplete = 0
		}
		h.localDB.Create(&uploadMate)
	}
	return uploadMate.Id
}

func (h *Handler) CreateOrIgnoreDir(dirName string, dirPath string, userId int) int {
	var dirMate models.Directory
	h.localDB.Where(&models.Directory{
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
		h.localDB.Create(&dirMate)
	}
	return dirMate.Id
}

func (h *Handler) CreateZXiPath(path string, userId int) {
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

func (h *Handler) CreateOrIgnoreFile(hash string, size int) int {
	var fileMate models.File
	h.localDB.Where(&models.File{
		Recycled: "N",
		Hash:     hash,
	}).First(&fileMate)
	if fileMate.Id == 0 {
		fileMate.Hash = hash
		fileMate.Path = filepath.Join(FilePath, hash)
		fileMate.Size = size
		h.localDB.Create(&fileMate)
	}
	return fileMate.Id
}

func (h *Handler) CreateUserFile(hash string, fileName string, userId int, path string) int {
	var dirPath, dirName string
	var fileMate models.File
	var dirMate models.Directory
	if path == "" {
		dirName, dirPath = "/", "-"
	} else {
		dirPath, dirName = h.PathSplit(path)
	}
	h.localDB.Where(&models.File{
		Recycled: "N",
		Hash:     hash,
	}).First(&fileMate)
	h.localDB.Where(&models.Directory{
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
	h.localDB.Create(&userFileMate)
	return userFileMate.Id
}

func (h *Handler) CreateOrIgnoreLocalPath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0777)
	}
	return err
}

func (h *Handler) GetProgress(uploadId int) float64 {
	var progress float64
	var totalNum, completeNum int
	var blockList []models.UploadBlock

	if h.GetUploadInfo(uploadId).IsComplete == 1 {
		return 100
	}

	h.localDB.Where(&models.UploadBlock{
		Recycled: "N",
		UploadId: uploadId,
	}).Find(&blockList)
	for _, blockMate := range blockList {
		if blockMate.IsComplete == 1 {
			completeNum++
		}
		totalNum++
	}
	if totalNum != 0 {
		progress, _ = strconv.ParseFloat(
			fmt.Sprintf("%.2f", float64(completeNum)/float64(totalNum)*100),
			64,
		)
	} else {
		progress = 0
	}
	if progress == 100 {
		fileMate := h.GetFileInfoByUploadId(uploadId)
		uploadMate := h.GetUploadInfo(uploadId)
		fileMate.IsComplete = 1
		uploadMate.IsComplete = 1
		uploadMate.Uploading = 0
		h.localDB.Save(fileMate)
		h.localDB.Save(uploadMate)
	}
	return progress
}

func (h *Handler) GetUploadTable(userId int, page int, size int) ([]ShowUploadTable, int) {
	var result []ShowUploadTable
	var count int
	var uploadList []models.Upload
	start := (page - 1) * size
	h.localDB.Where(&models.Upload{
		Recycled:   "N",
		UserInfoId: userId,
	}).Order("id desc").Limit(size).Offset(start).Find(&uploadList)
	h.localDB.Model(&models.Upload{}).Where(&models.Upload{
		Recycled:   "N",
		UserInfoId: userId,
	}).Count(&count)
	for _, uploadInfo := range uploadList {
		var progress float64
		var userFileMate models.UserFile
		h.localDB.Where(&models.UserFile{
			UserInfoId: userId,
			FileId:     uploadInfo.FileId,
		}).First(&userFileMate)
		h.localDB.Model(&userFileMate).Related(&userFileMate.File)
		if uploadInfo.IsComplete == 1 {
			progress = 100
		} else {
			var totalNum, completeNum int
			var blockList []models.UploadBlock
			h.localDB.Where(&models.UploadBlock{
				Recycled: "N",
				UploadId: uploadInfo.Id,
			}).Find(&blockList)
			for _, blockMate := range blockList {
				if blockMate.IsComplete == 1 {
					completeNum++
				}
				totalNum++
			}
			if totalNum != 0 {
				progress, _ = strconv.ParseFloat(
					fmt.Sprintf("%.2f", float64(completeNum)/float64(totalNum)*100),
					64,
				)
			} else {
				progress = 0
			}
		}
		result = append(result, ShowUploadTable{
			Id:         uploadInfo.Id,
			Uploading:  uploadInfo.Uploading,
			IsComplete: uploadInfo.IsComplete,
			Name:       userFileMate.Name,
			LocalPath:  uploadInfo.LocalPath,
			Size:       userFileMate.File.Size,
			Progress:   progress,
		})
	}
	return result, count
}

func (h *Handler) GetUploadInfo(uploadId int) models.Upload {
	var uploadMate models.Upload
	h.localDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	h.localDB.Model(&uploadMate).Related(&uploadMate.Block)
	return uploadMate
}

func (h *Handler) GetBlockListByUploadId(uploadId int) []models.UploadBlock {
	var uploadBlockList []models.UploadBlock
	h.localDB.Where(&models.UploadBlock{
		Recycled: "N",
		UploadId: uploadId,
	}).Find(&uploadBlockList)
	return uploadBlockList
}

func (h *Handler) GetFileInfoByUploadId(uploadId int) models.File {
	var uploadMate models.Upload
	var fileMate models.File
	h.localDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	h.localDB.Model(&uploadMate).Related(&fileMate)
	return fileMate
}

func (h *Handler) GetUploadBlockInfo(blockId int) models.UploadBlock {
	var uploadBlockMate models.UploadBlock
	h.localDB.Where(&models.UploadBlock{
		Recycled: "N",
		Id:       blockId,
	}).First(&uploadBlockMate)
	return uploadBlockMate
}

func (h *Handler) UpdateFileComplete(uploadId int, state int) {
	var uploadMate models.Upload
	var fileMate models.File
	h.localDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	h.localDB.Model(&uploadMate).Related(&fileMate)
	uploadMate.IsComplete = state
	fileMate.IsComplete = state
	h.localDB.Save(&uploadMate)
	h.localDB.Save(&fileMate)
}

func (h *Handler) UpdateBlockComplete(blockId int, state int) {
	var uploadBlockMate models.UploadBlock
	h.localDB.Where(&models.UploadBlock{
		Recycled: "N",
		Id:       blockId,
	}).First(&uploadBlockMate)
	uploadBlockMate.IsComplete = state
	h.localDB.Save(&uploadBlockMate)
}

func (h *Handler) UpdateUploading(uploadId int, state int) {
	var uploadMate models.Upload
	h.localDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	uploadMate.Uploading = state
	h.localDB.Save(&uploadMate)
}

func (h *Handler) DeleteUpload(uploadId int) {
	var uploadMate models.Upload
	h.localDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	uploadMate.Recycled = "Y"
	h.localDB.Save(uploadMate)

	h.localDB.Where(&models.UploadBlock{
		UploadId: uploadId,
	}).Delete(models.UploadBlock{})
}

func (h *Handler) SaveFile(data []byte, path string) error {
	fo, err := os.Create(path)
	if err == nil {
		_, err = fo.Write(data)
	}
	return err
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
