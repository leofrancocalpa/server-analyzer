package rest

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/leofrancocalpa/server-analyzer/pkg/http/client"
	"github.com/leofrancocalpa/server-analyzer/pkg/repository"

	"github.com/valyala/fasthttp"
)

const (
	//APISSLLABS f
	APISSLLABS = "https://api.ssllabs.com/api/v3/analyze?host="

	PORT = ":8888"

	APIBASE     = "/api/v1"
	URL_SERVERS = APIBASE + "/servers"
	URL_ANALYZE = URL_SERVERS + "/analyze"
)

var repo repository.Repository

func index(ctx *fasthttp.RequestCtx) {

	fmt.Fprintln(ctx, "http://localhost:8888/services \nhttp://localhost:8888/services/analyze?host=<server name>")
}

func analyzeRequestHandler(ctx *fasthttp.RequestCtx) {
	host := ctx.QueryArgs().Peek("host")
	fmt.Println("analyzinfg host: ", string(host))

	var url string = APISSLLABS + string(host)
	resultServerInfo := client.DoRequest(url, string(host)) //Getting info from SSLLABS, whois

	analysis := analyzeServer(resultServerInfo, repo) //looking for changes on servers and preparing the response

	json.NewEncoder(ctx.Response.BodyWriter()).Encode(analysis)

}

func listServersRequestHandler(ctx *fasthttp.RequestCtx) {
	response := retrieveHistoryFromDB(repo)
	json.NewEncoder(ctx.Response.BodyWriter()).Encode(response)
}

// StartServer init rest service
func StartServer() {

	repo = repository.NewRepository()
	router := fasthttprouter.New()

	router.GET("/", index)
	router.GET(URL_ANALYZE, analyzeRequestHandler)
	router.GET(URL_SERVERS, listServersRequestHandler)

	fmt.Println("****** SERVER UP *********")
	log.Fatal(fasthttp.ListenAndServe(PORT, Cors(router.Handler)))

}
