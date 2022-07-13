package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"watcharis/go-directus/directus"

	"github.com/gorilla/mux"
)

const (
	defultPort = "8022"
)

func GraceFullShutdownAndRunServer() {
	var wg sync.WaitGroup

	handler, err := setupServer()
	if err != nil {
		log.Printf("[ Error setUpServer ] => %+v", err.Error())
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", defultPort),
		Handler: handler,
	}

	go func(srv *http.Server) {
		log.Printf("[ GO HTTP START SERVER ] %+v\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[ Error server close ] => %+v\n", err)
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

func setupServer() (http.Handler, error) {

	ctx := context.Background()
	router := mux.NewRouter()

	directusService := directus.NewService()
	directusEndpoint := directus.NewEndpoint(directusService)
	directusTransport := directus.NewTransports(directusEndpoint)

	//---- router -----
	directusRouter := router.PathPrefix("/directus").Subrouter()
	directusRouter.Handle("/fetch", directusTransport.CallDirectus(ctx)).Methods(http.MethodGet)
	directusRouter.Handle("/test", directusTransport.FetchDataFromDirectus(ctx)).Methods(http.MethodGet)

	return router, nil
}

func main() {
	GraceFullShutdownAndRunServer()
}
