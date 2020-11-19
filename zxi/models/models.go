package models

import "time"

//type UserInfo struct {
//	Id   int64  `db:"id"`
//	Name string `db:"name"`
//	User string `db:"user"`
//	Pwd  string `db:"pwd"`
//}
//
//type Directory struct {
//	Id       int64    `db:"id"`
//	Name     string   `db:"name"`
//	Path     string   `db:"path"`
//	UserInfo UserInfo `db:"user_id"`
//	IsKey    int      `db:"is_key"`
//}
//
//type File struct {
//	Id         int64  `db:"id"`
//	Hash       string `db:"hash"`
//	Path       string `db:"path"`
//	Size       int64  `db:"size"`
//	IsComplete int    `db:"is_complete"`
//}
//
//type UserFile struct {
//	Id        int64     `db:"id"`
//	Name      string    `db:"name"`
//	File      File      `db:"file_id"`
//	Directory Directory `db:"dir_id"`
//	UserInfo  UserInfo  `db:"user_id"`
//	IsKey     int       `db:"is_key"`
//}
//
//type Upload struct {
//	Id         int64    `db:"id"`
//	File       File     `db:"file_id"`
//	UserInfo   UserInfo `db:"user_id"`
//	LocalPath  string   `db:"local_path"`
//	BlockSize  int      `db:"block_size"`
//	Uploading  int      `db:"uploading"`
//	IsComplete int      `db:"is_complete"`
//}
//
//type UploadBlock struct {
//	Id         int64  `db:"id"`
//	Upload     Upload `db:"upload_id"`
//	Offset     int64  `db:"Offset"`
//	Size       int64  `db:"size"`
//	IsComplete int    `db:"is_complete"`
//}
//
//type Download struct {
//	Id         int64    `db:"id"`
//	File       File     `db:"file_id"`
//	UserInfo   UserInfo `db:"user_id"`
//	LocalPath  string   `db:"local_path"`
//	BlockSize  int      `db:"block_size"`
//	IsComplete int      `db:"is_complete"`
//}
//
//type DownloadBlock struct {
//	Id         int64  `db:"id"`
//	Download   Upload `db:"upload_id"`
//	Offset     int64  `db:"Offset"`
//	Size       int64  `db:"size"`
//	IsComplete int    `db:"is_complete"`
//}

type UserInfo struct {
	Id        int       `gorm:"primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	User      string    `gorm:"type:varchar(255);not null"`
	Pwd       string    `gorm:"type:varchar(255);not null"`
	Recycled  string    `gorm:"type:varchar(255);not null;default:'N'"`
	UpdatedAt time.Time `gorm:"column:utime"`
	CreatedAt time.Time `gorm:"column:ctime"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

type Directory struct {
	Id        int       `gorm:"primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Path      string    `gorm:"type:varchar(255);not null"`
	IsKey     int       `gorm:"type:int(255);not null"`
	Recycled  string    `gorm:"type:varchar(255);not null;default:'N'"`
	UpdatedAt time.Time `gorm:"column:utime"`
	CreatedAt time.Time `gorm:"column:ctime"`

	UserInfo   UserInfo `gorm:"ForeignKey:UserInfoId;AssociationForeignKey:Id"`
	UserInfoId int      `gorm:"column:user_id"`
}

func (Directory) TableName() string {
	return "directory"
}

type File struct {
	Id         int       `gorm:"primary_key"`
	Hash       string    `gorm:"type:varchar(255);not null"`
	Path       string    `gorm:"type:varchar(255);not null"`
	Size       int       `gorm:"type:int(255);not null"`
	IsComplete int       `gorm:"type:int(255);not null"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'"`
	UpdatedAt  time.Time `gorm:"column:utime"`
	CreatedAt  time.Time `gorm:"column:ctime"`
}

func (File) TableName() string {
	return "file"
}

type UserFile struct {
	Id        int       `gorm:"primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	IsKey     int       `gorm:"type:int(255);not null"`
	Recycled  string    `gorm:"type:varchar(255);not null;default:'N'"`
	UpdatedAt time.Time `gorm:"column:utime"`
	CreatedAt time.Time `gorm:"column:ctime"`

	File        File      `gorm:"ForeignKey:FileId;AssociationForeignKey:Id"`
	Directory   Directory `gorm:"ForeignKey:DirectoryId;AssociationForeignKey:Id"`
	UserInfo    UserInfo  `gorm:"ForeignKey:UserInfoId;AssociationForeignKey:Id"`
	FileId      int       `gorm:"column:file_id"`
	DirectoryId int       `gorm:"column:dir_id"`
	UserInfoId  int       `gorm:"column:user_id"`
}

func (UserFile) TableName() string {
	return "user_file"
}

type Upload struct {
	Id         int       `gorm:"primary_key"`
	LocalPath  string    `gorm:"type:varchar(255);not null"`
	BlockSize  int       `gorm:"type:int(255);not null"`
	Uploading  int       `gorm:"type:int(255);not null"`
	IsComplete int       `gorm:"type:int(255);not null"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'"`
	UpdatedAt  time.Time `gorm:"column:utime"`
	CreatedAt  time.Time `gorm:"column:ctime"`

	File       File     `gorm:"ForeignKey:FileId;AssociationForeignKey:Id"`
	UserInfo   UserInfo `gorm:"ForeignKey:UserInfoId;AssociationForeignKey:Id"`
	FileId     int      `gorm:"column:file_id"`
	UserInfoId int      `gorm:"column:user_id"`
}

func (Upload) TableName() string {
	return "upload"
}

type UploadBlock struct {
	Id         int       `gorm:"primary_key"`
	Offset     int       `gorm:"type:int(255);not null"`
	Size       int       `gorm:"type:varchar(255);not null"`
	IsComplete int       `gorm:"type:int(255);not null"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'"`
	UpdatedAt  time.Time `gorm:"column:utime"`
	CreatedAt  time.Time `gorm:"column:ctime"`

	Upload   Upload `gorm:"ForeignKey:UploadId;AssociationForeignKey:Id"`
	UploadId int    `gorm:"column:upload_id"`
}

func (UploadBlock) TableName() string {
	return "upload_block"
}

type Download struct {
	Id         int       `gorm:"primary_key"`
	LocalPath  string    `gorm:"type:varchar(255);not null"`
	BlockSize  int       `gorm:"type:int(255);not null"`
	IsComplete int       `gorm:"type:int(255);not null"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'"`
	UpdatedAt  time.Time `gorm:"column:utime"`
	CreatedAt  time.Time `gorm:"column:ctime"`

	File       File     `gorm:"ForeignKey:FileId;AssociationForeignKey:Id"`
	UserInfo   UserInfo `gorm:"ForeignKey:UserInfoId;AssociationForeignKey:Id"`
	FileId     int      `gorm:"column:file_id"`
	UserInfoId int      `gorm:"column:user_id"`
}

func (Download) TableName() string {
	return "download"
}

type DownloadBlock struct {
	Id         int       `gorm:"primary_key"`
	Offset     int       `gorm:"type:int(255);not null"`
	Size       int       `gorm:"type:varchar(255);not null"`
	IsComplete int       `gorm:"type:int(255);not null"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'"`
	UpdatedAt  time.Time `gorm:"column:utime"`
	CreatedAt  time.Time `gorm:"column:ctime"`

	Download   Upload `db:"upload_id"`
	DownloadId int    `gorm:"column:download_id"`
}

func (DownloadBlock) TableName() string {
	return "download_block"
}
