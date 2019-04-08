package controllers

import (
	"encoding/json"
	"github.com/Betchika99/tp_db_project/database"
	"github.com/Betchika99/tp_db_project/models"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

const (
	authorNotFound = "Can't find thread author by nickname: "
	forumBySlugNotFound = "Can't find thread forum by slug: "
	threadBySlugNotFound = "Can't find thread by slug: "
	threadByIDNotFound = "Can't find thread with id: "
)

func ThreadCreate(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/forum/:slug/create")

	thread := models.Thread{}
	forumSlug := ctx.UserValue("slug").(string)

	err := json.Unmarshal(ctx.PostBody(), &thread)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	forum, err := database.SelectForumBySlug(forumSlug)
	if err != nil {
		log.Println("ERROR is", err.Error())

		message := forumBySlugNotFound + forumSlug

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
	thread.Forum = forum.Slug

	_, err = database.SelectOneUser(thread.Author)
	if err != nil {
		log.Println("ERROR is", err.Error())

		message := authorNotFound + thread.Author

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

	existThread := models.Thread{}
	if thread.Slug != "" {
		existThread, err = database.SelectThreadBySlug(thread.Slug)
		if err != nil && err != pgx.ErrNoRows {
			log.Println("ERROR is", err.Error())
			return
		}
	}

	if existThread.Id != 0 {
		jsonBody, err := json.Marshal(existThread)
		if err != nil {
			log.Println(" json marshal fail")
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusConflict)
		ctx.Response.SetBody(jsonBody)
		return
	}

	data, err := database.InsertThread(thread)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	err = database.UpdateForumField(thread.Forum, "threads")
	if err != nil {
		log.Println("ERROR is", err.Error())
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

func ThreadVote(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/thread/:slug_or_id/vote")

	slugOrId := ctx.UserValue("slug_or_id").(string)

	var vote models.Vote
	err := json.Unmarshal(ctx.PostBody(), &vote)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	var threadID int
	var threadModel models.Thread
	if threadID, err = strconv.Atoi(slugOrId); err != nil {
		threadModel, err = database.SelectThreadBySlug(slugOrId)
		if err != nil {
			if err == pgx.ErrNoRows {
				message := threadBySlugNotFound + slugOrId

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
		threadID = threadModel.Id
	} else {
		threadModel, err = database.SelectThreadById(threadID)
		if err != nil {
			if err == pgx.ErrNoRows {
				message := threadByIDNotFound + strconv.FormatInt(int64(threadID), 10)

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
		threadID = threadModel.Id
	}
	if threadID == 0 {
		log.Println("There are not any thread with this slug or id")
		return
	}

	_, err = database.SelectOneUser(vote.Nickname)
	if err != nil {
		if err == pgx.ErrNoRows {
			jsonBody, err := json.Marshal(models.ModelError{userNotFound + vote.Nickname})
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

	existedVote, err := database.SelectVote(vote, threadID)
	if err != nil {
		log.Println("ERROR is", err.Error())
		if err == pgx.ErrNoRows {
			_, err = database.InsertVote(vote, threadID)
			threadModel.Votes += vote.Voice
			if err != nil {
				log.Println("ERROR is", err.Error())
				return
			}
		} else {
			return
		}
	}

	if existedVote.Voice != vote.Voice {
		_, err = database.UpdateVote(vote, threadID)
		if err != nil {
			log.Println("ERROR is", err.Error())
			return
		}
		threadModel.Votes = threadModel.Votes - existedVote.Voice + vote.Voice
	}


	data, err := database.UpdateThreadVote(threadModel)
	if err != nil {
		log.Println("ERROR is", err.Error())
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

func ThreadGetOne(ctx *fasthttp.RequestCtx) {
	log.Println("GET /api/thread/:slug_or_id/details")

	slugOrId := ctx.UserValue("slug_or_id").(string)

	var threadModel models.Thread
	if threadID, err := strconv.Atoi(slugOrId); err != nil {
		threadModel, err = database.SelectThreadBySlug(slugOrId)
		if err != nil {
			if err == pgx.ErrNoRows {
				message := threadBySlugNotFound + slugOrId

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
	} else {
		threadModel, err = database.SelectThreadById(threadID)
		if err != nil {
			if err == pgx.ErrNoRows {
				message := threadByIDNotFound + strconv.FormatInt(int64(threadID), 10)

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
	}

	jsonBody, err := json.Marshal(threadModel)
	if err != nil {
		log.Println(" json marshal fail")
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(jsonBody)
}

// TODO: make this shit
func ThreadGetPosts(ctx *fasthttp.RequestCtx) {
	log.Println("GET /api/thread/:slug_or_id/posts")

	slugOrId := ctx.UserValue("slug_or_id").(string)

	var threadModel models.Thread
	if threadID, err := strconv.Atoi(slugOrId); err != nil {
		threadModel, err = database.SelectThreadBySlug(slugOrId)
		if err != nil {
			log.Println("ERROR is", err.Error())
			if err == pgx.ErrNoRows {
				message := threadBySlugNotFound + slugOrId

				jsonBody, err := json.Marshal(models.ModelError{message})
				if err != nil {
					log.Println("ERROR is", err.Error())
					return
				}

				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				ctx.Response.SetBody(jsonBody)
				return
			}
			return
		}
	} else {
		threadModel, err = database.SelectThreadById(threadID)
		if err != nil {
			log.Println("ERROR is", err.Error())
			if err == pgx.ErrNoRows {
				message := threadByIDNotFound + strconv.FormatInt(int64(threadID), 10)

				jsonBody, err := json.Marshal(models.ModelError{message})
				if err != nil {
					log.Println("ERROR is", err.Error())
					return
				}

				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				ctx.Response.SetBody(jsonBody)
				return
			}
			return
		}
	}

	limit := string(ctx.FormValue("limit"))
	sort := string(ctx.FormValue("sort"))
	since := string(ctx.FormValue("since"))
	desc := string(ctx.FormValue("desc"))
	//sinceString := "AND created_at >= $2"

	//
	//if since == "" {
	//	sinceString = ""
	//}
	//

	limitString := ""
	if limit != "" {
		limitString = " LIMIT " + limit
	}

	data, err := database.SelectThreadPosts(threadModel.Id, limitString, desc, sort, since)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	//data, err := database.SelectThreadPosts(threadModel.Id)
	//if err != nil {
	//	log.Println("ERROR is", err.Error())
	//	return
	//}

	if data == nil {
		log.Println("There are not any posts in this thread")
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

func ThreadUpdate(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/thread/:slug_or_id/details")

	thread := models.Thread{}
	slugOrId := ctx.UserValue("slug_or_id").(string)
	err := json.Unmarshal(ctx.PostBody(), &thread)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	if threadID, err := strconv.Atoi(slugOrId); err != nil {
		threadModel, err := database.SelectThreadBySlug(slugOrId)
		if err != nil {
			if err == pgx.ErrNoRows {
				message := threadBySlugNotFound + slugOrId

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
		thread.Id = threadModel.Id
	} else {
		_, err = database.SelectThreadById(threadID)
		if err != nil {
			if err == pgx.ErrNoRows {
				message := threadByIDNotFound + strconv.FormatInt(int64(threadID), 10)

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
		thread.Id = threadID
	}

	data, err := database.UpdateThread(thread)

	jsonBody, err := json.Marshal(data)
	if err != nil {
		log.Println(" json marshal fail")
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(jsonBody)
}
