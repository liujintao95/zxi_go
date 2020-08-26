package models

type UserInfo struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	User     string `db:"user"`
	Pwd      string `db:"pwd"`
}

type Directory struct {
	Id       int64    `db:"id"`
	Name     string   `db:"name"`
	Fid      int64    `db:"fid"`
	UserInfo UserInfo `db:"user_id"`
	IsKey    int      `db:"is_key"`
}

type File struct {
	Id         int64  `db:"id"`
	Hash       string `db:"hash"`
	Path       string `db:"path"`
	Size       int64  `db:"size"`
	IsComplete int    `db:"is_complete"`
}

type UserFile struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	File      File      `db:"file_id"`
	Directory Directory `db:"dir_id"`
	IsKey     int       `db:"is_key"`
}

type Upload struct {
	Id         int64    `db:"id"`
	File       File     `db:"file_id"`
	UserInfo   UserInfo `db:"user_id"`
	LocalPath  string   `db:"local_path"`
	BlockSize  int      `db:"block_size"`
	IsComplete int      `db:"is_complete"`
}

type UploadBlock struct {
	Id         int64  `db:"id"`
	Upload     Upload `db:"upload_id"`
	Offset     int64  `db:"Offset"`
	Size       int64  `db:"size"`
	IsComplete int    `db:"is_complete"`
}

type Download struct {
	Id         int64    `db:"id"`
	File       File     `db:"file_id"`
	UserInfo   UserInfo `db:"user_id"`
	LocalPath  string   `db:"local_path"`
	BlockSize  int      `db:"block_size"`
	IsComplete int      `db:"is_complete"`
}

type DownloadBlock struct {
	Id         int64  `db:"id"`
	Download   Upload `db:"upload_id"`
	Offset     int64  `db:"Offset"`
	Size       int64  `db:"size"`
	IsComplete int    `db:"is_complete"`
}
