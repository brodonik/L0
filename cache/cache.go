package cache

import (
	"database/sql"
	"fmt"
	"log"
	"main/model"
)

var OrdersCache = make(map[int]model.Order)

func CacheOrder(order model.Order, db *sql.DB) {
	var orderId int
	err := db.QueryRow("SELECT id FROM delivery WHERE customer_id = $1", order.Customer_id).Scan(&orderId)
	if err != nil {
		log.Printf("Error fetching order id from database: %v", err)
		return
	}

	OrdersCache[orderId] = order
	fmt.Println("Order cached")
}

func RestoreCacheFromDB(db *sql.DB) {
	rows, err := db.Query("SELECT locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM delivery")
	if err != nil {
		log.Fatalf("Error restoring cache from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service, &order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard); err != nil {
			log.Printf("Error scanning order: %v", err)
			continue
		}
		CacheOrder(order, db)
	}
}
