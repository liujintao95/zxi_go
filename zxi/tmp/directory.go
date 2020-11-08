package tmp

import (
	"zxi_go/core"
	"zxi_go/zxi/models"
)

type DirectoryManager struct {
	table string
}

func NewDirectoryManager() *DirectoryManager {
	return &DirectoryManager{table: "directory"}
}

func (d *DirectoryManager) GetByUserIdPathName(userId int64, path string, name string) (models.Directory, error) {
	dirMate := new(models.Directory)
	sql := `
		SELECT id, name, path, is_key
		FROM directory
		WHERE recycled = 'N'
		AND user_id = ?
		AND path = ?
		AND name = ?
	`
	row := core.Conn.QueryRow(sql, userId, path, name)
	err := row.Scan(
		&dirMate.Id, &dirMate.Name, &dirMate.Path, &dirMate.IsKey,
	)
	return *dirMate, err
}

func (d *DirectoryManager) GetRootByUserId(userId int64) ([]models.Directory, error) {
	var dirList []models.Directory
	sql := `
		SELECT id, name, path, is_key
		FROM directory
		WHERE recycled = 'N'
		AND path = '\\'
		AND user_id = ?
	`
	rows, err := core.Conn.Query(sql, userId)
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

func (d *DirectoryManager) GetByDirId(id int64) (models.Directory, error) {
	dirMate := new(models.Directory)
	sql := `
		SELECT id, name, path, is_key
		FROM directory
		WHERE recycled = 'N'
		AND id = ?
	`
	row := core.Conn.QueryRow(sql, id)
	err := row.Scan(
		&dirMate.Id, &dirMate.Name, &dirMate.Path, &dirMate.IsKey,
	)
	return *dirMate, err
}

func (d *DirectoryManager) GetListByUserIdPath(userId int64, path string) ([]models.Directory, error) {
	var dirList []models.Directory
	sql := `
		SELECT id, name, path, is_key
		FROM directory
		WHERE recycled = 'N'
		AND user_id = ?
		AND path = ?
	`
	rows, err := core.Conn.Query(sql, userId, path)
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
	res, err := core.Conn.Exec(
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
	_, err := core.Conn.Exec(
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
	_, err := core.Conn.Exec(sql, name, id)
	return err
}

func (d *DirectoryManager) UpdateKey(key int64, id int64) error {
	sql := `
		UPDATE directory 
		SET is_key = ?
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, key, id)
	return err
}

func (d *DirectoryManager) UpdateFId(FId int64, id int64) error {
	sql := `
		UPDATE directory 
		SET fid = ?
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, FId, id)
	return err
}

func (d *DirectoryManager) DelById(id int64) error {
	sql := `
		UPDATE directory 
		SET recycled = 'Y'
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, id)
	return err
}
