package consumer

import (
	"encoding/json"
	"log"
	"main/cache"
	"main/model"
	"main/storage"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	Consumer sarama.Consumer
	Db       *storage.Db
	Ch       *cache.OrdersCache
}

func NewConsumer(consumer sarama.Consumer, db *storage.Db, ch *cache.OrdersCache) *Consumer {
	return &Consumer{Consumer: consumer, Db: db, Ch: ch}
}

func (c *Consumer) SubscribeToOrder() {
	partitionConsumer, err := c.Consumer.ConsumePartition("orders", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case message := <-partitionConsumer.Messages():
			var order model.Order
			if err := json.Unmarshal(message.Value, &order); err != nil {
				log.Printf("Error unmarshalling order: %v", err)
				continue
			}

			if order.Order_uid == "" || order.Customer_id == "" {
				log.Printf("Received order with invalid critical data: %+v\n", order)
				continue
			}

			if order.Locale == "" {
				order.Locale = "en"
			}
			if order.Internal_signature == "" {
				order.Internal_signature = ""
			}
			if order.Delivery_service == "" {
				order.Delivery_service = "default_service"
			}
			if order.Shardkey == "" {
				order.Shardkey = "0"
			}
			if order.Date_created == "" {
				order.Date_created = time.Now().Format(time.RFC3339)
			}
			if order.Oof_shard == "" {
				order.Oof_shard = "0"
			}

			var exists bool
			err := c.Db.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)", order.Order_uid).Scan(&exists)
			if err != nil {
				log.Printf("Error checking if order exists: %v", err)
				continue
			}
			if exists {
				log.Printf("Order already exists: %+v\n", order)
				continue
			}

			c.Db.SaveOrderToDB(order)
			c.Ch.CacheOrder(order)

		case err := <-partitionConsumer.Errors():
			log.Printf("Error from consumer: %v", err)
		}
	}
}
