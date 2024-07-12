package main

import (
	elastic "ex01/db"
	"ex01/server"
)

func main() {
	// client := elastic.NewClient()
	server := server.NewServer(elastic.NewClient())

	server.Run()
}
