package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetOpenIdByUserId(ctx context.Context, userId interface{}) (string, error)
	GetById() error
	GetByPhoneNumber(string) error
	CheckIfUserExistsByPhoneNumber(string) bool
	GetIdByPhoneNumber(string) UserID
	CreateUser(CreateUser) (int64, error)
}

type user struct {
	db *sqlx.DB
}

func newUser(dbc *sqlx.DB) *user {
	return &user{dbc}
}

func (u *user) GetById() error {
	return nil
}

func (u *user) GetByPhoneNumber(phoneNumber string) error {
	fmt.Println("Get user by phone number.")
	// getting logic here...
	return nil
}

func (u *user) CheckIfUserExistsByPhoneNumber(phoneNumber string) bool {
	query := "SELECT 1 AS `exist` FROM `members` WHERE `phone_number` = ?"
	var exist int8
	u.db.Get(&exist, query, phoneNumber)
	return exist == 1
}

func (u *user) GetIdByPhoneNumber(phoneNumber string) UserID {
	var userId UserID
	query := "SELECT id FROM `members` WHERE `phone_number` = ?"
	u.db.Get(&userId, query, phoneNumber)
	return userId
}

func (u *user) CreateUser(newUser CreateUser) (int64, error) {
	query := "INSERT INTO `members`(`phone_number`, `username`) VALUES(?, ?)"
	result, err := u.db.Exec(query, newUser.PhoneNumber, newUser.Username)
	if err != nil {
		return 0, errors.New("创建数据时发生错误")
	}
	return result.LastInsertId()
}

func (u *user) GetOpenIdByUserId(ctx context.Context, userId interface{}) (openId string, err error) {
	query := "SELECT `open_id` AS `openId` FROM `member_binds` WHERE `user_id` = ?"
	return openId, u.db.GetContext(ctx, &openId, query, userId)
}
