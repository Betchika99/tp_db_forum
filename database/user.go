package database

import (
	"fmt"
	"github.com/Betchika99/tp_db_project/models"
)

const (
	sqlInsertUser = `
	INSERT INTO users ("nickname", "fullname", "about", "email")
	VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

 	sqlSelectUsers = `SELECT nickname, fullname, email, about FROM users WHERE nickname = $1 OR email = $2`

	sqlSelectUser = `SELECT nickname, fullname, email, about FROM users WHERE nickname = $1`

	sqlUpdateUser = `
	UPDATE users
	SET fullname = COALESCE(NULLIF($1, ''), fullname),
		about    = COALESCE(NULLIF($2, ''), about),
		email    = COALESCE(NULLIF($3, ''), email)
	WHERE nickname = $4
	RETURNING fullname, about, email, nickname`
)

func CreateUser(user models.User) error {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	data, err := transaction.Exec(sqlInsertUser, user.Nickname, user.Fullname, user.About, user.Email)
	if data.RowsAffected() == 0 {
		return fmt.Errorf("User existed!")
	}
	if err != nil {
		return err
	}

	if err = transaction.Commit(); err != nil {
		return err
	}

	return nil
}

func SelectUsersByNickAndEmail(nick, email string) (models.Users, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	data, err := transaction.Query(sqlSelectUsers, nick, email)
	if err != nil {
		return nil, err
	}

	var users models.Users
	for data.Next() {
		user := models.User{}
		err := data.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func SelectOneUser(nick string) (models.User, error) {
	user := models.User{}
	transaction, err := GetConnect().Begin()
	if err != nil {
		return user, err
	}
	defer transaction.Rollback()

	if err != nil {
		return user, err
	}

	data := transaction.QueryRow(sqlSelectUser, nick)

	err = data.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
	if err != nil {
		return user, err
	}

	if err = transaction.Commit(); err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUser(user models.User) (models.User, error) {
	transaction, err := GetConnect().Begin()
	defer transaction.Rollback()

	if err != nil {
		return user, err
	}

	data := transaction.QueryRow(sqlUpdateUser, user.Fullname, user.About, user.Email, user.Nickname)

	userUpdated := models.User{}

	err = data.Scan(&userUpdated.Fullname, &userUpdated.About, &userUpdated.Email, &userUpdated.Nickname)
	if err != nil {
		return user, err
	}

	if err = transaction.Commit(); err != nil {
		return userUpdated, err
	}

	return userUpdated, nil
}