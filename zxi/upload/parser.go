package upload

type ShowUpload struct {
	Id         int
	Name       string
	Size       string
	Progress   float64
	Uploading  int
	IsComplete int
	LocalPath  string
}
