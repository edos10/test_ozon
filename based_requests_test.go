package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestShortUrlHandlerInMemory(t *testing.T) {
	db := NewInMemoryDatabase()
	defer db.Close()

	// Инициализация базы данных

	handler := SendUrlHandler(db)

	requestBody := SendRequest{
		URL: "https://www.example.com",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Создание временного HTTP-сервера
	server := httptest.NewServer(handler)
	defer server.Close()

	// Отправка запроса к временному серверу
	time.Sleep(1 * time.Second)
	resp, err := http.Post(server.URL+"/send", "application/json", bytes.NewBuffer(body))
	time.Sleep(1 * time.Second)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200; got %d", resp.StatusCode)
	}

	// Проверка тела ответа
	var response SendResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Проверка корректности сокращенной ссылки
	if len(response.ShortURL) != sizeUrl {
		t.Errorf("Expected shortened URL length %d; got %d", sizeUrl, len(response.ShortURL))
	}

	// Проверка, что сокращенная ссылка действительно сохранена в базе данных
	originalURL, err := db.GetUrl(response.ShortURL)
	if err != nil {
		t.Fatalf("Failed to get original URL from database: %v", err)
	}

	if originalURL != requestBody.URL {
		t.Errorf("Expected original URL %s; got %s", requestBody.URL, originalURL)
	}
}

func TestShortUrlHandlerPostgres(t *testing.T) {
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "postgres")
	_ = os.Setenv("DB_PASSWORD", "default")
	_ = os.Setenv("DB_NAME", "urls")
	_ = os.Setenv("STORAGE", "postgres")
	db := NewPostgresDatabase()
	defer db.Close()
	//handler := SendUrlHandler(db)

}
