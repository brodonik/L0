package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"main/model"
	"main/storage"
)

type OrdersCache struct {
	ch map[string]model.Order
	db *storage.Db
}

func NewOrdersCache(db *storage.Db) *OrdersCache {
	return &OrdersCache{ch: make(map[string]model.Order), db: db}
}

func (o *OrdersCache) CacheOrder(order model.Order) {
	o.ch[order.Order_uid] = order
	fmt.Println("Order cached")
}

func (r *OrdersCache) RestoreCacheFromDB() {
	rows, err := r.db.Db.Query(`SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders`)
	if err != nil {
		log.Fatalf("Error restoring cache from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		var deliveryJSON, paymentJSON, itemsJSON string

		if err := rows.Scan(&order.Order_uid, &order.Track_number, &order.Entry, &deliveryJSON, &paymentJSON, &itemsJSON,
			&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
			&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard); err != nil {
			log.Printf("Error scanning order: %v", err)
			continue
		}

		json.Unmarshal([]byte(deliveryJSON), &order.Delivery)
		json.Unmarshal([]byte(paymentJSON), &order.Payment)
		json.Unmarshal([]byte(itemsJSON), &order.Items)

		r.CacheOrder(order)
	}
}

func (o *OrdersCache) GetOrderByUid(orderUid string) (*model.Order, bool) {
	order, ok := o.ch[orderUid]
	return &order, ok
}
