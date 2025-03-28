package main

import (
	"huinya/db"
	"huinya/handlers"
	"huinya/middleware"
	"net/http"
)

func main() {
	db.DatabaseConfiguration()

	http.HandleFunc("/auth/register", handlers.CreateUser)
	http.HandleFunc("/auth/login", handlers.LoginUser)
	http.Handle("/auth/profile", middleware.JWTMiddleware(http.HandlerFunc(handlers.UserProfile)))
	http.Handle("/api/v1/upload", middleware.JWTMiddleware(http.HandlerFunc(handlers.UploadPost)))
	http.ListenAndServe(":8080", nil)

}
