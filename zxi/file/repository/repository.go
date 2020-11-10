package repository

import (
	"zxi_go/core"
	"zxi_go/zxi/models"
)

func CreateDirInfo(dirMate models.Directory) (int64, error) {
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

func CreateFileInfo(fileMate models.File) (int64, error) {
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

func CreateUploadInfo(uploadMate models.Upload) (int64, error) {
	sql := `
		INSERT INTO upload(
			file_id, user_id, local_path, block_size, is_complete
		)
		VALUES(?, ?, ?, ?, ?)
	`
	res, err := core.Conn.Exec(
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

func CreateUploadBlock(uploadBlockMate models.UploadBlock) (int64, error) {
	sql := `
		INSERT INTO upload_block(
			upload_id, offset, size, is_complete
		)
		VALUES(?, ?, ?, ?)
	`
	res, err := core.Conn.Exec(
		sql,
		uploadBlockMate.Upload.Id, uploadBlockMate.Offset,
		uploadBlockMate.Size, uploadBlockMate.IsComplete,
	)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func CreateUserFileInfo(userFileMate models.UserFile) (int64, error) {
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

func GetUploadInfoByUserIdFileId(userId int64, fileId int64) (models.Upload, error) {
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
	rows := core.Conn.QueryRow(sql, userId, fileId)
	err := rows.Scan(
		&uploadMate.Id, &uploadMate.UserInfo.Name, &uploadMate.UserInfo.User,
		&uploadMate.File.Hash, &uploadMate.File.Size, &uploadMate.File.Path,
		&uploadMate.File.IsComplete, &uploadMate.LocalPath, &uploadMate.BlockSize,
		&uploadMate.IsComplete,
	)
	return *uploadMate, err
}

func GetDirInfoByUserIdPathName(userId int64, path string, name string) (models.Directory, error) {
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

func GetFileInfoByHash(hash string) (models.File, error) {
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

func GetRootDirListByUserId(userId int64) ([]models.Directory, error) {
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

func GetRootFileListByUserId(userId int64) ([]models.UserFile, error) {
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

func GetDirListByUserIdPath(userId int64, path string) ([]models.Directory, error) {
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

func GetFileListByDirId(dirId int64) ([]models.UserFile, error) {
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

func GetUploadListByUserId(userId int64) ([]models.Upload, error) {
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
	rows, err := core.Conn.Query(sql, userId)
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

func GetUploadInfoById(id int64) (models.Upload, error) {
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
	rows := core.Conn.QueryRow(sql, id)
	err := rows.Scan(
		&uploadMate.Id, &uploadMate.UserInfo.Name, &uploadMate.UserInfo.User,
		&uploadMate.File.Hash, &uploadMate.File.Size, &uploadMate.File.Path,
		&uploadMate.File.IsComplete, &uploadMate.LocalPath, &uploadMate.BlockSize,
		&uploadMate.IsComplete,
	)
	return *uploadMate, err
}

func UpdateUploadComplete(complete int, id int64) error {
	sql := `
		UPDATE upload 
		SET is_complete = ?
		WHERE id = ?
	`
	_, err := core.Conn.Exec(sql, complete, id)
	return err
}

func GetUserFileInfoByUserIdFileId(userId int64, fileId int64) (models.Upload, error) {
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
	rows := core.Conn.QueryRow(sql, userId, fileId)
	err := rows.Scan(
		&uploadMate.Id, &uploadMate.UserInfo.Name, &uploadMate.UserInfo.User,
		&uploadMate.File.Hash, &uploadMate.File.Size, &uploadMate.File.Path,
		&uploadMate.File.IsComplete, &uploadMate.LocalPath, &uploadMate.BlockSize,
		&uploadMate.IsComplete,
	)
	return *uploadMate, err
}
