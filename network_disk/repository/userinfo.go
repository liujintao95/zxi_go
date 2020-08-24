package repository

import (
	"zxi_network_disk_go/network_disk/models"
	"zxi_network_disk_go/utils"
)

type UserManager struct {
	table string
}

func NewUserManager() *UserManager {
	return &UserManager{table: "user_info"}
}

func (u *UserManager) GetSqlByUser(user string) (models.UserInfo, error) {
	userMate := new(models.UserInfo)
	sql := `
		SELECT id, name, user, pwd, recycled
		FROM user_info 
		WHERE user = ?
		AND recycled = 'N'
	`
	rows := utils.Conn.QueryRow(sql, user)
	err := rows.Scan(
		userMate.Id, userMate.Name, userMate.User,
		userMate.Pwd, userMate.Recycled,
	)
	return *userMate, err
}

func (u *UserManager) Create(userMate models.UserInfo) (int64, error) {
	sql := `
		INSERT INTO user_info(
			name, user, pwd
		) 
		VALUES (?, ?, ?)
	`
	res, err := utils.Conn.Exec(
		sql,
		userMate.User, userMate.User, userMate.Pwd,
	)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (u *UserManager) Update(userMate models.UserInfo) error {
	sql := `
		UPDATE user_info 
		SET name = ?, pwd = ?
		WHERE user=?
	`
	_, err := utils.Conn.Exec(
		sql,
		userMate.Name, userMate.Pwd, userMate.User,
	)
	return err
}

func (u *UserManager) DelByUser(user string) error {
	sql := `
		UPDATE user_info 
		SET recycled = 'Y'
		WHERE user = ?
	`
	_, err := utils.Conn.Exec(sql, user)
	return err
}
