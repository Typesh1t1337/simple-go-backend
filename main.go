package main

import (
	"huinya/db"
	"huinya/handlers"
	"net/http"
)

func main() {
	db.DatabaseConfiguration()

	http.HandleFunc("/auth/register", handlers.CreateUser)
	http.HandleFunc("/auth/login", handlers.LoginUser)
	http.ListenAndServe(":8080", nil)

}
