package database

import (
	"github.com/Betchika99/tp_db_project/models"
	"github.com/jackc/pgx"
	"log"
	"strconv"
	"strings"
)

const (
	sqlInsertThread = `
	INSERT INTO threads (
			slug, forum_id, author_id, created_at, message, title ) 
		VALUES ( $1, $2, $3, $4, $5, $6)
		RETURNING id, slug, forum_id, author_id, created_at, message, title `

	sqlSelectThreadBySlug = `SELECT forum_id, slug, author_id, title, message, created_at, id, votes
							 FROM threads WHERE slug = $1`

	sqlSelectThreadById = `SELECT forum_id, slug, author_id, title, message, created_at, id, votes
							 FROM threads WHERE id = $1`

	sqlUpdateThreadVote = `UPDATE threads
				       	   SET votes = $1
					       WHERE id = $2
					       RETURNING *`

	sqlUpdateThread = `UPDATE threads
					   SET title = COALESCE(NULLIF($1, ''), title),
						   message = COALESCE(NULLIF($2, ''), message)
					   WHERE id = $3
					   RETURNING *`

	sqlSelectThreadPosts = `SELECT id, author_id, parent_id, message, forum_id, thread_id, created_at, is_edited
							FROM posts WHERE thread_id = $1 `

	sqlSelectThreadPostsTreeBegin = `WITH RECURSIVE tree (id, path) AS (
							(SELECT id, array[id]
							FROM posts 
							WHERE parent_id = 0 AND thread_id = $1
							ORDER BY created_at, id `
	sqlSelectThreadPostsTreeEnd = `)
							UNION ALL
							SELECT p.id, array_append(path, p.id) FROM posts AS p
							JOIN tree ON tree.id = p.parent_id
							)
							SELECT p.id, p.author_id, p.parent_id, p.message, p.forum_id, p.thread_id, p.created_at, p.is_edited
							FROM tree
							   JOIN posts AS p ON tree.id = p.id
							ORDER BY tree.path`

	sqlSelectThreadPostsTreeDescBegin = `
							WITH RECURSIVE tree (id, path) AS (
							(SELECT id, array[-id] AS path
							FROM posts 
							WHERE parent_id = 0 AND thread_id = $1 
							ORDER BY path `
	sqlSelectThreadPostsTreeDescEnd = ` )
							UNION ALL
							SELECT p.id, array_append(path, p.id) FROM posts AS p
							JOIN tree ON tree.id = p.parent_id
							)
							SELECT p.id, p.author_id, p.parent_id, p.message, p.forum_id, p.thread_id, p.created_at, p.is_edited
							FROM tree 
						        JOIN posts AS p ON tree.id = p.id
							ORDER BY tree.path`

	sqlSelectThreadPostsTreeSince = `
							SELECT id, author_id, parent_id, message, forum_id, thread_id, created_at, is_edited
							FROM posts_tree WHERE num > (SELECT num FROM posts_tree WHERE id = $1) ORDER BY num
							`
	sqlSelectThreadPostsTreeSinceDesc = `
							SELECT id, author_id, parent_id, message, forum_id, thread_id, created_at, is_edited
							FROM posts_tree WHERE num < (SELECT num FROM posts_tree WHERE id = $1) ORDER BY num DESC
							`

	sqlSelectThreadPostsParentTreeSinceDesc = `
							SELECT id, author_id, parent_id, message, forum_id, thread_id, created_at, is_edited
							FROM posts_tree WHERE pnum < (SELECT pnum FROM posts_tree WHERE id = $1) AND pnum >= (SELECT pnum - $2 FROM posts_tree WHERE id = $1) ORDER BY num
							`

	sqlMakeTempTable = `CREATE TEMPORARY TABLE IF NOT EXISTS posts_tree (id BIGINT, author_id CITEXT, parent_id INT, message TEXT, forum_id CITEXT, thread_id INT, created_at TIMESTAMPTZ, is_edited BOOLEAN, num BIGINT, pnum BIGINT)`

	sqlInsertTempTable = `
							INSERT INTO posts_tree 
							WITH RECURSIVE tree (id, path) AS (
							(SELECT id, array[id]
							FROM posts 
							WHERE parent_id = 0 AND thread_id = $1
							ORDER BY created_at, id)
							UNION ALL
							SELECT p.id, array_append(path, p.id) FROM posts AS p
							JOIN tree ON tree.id = p.parent_id
							)
							SELECT p.id, p.author_id, p.parent_id, p.message, p.forum_id, p.thread_id, p.created_at, p.is_edited, ROW_NUMBER() OVER (ORDER BY path) as num
							FROM tree
							   JOIN posts AS p ON tree.id = p.id
							ORDER BY tree.path
							`

	sqlInsertParentTempTable = `
							INSERT INTO posts_tree 
							WITH RECURSIVE tree (id, path, pnum) AS (
							(SELECT id, array[-id] AS path, ROW_NUMBER() OVER (ORDER BY id) as pnum
							FROM posts 
							WHERE parent_id = 0 AND thread_id = $1 
							ORDER BY path DESC)
							UNION ALL
							SELECT p.id, array_append(path, p.id), pnum FROM posts AS p
							JOIN tree ON tree.id = p.parent_id
							)
							SELECT p.id, p.author_id, p.parent_id, p.message, p.forum_id, p.thread_id, p.created_at, p.is_edited, ROW_NUMBER() OVER (ORDER BY path) as num, tree.pnum
							FROM tree 
						        JOIN posts AS p ON tree.id = p.id
							ORDER BY tree.path
	`

	sqlTruncateTempTable = `TRUNCATE TABLE posts_tree`

	//sqlSelectFromTree = `SELECT * FROM posts WHERE id > $2 AND id IN (`
)

func InsertThread(thread models.Thread) (models.Thread, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return thread, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlInsertThread,
		thread.Slug,
		thread.Forum,
		thread.Author,
		thread.Created,
		thread.Message,
		thread.Title)

	threadGot := models.Thread{}

	err = data.Scan(&threadGot.Id,
		&threadGot.Slug,
		&threadGot.Forum,
		&threadGot.Author,
		&threadGot.Created,
		&threadGot.Message,
		&threadGot.Title,
	)

	if err != nil {
		return thread, err
	}

	if err = transaction.Commit(); err != nil {
		return thread, err
	}

	return threadGot, nil
}

func SelectThreadBySlug(slug string) (models.Thread, error) {
	thread := models.Thread{}
	transaction, err := GetConnect().Begin()
	if err != nil {
		return thread, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlSelectThreadBySlug, slug)

	err = data.Scan(&thread.Forum,
		&thread.Slug,
		&thread.Author,
		&thread.Title,
		&thread.Message,
		&thread.Created,
		&thread.Id,
		&thread.Votes,
	)
	if err != nil {
		return thread, err
	}

	if err = transaction.Commit(); err != nil {
		return thread, err
	}

	return thread, nil
}

func SelectThreadById(id int) (models.Thread, error) {
	thread := models.Thread{}
	transaction, err := GetConnect().Begin()
	if err != nil {
		return thread, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlSelectThreadById, id)

	err = data.Scan(&thread.Forum,
		&thread.Slug,
		&thread.Author,
		&thread.Title,
		&thread.Message,
		&thread.Created,
		&thread.Id,
		&thread.Votes,
	)
	if err != nil {
		return thread, err
	}

	if err = transaction.Commit(); err != nil {
		return thread, err
	}

	return thread, nil
}

func UpdateThread(thread models.Thread) (models.Thread, error) {
	transaction, err := GetConnect().Begin()
	defer transaction.Rollback()
	if err != nil {
		return thread, err
	}

	data := transaction.QueryRow(sqlUpdateThread, thread.Title, thread.Message, thread.Id)

	threadUpdated := models.Thread{}

	err = data.Scan(&threadUpdated.Id,
		&threadUpdated.Slug,
		&threadUpdated.Title,
		&threadUpdated.Message,
		&threadUpdated.Votes,
		&threadUpdated.Created,
		&threadUpdated.Forum,
		&threadUpdated.Author,
	)
	if err != nil {
		return thread, err
	}

	if err = transaction.Commit(); err != nil {
		return thread, err
	}

	return threadUpdated, nil
}

func UpdateThreadVote(thread models.Thread) (models.Thread, error) {
	transaction, err := GetConnect().Begin()
	defer transaction.Rollback()
	if err != nil {
		return thread, err
	}

	data := transaction.QueryRow(sqlUpdateThreadVote, thread.Votes, thread.Id)

	threadUpdated := models.Thread{}

	err = data.Scan(&threadUpdated.Id,
		&threadUpdated.Slug,
		&threadUpdated.Title,
		&threadUpdated.Message,
		&threadUpdated.Votes,
		&threadUpdated.Created,
		&threadUpdated.Forum,
		&threadUpdated.Author,
	)
	if err != nil {
		return thread, err
	}

	if err = transaction.Commit(); err != nil {
		return thread, err
	}

	return threadUpdated, nil
}

func SelectThreadPosts(threadID int, limit, desc, sortFlag, since string) (models.Posts, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	sortString := " ORDER BY id"
	sinceString := "AND id > $2"

	var query string
	var data *pgx.Rows
	switch sortFlag {
	case "tree":
		sinceString = " WHERE p.id > $2 "
		if desc == "true" {
			sortString = " DESC"
			sinceString = ""
		} else {
			sortString = ""
		}
		if since == "" {
			sinceString = ""
			query = sqlSelectThreadPostsTreeBegin + sqlSelectThreadPostsTreeEnd + sortString + limit
		} else {
			if desc == "true" {
				query = sqlSelectThreadPostsTreeSinceDesc + limit
			} else {
				query = sqlSelectThreadPostsTreeSince + limit
			}
		}
	case "parent_tree":
		if desc == "true" {
			if since == "" {
				sortString = ""
				query = sqlSelectThreadPostsTreeDescBegin + limit + sqlSelectThreadPostsTreeDescEnd
			} else {
				query = sqlSelectThreadPostsParentTreeSinceDesc
			}
		} else {
			if since == "" {
				sortString = ""
				query = sqlSelectThreadPostsTreeBegin + limit + sqlSelectThreadPostsTreeEnd
			} else {
				query = sqlSelectThreadPostsTreeSince
			}
		}
	default:
		if desc == "true" {
			sortString = "ORDER BY id DESC"
			sinceString = "AND id < $2"
		}
		if since == "" {
			sinceString = ""
		}
		query = sqlSelectThreadPosts + sinceString + sortString + limit
	}
	log.Println(query)
	if since != "" {
		sinceInt, err := strconv.Atoi(since)
		if err != nil {
			return nil, err
		}
		if sortFlag == "tree" || sortFlag == "parent_tree" {
			_, err = transaction.Exec(sqlMakeTempTable)
			if err != nil {
				log.Println("make")
				return nil, err
			}
			_, err = transaction.Exec(sqlTruncateTempTable)
			if err != nil {
				log.Println("drop")
				return nil, err
			}

			if sortFlag == "parent_tree" && desc == "true" {
				_, err = transaction.Exec(sqlInsertParentTempTable, threadID)
			} else {
				_, err = transaction.Exec(sqlInsertTempTable, threadID)
			}
			if err != nil {
				log.Println("ins")
				return nil, err
			}
			if sortFlag == "parent_tree" && desc == "true" {
				limits := strings.Split(limit, " ")
				data, err = transaction.Query(query, sinceInt, limits[len(limits) - 1])
			} else {
				data, err = transaction.Query(query, sinceInt)
			}
			if err != nil {
				log.Println("select")
				return nil, err
			}
		} else {
			data, err = transaction.Query(query, threadID, sinceInt)
		}
	} else {
		data, err = transaction.Query(query, threadID)
	}
	//if (sinceString == "") {
	//	query = sqlSelectThreadsBySlug + sort + limit
	//	//log.Println(query)
	//	data, err = transaction.Query(query, threadID)
	//} else {
	//	sinceTime, err := time.Parse(time.RFC3339, since)
	//	if err != nil {
	//		return nil, err
	//	}
	//	query = sqlSelectThreadsBySlug + sinceString + sort + limit
	//	//log.Println(query)
	//	data, err = transaction.Query(query, threadID, sinceTime)
	//}
	if err != nil {
		return nil, err
	}

	var posts models.Posts
	for data.Next() {
		post := models.Post{}
		err := data.Scan(
			&post.Id,
			&post.Author,
			&post.Parent,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.IsEdited,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = transaction.Commit(); err != nil {
		return posts, err
	}

	return posts, nil

}
