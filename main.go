package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"watcharis/go-directus/directus"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	defultPort = "8022"
)

func GraceFullShutdownAndRunServer() {
	var wg sync.WaitGroup

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	handler, err := setupServer(sugarLogger)
	if err != nil {
		log.Printf("[ Error setUpServer ] => %+v", err.Error())
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", defultPort),
		Handler: handler,
	}

	go func(srv *http.Server) {
		sugarLogger.Info("[ GO start http server port] ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errServer := errors.New("error server close")
			sugarLogger.With("error", errServer)
			return
		}
	}(srv)

	sigint := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	sigs := <-sigint
	wg.Add(1)
	go func(sigs os.Signal) {
		defer wg.Done()
		fmt.Printf("[SIGNAL STATUS] => %+v\n", sigs)
	}(sigs)
	wg.Wait()
}

func setupServer(sugarLogger *zap.SugaredLogger) (http.Handler, error) {

	ctx := context.Background()
	router := mux.NewRouter()

	directusService := directus.NewService()
	directusEndpoint := directus.NewEndpoint(directusService)
	directusTransport := directus.NewTransports(directusEndpoint)

	//---- router -----
	directusRouter := router.PathPrefix("/directus").Subrouter()
	directusRouter.Handle("/fetch", directusTransport.CallDirectus(ctx, sugarLogger)).Methods(http.MethodGet)
	directusRouter.Handle("/test", directusTransport.FetchDataFromDirectus(ctx, sugarLogger)).Methods(http.MethodGet)

	return router, nil
}

func main() {
	GraceFullShutdownAndRunServer()
}
