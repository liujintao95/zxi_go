package handler

import (
	"strings"
)

type DirTree struct {
	id       int64
	name     string
	DirChild []DirTree
}


func PathSplit(path string) (string, string) {
	i := strings.LastIndex(path, `\`)
	if path[:i] != ""{
		return path[:i], path[i+1:]
	} else {
		return `\`, path[i+1:]
	}

}
