package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"os"
)

type InMemoryDatabase struct {
	shortToOriginal map[string]string
}

type Database interface {
	Save(shortURL, originalURL string) error
	Get(shortURL string) (string, error)
	io.Closer
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		shortToOriginal: make(map[string]string),
	}
}

func (db *InMemoryDatabase) Save(shortURL, originalURL string) error {
	for _, v := range db.shortToOriginal {
		if v == originalURL {
			return fmt.Errorf("Original URL already exists")
		}
	}

	_, ok := db.shortToOriginal[shortURL]
	if ok {
		return fmt.Errorf("Short URL already exists")
	}

	db.shortToOriginal[shortURL] = originalURL

	return nil
}

func (db *InMemoryDatabase) Get(shortURL string) (string, error) {
	fmt.Println(shortURL, db.shortToOriginal)
	originalURL, ok := db.shortToOriginal[shortURL]
	if !ok {
		return "", fmt.Errorf("Short URL not found")
	}
	return originalURL, nil
}

func (db *InMemoryDatabase) Close() error {
	return nil
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase() *PostgresDatabase {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}
	return &PostgresDatabase{
		db: db,
	}
}

func (db *PostgresDatabase) Save(shortURL, originalURL string) error {
	if _, err := db.db.Exec("SELECT * FROM urls WHERE 'original_url' = $1"); err != nil {
		return fmt.Errorf("Original URL already exists!")
	}
	_, err := db.db.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)
	return err
}

func (db *PostgresDatabase) Get(shortURL string) (string, error) {
	var originalURL string
	err := db.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Short URL not found")
		}
		return "", err
	}
	return originalURL, nil
}

func (db *PostgresDatabase) Close() error {
	return db.db.Close()
}
