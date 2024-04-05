package cache

import (
	"fmt"
	"log"
	"main/model"
	"main/storage"
)

type OrdersCache struct {
	ch map[int]model.Order
	db *storage.Db
}

func NewOrdersCache(db *storage.Db) *OrdersCache {
	return &OrdersCache{ch: make(map[int]model.Order), db: db}
}

func (o *OrdersCache) CacheOrder(order model.Order) {
	var orderId int
	err := o.db.Db.QueryRow("SELECT id FROM delivery WHERE customer_id = $1", order.Customer_id).Scan(&orderId)
	if err != nil {
		log.Printf("Error fetching order id from database: %v", err)
		return
	}

	o.ch[orderId] = order

	fmt.Println("Order cached")
}

func (r *OrdersCache) RestoreCacheFromDB() {
	rows, err := r.db.Db.Query("SELECT locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM delivery")
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
		r.CacheOrder(order)

	}
}

func (o *OrdersCache) GetOrderById(orderId int) (*model.Order, bool) {
	order, ok := o.ch[orderId]
	return &order, ok
}
