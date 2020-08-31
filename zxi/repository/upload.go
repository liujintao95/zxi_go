package repository

import (
	"zxi_network_disk_go/utils"
	"zxi_network_disk_go/zxi/models"
)

type UploadManager struct {
	table string
}

func NewUploadManager() *UploadManager {
	return &UploadManager{table: "upload"}
}

func (u *UploadManager) GetByUserIdFileId(userId int64, fileId int64) (models.Upload, error) {
	uploadMate := new(models.Upload)
	sql := `
		SELECT up.id, ui.name, ui.user, fi.hash, fi.size, fi.path,
		fi.is_complete, local_path, block_size, up.is_complete
		FROM upload AS up
		INNER JOIN file AS fi
		ON up.file_id = fi.id
		INNER JOIN user_info AS ui
		ON up.user_id = ui.id
		WHERE up.recycled = 'N'
		AND up.user_id = ?
		AND up.file_id = ?
	`
	rows := utils.Conn.QueryRow(sql, userId, fileId)
	err := rows.Scan(
		&uploadMate.Id, &uploadMate.UserInfo.Name, &uploadMate.UserInfo.User,
		&uploadMate.File.Hash, &uploadMate.File.Size, &uploadMate.File.Path,
		&uploadMate.File.IsComplete, &uploadMate.LocalPath, &uploadMate.BlockSize,
		&uploadMate.IsComplete,
	)
	return *uploadMate, err
}

func (u *UploadManager) GetById(id int64) (models.Upload, error) {
	uploadMate := new(models.Upload)
	sql := `
		SELECT up.id, ui.name, ui.user, fi.hash, fi.size, fi.path,
		fi.is_complete, local_path, block_size, up.is_complete
		FROM upload AS up
		INNER JOIN file AS fi
		ON up.file_id = fi.id
		INNER JOIN user_info AS ui
		ON up.user_id = ui.id
		WHERE up.recycled = 'N'
		AND up.id = ?
	`
	rows := utils.Conn.QueryRow(sql, id)
	err := rows.Scan(
		&uploadMate.Id, &uploadMate.UserInfo.Name, &uploadMate.UserInfo.User,
		&uploadMate.File.Hash, &uploadMate.File.Size, &uploadMate.File.Path,
		&uploadMate.File.IsComplete, &uploadMate.LocalPath, &uploadMate.BlockSize,
		&uploadMate.IsComplete,
	)
	return *uploadMate, err
}

func (u *UploadManager) GetListByUserId(userId int64) ([]models.Upload, error) {
	var uploadList []models.Upload
	sql := `
		SELECT up.id, ui.name, ui.user, fi.hash, fi.size, fi.path,
		fi.is_complete, local_path, block_size, up.is_complete
		FROM upload AS up
		INNER JOIN file AS fi
		ON up.file_id = fi.id
		INNER JOIN user_info AS ui
		ON up.user_id = ui.id
		WHERE up.recycled = 'N'
		AND up.user_id = ?
	`
	rows, err := utils.Conn.Query(sql, userId)
	if err != nil {
		return uploadList, err
	}
	for rows.Next() {
		uploadMate := new(models.Upload)
		_ = rows.Scan(
			&uploadMate.Id, &uploadMate.UserInfo.Name, &uploadMate.UserInfo.User,
			&uploadMate.File.Hash, &uploadMate.File.Size, &uploadMate.File.Path,
			&uploadMate.File.IsComplete, &uploadMate.LocalPath, &uploadMate.BlockSize,
			&uploadMate.IsComplete,
		)
		uploadList = append(uploadList, *uploadMate)
	}
	return uploadList, err
}

func (u *UploadManager) Create(uploadMate models.Upload) (int64, error) {
	sql := `
		INSERT INTO upload(
			file_id, user_id, local_path, block_size, is_complete
		)
		VALUES(?, ?, ?, ?, ?)
	`
	res, err := utils.Conn.Exec(
		sql,
		uploadMate.File.Id, uploadMate.UserInfo.Id, uploadMate.LocalPath,
		uploadMate.BlockSize, uploadMate.IsComplete,
	)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (u *UploadManager) Update(uploadMate models.Upload) error {
	sql := `
		UPDATE upload 
		SET local_path = ?, block_size = ?, is_complete = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(
		sql,
		uploadMate.LocalPath, uploadMate.BlockSize, uploadMate.IsComplete,
		uploadMate.Id,
	)
	return err
}

func (u *UploadManager) UpdateComplete(complete int, id int64) error {
	sql := `
		UPDATE upload 
		SET is_complete = ?
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, complete, id)
	return err
}

func (u *UploadManager) DeleteById(id int64) error {
	sql := `
		UPDATE upload 
		SET recycled = 'Y'
		WHERE id = ?
	`
	_, err := utils.Conn.Exec(sql, id)
	return err
}