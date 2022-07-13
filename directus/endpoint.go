package directus

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.uber.org/zap"
)

type WebhookEndpoint interface {
	FetchDataFromDirectusEndpoint(ctx context.Context, logger *zap.SugaredLogger) http.HandlerFunc
	CallDirectusEndpoint(ctx context.Context, logger *zap.SugaredLogger) http.HandlerFunc
}

type Endpoint struct {
	Whs WebhookServices
}

func NewEndpoint(whs WebhookServices) *Endpoint {
	return &Endpoint{
		Whs: whs,
	}
}

func (e *Endpoint) FetchDataFromDirectusEndpoint(ctx context.Context, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger.Info(
			zap.Any("url", r.URL),
			zap.String("method", r.Method),
		)
		w.Header().Set("Content-Type", "application/json")

		result, err := e.Whs.FetchDataFromDirectusService(ctx, logger)
		if err != nil {
			logger.With("error", err).Error(err)
			res := map[string]interface{}{
				"success": false,
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

func (e *Endpoint) CallDirectusEndpoint(ctx context.Context, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		result, err := e.Whs.CallDirectusService(ctx, logger)
		if err != nil {
			log.Printf("[ERROR] CallDirectusEndpoint service => %+v\n", err)
			res := map[string]interface{}{
				"success": false,
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
