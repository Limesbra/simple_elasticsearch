package main

import (
	elastic "ex03/db"
	"ex03/server"
)

func main() {
	server := server.NewServer(elastic.NewClient())

	server.Run()
}
