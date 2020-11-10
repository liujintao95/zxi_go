package handlers

import (
	"path/filepath"
	"strings"
	consts "zxi_go/zxi/file/const"
	"zxi_go/zxi/file/repository"
	"zxi_go/zxi/models"
)

func PathSplit(path string) (string, string) {
	i := strings.LastIndex(path, `\`)
	if path[:i] != "" {
		return path[:i], path[i+1:]
	} else {
		return `\`, path[i+1:]
	}
}

func AbsolutePathToRelativePath(absolutePath string, rootPath string) string {
	prefix, _ := PathSplit(rootPath)
	return strings.Replace(absolutePath, prefix, ``, 1)
}

func CreateOrIgnoreFile(hash string, size int64) (int64, error) {
	fileMate, err := repository.GetFileInfoByHash(hash)
	if err != nil {
		return 0, err
	}
	if fileMate.Id != 0 {
		return fileMate.Id, nil
	} else {
		return repository.CreateFileInfo(models.File{
			Hash: hash,
			Path: filepath.Join(consts.SAVE_PATH, hash),
			Size: size,
		})
	}
}

func CreateOrIgnoreUpload(hash string, path string, userId int64) (int64, error) {
	fileMate, err := repository.GetFileInfoByHash(hash)
	if err != nil {
		return 0, err
	}
	if fileMate.IsComplete != 1 {
		// 网盘上没有该文件，需要上传
		uploadMate, err := repository.GetUploadInfoByUserIdFileId(userId, fileMate.Id)
		if err != nil {
			return 0, err
		}
		if uploadMate.Id == 0 {
			return repository.CreateUploadInfo(models.Upload{
				File:      models.File{Id: fileMate.Id},
				UserInfo:  models.UserInfo{Id: userId},
				BlockSize: consts.BLOCK_SIZE,
				LocalPath: path,
			})
		}
	}
	return 0, nil
}

func CreateUploadBlock(uploadId int64, size int64) error {
	blockNum := (size / consts.BLOCK_SIZE) + 1
	for i := int64(0); i < blockNum; i++ {
		var blockSize int64
		if i+1 == blockNum {
			blockSize = size % consts.BLOCK_SIZE
		} else {
			blockSize = consts.BLOCK_SIZE
		}
		_, err := repository.CreateUploadBlock(models.UploadBlock{
			Upload: models.Upload{Id: uploadId},
			Offset: i * consts.BLOCK_SIZE,
			Size:   blockSize,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateUserFile(hash string, fileName string, userId int64, dirId int64) (int64, error) {
	fileMate, err := repository.GetFileInfoByHash(hash)
	if err != nil {
		return 0, err
	}
	return repository.CreateUserFileInfo(models.UserFile{
		Name:      fileName,
		File:      models.File{Id: fileMate.Id},
		Directory: models.Directory{Id: dirId},
		UserInfo:  models.UserInfo{Id: userId},
	})
}

func CreateOrIgnoreDir(dirName string, dirPath string, userId int64) (int64, error) {
	dirMate, err := repository.GetDirInfoByUserIdPathName(userId, dirPath, dirName)
	if err != nil {
		return 0, err
	}

	if dirMate.Id == 0 {
		return repository.CreateDirInfo(models.Directory{
			Name:     dirName,
			Path:     dirPath,
			UserInfo: models.UserInfo{Id: userId},
		})
	}
	return 0, nil
}

func CreatePath(relativePath string, userId int64) (int64, error) {
	var dirId int64
	// 绝对路径转相对路径

	// 分割路径与文件
	dirPath, _ := PathSplit(relativePath)
	for {
		dirPath, dirName := PathSplit(dirPath)
		lastId, err := CreateOrIgnoreDir(dirName, dirPath, userId)
		if err != nil {
			return 0, err
		}
		if dirId == 0 {
			dirId = lastId
		}
		if dirPath == `\` {
			break //跳出循环
		}
	}
	return dirId, nil
}

func GetDirList(path string, userId int64) ([]models.Directory, error) {
	if path == "" {
		return repository.GetRootDirListByUserId(userId)
	} else {
		return repository.GetDirListByUserIdPath(userId, path)
	}
}

func GetFileList(dirId int64, userId int64) ([]models.UserFile, error) {
	if dirId == 0 {
		return repository.GetRootFileListByUserId(userId)
	} else {
		return repository.GetFileListByDirId(dirId)
	}
}

func GetUploadList(userId int64) ([]models.Upload, error) {
	return repository.GetUploadListByUserId(userId)
}

func GetUploadInfo(uploadId int64) (models.Upload, error) {
	return repository.GetUploadInfoById(uploadId)
}

func SetIsComplete(uploadId int64, state int) error {
	return repository.UpdateUploadComplete(state, uploadId)
}