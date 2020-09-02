package handler

import (
	"strings"
)

func PathSplit(path string) (string, string) {
	i := strings.LastIndex(path, `\`)
	if path[:i] != "" {
		return path[:i], path[i+1:]
	} else {
		return `\`, path[i+1:]
	}
}
