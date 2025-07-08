package main

import (
	"fmt"
	"main/Database"
	"main/Routes"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	Database.Database()
	fmt.Printf("Starting server at port 8080\n")
	newRouter := Routes.PublicRoutes()
	http.ListenAndServe(":8080", newRouter)
	fmt.Printf("Server running on port 8080\n")
}
