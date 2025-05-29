package main

import (
	"main/Database"
	"main/Routes"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	Database.Database()

	newRouter := Routes.PublicRoutes()
	http.ListenAndServe(":8080", newRouter)
}
