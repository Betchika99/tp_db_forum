package controllers

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

func Index(ctx *fasthttp.RequestCtx) {
	log.Println("GET /")
	fmt.Fprintf(ctx, "hello!!!!")
}
