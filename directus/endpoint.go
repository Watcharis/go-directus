package directus

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type WebhookEndpoint interface {
	FetchDataFromDirectusEndpoint(ctx context.Context) http.HandlerFunc
	CallDirectusEndpoint(ctx context.Context) http.HandlerFunc
}

type Endpoint struct {
	Whs WebhookServices
}

func NewEndpoint(whs WebhookServices) *Endpoint {
	return &Endpoint{
		Whs: whs,
	}
}

func (e *Endpoint) FetchDataFromDirectusEndpoint(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		result, err := e.Whs.FetchDataFromDirectusService(ctx)
		if err != nil {
			log.Printf("[ERROR] FetchDataFromDirectusEndpoint service => %+v\n", err)
			res := map[string]interface{}{
				"success": true,
				"data":    result,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}

		res := map[string]interface{}{
			"success": true,
			"data":    result,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (e *Endpoint) CallDirectusEndpoint(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		result, err := e.Whs.CallDirectusService(ctx)
		if err != nil {
			log.Printf("[ERROR] CallDirectusEndpoint service => %+v\n", err)
			res := map[string]interface{}{
				"success": true,
				"data":    result,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		res := map[string]interface{}{
			"success": true,
			"data":    result,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
