package main

import (
	"fmt"

	service "github.com/leofrancocalpa/server-analyzer/pkg/http/rest"
)

func main() {
	fmt.Println("******* STARTING SERVICE *******")
	service.StartServer()
}
