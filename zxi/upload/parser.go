package upload

import (
	"zxi_go/zxi/models"
)

type ShowUploadTable struct {
	Id int
	Uploading int
	IsComplete int
	LocalPath string
	Name string
	Size string
	Progress   float64
}

type ShowUploadInfo struct {
	Id         int
	LocalPath  string
	BlockSize  int
	Uploading  int
	IsComplete int
	FileId     int
	UserInfoId int
	BlockList []models.UploadBlock
}
