package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Karanth1r3/l_0/internal/config"
	"github.com/Karanth1r3/l_0/internal/service"
	"github.com/Karanth1r3/l_0/internal/storage"
	"github.com/Karanth1r3/l_0/internal/utils"
)

func main() {
	// trying to initialize config based on provided file
	cfg, err := config.Parse("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	// []byte data from page to feed to the server (i guess)
	htmlPage, err := os.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
	// setting db connection was put to an internal module for convinience
	dbConn, err := utils.ConnectDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { // on main goroutine executed cast conn close check
		if err := dbConn.Close(); err != nil {
			log.Println(err)
		}
	}()
	// connect to nats streaming service using config params & defer it's closure on end
	natsConn, err := utils.ConnectNats(&cfg.NATS)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.NATS.Host, cfg.NATS.Port)
	defer natsConn.Close()

	stanConn, err := utils.ConnectStan(natsConn, &cfg.STAN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := stanConn.Close(); err != nil {
			log.Println(err)
		}
	}()

	storage := storage.NewStorage(dbConn)
	service, err := service.NewService(storage, stanConn, cfg.Service, htmlPage)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(service.HandleIndex))
	mux.Handle("/get", http.HandlerFunc(service.HandleGet))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: mux,
	}

	go func() {
		if err = httpServer.ListenAndServe(); err != nil {
			log.Print("listen and serve err", err)
		}
	}()

	// stop-Channel && channel to check if cleanup goroutine executed (graceful shutdown)
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt) // on interrupt/termination, signal shall be notified
	go func() {
		for range signalChan {
			log.Printf("\nReceived an interrupt...\n\n")
			_ = httpServer.Shutdown(context.Background())
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
