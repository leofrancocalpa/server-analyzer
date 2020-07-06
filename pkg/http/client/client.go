package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/likexian/whois-go"
	"github.com/valyala/fasthttp"
)

var grades = []string{"A+", "A", "B", "C", "D", "E", "F"}

// DoRequest GET http request to SSLLABS - WHOIS
func DoRequest(url string, host string) ResultServerInfo {

	reslt := getSSLabIfo(url)

	//adding whois info for each endpoint
	for i := 0; i < len(reslt.Enpoints); i++ {
		data := whoisRequest(reslt.Enpoints[i].IPAddress)
		if len(data) > 1 {
			reslt.Enpoints[i].Owner = data[0]
			reslt.Enpoints[i].Country = data[1]
		} else {
			reslt.Enpoints[i].Country = data[0]
		}
	}
	//getting the lowest SSL Grade
	reslt.SSLGrade = lowestSSLGrade(reslt.Enpoints)

	if strings.EqualFold(reslt.Status, "READY") {
		reslt.IsDown = false
	} else {
		reslt.IsDown = true
	}

	title, logo := getHTMLInfo(host)

	reslt.Title = title
	reslt.Logo = logo

	//fmt.Println(reslt)
	fmt.Println("==== fiish client requests ====")

	return reslt
}

//Get server info from API_SSLLAB
func getSSLabIfo(url string) ResultServerInfo {

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(req, resp)

	bodyBytes := resp.Body()

	var reslt ResultServerInfo

	err := json.Unmarshal(bodyBytes, &reslt)
	if err != nil {
		fmt.Println("[ERROR] Failure unmarshal")
	}

	return reslt
}

//Get title and favicon from host's website
func getHTMLInfo(hostname string) (string, string) {
	var url string
	if strings.Contains(hostname, "http") {
		url = hostname
	} else {
		url = "https://" + hostname
	}
	title, logo := GetTitleAndIcon(url)
	return title, logo
}

//Get the country and owner using whois command
func whoisRequest(ipaddress string) [2]string {
	fmt.Println("whois: ", ipaddress)
	owner := true
	country := true
	result, err := whois.Whois(ipaddress)
	if err != nil {
		fmt.Println("[ERROR] Failure getting whois request")
	}
	//fmt.Println(result)

	lines := strings.Split(result, "\n")
	var data [2]string
	for _, value := range lines {
		if owner && (strings.Contains(value, "OrgName:") || strings.Contains(value, "owner:") || strings.Contains(value, "org-name:")) {
			line := strings.Split(value, ": ")
			data[0] = strings.TrimSpace(line[1])
			owner = false
		} else if country && (strings.Contains(value, "Country:") || strings.Contains(value, "country")) {
			line := strings.Split(value, ": ")
			data[1] = strings.TrimSpace(line[1])
			country = false
		}
	}
	fmt.Println(data)
	return data
}

func lowestSSLGrade(endpoints []Enpoint) string {
	fmt.Println("Calculing lowest ssl grade")
	result := "A+"
	index := 0
	for _, value := range endpoints {
		i := indexOf(value.Grade, grades)
		if i > index {
			result = value.Grade
			index = i
		}
	}
	return result
}

func indexOf(element string, data []string) int {
	for i, value := range data {
		if strings.EqualFold(element, value) {
			return i
		}
	}
	return -1
}
