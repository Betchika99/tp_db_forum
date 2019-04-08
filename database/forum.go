package database

import (
	"github.com/Betchika99/tp_db_project/models"
	"github.com/jackc/pgx"
	"time"
)

const (
	sqlInsertForum = `INSERT INTO forums (slug, title, user_id)
  					  VALUES ($1, $2, (
							SELECT nickname FROM users WHERE nickname = $3)
					  )
			          RETURNING slug, title, user_id`

	sqlSelectForumBySlug = `SELECT * FROM forums
							WHERE slug = $1`

	sqlSelectThreadsBySlug = `SELECT * FROM threads 
							  WHERE forum_id = $1 `

	sqlSelectForumUsers = `SELECT DISTINCT u.nickname, u.about, u.email, u.fullname 
						   FROM users AS u
						   LEFT JOIN posts AS p ON u.nickname = p.author_id
						   LEFT JOIN threads AS th ON u.nickname =  th.author_id
						   WHERE (p.forum_id = $1 OR th.forum_id = $1)`
)

func CreateForum(forum models.Forum) (models.Forum, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return forum, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlInsertForum, forum.Slug, forum.Title, forum.User)

	forumGot := models.Forum{}

	err = data.Scan(&forumGot.Slug, &forumGot.Title, &forumGot.User)
	if err != nil {
		return forum, err
	}

	if err = transaction.Commit(); err != nil {
		return forum, err
	}

	return forumGot, nil
}

func SelectForumBySlug(slug string) (models.Forum, error) {
	forum := models.Forum{}
	transaction, err := GetConnect().Begin()
	if err != nil {
		return forum, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlSelectForumBySlug, slug)

	err = data.Scan(&forum.Slug,
				    &forum.Posts,
				    &forum.Threads,
				    &forum.Title,
				    &forum.User,
				    )
	if err != nil {
		return forum, err
	}

	if err = transaction.Commit(); err != nil {
		return forum, err
	}

	return forum, nil
}

func SelectThreadsBySlug(slug, limit, sort, sinceString, since string) (models.Threads, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	var query string
	var data *pgx.Rows
	if (sinceString == "") {
		query = sqlSelectThreadsBySlug + sort + limit
		//log.Println(query)
		data, err = transaction.Query(query, slug)
	} else {
		sinceTime, err := time.Parse(time.RFC3339, since)
		if err != nil {
			return nil, err
		}
		query = sqlSelectThreadsBySlug + sinceString + sort + limit
		//log.Println(query)
		data, err = transaction.Query(query, slug, sinceTime)
	}
	if err != nil {
		return nil, err
	}

	var threads models.Threads
	for data.Next() {
		thread := models.Thread{}
		err := data.Scan(&thread.Id,
						 &thread.Slug,
						 &thread.Title,
						 &thread.Message,
						 &thread.Votes,
						 &thread.Created,
						 &thread.Forum,
						 &thread.Author,
						 )
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}

func SelectForumUsers(slug, limit, sort, sinceString, since string) (models.Users, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	var query string
	var data *pgx.Rows
	if sinceString == "" {
		query = sqlSelectForumUsers + sort + limit
		data, err = transaction.Query(query, slug)
	} else {
		query = sqlSelectForumUsers + sinceString + sort + limit
		data, err = transaction.Query(query, slug, since)
	}
	if err != nil {
		return nil, err
	}

	var users models.Users
	for data.Next() {
		user := models.User{}
		err := data.Scan(&user.Nickname,
			&user.About,
			&user.Email,
			&user.Fullname,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
