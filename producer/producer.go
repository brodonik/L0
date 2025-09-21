package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/model"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	order := model.Order{
		Order_uid:    "b563feb7b2b84b6test",
		Track_number: "WBILMTESTTRACK",
		Entry:        "WBIL",
		Delivery: model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:   "b563feb7b2b84b6test",
			Request_id:    "",
			Currency:      "USD",
			Provider:      "wbpay",
			Amount:        1817,
			Payment_dt:    1637907727,
			Bank:          "alpha",
			Delivery_cost: 1500,
			Goods_total:   317,
			Custom_fee:    0,
		},
		Items: []model.Item{
			{
				Chrt_id:      9934930,
				Track_number: "WBILMTESTTRACK",
				Price:        453,
				Rid:          "ab4219087a764ae0btest",
				Name:         "Mascaras",
				Sale:         30,
				Size:         "0",
				Total_price:  317,
				Nm_id:        2389212,
				Brand:        "Vivienne Sabo",
				Status:       202,
			},
		},
		Locale:             "en",
		Internal_signature: "",
		Customer_id:        "test",
		Delivery_service:   "meest",
		Shardkey:           "9",
		Sm_id:              99,
		Date_created:       "2021-11-26T06:22:19Z",
		Oof_shard:          "1",
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Error marshalling order: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.StringEncoder(orderJSON),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	fmt.Printf("Order UID: %s\n", order.Order_uid)
}
