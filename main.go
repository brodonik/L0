package main

import (
	"fmt"
	"log"
	"main/cache"
	serverhttp "main/serverhttp"

	"main/storage"
	"main/subscriber"
	"net/http"

	"github.com/nats-io/stan.go"
)

func main() {
	db := storage.NewDb()
	defer db.Db.Close()

	sc, err := stan.Connect("test-cluster", "client-123", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming Server: %v", err)
	}
	defer sc.Close()
	ch := cache.NewOrdersCache(db)
	s := subscriber.NewSubscriber(&sc, db, ch)
	s.SubscribeToOrder()

	ch.RestoreCacheFromDB()

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		if idParam := r.URL.Query().Get("id"); idParam != "" {
			serverhttp.GetOrderHandler(ch).ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, "./ui/index.html")
	})

	fmt.Println("Starting HTTP server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
