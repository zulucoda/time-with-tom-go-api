package main

import (
	"github.com/zulucoda/time-with-tom-go-api/booking"
	"net/http"
)

const apiBasePath = "/api"

func main() {

	booking.SetupRoutes(apiBasePath)
	// port, ServeMux nil is the default
	http.ListenAndServe(":5000", nil)
}
