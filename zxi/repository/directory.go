package repository

import (
	"zxi_network_disk_go/utils"
	"zxi_network_disk_go/zxi/models"
)

type DirectoryManager struct {
	table string
}

func NewDirectoryManager() *DirectoryManager {
	return &DirectoryManager{table: "directory"}
}

func (d *DirectoryManager) GetByUserIdPath(userId int64, path string) (models.Directory, error) {
	dirMate := new(models.Directory)
	sql := `
		SELECT id, name, path, is_key
		FROM directory
		WHERE recycled = 'N'
		AND user_id = ?
		AND path = ?
	`
	row := utils.Conn.QueryRow(sql, userId, path)
	err := row.Scan(
		&dirMate.Id, &dirMate.Name, &dirMate.Path, &dirMate.IsKey,
	)
	return *dirMate, err
}

func (d *DirectoryManager) GetRootByUserId(userId int64) (models.Directory, error) {
	dirMate := new(models.Directory)
	sql := `
		SELECT id, name, path, is_key
		FROM directory
		WHERE recycled = 'N'
		AND fid = -1
		AND user_id = ?
	`
	row := utils.Conn.QueryRow(sql, userId)
	err := row.Scan(
		&dirMate.Id, &dirMate.Name, &dirMate.Path, &dirMate.IsKey,
	)
	return *dirMate, err
}

func (d *DirectoryManager) GetByDirId(id int64) (models.Directory, error) {
	dirMate := new(models.Directory)
	sql := `
		SELECT id, name, path, is_key
		FROM directory
		WHERE recycled = 'N'
		AND fid = ?
	`
	row := utils.Conn.QueryRow(sql, id)
	err := row.Scan(
		&dirMate.Id, &dirMate.Name, &dirMate.Path, &dirMate.IsKey,
	)
	return *dirMate, err
}

func (d *DirectoryManager) GetListByFId(FId int64) ([]models.Directory, error) {
	var dirList []models.Directory
	sql := `
		SELECT id, name, path, is_key
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
			&dirMate.Id, &dirMate.Name, &dirMate.Path, &dirMate.IsKey,
		)
		dirList = append(dirList, *dirMate)
	}
	return dirList, err
}

func (d *DirectoryManager) Create(dirMate models.Directory) (int64, error) {
	sql := `
		INSERT INTO directory(
			name, path, is_key, user_id
		)
		VALUES(?, ?, ?, ?)
	`
	res, err := utils.Conn.Exec(
		sql,
		dirMate.Name, dirMate.Path, dirMate.IsKey, dirMate.UserInfo.Id,
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
		SET name = ?, path = ?, is_key = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(
		sql,
		dirMate.Name, dirMate.Path, dirMate.IsKey, dirMate.Id,
	)
	return err
}

func (d *DirectoryManager) UpdateName(name string, id int64) error {
	sql := `
		UPDATE directory 
		SET name = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, name, id)
	return err
}

func (d *DirectoryManager) UpdateKey(key int64, id int64) error {
	sql := `
		UPDATE directory 
		SET is_key = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, key, id)
	return err
}

func (d *DirectoryManager) UpdateFId(FId int64, id int64) error {
	sql := `
		UPDATE directory 
		SET fid = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, FId, id)
	return err
}

func (d *DirectoryManager) DelById(id int64) error {
	sql := `
		UPDATE directory 
		SET recycled = 'Y'
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, id)
	return err
}
