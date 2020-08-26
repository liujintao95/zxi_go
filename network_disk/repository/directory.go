package repository

import (
	"zxi_network_disk_go/network_disk/models"
	"zxi_network_disk_go/utils"
)

type DirectoryManager struct {
	table string
}

func NewDirectoryManager() *DirectoryManager {
	return &DirectoryManager{table: "directory"}
}

func (d *DirectoryManager) GetRootByUserId(userId int) (models.Directory, error) {
	dirMate := new(models.Directory)
	sql := `
		SELECT id, name, fid, is_key
		FROM directory
		WHERE recycled = 'N'
		AND fid = -1
		AND user_id = ?
	`
	row := utils.Conn.QueryRow(sql, userId)
	err := row.Scan(
		dirMate.Id, dirMate.Name, dirMate.Fid, dirMate.IsKey,
	)
	return *dirMate, err
}

func (d *DirectoryManager) GetByDirId(id int) (models.Directory, error) {
	dirMate := new(models.Directory)
	sql := `
		SELECT id, name, fid, is_key
		FROM directory
		WHERE recycled = 'N'
		AND fid = ?
	`
	row := utils.Conn.QueryRow(sql, id)
	err := row.Scan(
		dirMate.Id, dirMate.Name, dirMate.Fid, dirMate.IsKey,
	)
	return *dirMate, err
}

func (d *DirectoryManager) GetListByFId(FId int) ([]models.Directory, error) {
	var dirList []models.Directory
	sql := `
		SELECT id, name, fid, is_key
		FROM directory
		WHERE recycled = 'N'
		AND fid = ?
	`
	rows, err := utils.Conn.Query(sql, FId)
	if err != nil {
		return dirList, err
	}
	for rows.Next() {
		dirMate := new(models.Directory)
		_ = rows.Scan(
			dirMate.Id, dirMate.Name, dirMate.Fid, dirMate.IsKey,
		)
		dirList = append(dirList, *dirMate)
	}
	return dirList, err
}

func (d *DirectoryManager) Create(dirMate models.Directory) (int64, error) {
	sql := `
		INSERT INTO directory(
			name, fid, is_key, user_id
		)
		VALUES(?, ?, ?, ?)
	`
	res, err := utils.Conn.Exec(
		sql,
		dirMate.Name, dirMate.Fid, dirMate.IsKey, dirMate.UserInfo.Id,
	)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (d *DirectoryManager) Update(dirMate models.Directory) error {
	sql := `
		UPDATE directory 
		SET name = ?, fid = ?, is_key = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(
		sql,
		dirMate.Name, dirMate.Fid, dirMate.IsKey, dirMate.Id,
	)
	return err
}

func (d *DirectoryManager) UpdateName(name string, id int) error {
	sql := `
		UPDATE directory 
		SET name = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, name, id)
	return err
}

func (d *DirectoryManager) UpdateKey(key int, id int) error {
	sql := `
		UPDATE directory 
		SET is_key = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, key, id)
	return err
}

func (d *DirectoryManager) UpdateFId(FId int, id int) error {
	sql := `
		UPDATE directory 
		SET fid = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, FId, id)
	return err
}

func (d *DirectoryManager) DelById(id int) error {
	sql := `
		UPDATE directory 
		SET recycled = 'Y'
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, id)
	return err
}