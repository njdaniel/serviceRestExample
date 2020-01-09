package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// ====================================
	// App starting
	log.Printf("main: started")
	defer log.Printf("main: completed")

	// =====================================
	// Start API Service
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	api := http.Server{
		Addr:         "localhost:8001",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a buffered channel to listen for errors coming from listener.
	serverErrors := make(chan error, 1)

	// Start listener
	go func() {
		log.Printf("main: listener goroutine: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	//if err := http.ListenAndServe("localhost:8001", r); err != nil {
	//	log.Fatalf("error: listening and serving: %s", err)
	//}

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Site it up")
}