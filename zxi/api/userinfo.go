package api

import (
	"zxi_go/utils"
	"zxi_go/zxi/repository"
)

var errCheck = utils.ErrCheck
var fileManager = repository.NewFileManager()
var userFileManager = repository.NewUserFileManager()
var directoryManager = repository.NewDirectoryManager()
var uploadManager = repository.NewUploadManager()
