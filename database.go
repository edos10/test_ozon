package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
)

type InMemoryDatabase struct {
	shortToOriginal map[string]string
	currentString   string
}

type Database interface {
	SaveUrl(shortURL, originalURL string) error
	GetUrl(shortURL string) (string, error)
	SaveCurrentString(currentString string) error
	GetCurrentString() string
	InitializeCurrentString() error
	io.Closer
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		shortToOriginal: make(map[string]string),
	}
}

func (db *InMemoryDatabase) SaveUrl(shortURL, originalURL string) error {
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

func (db *InMemoryDatabase) GetUrl(shortURL string) (string, error) {
	fmt.Println(shortURL, db.shortToOriginal)
	originalURL, ok := db.shortToOriginal[shortURL]
	if !ok {
		return "", fmt.Errorf("Short URL not found")
	}
	return originalURL, nil
}

func (db *InMemoryDatabase) GetCurrentString() string {
	return db.currentString
}

func (db *InMemoryDatabase) SaveCurrentString(currentString string) error {
	db.currentString = currentString
	return nil
}

func (db *InMemoryDatabase) InitializeCurrentString() error {
	db.currentString = "aaaaaaaaaa"
	return nil
}

func (db *InMemoryDatabase) Close() error {
	return nil
}

type PostgresDatabase struct {
	db *sql.DB
}

func (db *PostgresDatabase) InitializeCurrentString() error {
	exists, err := checkTableExists(db.db, "genstring")
	if err != nil {
		return err
	}

	if !exists {
		if err := createTableForString(db.db); err != nil {
			return err
		}
		value := "aaaaaaaaaa"
		_, err := db.db.Exec("INSERT INTO genstring (currentstring) VALUES ($1)", value)
		if err != nil {
			return err
		}

		log.Println("Initialized currentString with value:", value)
	} else {
		log.Println("currentString already initialized")
	}

	return nil
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

func (db *PostgresDatabase) SaveUrl(shortURL, originalURL string) error {
	result, errQuery := db.db.Exec("SELECT * FROM urls WHERE original_url = $1", originalURL)
	if errQuery != nil {
		return errQuery
	}
	rowsAffected, errRows := result.RowsAffected()
	if errRows != nil {
		return errRows
	}
	if rowsAffected > 0 {
		return fmt.Errorf("Original URL already exists!")
	}
	_, err := db.db.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)
	return err
}

func (db *PostgresDatabase) GetUrl(shortURL string) (string, error) {
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

func (db *PostgresDatabase) GetCurrentString() string {
	var currentString string
	err := db.db.QueryRow("SELECT genstring FROM genstring LIMIT 1").Scan(&currentString)
	if err != nil {
		return fmt.Sprintf("failed to get current string: %w", err)
	}
	return currentString
}

func (db *PostgresDatabase) SaveCurrentString(currentString string) error {
	log.Printf(currentString)
	_, err := db.db.Exec("UPDATE genstring SET genstring = $1", currentString)
	if err != nil {
		return fmt.Errorf("failed to update current string: %w", err)
	}
	return nil
}

func (db *PostgresDatabase) Close() error {
	return db.db.Close()
}
