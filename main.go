package main

import (
	"fmt"
	mongodbDriver "go-practise/driver"
	"go-practise/handler"
	"net/http"
)

func main() {
	mongodbDriver.ConnectMongoDb("admin", "password")

	http.HandleFunc("/api/check-health", handler.CheckHealthHandler)
	http.HandleFunc("/api/user/register", handler.Register)
	http.HandleFunc("/api/user/get", handler.GetUser)
	http.HandleFunc("/api/user/login", handler.Login)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
