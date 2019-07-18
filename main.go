package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"github.com/ulule/limiter/drivers/store/memory"
)

func main() {

	// Define a limit rate to 4 requests per hour.
	rate, err := limiter.NewRateFromFormatted("4-H")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a store with the redis client.
	store := memory.NewStoreWithOptions(limiter.StoreOptions{
		Prefix:   "limiter_chi_example",
		MaxRetry: 3,
	})

	// Create a new middleware with the limiter instance.
	middleware := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))

	// Launch a simple chi server.
	router := chi.NewRouter()
	router.Use(middleware.Handler)
	router.Get("/", index)
	fmt.Println("Server is running on port 7777...")
	log.Fatal(http.ListenAndServe(":7777", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"message": "ok"}`))
}
