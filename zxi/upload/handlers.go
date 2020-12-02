package upload

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	. "zxi_go/core"
	"zxi_go/zxi/models"
)

type Handlers struct {
}

var handlers = Handlers{}

func (h *Handlers) CreateOrIgnoreUploadBlock(uploadId int, size int) {
	var uploadMate models.Upload
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	if uploadMate.IsComplete == 0 {
		blockNum := (size / BLOCK_SIZE) + 1
		for i := 0; i < blockNum; i++ {
			var blockSize int
			var uploadBlockMate models.UploadBlock
			if i+1 == blockNum {
				blockSize = size % BLOCK_SIZE
			} else {
				blockSize = BLOCK_SIZE
			}
			LocalDB.Where(&models.UploadBlock{
				Recycled: "N",
				UploadId: uploadId,
				Offset:   i,
				Size:     blockSize,
			}).First(&uploadBlockMate)
			if uploadBlockMate.Id == 0 {
				uploadBlockMate.UploadId = uploadId
				uploadBlockMate.Offset = i
				uploadBlockMate.Size = blockSize
				LocalDB.Create(&uploadBlockMate)
			}
		}
	}
}

func (h *Handlers) CreateOrIgnoreUpload(hash string, path string, userId int) int {
	var fileMate models.File
	var uploadMate models.Upload
	LocalDB.Where(&models.File{
		Recycled: "N",
		Hash:     hash,
	}).First(&fileMate)
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		FileId:   fileMate.Id,
		UserInfoId: userId,
	}).First(&uploadMate)
	if uploadMate.Id == 0 {
		uploadMate.FileId = fileMate.Id
		uploadMate.UserInfoId = userId
		uploadMate.BlockSize = BLOCK_SIZE
		uploadMate.LocalPath = path
		if fileMate.IsComplete == 1 {
			uploadMate.Uploading = 0
			uploadMate.IsComplete = 1
		} else {
			uploadMate.Uploading = 1
			uploadMate.IsComplete = 0
		}
		LocalDB.Create(&uploadMate)
	} else if uploadMate.Recycled == "Y"{
		uploadMate.Recycled = "N"
		uploadMate.LocalPath = path
		uploadMate.BlockSize = BLOCK_SIZE
		if fileMate.IsComplete == 1 {
			uploadMate.Uploading = 0
			uploadMate.IsComplete = 1
		} else {
			uploadMate.Uploading = 1
			uploadMate.IsComplete = 0
		}
		LocalDB.Save(uploadMate)
	}
	return uploadMate.Id
}

func (h *Handlers) CreateOrIgnorePath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0777)
	}
	return err
}

func (h *Handlers) GetUploadList(userId int, page int, size int) ([]ShowUploadTable, int) {
	var count int
	var uploadList []models.Upload
	var result []ShowUploadTable
	start := (page - 1) * size
	LocalDB.Where(&models.Upload{
		Recycled:   "N",
		UserInfoId: userId,
	}).Order("id desc").Limit(size).Offset(start).Find(&uploadList)
	LocalDB.Model(&models.Upload{}).Where(&models.Upload{
		Recycled:   "N",
		UserInfoId: userId,
	}).Count(&count)
	for _, uploadInfo := range uploadList {
		var progress float64
		var fileMate models.File
		LocalDB.Where(&models.File{
			Id: uploadInfo.FileId,
		}).First(&fileMate)
		if uploadInfo.IsComplete == 1 {
			progress = 100
		} else {
			var totalNum, completeNum int
			var blockList []models.UploadBlock
			LocalDB.Where(&models.UploadBlock{
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
			LocalPath:  uploadInfo.LocalPath,
			Name:       GetFileByPath(uploadInfo.LocalPath),
			SizeFmt:       StrSize(fileMate.Size),
			Size:       fileMate.Size,
			Progress:   progress,
		})
	}
	return result, count
}

func (h *Handlers) GetProgress(uploadId int) float64 {
	var progress float64
	var totalNum, completeNum int
	var blockList []models.UploadBlock

	LocalDB.Where(&models.UploadBlock{
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
	return progress
}

func (h *Handlers) GetUploadInfo(uploadId int) ShowUploadInfo {
	var uploadMate models.Upload
	var uploadBlockList []models.UploadBlock
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	LocalDB.Where(&models.UploadBlock{
		Recycled: "N",
		UploadId: uploadId,
	}).Find(&uploadBlockList)
	return ShowUploadInfo{
		Id:         uploadMate.Id,
		LocalPath:  uploadMate.LocalPath,
		BlockSize:  uploadMate.BlockSize,
		Uploading:  uploadMate.Uploading,
		IsComplete: uploadMate.IsComplete,
		FileId:     uploadMate.FileId,
		UserInfoId: uploadMate.UserInfoId,
		BlockList:  uploadBlockList,
	}
}

func (h *Handlers) GetFileInfoByUploadId(uploadId int) models.File {
	var uploadMate models.Upload
	var fileMate models.File
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	LocalDB.Where(&models.File{
		Recycled: "N",
		Id:       uploadMate.FileId,
	}).First(&fileMate)
	return fileMate
}

func (h *Handlers) GetUploadBlockInfo(blockId int) models.UploadBlock {
	var uploadBlockMate models.UploadBlock
	LocalDB.Where(&models.UploadBlock{
		Recycled: "N",
		Id:       blockId,
	}).First(&uploadBlockMate)
	return uploadBlockMate
}

func (h *Handlers) UpdateFileComplete(uploadId int, state int) {
	var uploadMate models.Upload
	var fileMate models.File
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	LocalDB.Where(&models.File{
		Recycled: "N",
		Id:       uploadMate.FileId,
	}).First(&fileMate)
	uploadMate.IsComplete = state
	fileMate.IsComplete = state
	LocalDB.Save(&uploadMate)
	LocalDB.Save(&fileMate)
}

func (h *Handlers) UpdateBlockComplete(blockId int, state int) {
	var uploadBlockMate models.UploadBlock
	LocalDB.Where(&models.UploadBlock{
		Recycled: "N",
		Id:       blockId,
	}).First(&uploadBlockMate)
	uploadBlockMate.IsComplete = state
	LocalDB.Save(&uploadBlockMate)
}

func (h *Handlers) UpdateUploading(uploadId int, state int) {
	var uploadMate models.Upload
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	uploadMate.Uploading = state
	LocalDB.Save(&uploadMate)
}

func (h *Handlers) DeleteUpload(uploadId int) {
	var uploadMate models.Upload
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(&uploadMate)
	uploadMate.Recycled = "Y"
	LocalDB.Save(uploadMate)

	LocalDB.Where(&models.UploadBlock{
		UploadId: uploadId,
	}).Delete(models.UploadBlock{})
}

func (h *Handlers) SaveFile(data []byte, path string) error {
	fo, err := os.Create(path)
	if err == nil {
		_, err = fo.Write(data)
	}
	return err
}



func GetFileByPath(path string) string {
	path = strings.Replace(path, `\`, `/`, -1)
	if path == "" {
		return "."
	}
	for len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	if i := strings.LastIndex(path, "/"); i >= 0 {
		path = path[i+1:]
	}
	if path == "" {
		return "/"
	}
	return path
}

func StrSize(size int) string {
	floatSize := float64(size)
	unitList := [5]string{"B", "KB", "MB", "GB", "TB"}
	index := 0
	for floatSize > 1024 {
		index++
		floatSize, _ = strconv.ParseFloat(
			fmt.Sprintf("%.2f", floatSize/float64(1024)),
			64,
		)
	}
	return fmt.Sprintf("%.2f%s", floatSize, unitList[index])
}
