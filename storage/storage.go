package storage

import (
	"database/sql"
	"log"
	"main/model"

	_ "github.com/lib/pq"
)

func ConnectToDB() *sql.DB {

	connStr := "user=postgres dbname=dbl0 password=root sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	return db
}

func SaveOrderToDB(db *sql.DB, order model.Order) {

	_, err := db.Exec("INSERT INTO delivery (locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		log.Printf("Error saving order to database: %v", err)
	} else {
		log.Println("Order successfully saved to database")
	}
}
