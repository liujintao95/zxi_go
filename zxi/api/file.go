package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
	"zxi_network_disk_go/conf"
	"zxi_network_disk_go/zxi/handler"
	"zxi_network_disk_go/zxi/models"
	"zxi_network_disk_go/zxi/parser"
)

func SaveFileInfo(g *gin.Context) {
	filesJsonList := g.PostFormArray("files_json")
	rootPath := g.PostForm("root")
	userInter, _ := g.Get("userInfo")
	userMate := userInter.(models.UserInfo)
	prefix, _ := handler.PathSplit(rootPath)

	for _, fileJson := range filesJsonList {
		// json数据反序列化
		postFile := new(parser.PostFileInfo)
		err := json.Unmarshal([]byte(fileJson), postFile)
		errCheck(g, err, "文件属性解析错误", http.StatusInternalServerError)

		// 保存目录信息
		var dirId int64
		filePath := strings.Replace(postFile.Path, prefix, ``, 1)
		dirPath, dirName := handler.PathSplit(filePath)
		for {
			dirPath, dirName = handler.PathSplit(dirPath)
			dirMate, err := directoryManager.GetByUserIdPathName(userMate.Id, dirPath, dirName)
			errCheck(g, err, "数据库查询错误", http.StatusInternalServerError)

			if dirMate.Id == 0 {
				lastId, err := directoryManager.Create(models.Directory{
					Name:     dirName,
					Path:     dirPath,
					UserInfo: userMate,
				})
				errCheck(g, err, "数据库写入错误", http.StatusInternalServerError)
				if dirId == 0{
					dirId = lastId
				}
			} else {
				if dirId == 0{
					dirId = dirMate.Id
				}
			}

			if dirPath == `\` {
				break //跳出循环
			}
		}

		// 保存文件信息以及处理需要上传的文件
		fileMate, _ := fileManager.GetByHash(postFile.Hash)
		var fileId int64
		if fileMate.Id != 0 {
			fileId = fileMate.Id
		} else {
			fileId, err = fileManager.Create(models.File{
				Hash: postFile.Hash,
				Path: filepath.Join(conf.SAVE_PATH, postFile.Hash),
				Size: postFile.Size,
			})
			errCheck(g, err, "数据库写入错误", http.StatusInternalServerError)
		}
		if fileMate.IsComplete != 1{
			// 网盘上没有该文件，需要上传
			uploadMate, _ := uploadManager.GetByUserIdFileId(userMate.Id, fileId)
			if uploadMate.Id == 0 {
				_, err = uploadManager.Create(models.Upload{
					File: models.File{Id: fileId},
					UserInfo: models.UserInfo{Id: userMate.Id},
					BlockSize: conf.BLOCK_SIZE,
					LocalPath: postFile.Path,
				})
			}
		}
		_, err = userFileManager.Create(models.UserFile{
			Name:      postFile.Name,
			File:      models.File{Id: fileId},
			Directory: models.Directory{Id: dirId},
		})
		errCheck(g, err, "数据库写入错误", http.StatusInternalServerError)
	}

	g.JSON(http.StatusOK, gin.H{
		"errmsg": "ok",
		"data":   "",
	})
}
