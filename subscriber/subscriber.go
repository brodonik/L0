package subscriber

import (
	"database/sql"
	"encoding/json"
	"log"
	"main/cache"
	"main/model"
	"main/storage"
	"time"

	"github.com/nats-io/stan.go"
)

func SubscribeToOrder(sc stan.Conn, db *sql.DB) {

	_, err := sc.Subscribe("order", func(m *stan.Msg) {

		var order model.Order
		if err := json.Unmarshal(m.Data, &order); err != nil {
			log.Printf("Error unmarshalling order: %v", err)
			return
		}

		if order.Customer_id == "" || order.Sm_id == 0 {
			log.Printf("Received order with invalid critical data: %+v\n", order)
			return
		}

		if order.Locale == "" {
			order.Locale = "default_locale"
		}
		if order.Internal_signature == "" {
			order.Internal_signature = "default_signature"
		}
		if order.Delivery_service == "" {
			order.Delivery_service = "default_service"
		}
		if order.Shardkey == "" {
			order.Shardkey = "default_shardkey"
		}
		if order.Date_created == "" {
			order.Date_created = time.Now().Format(time.RFC3339)
		}
		if order.Oof_shard == "" {
			order.Oof_shard = "default_oof_shard"
		}

		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM delivery WHERE customer_id = $1 AND sm_id = $2)", order.Customer_id, order.Sm_id).Scan(&exists)
		if err != nil {
			log.Printf("Error checking if order exists: %v", err)
			return
		}
		if exists {
			log.Printf("Order already exists: %+v\n", order)
			return
		}

		storage.SaveOrderToDB(db, order)
		cache.CacheOrder(order, db)
		m.Ack()
	}, stan.DeliverAllAvailable(), stan.DurableName("order-subscription"), stan.SetManualAckMode(), stan.AckWait(time.Second*30))

	if err != nil {
		log.Fatalf("Error subscribing to order: %v", err)
	}
}
