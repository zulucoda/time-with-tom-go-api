package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LogStartAndEndTimeHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler; middleware start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("middleware finished; %s\n", time.Since(start))
	})
}
