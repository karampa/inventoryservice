package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/karampa/inventoryservice/database"
	"github.com/karampa/inventoryservice/product"
)

const apiBasePath = "/api"

func main() {
	database.SetupDBConnection()
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
