package upload

import "zxi_go/zxi/models"

type ShowUploadTable struct {
	*models.Upload
	Name string
	Size string
	Progress   float64
}

type ShowUploadInfo struct {
	*models.Upload
	BlockList []models.UploadBlock
}
