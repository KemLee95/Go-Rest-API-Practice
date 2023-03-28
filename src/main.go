package main

import (
	"log"

	config "github.com/kemlee/go-rest-api-practise/config"
	server "github.com/kemlee/go-rest-api-practise/server"
)

func main() {
	if _, err := config.GetAPIConfig(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	server := server.NewServer()

	if err := server.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
