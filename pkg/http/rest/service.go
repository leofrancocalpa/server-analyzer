package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/leofrancocalpa/server-analyzer/pkg/http/client"
	"github.com/leofrancocalpa/server-analyzer/pkg/repository"
	"github.com/leofrancocalpa/server-analyzer/pkg/repository/model"

	"github.com/valyala/fasthttp"
)

const (
	//APISSLLABS f
	APISSLLABS = "https://api.ssllabs.com/api/v3/analyze?host="
)

var repo repository.Repository

func index(ctx *fasthttp.RequestCtx) {

	fmt.Fprintln(ctx, "http://localhost:8888/services \n http://localhost:8888/services/analyze?host=<server name>")
}

func analyzeRequestHandler(ctx *fasthttp.RequestCtx) {
	host := ctx.QueryArgs().Peek("host")
	fmt.Println("analyzinfg host: ", string(host))

	var url string = APISSLLABS + string(host)
	resultServerInfo := client.DoRequest(url, string(host))

	var analysis model.Analysis
	analysis.ServerChanged = false
	analysis.PreviousSSLGrade = resultServerInfo.SSLGrade
	analysis.IsDown = resultServerInfo.IsDown
	analysis.Title = resultServerInfo.Title
	analysis.Logo = resultServerInfo.Logo
	analysis.SSLGradre = resultServerInfo.SSLGrade
	if strings.EqualFold("READY", resultServerInfo.Status) {
		analysis.IsDown = false
	} else {
		analysis.IsDown = true
	}
	var servers []model.Server

	for _, value := range resultServerInfo.Enpoints {
		var serv model.Server
		serv.Address = value.IPAddress
		serv.Country = value.Country
		serv.Owner = value.Owner
		serv.SSLGrade = value.Grade

		servers = append(servers, serv)
	}

	for i := 0; i < len(resultServerInfo.Enpoints); i++ {
		var serv model.Server
		serv.Address = resultServerInfo.Enpoints[i].IPAddress
		serv.Country = resultServerInfo.Enpoints[i].Country
		serv.Owner = resultServerInfo.Enpoints[i].Owner
		serv.SSLGrade = resultServerInfo.Enpoints[i].Grade

		servers = append(servers, serv)
	}

	analysis.Servers = servers
	fmt.Println(">>> ", servers)
	json.NewEncoder(ctx.Response.BodyWriter()).Encode(analysis)

	jsonData, _ := json.Marshal(analysis)
	dataToPersist := model.ServerInfo{
		DNS:  resultServerInfo.Host,
		Data: jsonData,
	}
	fmt.Println("### >", string(jsonData))
	repo.CreateServerInfo(dataToPersist)

}

func listServersRequestHandler(ctx *fasthttp.RequestCtx) {
	serversQueried, _ := repo.FetchServersInfo()
	var response model.ItemsToResponse
	var items []model.ServersQueried
	for _, value := range serversQueried.Items {
		var analysis model.Analysis
		err := json.Unmarshal(value.Data, &analysis)
		if err != nil {
			panic("[ERROR] Failure response servers queried")
		}
		var servQueried model.ServersQueried
		servQueried.DNS = value.DNS
		servQueried.Data = analysis
		items = append(items, servQueried)
	}
	response.Items = items
	json.NewEncoder(ctx.Response.BodyWriter()).Encode(response)
}

// StartServer init rest service
func StartServer() {

	repo = repository.NewRepository()
	router := fasthttprouter.New()

	router.GET("/", index)
	router.GET("/servers/analyze", analyzeRequestHandler)
	router.GET("/servers", listServersRequestHandler)

	fmt.Println("****** SERVER UP *********")
	log.Fatal(fasthttp.ListenAndServe(":8888", Cors(router.Handler)))

}
