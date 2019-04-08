package database

import (
	"fmt"
	"github.com/Betchika99/tp_db_project/models"
	"github.com/jackc/pgx"
	"strconv"
)

const (
	sqlInsertPost = `INSERT INTO posts (author_id, forum_id, message, parent_id, thread_id)
					  VALUES ($1, $2, $3, $4, $5)
					  RETURNING id, author_id, forum_id, created_at, message, parent_id, thread_id`

	sqlSelectPostByID = `SELECT id, author_id, forum_id, created_at, message, is_edited, thread_id
						 FROM posts WHERE id = $1`

	sqlUpdatePost = `UPDATE posts
					 SET message = COALESCE(NULLIF($1, ''), message),
					     is_edited = true
					 WHERE id = $2
				     RETURNING id, author_id, forum_id, created_at, message, is_edited, thread_id`
)

func InsertPosts(posts models.Posts, slugOrId string) (models.Posts, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	var thread int
	var forum string
	if thread, err = strconv.Atoi(slugOrId); err != nil {
		threadModel, err := SelectThreadBySlug(slugOrId)
		if err != nil {
			message := "Can't find post thread by slug: " + slugOrId
			return nil, fmt.Errorf(message)
		}
		thread = threadModel.Id
		forum = threadModel.Forum
	} else {
		threadModel, err := SelectThreadById(thread)
		if err != nil {
			message := "Can't find post thread by id: " + strconv.FormatInt(int64(thread), 10)
			return nil, fmt.Errorf(message)
		}
		forum = threadModel.Forum
	}

	if forum == "" {
		return nil, fmt.Errorf("There are not any forum to the thread with this slug or id")
	}

	var postsGot models.Posts
	for _, post := range posts {
		if post.Parent != 0 {
			existedParent, err := SelectPostByID(post.Parent)
			if err != nil {
				return nil, err
			}
			if existedParent.Thread != thread {
				return nil, pgx.ErrNoRows
			}
		}

		if _, err = SelectOneUser(post.Author); err != nil {
			message := "Can't find post author by nickname: " + post.Author
			return nil, fmt.Errorf(message)
		}

		data := transaction.QueryRow(sqlInsertPost,
			post.Author,
			forum,
			post.Message,
			post.Parent,
			thread,
			)

		postGot := models.Post{}

		err = data.Scan(&postGot.Id,
			&postGot.Author,
			&postGot.Forum,
			&postGot.Created,
			&postGot.Message,
			&postGot.Parent,
			&postGot.Thread,
		)
		if err != nil {
			return nil, err
		}
		err = UpdateForumField(forum, "posts")
		postsGot = append(postsGot, postGot)
	}

	if err = transaction.Commit(); err != nil {
		return nil, err
	}

	return postsGot, nil
}

func SelectPostByID(postID int) (models.Post, error) {
	post := models.Post{}
	transaction, err := GetConnect().Begin()
	if err != nil {
		return post, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlSelectPostByID, postID)

	err = data.Scan(&post.Id,
		&post.Author,
		&post.Forum,
		&post.Created,
		&post.Message,
		&post.IsEdited,
		&post.Thread,
	)
	if err != nil {
		return post, err
	}

	if err = transaction.Commit(); err != nil {
		return post, err
	}

	return post, nil
}

func UpdatePost(post models.Post) (models.Post, error) {
	transaction, err := GetConnect().Begin()
	defer transaction.Rollback()
	if err != nil {
		return post, err
	}

	if len(post.Message) == 0 {
		return post, fmt.Errorf("Post message is empty")
	}

	data := transaction.QueryRow(sqlUpdatePost, post.Message, post.Id)

	postUpdated := models.Post{}
	err = data.Scan(&postUpdated.Id,
		&postUpdated.Author,
		&postUpdated.Forum,
		&postUpdated.Created,
		&postUpdated.Message,
		&postUpdated.IsEdited,
		&postUpdated.Thread,
	)
	if err != nil {
		return post, err
	}

	if err = transaction.Commit(); err != nil {
		return post, err
	}

	return postUpdated, nil
}