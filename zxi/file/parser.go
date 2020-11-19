package file

type DirInfo struct {
	Files []FileInfo `form:"files" json:"files"`
	Root  string     `form:"root" json:"root"`
}

type FileInfo struct {
	Size int    `form:"size" json:"size"`
	Hash string `form:"hash" json:"hash"`
	Name string `form:"name" json:"name"`
	Path string `form:"path" json:"path"`
}
