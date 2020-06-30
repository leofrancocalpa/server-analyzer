package rest

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// Index is the index handler
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "init method")
	fmt.Println("Req Index")
}

func startServer() {
	fmt.Println("I'm the server")
	router := fasthttprouter.New()
	router.GET("/", Index)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
