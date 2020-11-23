package upload

import (
	"fmt"
	"strconv"
	"strings"
	. "zxi_go/core"
	"zxi_go/zxi/models"
)

type Handlers struct {
}

var handlers = Handlers{}

func (h *Handlers) CreateOrIgnoreUploadBlock(uploadId int, size int) {
	uploadMate := new(models.Upload)
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(uploadMate)
	if uploadMate.IsComplete == 0 {
		blockNum := (size / BLOCK_SIZE) + 1
		for i := 0; i < blockNum; i++ {
			var blockSize int
			if i+1 == blockNum {
				blockSize = size % BLOCK_SIZE
			} else {
				blockSize = BLOCK_SIZE
			}

			uploadBlockMate := new(models.UploadBlock)
			LocalDB.Where(&models.UploadBlock{
				Recycled: "N",
				UploadId: uploadId,
				Offset:   i,
				Size:     blockSize,
			}).First(uploadBlockMate)
			if uploadBlockMate.Id == 0 {
				uploadBlockMate.UploadId = uploadId
				uploadBlockMate.Offset = i
				uploadBlockMate.Size = blockSize
				LocalDB.Create(uploadBlockMate)
			}
		}
	}
}

func (h *Handlers) CreateOrIgnoreUpload(hash string, path string, userId int) int {
	fileMate := new(models.File)
	uploadMate := new(models.Upload)

	LocalDB.Where(&models.File{
		Recycled: "N",
		Hash:     hash,
	}).First(fileMate)
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		FileId:   fileMate.Id,
	}).First(uploadMate)

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
		LocalDB.Create(uploadMate)
	}
	return uploadMate.Id
}

func (h *Handlers) GetUploadList(userId int, page int, size int) ([]ShowUpload, int) {
	var count int
	start := (page - 1) * size
	uploadList := new([]models.Upload)

	LocalDB.Where(&models.Upload{
		Recycled:   "N",
		UserInfoId: userId,
	}).Order("id desc").Limit(size).Offset(start).Find(uploadList)
	LocalDB.Model(&models.Upload{}).Where(&models.Upload{
		Recycled:   "N",
		UserInfoId: userId,
	}).Count(&count)

	var result []ShowUpload
	for _, uploadInfo := range *uploadList {
		fileMate := new(models.File)
		LocalDB.Where(&models.File{
			Id: uploadInfo.FileId,
		}).First(fileMate)
		if uploadInfo.IsComplete == 1 {
			result = append(result, ShowUpload{
				Id:         uploadInfo.Id,
				Uploading:  uploadInfo.Uploading,
				IsComplete: uploadInfo.IsComplete,
				LocalPath:  uploadInfo.LocalPath,
				Name:       GetFileByPath(uploadInfo.LocalPath),
				Size:       StrSize(fileMate.Size),
				Progress:   100,
			})
		} else {
			totalNum := 0
			completeNum := 0
			blockList := new([]models.UploadBlock)
			LocalDB.Where(&models.UploadBlock{
				Recycled: "N",
				UploadId: uploadInfo.Id,
			})

			for _, blockMate := range *blockList {
				if blockMate.IsComplete == 1 {
					completeNum++
				}
				totalNum++
			}
			var progress float64
			if totalNum != 0 {
				progress, _ = strconv.ParseFloat(
					fmt.Sprintf("%.2f", float64(completeNum)/float64(totalNum)*100),
					64,
				)
			} else {
				progress = 0
			}
			result = append(result, ShowUpload{
				Id:         uploadInfo.Id,
				Uploading:  uploadInfo.Uploading,
				IsComplete: uploadInfo.IsComplete,
				LocalPath:  uploadInfo.LocalPath,
				Name:       GetFileByPath(uploadInfo.LocalPath),
				Size:       StrSize(fileMate.Size),
				Progress:   progress,
			})
		}
	}
	return result, count
}

func (h *Handlers) GetUploadInfo(uploadId int) models.Upload {
	uploadMate := new(models.Upload)
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(uploadMate)
	return *uploadMate
}

func (h *Handlers) GetFileInfoByUploadId(uploadId int) models.File {
	uploadMate := new(models.Upload)
	fileMate := new(models.File)
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(uploadMate)
	LocalDB.Where(&models.File{
		Recycled: "N",
		Id:       uploadMate.FileId,
	}).First(fileMate)
	return *fileMate
}

func (h *Handlers) SetIsComplete(uploadId int, state int) {
	uploadMate := new(models.Upload)
	fileMate := new(models.File)
	LocalDB.Where(&models.Upload{
		Recycled: "N",
		Id:       uploadId,
	}).First(uploadMate)
	LocalDB.Where(&models.File{
		Recycled: "N",
		Id:       uploadMate.FileId,
	}).First(fileMate)
	uploadMate.IsComplete = state
	fileMate.IsComplete = state
	LocalDB.Save(uploadMate)
	LocalDB.Save(fileMate)
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
