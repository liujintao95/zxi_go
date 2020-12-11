package download

import "zxi_go/core/models"

type ShowDownloadTable struct {
	Id          int     `json:"id"`
	Downloading int     `json:"downloading"`
	IsComplete  int     `json:"is_complete"`
	Name   string  `json:"name"`
	LocalPath   string  `json:"local_path"`
	Size        int     `json:"size"`
	Progress    float64 `json:"progress"`
}

type ShowDownloadInfo struct {
    Download    models.Download        `json:"download_map"`
    BlockList []models.DownloadBlock `json:"block_list"`
}