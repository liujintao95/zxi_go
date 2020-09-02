package field

import "zxi_network_disk_go/zxi/models"

type ErrField struct {
	Error string `json:"err"`
}

type ShowFilesField struct {
	Error string `json:"err"`
	DirList []models.Directory `json:"dir_list"`
	FileList []models.UserFile `json:"file_list"`
}
