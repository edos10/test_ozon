package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/edos10/test_ozon/greetings"
	"github.com/gorilla/mux"
)

func InitializeCurrentString(db *sql.DB) error {
	exists, err := checkTableExists(db, "genstring")
	if err != nil {
		return err
	}

	if !exists {
		if err := createTableForString(db); err != nil {
			return err
		}
		value := "aaaaaaaaaa"
		_, err := db.Exec("INSERT INTO genstring (currentstring) VALUES ($1)", value)
		if err != nil {
			return err
		}

		log.Println("Initialized currentString with value:", value)
	} else {
		log.Println("currentString already initialized")
	}

	return nil
}

func InitializeForUrls(db *sql.DB) error {
	tableName := "urls"

	exists, err := checkTableExists(db, tableName)
	if err != nil {
		return err
	}

	if !exists {
		err := createTableForUrls(db)
		if err != nil {
			return err
		}
		log.Printf("Table '%s' created\n", tableName)
	} else {
		log.Printf("Table '%s' already exists\n", tableName)
	}

	return nil
}

func checkTableExists(db *sql.DB, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		)
	`

	var exists bool
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func createTableForString(db *sql.DB) error {
	query := `
		CREATE TABLE genstring (
			currentstring TEXT
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func createTableForUrls(db *sql.DB) error {
	query := `
		CREATE TABLE urls (
			original_url TEXT,
			short_url TEXT
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// как вариант, можно в main запускать порядка 10-20 горутин

func worker(wg *sync.WaitGroup, jobQueue <-chan *http.Request, handler http.Handler) {
	for req := range jobQueue {
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		log.Printf("Request: %s %s | Response: %d %s\n", req.Method, req.URL.Path, res.Code, res.Body.String())
	}
	wg.Done()
}

func InitDbForLocal() {
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "postgres")
	_ = os.Setenv("DB_PASSWORD", "default")
	_ = os.Setenv("DB_NAME", "urls")
}

func main() {
	// инициализируем мапы для работы
	makeMaps()
	log.Println("HER")
	log.Println(greetings.Hello("Ed"))
	// InitDbForLocal()
	if len(os.Args) == 2 {
		log.Println("Start program with argument storage")
		_ = os.Setenv("STORAGE", os.Args[1])
	} else {
		log.Fatalf("Error in cmd line, restart program with right arguments, please.")
		return
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	storageType := os.Getenv("STORAGE")

	if storageType == "" {
		log.Fatal("Wrong input parameters, restart service please with parameter database - in-memory or postgres")
	}

	var db Database
	log.Printf("Service storage is %s", storageType)
	switch storageType {
	case "in-memory":
		db = NewInMemoryDatabase()
	case "postgres":
		db = NewPostgresDatabase()
		dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)

		dbForInit, err := sql.Open("postgres", dataSourceName)
		if err != nil {
			panic(err)
		}

		if err := InitializeForUrls(dbForInit); err != nil {
			panic(err)
		}

	default:
		log.Fatal("Invalid storage type")
	}

	if err := db.InitializeCurrentString(); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/send", SendUrlHandler(db)).Methods("POST")
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

	log.Println("Service successfully stopped")
}
