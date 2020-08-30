package handler

import (
	"strings"
	"zxi_network_disk_go/zxi/models"
)

type DirTree struct {
	id       int64
	name     string
	DirChild []DirTree
}

func SaveFile(file models.File) {

}

func createDirTree(dirList []models.Directory) {
	var dirTree map[string]DirTree
	for _, dirMate := range dirList {
		if dirMate.Name != "" {
			dirTree[dirMate.Name] = DirTree{
				id:   dirMate.Id,
				name: dirMate.Name,
			}
		}
	}
}

func PathSplit(path string) (string, string) {
	i := strings.LastIndex(path, `\`)
	if path[:i] != ""{
		return path[:i], path[i+1:]
	} else {
		return `\`, path[i+1:]
	}

}
