package main

import (
	db "github.com/Betchika99/tp_db_project/database"
	r "github.com/Betchika99/tp_db_project/router"
	"github.com/valyala/fasthttp"
	"log"
)

const port = ":5000"

//var addr = flag.String("addr", "127.0.0.1:5000", "TCP address to listen to for incoming connections")

//
//func handler(ctx *fasthttp.RequestCtx) {
//	ctx.WriteString("Hello, world!\n")
//	db.Connect(psqlURI)
//}

func main() {
	defer db.CloseConnect()

	err := db.OpenConnect()
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	err = db.LoadSchema()
	if err != nil {
		log.Println("ERROR is", err.Error())
		return
	}

	//db.Connect(psqlURI)

	router := r.InitRouter()
	log.Fatal(fasthttp.ListenAndServe(port, router.Handler))

	log.Printf("Server started")

	//flag.Parse()
	//s := fasthttp.Server{
	//	Handler: handler,
	//}
	//err := s.ListenAndServe(*addr)
	//if err != nil {
	//	log.Fatalf("error in ListenAndServe: %s", err)
	//}
}
