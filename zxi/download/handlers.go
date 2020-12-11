package download

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"strconv"
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

func (h *Handler) CreateOrIgnoreDownloadBlock(downloadId int) {
	var downloadMate models.Download
	h.localDB.Where(&models.Download{
		Recycled: "N",
		Id:       downloadId,
	}).First(&downloadMate)
	h.localDB.Model(&downloadMate).Related(&downloadMate.File)

	if downloadMate.IsComplete == 0 {
		blockNum := (downloadMate.File.Size / downloadMate.BlockSize) + 1
		for i := 0; i < blockNum; i++ {
			var blockSize int
			var downloadBlockMate models.DownloadBlock
			if i+1 == blockNum {
				blockSize = downloadMate.File.Size % downloadMate.BlockSize
			} else {
				blockSize = downloadMate.BlockSize
			}
			h.localDB.Where(&models.DownloadBlock{
				Recycled: "N",
				DownloadId: downloadId,
				Offset:   i,
				Size:     blockSize,
			}).First(&downloadBlockMate)
			if downloadBlockMate.Id == 0 {
				downloadBlockMate.DownloadId = downloadId
				downloadBlockMate.Offset = i
				downloadBlockMate.Size = blockSize
				h.localDB.Create(&downloadBlockMate)
			}
		}
	}
}

func (h *Handler) CreateOrIgnoreDownload(fileId int, userId int, localPath string) int {
	var downloadMate models.Download
	h.localDB.Where(&models.Download{
		Recycled:   "N",
		FileId:       fileId,
		UserInfoId:       userId,
	}).First(&downloadMate)
	if downloadMate.Id == 0 {
		var uploadMate models.Upload
		h.localDB.Where(&models.Upload{
			Recycled:   "N",
			FileId:       fileId,
			UserInfoId:       userId,
		}).First(&uploadMate)

		downloadMate.FileId = fileId
		downloadMate.UserInfoId = userId
		downloadMate.LocalPath = localPath
		downloadMate.BlockSize = uploadMate.BlockSize
		downloadMate.Downloading = 1
		downloadMate.IsComplete = 0
		h.localDB.Create(&downloadMate)

	}
	return downloadMate.Id
}

func (h *Handler) GetDownloadTable(userId int, page int, size int) ([]ShowDownloadTable, int) {
	var result []ShowDownloadTable
	var count int
	var downloadList []models.Download
	start := (page - 1) * size
	h.localDB.Where(&models.Download{
		Recycled:   "N",
		UserInfoId: userId,
	}).Order("id desc").Limit(size).Offset(start).Find(&downloadList)
	h.localDB.Model(&models.Download{}).Where(&models.Download{
		Recycled:   "N",
		UserInfoId: userId,
	}).Count(&count)
	for _, downloadInfo := range downloadList {
		var progress float64
		var userFileMate models.UserFile

		h.localDB.Where(&models.UserFile{
			UserInfoId: userId,
			FileId: downloadInfo.FileId,
		}).First(&userFileMate)
		h.localDB.Model(&userFileMate).Related(&userFileMate.File)

		if downloadInfo.IsComplete == 1 {
			progress = 100
		} else {
			var totalNum, completeNum int
			var blockList []models.DownloadBlock
			h.localDB.Where(&models.DownloadBlock{
				Recycled:   "N",
				DownloadId: downloadInfo.Id,
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
		result = append(result, ShowDownloadTable{
			Id:          downloadInfo.Id,
			Downloading: downloadInfo.Downloading,
			IsComplete:  downloadInfo.IsComplete,
			Name: userFileMate.Name,
			LocalPath:   downloadInfo.LocalPath,
			Size:        userFileMate.File.Size,
			Progress:    progress,
		})
	}
	return result, count
}

func (h *Handler) GetDownloadInfo(downloadId int) models.Download {
    var downloadMate models.Download
    h.localDB.Where(&models.Download{
        Recycled: "N",
        Id:       downloadId,
    }).First(&downloadMate)
    h.localDB.Model(&downloadMate).Related(&downloadMate.Block)
    return downloadMate
}

func (h *Handler) GetProgress(downloadId int) float64 {
	var progress float64
	var totalNum, completeNum int
	var blockList []models.DownloadBlock

	h.localDB.Where(&models.DownloadBlock{
		Recycled: "N",
		DownloadId: downloadId,
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

func (h *Handler) GetFileInfoByDownloadId(downloadId int) models.File {
	var downloadMate models.Download
	var fileMate models.File
	h.localDB.Where(&models.Download{
		Recycled: "N",
		Id:       downloadId,
	}).First(&downloadMate)
	h.localDB.Model(&downloadMate).First(&fileMate)
	return fileMate
}

func (h *Handler) UpdateBlockComplete(blockId int, state int) {
	var downloadBlockMate models.DownloadBlock
	h.localDB.Where(&models.DownloadBlock{
		Recycled: "N",
		Id:       blockId,
	}).First(&downloadBlockMate)
	downloadBlockMate.IsComplete = state
	h.localDB.Save(&downloadBlockMate)
}

func (h *Handler) UpdateDownloading(downloadId int, state int) {
	var downloadMate models.Download
	h.localDB.Where(&models.Download{
		Recycled: "N",
		Id:       downloadId,
	}).First(&downloadMate)
	downloadMate.Downloading = state
	h.localDB.Save(&downloadMate)
}

func (h *Handler) DeleteDownload(downloadId int) {
	var downloadMate models.Download
	h.localDB.Where(&models.Download{
		Recycled: "N",
		Id:       downloadId,
	}).First(&downloadMate)
	downloadMate.Recycled = "Y"
	h.localDB.Save(downloadMate)

	h.localDB.Where(&models.DownloadBlock{
		DownloadId: downloadId,
	}).Delete(models.DownloadBlock{})
}

func (h *Handler) readFile(path string) (string, error){
	f, err := ioutil.ReadFile(path)
	return string(f), err
}