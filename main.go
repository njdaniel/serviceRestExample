package main

import (
	"context"
	"database/sql"
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


func ConnectDB() *sql.DB {
	log.Println("Connecting to db")
	defer log.Println("Connected to db")
	host := "localhost:5432"
	// TODO: remove hardcoded pwd
	password := "password"
	dbsource := fmt.Sprintf("postgres://postgres:%s@%s/postgres?sslmode=disable", password, host)
	db, err := sql.Open("postgres", dbsource)
	if err != nil {
		log.Printf("error opening db: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Printf("error contacting db %v", err)
	}
	return db
}

type Item struct {
	Name string `json:"name"`
	Cost string `json:"cost"`
	Description string `json:"description"`
}
type Items []Item

func ListItems(w http.ResponseWriter, r *http.Request)  {
	const query = `SELECT * FROM items`
	db := ConnectDB()
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("error: could not query: %v \n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items Items
	for rows.Next() {
		item := Item{}
		err := rows.Scan(&item.Name, &item.Cost, &item.Description)
		if err != nil {
			log.Printf("error: could not scan row into item, %v \n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	data, err := json.Marshal(items)
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