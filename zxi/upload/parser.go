package upload

import (
	"zxi_go/zxi/models"
)

type ShowUploadTable struct {
	Id         int     `json:"id"`
	Uploading  int     `json:"uploading"`
	IsComplete int     `json:"is_complete"`
	LocalPath  string  `json:"local_path"`
	Size       int     `json:"size"`
	Progress   float64 `json:"progress"`
}

type ShowUploadInfo struct {
	Upload    models.Upload        `json:"upload_map"`
	BlockList []models.UploadBlock `json:"block_list"`
}
