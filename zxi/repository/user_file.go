package repository

import (
    "zxi_go/core"
    "zxi_go/zxi/models"
)

type UserFileManager struct {
	table string
}

func NewUserFileManager() *UserFileManager {
	return &UserFileManager{table: "file_dir"}
}

func (f *UserFileManager) GetRootByUserId(userId int64) ([]models.UserFile, error) {
	var userFileList []models.UserFile
	sql := `
		SELECT uf.id AS id, uf.name AS file_name, uf.is_key AS file_key,
		f.id AS file_id, f.hash, f.size, f.path
		FROM user_file AS uf
		INNER JOIN file AS f
		ON uf.file_id = f.id
		WHERE uf.recycled = 'N'
		AND f.recycled = 'N'
		AND f.is_complete = 1
		AND uf.dir_id = 0
		AND uf.user_id = ?
	`
	rows, err := core.Conn.Query(sql, userId)
	if err != nil {
		return userFileList, err
	}
	for rows.Next() {
		userFileMate := new(models.UserFile)
		_ = rows.Scan(
			&userFileMate.Id, &userFileMate.Name, &userFileMate.IsKey, &userFileMate.File.Id,
			&userFileMate.File.Hash, &userFileMate.File.Size, &userFileMate.File.Path,
		)
		userFileList = append(userFileList, *userFileMate)
	}
	return userFileList, err
}

func (f *UserFileManager) GetListByDirId(dirId int64) ([]models.UserFile, error) {
	var userFileList []models.UserFile
	sql := `
		SELECT d.id AS dir_id, d.name AS dir_name, d.is_key AS dir_key,
		uf.id AS map_id, uf.name AS file_name, uf.is_key AS file_key,
		f.id AS file_id, f.hash, f.size, f.path
		FROM user_file AS uf
		INNER JOIN directory AS d
		ON uf.dir_id = d.id
		INNER JOIN file AS f
		ON uf.file_id = f.id
		WHERE d.recycled = 'N'
		AND uf.recycled = 'N'
		AND f.recycled = 'N'
		AND f.is_complete = 1
		AND uf.dir_id = ?
	`
	rows, err := core.Conn.Query(sql, dirId)
	if err != nil {
		return userFileList, err
	}
	for rows.Next() {
		userFileMate := new(models.UserFile)
		_ = rows.Scan(
			&userFileMate.Directory.Id, &userFileMate.Directory.Name, &userFileMate.Directory.IsKey,
			&userFileMate.Id, &userFileMate.Name, &userFileMate.IsKey, &userFileMate.File.Id,
			&userFileMate.File.Hash, &userFileMate.File.Size, &userFileMate.File.Path,
		)
		userFileList = append(userFileList, *userFileMate)
	}
	return userFileList, err
}

func (f *UserFileManager) Create(userFileMate models.UserFile) (int64, error) {
	sql := `
		INSERT INTO user_file(
			file_id, dir_id, user_id, name, is_key
		)
		VALUES(?, ?, ?, ?, ?)
	`
	res, err := core.Conn.Exec(
		sql,
		userFileMate.File.Id, userFileMate.Directory.Id, userFileMate.UserInfo.Id,
		userFileMate.Name, userFileMate.IsKey,
	)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (f *UserFileManager) Update(userFileMate models.UserFile) error {
	sql := `
		UPDATE user_file 
		SET dir_id = ?, name = ?, is_key = ?
		WHERE id = ?
	`
	_, err := core.Conn.Exec(
		sql,
		userFileMate.Directory.Id, userFileMate.Name, userFileMate.IsKey, userFileMate.Id,
	)
	return err
}

func (f *UserFileManager) UpdateName(name string, id int64) error {
	sql := `
		UPDATE user_file 
		SET name = ?
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, name, id)
	return err
}

func (f *UserFileManager) UpdateKey(key int, id int64) error {
	sql := `
		UPDATE user_file 
		SET is_key = ?
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, key, id)
	return err
}

func (f *UserFileManager) UpdateDirId(dirId int, id int64) error {
	sql := `
		UPDATE user_file 
		SET dir_id = ?
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, dirId, id)
	return err
}

func (f *UserFileManager) DelById(id int64) error {
	sql := `
		UPDATE user_file 
		SET recycled = 'Y'
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, id)
	return err
}
