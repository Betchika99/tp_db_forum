package controllers

import (
	"encoding/json"
	"github.com/Betchika99/tp_db_project/database"
	"github.com/valyala/fasthttp"
	"log"
)

func Status(ctx *fasthttp.RequestCtx) {
	log.Println("GET /api/service/status")

	data, err := database.SelectCounts()

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

func Clear(ctx *fasthttp.RequestCtx) {
	log.Println("POST /api/service/clear")

	err := database.DeleteAll()
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
