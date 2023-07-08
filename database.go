package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"io"
)

type RedisDatabase struct {
	client *redis.Client
}

func NewRedisDatabase() *RedisDatabase {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisDatabase{
		client: client,
	}
}

func (db *RedisDatabase) Save(shortURL, originalURL string) error {
	ctx := context.Background()
	err := db.client.Set(ctx, shortURL, originalURL, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}

func (db *RedisDatabase) Get(shortURL string) (string, error) {
	ctx := context.Background()
	val, err := db.client.Get(ctx, shortURL).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("short URL not found")
		}
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	return val, nil
}

func (db *RedisDatabase) Close() error {
	err := db.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	return nil
}

type Database interface {
	Save(shortURL, originalURL string) error
	Get(shortURL string) (string, error)
	io.Closer
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase() *PostgresDatabase {
	db, err := sql.Open("postgres", "postgres://postgres:default@localhost:5432/urls?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return &PostgresDatabase{
		db: db,
	}
}

func (db *PostgresDatabase) Save(shortURL, originalURL string) error {
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
