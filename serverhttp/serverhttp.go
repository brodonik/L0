package serverhttp

import (
	"encoding/json"
	"main/cache"
	"net/http"
	"strconv"
)

func GetOrderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}

		order, ok := cache.OrdersCache[id]
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
