package repository

import (
	"zxi_network_disk_go/network_disk/models"
	"zxi_network_disk_go/utils"
)

type FileManager struct {
	table string
}

func NewFileManager() *FileManager {
	return &FileManager{table: "user_info"}
}

func (f *FileManager) GetById(id int) (models.File, error) {
	fileMate := new(models.File)
	sql := `
		SELECT id, hash, path, size, is_complete, utime, ctime
		FROM file
		WHERE recycled = 'N'
		AND id = ?
	`
	rows := utils.Conn.QueryRow(sql, id)
	err := rows.Scan(
		fileMate.Id, fileMate.Hash, fileMate.Path, fileMate.Size,
		fileMate.IsComplete, fileMate.Utime, fileMate.Utime,
	)
	return *fileMate, err
}

func (f *FileManager) GetByHash(hash string) (models.File, error) {
	fileMate := new(models.File)
	sql := `
		SELECT id, hash, path, size, is_complete, utime, ctime
		FROM file
		WHERE recycled = 'N'
		AND hash = ?
	`
	rows := utils.Conn.QueryRow(sql, hash)
	err := rows.Scan(
		fileMate.Id, fileMate.Hash, fileMate.Path, fileMate.Size,
		fileMate.IsComplete, fileMate.Utime, fileMate.Utime,
	)
	return *fileMate, err
}

func (f *FileManager) Create(fileMate models.File) (int64, error) {
	sql := `
		INSERT INTO file(
			hash, path, size, is_complete
		)
		VALUES(?, ?, ?, ?)
	`
	res, err := utils.Conn.Exec(
		sql,
		fileMate.Hash, fileMate.Path, fileMate.Size, fileMate.IsComplete,
	)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (f *FileManager) Update(fileMate models.File) error {
	sql := `
		UPDATE file 
		SET path = ?, size = ?, is_complete = ?
		WHERE hash = ?
	`
	_, err := utils.Conn.Exec(
		sql,
		fileMate.Path, fileMate.Size,
		fileMate.IsComplete, fileMate.Hash,
	)
	return err
}

func (f *FileManager) DelByHash(hash string) error {
	sql := `
		UPDATE file 
		SET recycled = 'Y'
		WHERE hash = ?
	`
	_, err := utils.Conn.Exec(sql, hash)
	return err
}

