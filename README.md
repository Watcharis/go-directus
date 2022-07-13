# install project & run project

    1. clone project
        * https://github.com/Watcharis/go-directus.git

    2. install package go
        * go mod download
        * go mod tidy

    3. run project use command
        * go run main.go

# install directus
    docker-compose up -d

# Present graceful shutdown in Go

```
func GraceFullShutdownAndRunServer() {
	var wg sync.WaitGroup

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	handler, err := setupServer(sugarLogger)
	if err != nil {
		sugarLogger.Infof("[ Error setUpServer ] => %+v\n", err)
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
```


