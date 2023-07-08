package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(wg *sync.WaitGroup, jobQueue <-chan *http.Request, handler http.Handler) {
	for req := range jobQueue {
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		log.Printf("Request: %s %s | Response: %d %s\n", req.Method, req.URL.Path, res.Code, res.Body.String())
	}
	wg.Done()
}

func main() {
	makeMaps()

	currentString := "aaaaaaaaaa"

	storageType := os.Args[1]
	if len(os.Args) != 2 {
		log.Fatal("Wrong input parameters, restart service please with parameter database - redis or postgres")
	}

	var db Database

	switch storageType {
	case "in-memory":
		db = NewInMemoryDatabase()
	case "postgres":
		db = NewPostgresDatabase()
	default:
		log.Fatal("Invalid storage type")
	}

	r := mux.NewRouter()
	r.HandleFunc("/send", SendUrlHandler(db, &currentString)).Methods("POST")
	r.HandleFunc("/get", GetUrlHandler(db)).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server successfully stopped")
}
