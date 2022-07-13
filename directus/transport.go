package directus

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type WebhookTransports interface {
	FetchDataFromDirectus(ctx context.Context, logger *zap.SugaredLogger) http.Handler
	CallDirectus(ctx context.Context, logger *zap.SugaredLogger) http.Handler
}

type Transports struct {
	whe WebhookEndpoint
}

func NewTransports(whe WebhookEndpoint) *Transports {
	return &Transports{
		whe: whe,
	}
}

func (w *Transports) FetchDataFromDirectus(ctx context.Context, logger *zap.SugaredLogger) http.Handler {
	return w.whe.FetchDataFromDirectusEndpoint(ctx, logger)
}

func (w *Transports) CallDirectus(ctx context.Context, logger *zap.SugaredLogger) http.Handler {
	return w.whe.CallDirectusEndpoint(ctx, logger)
}
