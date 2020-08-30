package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

	var uploadPath []string

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
			dirMate, err := directoryManager.GetByUserIdPath(userMate.Id, dirPath)
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
		fileMate, err := fileManager.GetByHash(postFile.Hash)
		errCheck(g, err, "数据库查询错误", http.StatusInternalServerError)
		var fileId int64
		if fileMate.Id != 0 {
			// 如果有相关文件数据则触发秒传
			fileId = fileMate.Id
		} else {
			// 如果没有相关文件数据则需要上传
			uploadPath = append(uploadPath, postFile.Path)
			// 创建上传任务
			fileId, err = fileManager.Create(models.File{
				Hash: postFile.Hash,
				Path: postFile.Path,
				Size: postFile.Size,
			})
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
