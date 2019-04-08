package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/Betchika99/tp_db_project/database"
	"github.com/Betchika99/tp_db_project/models"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"log"
)

const (
	forumNotFound = "Can't find forum with slug: "
)

func ForumCreate(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/forum/create")

	forum := models.Forum{}

	err := json.Unmarshal(ctx.PostBody(), &forum)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	user, err := database.SelectOneUser(forum.User)
	if err != nil {
		log.Println("ERROR is", err.Error())

		message := userNotFound + user.Nickname

		jsonBody, err := json.Marshal(models.ModelError{message})
		if err != nil {
			log.Println(" json marshal fail")
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(jsonBody)
		return
	}

	data, err := database.CreateForum(forum)
	if err != nil {
		log.Println("ERROR is", err.Error())

		forumSelected, err := database.SelectForumBySlug(forum.Slug)
		if err != nil {
			log.Println("ERROR is", err.Error())
			return
		}

		jsonBody, err := json.Marshal(forumSelected)
		if err != nil {
			log.Println(" json marshal fail")
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusConflict)
		ctx.Response.SetBody(jsonBody)
		return
	}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		log.Println(" json marshal fail")
		return

	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Response.SetBody(jsonBody)
}

func ForumGetOne(ctx *fasthttp.RequestCtx) {
	log.Println("GET /api/forum/:slug/details")

	slug := ctx.UserValue("slug").(string)

	data, err := database.SelectForumBySlug(slug)

	if err != nil {
		log.Println("ERROR is", err.Error())

		if err == pgx.ErrNoRows {
			message := forumNotFound + slug

			jsonBody, err := json.Marshal(models.ModelError{message})
			if err != nil {
				log.Println(" json marshal fail")
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBody(jsonBody)
			return
		} else {
			return
		}
	}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		log.Println(" json marshal fail")
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(jsonBody)
}

func ForumGetThreads(ctx *fasthttp.RequestCtx) {
	log.Println("GET /api/forum/:slug/threads")

	slug := ctx.UserValue("slug").(string)
	_, err := database.SelectForumBySlug(slug)
	if err != nil {
		log.Println("ERROR is", err.Error())

		if err == pgx.ErrNoRows {
			message := forumNotFound + slug

			jsonBody, err := json.Marshal(models.ModelError{message})
			if err != nil {
				log.Println(" json marshal fail")
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBody(jsonBody)
			return
		} else {
			return
		}
	}

	limit := string(ctx.FormValue("limit"))
	desc := ctx.FormValue("desc")
	since := string(ctx.FormValue("since"))
	sinceString := "AND created_at >= $2"

	sortString := " ORDER BY created_at ASC"
	if bytes.Equal([]byte("true"), desc) {
		sortString = " ORDER BY created_at DESC"
		sinceString = "AND created_at <= $2"
	}

	if since == "" {
		sinceString = ""
	}

	limitString := ""
	if limit != "" {
		limitString = " LIMIT " + limit
	}

	data, err := database.SelectThreadsBySlug(slug, limitString, sortString, sinceString, since)

	if err != nil {
		log.Println("ERROR is", err.Error())

		if err == pgx.ErrNoRows {
			message := forumNotFound + slug

			jsonBody, err := json.Marshal(models.ModelError{message})
			if err != nil {
				log.Println(" json marshal fail")
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBody(jsonBody)
			return
		} else {
			return
		}
	}
	if data == nil {
		log.Println("There are not any threads in this forum")
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBody([]byte{91, 93})
			return
	}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		log.Println(" json marshal fail")
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(jsonBody)
}

func ForumGetUsers(ctx *fasthttp.RequestCtx) {
	log.Println("GET /api/forum/:slug/threads")

	slug := ctx.UserValue("slug").(string)
	_, err := database.SelectForumBySlug(slug)
	if err != nil {
		log.Println("ERROR is", err.Error())

		if err == pgx.ErrNoRows {
			message := forumNotFound + slug

			jsonBody, err := json.Marshal(models.ModelError{message})
			if err != nil {
				log.Println(" json marshal fail")
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBody(jsonBody)
			return
		} else {
			return
		}
	}

	limit := string(ctx.FormValue("limit"))
	desc := ctx.FormValue("desc")
	since := string(ctx.FormValue("since"))
	sinceString := "AND u.nickname > $2"

	sortString := " ORDER BY u.nickname ASC"
	if bytes.Equal([]byte("true"), desc) {
		sortString = " ORDER BY u.nickname DESC"
		sinceString = "AND u.nickname < $2"
	}

	if since == "" {
		sinceString = ""
	}

	limitString := ""
	if limit != "" {
		limitString = " LIMIT " + limit
	}

	data, err := database.SelectForumUsers(slug, limitString, sortString, sinceString, since)

	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}
	if data == nil {
		log.Println("There are not any threads in this forum")
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody([]byte{91, 93})
		return
	}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		log.Println(" json marshal fail")
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(jsonBody)
}