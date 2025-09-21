package storage

import (
	"database/sql"
	"encoding/json"
	"log"
	"main/model"

	_ "github.com/lib/pq"
)

type Db struct {
	Db *sql.DB
}

func NewDb() *Db {
	connStr := "user=postgres dbname=dbl0 password=root sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	return &Db{Db: db}
}

func (d *Db) SaveOrderToDB(order model.Order) {
	tx, err := d.Db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return
	}
	defer tx.Rollback()

	deliveryJSON, _ := json.Marshal(order.Delivery)
	paymentJSON, _ := json.Marshal(order.Payment)
	itemsJSON, _ := json.Marshal(order.Items)

	_, err = tx.Exec(`INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		order.Order_uid, order.Track_number, order.Entry, deliveryJSON, paymentJSON, itemsJSON,
		order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service,
		order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		log.Printf("Error saving order to database: %v", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}

	log.Println("Order successfully saved to database")
}
