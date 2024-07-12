package main

import (
	elastic "ex02/db"
	"ex02/server"
)

func main() {
	// client := elastic.NewClient()
	server := server.NewServer(elastic.NewClient())

	server.Run()
}
