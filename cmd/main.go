package main

import (
	"fmt"
	service "github.com/leofrancocalpa/server-analyzer/pkg/rest"
)

func main() {
	fmt.Println("Hello dude")
	service.startServer()
}
