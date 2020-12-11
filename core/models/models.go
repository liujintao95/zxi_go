package models

import "time"

type UserInfo struct {
	Id        int       `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	User      string    `gorm:"type:varchar(255);not null" json:"user"`
	Pwd       string    `gorm:"type:varchar(255);not null" json:"pwd"`
	Recycled  string    `gorm:"type:varchar(255);not null;default:'N'" json:"recycled"`
	UpdatedAt time.Time `gorm:"column:utime" json:"utime"`
	CreatedAt time.Time `gorm:"column:ctime" json:"ctime"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

type Directory struct {
	Id        int       `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Path      string    `gorm:"type:varchar(255);not null" json:"path"`
	IsKey     int       `gorm:"type:int(255);not null" json:"is_key"`
	Recycled  string    `gorm:"type:varchar(255);not null;default:'N'" json:"recycled"`
	UpdatedAt time.Time `gorm:"column:utime" json:"utime"`
	CreatedAt time.Time `gorm:"column:ctime" json:"ctime"`

	UserInfo   UserInfo `gorm:"ForeignKey:UserInfoId;AssociationForeignKey:Id"`
	UserInfoId int      `gorm:"column:user_id" json:"user_id"`
}

func (Directory) TableName() string {
	return "directory"
}

type File struct {
	Id         int       `gorm:"primary_key" json:"id"`
	Hash       string    `gorm:"type:varchar(255);not null" json:"hash"`
	Path       string    `gorm:"type:varchar(255);not null" json:"path"`
	Size       int       `gorm:"type:int(255);not null" json:"size"`
	IsComplete int       `gorm:"type:int(255);not null" json:"is_complete"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'" json:"recycled"`
	UpdatedAt  time.Time `gorm:"column:utime" json:"utime"`
	CreatedAt  time.Time `gorm:"column:ctime" json:"ctime"`
}

func (File) TableName() string {
	return "file"
}

type UserFile struct {
	Id        int       `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	IsKey     int       `gorm:"type:int(255);not null" json:"is_key"`
	Recycled  string    `gorm:"type:varchar(255);not null;default:'N'" json:"recycled"`
	UpdatedAt time.Time `gorm:"column:utime" json:"utime"`
	CreatedAt time.Time `gorm:"column:ctime" json:"ctime"`

	File        File      `gorm:"ForeignKey:FileId;AssociationForeignKey:Id"`
	Directory   Directory `gorm:"ForeignKey:DirectoryId;AssociationForeignKey:Id"`
	UserInfo    UserInfo  `gorm:"ForeignKey:UserInfoId;AssociationForeignKey:Id"`
	FileId      int       `gorm:"column:file_id" json:"file_id"`
	DirectoryId int       `gorm:"column:dir_id" json:"dir_id"`
	UserInfoId  int       `gorm:"column:user_id" json:"user_id"`
}

func (UserFile) TableName() string {
	return "user_file"
}

type Upload struct {
	Id         int       `gorm:"primary_key" json:"id"`
	LocalPath  string    `gorm:"type:varchar(255);not null" json:"local_path"`
	BlockSize  int       `gorm:"type:int(255);not null" json:"block_size"`
	Uploading  int       `gorm:"type:int(255);not null" json:"uploading"`
	IsComplete int       `gorm:"type:int(255);not null" json:"is_complete"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'" json:"recycled"`
	UpdatedAt  time.Time `gorm:"column:utime" json:"utime"`
	CreatedAt  time.Time `gorm:"column:ctime" json:"ctime"`

	Block []UploadBlock `json:"block_list"`
	File       File     `gorm:"ForeignKey:FileId;AssociationForeignKey:Id"`
	UserInfo   UserInfo `gorm:"ForeignKey:UserInfoId;AssociationForeignKey:Id"`
	FileId     int      `gorm:"column:file_id" json:"file_id"`
	UserInfoId int      `gorm:"column:user_id" json:"user_id"`
}

func (Upload) TableName() string {
	return "upload"
}

type UploadBlock struct {
	Id         int       `gorm:"primary_key" json:"id"`
	Offset     int       `gorm:"type:int(255);not null" json:"offset"`
	Size       int       `gorm:"type:varchar(255);not null" json:"size"`
	IsComplete int       `gorm:"type:int(255);not null" json:"is_complete"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'" json:"recycled"`
	UpdatedAt  time.Time `gorm:"column:utime" json:"utime"`
	CreatedAt  time.Time `gorm:"column:ctime" json:"ctime"`

	Upload   Upload `gorm:"ForeignKey:UploadId;AssociationForeignKey:Id"`
	UploadId int    `gorm:"column:upload_id" json:"upload_id"`
}

func (UploadBlock) TableName() string {
	return "upload_block"
}

type Download struct {
	Id          int       `gorm:"primary_key" json:"id"`
	LocalPath   string    `gorm:"type:varchar(255);not null" json:"local_path"`
	BlockSize   int       `gorm:"type:int(255);not null" json:"block_size"`
	Downloading int       `gorm:"type:int(255);not null" json:"downloading"`
	IsComplete  int       `gorm:"type:int(255);not null" json:"is_complete"`
	Recycled    string    `gorm:"type:varchar(255);not null;default:'N'" json:"recycled"`
	UpdatedAt   time.Time `gorm:"column:utime" json:"utime"`
	CreatedAt   time.Time `gorm:"column:ctime" json:"ctime"`

	Block []DownloadBlock  `json:"block_list"`
	File       File     `gorm:"ForeignKey:FileId;AssociationForeignKey:Id"`
	UserInfo   UserInfo `gorm:"ForeignKey:UserInfoId;AssociationForeignKey:Id"`
	FileId     int      `gorm:"column:file_id" json:"file_id"`
	UserInfoId int      `gorm:"column:user_id" json:"user_id"`
}

func (Download) TableName() string {
	return "download"
}

type DownloadBlock struct {
	Id         int       `gorm:"primary_key" json:"id"`
	Offset     int       `gorm:"type:int(255);not null" json:"offset"`
	Size       int       `gorm:"type:varchar(255);not null" json:"size"`
	IsComplete int       `gorm:"type:int(255);not null" json:"is_complete"`
	Recycled   string    `gorm:"type:varchar(255);not null;default:'N'" json:"recycled"`
	UpdatedAt  time.Time `gorm:"column:utime" json:"utime"`
	CreatedAt  time.Time `gorm:"column:ctime" json:"ctime"`

	Download   Upload `db:"upload_id"`
	DownloadId int    `gorm:"column:download_id" json:"download_id"`
}

func (DownloadBlock) TableName() string {
	return "download_block"
}
