package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Betchika99/tp_db_project/database"
	"github.com/Betchika99/tp_db_project/models"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"log"
)

const (
	userNotFound = "Can't find user by nickname: "
	emailRegistred = "This email is already registered by user: "
)
func UserCreate(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/user/:name/create")

	user := models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	err = database.CreateUser(user)
	if err != nil {
		log.Println("ERROR is", err.Error())

		users, err := database.SelectUsersByNickAndEmail(user.Nickname, user.Email)
		if err != nil {
			log.Println("ERROR is", err.Error())
			return
		}
		if len(users) == 0 {
			err := fmt.Errorf("There are not users")
			log.Println("ERROR is", err.Error())
			return
		}

		jsonBody, err := json.Marshal(users)
		if err != nil {
			log.Println(" json marshal fail")
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusConflict)
		ctx.Response.SetBody(jsonBody)
		return
	}

	jsonBody, err := json.Marshal(user)
	if err != nil {
		log.Println(" json marshal fail")
		return

	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Response.SetBody(jsonBody)
}

func UserGetOne(ctx *fasthttp.RequestCtx) {
	log.Println("GET /api/user/:name/profile")

	nickname := ctx.UserValue("nickname").(string)

	data, err := database.SelectOneUser(nickname)

	if err != nil {
		log.Println("ERROR is", err.Error())

		if err == pgx.ErrNoRows {
			message := userNotFound + nickname

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

func UserUpdate(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/user/:name/profile")

	user := models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	data, err := database.UpdateUser(user)
	if err != nil {
		log.Println("ERROR is", err.Error())

		if err == pgx.ErrNoRows {
			message := userNotFound + user.Nickname

			jsonBody, err := json.Marshal(models.ModelError{message})
			if err != nil {
				log.Println(" json marshal fail")
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBody(jsonBody)
			return
		}

		message := emailRegistred

		jsonBody, err := json.Marshal(models.ModelError{message})
		if err != nil {
			log.Println(" json marshal fail")
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusConflict)
		ctx.SetBody(jsonBody)
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