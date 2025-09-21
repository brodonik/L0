package serverhttp

import (
	"encoding/json"
	"main/cache"
	"net/http"
	"strings"
)

func GetOrderHandler(ch *cache.OrdersCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/order/")
		if path == "" {
			http.Error(w, "Missing order_uid parameter", http.StatusBadRequest)
			return
		}

		order, ok := ch.GetOrderByUid(path)
		if !ok {
			http.Error(w, "Order not found in cache", http.StatusNotFound)
			return
		}

		jsonData, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
