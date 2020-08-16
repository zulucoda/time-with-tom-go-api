package main

import (
	"github.com/zulucoda/time-with-tom-go-api/booking"
	"github.com/zulucoda/time-with-tom-go-api/database"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
    database.SetupDatabase()
	booking.SetupRoutes(apiBasePath)
	// port, ServeMux nil is the default
	http.ListenAndServe(":5000", nil)
}
