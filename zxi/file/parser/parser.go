package parser

type DirInfo struct {
	Files []FileInfo `form:"files" json:"files"`
	Root string `form:"root" json:"root"`
}

type FileInfo struct {
	Size int64 `form:"size" json:"size"`
	Hash string `form:"hash" json:"hash"`
	Name string `form:"name" json:"name"`
	Path string `form:"path" json:"path"`
}

type UploadFile struct {
	Id int64 `form:"id" json:"id"`
}

type ShowUpload struct {
	Id       int64
	Name     string
	Size     string
	Progress float64
}


