package api

import (
	"zxi_network_disk_go/utils"
	"zxi_network_disk_go/zxi/repository"
)

var errCheck = utils.ErrCheck
var fileManager = repository.NewFileManager()
var userFileManager = repository.NewUserFileManager()
var directoryManager = repository.NewDirectoryManager()
