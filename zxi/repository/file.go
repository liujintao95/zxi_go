package repository

import (
    "zxi_go/core"
    "zxi_go/zxi/models"
)

type FileManager struct {
	table string
}

func NewFileManager() *FileManager {
	return &FileManager{table: "file"}
}

func (f *FileManager) GetById(id int64) (models.File, error) {
	fileMate := new(models.File)
	sql := `
		SELECT id, hash, path, size, is_complete
		FROM file
		WHERE recycled = 'N'
		AND id = ?
	`
	row := core.Conn.QueryRow(sql, id)
	err := row.Scan(
		&fileMate.Id, &fileMate.Hash, &fileMate.Path, &fileMate.Size, &fileMate.IsComplete,
	)
	return *fileMate, err
}

func (f *FileManager) GetByHash(hash string) (models.File, error) {
	fileMate := new(models.File)
	sql := `
		SELECT id, hash, path, size, is_complete
		FROM file
		WHERE recycled = 'N'
		AND hash = ?
	`
	row := core.Conn.QueryRow(sql, hash)
	err := row.Scan(
		&fileMate.Id, &fileMate.Hash, &fileMate.Path, &fileMate.Size, &fileMate.IsComplete,
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
	res, err := core.Conn.Exec(
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
		WHERE id = ?
	`
	_, err := core.Conn.Exec(
		sql,
		fileMate.Path, fileMate.Size, fileMate.IsComplete, fileMate.Id,
	)
	return err
}

func (f *FileManager) UpdateComplete(complete int, id int64) error {
	sql := `
		UPDATE file 
		SET is_complete = ?
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, complete, id)
	return err
}

func (f *FileManager) DelByHash(hash string) error {
	sql := `
		UPDATE file 
		SET recycled = 'Y'
		WHERE hash = ?
	`
	_, err := core.Conn.Exec(sql, hash)
	return err
}

func (f *FileManager) DelById(id int64) error {
	sql := `
		UPDATE file 
		SET recycled = 'Y'
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, id)
	return err
}
