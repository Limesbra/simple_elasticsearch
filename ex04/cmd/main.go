package main

import (
	elastic "ex04/db"
	"ex04/server"
)

func main() {
	server := server.NewServer(elastic.NewClient())

	server.Run()
}
