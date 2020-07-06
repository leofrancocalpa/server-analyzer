package rest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/leofrancocalpa/server-analyzer/pkg/http/client"
	"github.com/leofrancocalpa/server-analyzer/pkg/repository"
	"github.com/leofrancocalpa/server-analyzer/pkg/repository/model"
)

func analyzeServer(resultServerInfo client.ResultServerInfo, repo repository.Repository) model.Analysis {
	var analysis model.Analysis
	analysis.ServerChanged = false                        //default
	analysis.PreviousSSLGrade = resultServerInfo.SSLGrade //default
	analysis.IsDown = resultServerInfo.IsDown
	analysis.Title = resultServerInfo.Title
	analysis.Logo = resultServerInfo.Logo
	analysis.SSLGradre = resultServerInfo.SSLGrade

	var servers []model.Server

	for _, value := range resultServerInfo.Enpoints {
		var serv model.Server
		serv.Address = value.IPAddress
		serv.Country = value.Country
		serv.Owner = value.Owner
		serv.SSLGrade = value.Grade

		servers = append(servers, serv)
	}

	analysis.Servers = servers

	existInDB := repo.ExistEntry(resultServerInfo.Host)

	if existInDB {
		lastUpdated, previousSSLGrade, serversP := repo.GetServerInfoFromDB(resultServerInfo.Host)

		nowLessAnHour := time.Now().Add(-1 * time.Hour)
		fmt.Println(nowLessAnHour)
		fmt.Println(lastUpdated)
		if lastUpdated.Before(nowLessAnHour) {
			jsonServers, _ := json.Marshal(analysis.Servers)

			var json1, json2 interface{}
			json.Unmarshal(serversP, &json1)
			json.Unmarshal(jsonServers, &json2)

			if !reflect.DeepEqual(json1, json2) {
				fmt.Println("EQUALS B*")
				analysis.ServerChanged = true

			}
			analysis.PreviousSSLGrade = previousSSLGrade

			jsonData, _ := json.Marshal(analysis)
			dataToUpdate := model.ServerInfo{
				DNS:  resultServerInfo.Host,
				Data: jsonData,
			}
			fmt.Println("UPDATING IN DB...")
			repo.UpdateEntry(dataToUpdate)
		}

	} else {
		jsonData, _ := json.Marshal(analysis)
		dataToPersist := model.ServerInfo{
			DNS:  resultServerInfo.Host,
			Data: jsonData,
		}

		repo.CreateServerInfo(dataToPersist)
	}

	return analysis
}

func retrieveHistoryFromDB(repo repository.Repository) model.ItemsToResponse {
	serversQueried, _ := repo.FetchServersInfo()
	var response model.ItemsToResponse
	var items []model.ServersQueried
	for _, value := range serversQueried.Items {
		var analysis model.Analysis
		err := json.Unmarshal(value.Data, &analysis)
		if err != nil {
			fmt.Println("[ERROR] retrieving History from DB")
		}
		var servQueried model.ServersQueried
		servQueried.DNS = value.DNS
		servQueried.Data = analysis
		items = append(items, servQueried)
	}
	response.Items = items

	return response
}
