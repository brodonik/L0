package main

import (
	"fmt"
	"log"
	"main/cache"
	"main/consumer"
	serverhttp "main/serverhttp"
	"main/storage"
	"net/http"

	"github.com/IBM/sarama"
)

func main() {
	db := storage.NewDb()
	defer db.Db.Close()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	saramaConsumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer saramaConsumer.Close()

	ch := cache.NewOrdersCache(db)
	ch.RestoreCacheFromDB()

	c := consumer.NewConsumer(saramaConsumer, db, ch)
	go c.SubscribeToOrder()

	http.HandleFunc("/order/", serverhttp.GetOrderHandler(ch))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
	})

	fmt.Println("Starting HTTP server on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
