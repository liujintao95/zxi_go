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
		var fileMate models.File
		h.localDB.Where(&models.File{
			Id: downloadInfo.FileId,
		}).First(&fileMate)
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
			LocalPath:   downloadInfo.LocalPath,
			Size:        fileMate.Size,
			Progress:    progress,
		})
	}
	return result, count
}

func (h *Handler) GetDownloadInfo(downloadId int) ShowDownloadInfo {
    var downloadMate models.Download
    var downloadBlockList []models.DownloadBlock
    h.localDB.Where(&models.Download{
        Recycled: "N",
        Id:       downloadId,
    }).First(&downloadMate)
    h.localDB.Where(&models.DownloadBlock{
        Recycled: "N",
        DownloadId: downloadId,
    }).Find(&downloadBlockList)
    return ShowDownloadInfo{
        Download:    downloadMate,
        BlockList: downloadBlockList,
    }
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
	h.localDB.Where(&models.File{
		Recycled: "N",
		Id:       downloadMate.FileId,
	}).First(&fileMate)
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