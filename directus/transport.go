package directus

import (
	"context"
	"net/http"
)

type WebhookTransports interface {
	FetchDataFromDirectus(ctx context.Context) http.Handler
	CallDirectus() http.Handler
}

type Transports struct {
	whe WebhookEndpoint
}

func NewTransports(whe WebhookEndpoint) *Transports {
	return &Transports{
		whe: whe,
	}
}

func (w *Transports) FetchDataFromDirectus(ctx context.Context) http.Handler {
	return w.whe.FetchDataFromDirectusEndpoint(ctx)
}

func (w *Transports) CallDirectus(ctx context.Context) http.Handler {
	return w.whe.CallDirectusEndpoint(ctx)
}
