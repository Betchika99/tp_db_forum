package controllers

import (
	"encoding/json"
	"github.com/Betchika99/tp_db_project/database"
	"github.com/Betchika99/tp_db_project/models"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
)

const (
	parentNotFound = "Parent post was created in another thread"
	postNotFound = "Can't find post with id: "
)

func PostsCreate(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/thread/:slug_or_id/create")

	slugOrId := ctx.UserValue("slug_or_id").(string)

	var posts models.Posts
	err := json.Unmarshal(ctx.PostBody(), &posts)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	data, err := database.InsertPosts(posts, slugOrId)
	if err != nil {
		log.Println("ERROR is", err.Error())
		if err == pgx.ErrNoRows {
			message := parentNotFound

			jsonBody, err := json.Marshal(models.ModelError{message})
			if err != nil {
				log.Println(" json marshal fail")
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusConflict)
			ctx.Response.SetBody(jsonBody)
			return
		} else {
			jsonBody, err := json.Marshal(models.ModelError{err.Error()})
			if err != nil {
				log.Println(" json marshal fail")
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Response.SetBody(jsonBody)
			return
		}
	}

	if data == nil {
		log.Println("Empty posts found")
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusCreated)
		ctx.SetBody([]byte{91, 93})
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

func PostGetOne(ctx *fasthttp.RequestCtx) {
	log.Println("GET /api/post/:id/details")

	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	related := ctx.FormValue("related")
	relatedParams := []string{"post"}
	if len(related) != 0 {
		relatedParams = append(relatedParams, strings.Split(string(related), ",")...)
	}

	var postResult models.PostsRelated
	for _, param := range relatedParams {
		switch param {
		case "post":
			post, err := database.SelectPostByID(id)
			if err != nil {
				log.Println("ERROR is", err.Error())

				if err == pgx.ErrNoRows {
					message := postNotFound + strconv.FormatInt(int64(id), 10)

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
			postResult.PostModel = &post
		case "user":
			author, err := database.SelectOneUser(postResult.PostModel.Author)
			if err != nil {
				log.Println("ERROR is", err.Error())
				return
			}
			postResult.AuthorModel = &author
		case "thread":
			thread, err := database.SelectThreadById(postResult.PostModel.Thread)
			if err != nil {
				log.Println("ERROR is", err.Error())
				return
			}
			postResult.ThreadModel = &thread
		case "forum":
			forum, err := database.SelectForumBySlug(postResult.PostModel.Forum)
			if err != nil {
				log.Println("ERROR is", err.Error())
				return
			}
			postResult.ForumModel = &forum
		}
	}

	jsonBody, err := json.Marshal(postResult)
	if err != nil {
		log.Println(" json marshal fail")
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(jsonBody)
}

func PostUpdate(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/post/:id/details")

	postID, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	existedPost, err := database.SelectPostByID(postID)
	if err != nil {
		log.Println("ERROR is", err.Error())
		if err == pgx.ErrNoRows {
			message := postNotFound + strconv.FormatInt(int64(postID), 10)

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
		return
	}

	post := models.Post{}
	post.Id = postID

	err = json.Unmarshal(ctx.PostBody(), &post)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	data := models.Post{}
	if post.Message == existedPost.Message {
		data = existedPost
	} else {
		data, err = database.UpdatePost(post)
		if err != nil {
			log.Println("ERROR is", err.Error())
			if err != pgx.ErrNoRows {
				data = existedPost
			} else {
				return
			}
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
