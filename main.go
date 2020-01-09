package main

import (
	"context"
	"encoding/json"
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
	r.HandleFunc("/v1/long", Long)
	r.HandleFunc("/v1/items", ListItems)
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

func Long(w http.ResponseWriter, r *http.Request)  {
	log.Println("start long request")
	defer log.Println("end long request")

	// testing for past 5sec timeout
	// after 5sec, the func would repeat.. not good
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		fmt.Printf("%d sec", i)
	}
	fmt.Fprintf(w, "Request: %s %s", r.Method, r.URL.Path)
}

type Item struct {
	Name string `json:"name"`
	Cost string `json:"cost"`
	Description string `json:"description"`
}

func ListItems(w http.ResponseWriter, r *http.Request)  {
	list := []Item{
		{Name: "backpack", Cost: "50ss", Description: "Carries 5 bulk"},
		{Name: "sack", Cost: "80bp", Description: "Carries 5 bulk, requires at least one hand"},
		{Name: "satchel", Cost: "30ss", Description: "Carries 2 bulk"},
		{Name: "compass", Cost: "2gc", Description: "shows north"},
		{Name: "50ft of rope", Cost: "20ss", Description: "hemp rope"},
		{Name: "spyglass", Cost: "5gc", Description: "4x magnification"},
		{Name: "crowbar", Cost: "18ss", Description: "tool to pry open"},
		{Name: "hammer", Cost: "15ss", Description: "tool to hammer"},
		{Name: "hourglass", Cost: "1gc", Description: "takes 1 hour for sand"},
		{Name: "fishing rod", Cost: "10ss", Description: "tool for fishing"},
		{Name: "1lb fish bait", Cost: "5bp", Description: "bugs to help catch fish"},
	}
	data, err := json.Marshal(list)
	if err != nil {
		log.Println("error: marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error: writing result", err)
	}
}