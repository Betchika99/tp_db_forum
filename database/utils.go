package database

import (
	"github.com/Betchika99/tp_db_project/models"
)

func UpdateForumField(slug, field string) error {
	transaction, err := GetConnect().Begin()
	defer transaction.Rollback()
	if err != nil {
		return err
	}

	query := `UPDATE forums SET ` + field + ` = ` + field + ` + 1` + ` WHERE slug = $1`
	_, err = transaction.Exec(query, slug)
	if err != nil {
		return err
	}

	if err = transaction.Commit(); err != nil {
		return err
	}

	return nil
}

func SelectCounts() (models.Status, error) {
	status := models.Status{}
	transaction, err := GetConnect().Begin()
	if err != nil {
		return status, err
	}
	defer transaction.Rollback()


	tables := [...]string{"forums", "posts", "threads", "users"}

	for i, table := range tables {
		data := transaction.QueryRow(`SELECT COUNT (*) FROM ` + table)
		switch i {
		case 0:
			err = data.Scan(&status.Forum)
			if err != nil {
				return status, err
			}
		case 1:
			err = data.Scan(&status.Post)
			if err != nil {
				return status, err
			}
		case 2:
			err = data.Scan(&status.Thread)
			if err != nil {
				return status, err
			}
		case 3:
			err = data.Scan(&status.User)
			if err != nil {
				return status, err
			}
		}
	}

	if err = transaction.Commit(); err != nil {
		return status, err
	}

	return status, nil
}

func DeleteAll() error {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	_, err = transaction.Exec("TRUNCATE users CASCADE")
	if err != nil {
		return err
	}

	if err = transaction.Commit(); err != nil {
		return err
	}

	return nil
}