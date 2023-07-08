package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testShortenUrlHandlerInMemory(t *testing.T) {
	db := NewInMemoryDatabase()
	defer db.Close()
	str := "aaaaaaaaaa"
	handler := SendUrlHandler(db, &str)

	// Создание тестового запроса
	requestBody := SendRequest{
		URL: "https://www.example.com",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/send", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Создание тестового ResponseWriter для записи ответа
	recorder := httptest.NewRecorder()

	// Выполнение запроса
	handler.ServeHTTP(recorder, req)

	// Проверка статуса ответа
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	// Проверка тела ответа
	var response SendResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Проверка корректности сокращенной ссылки
	if len(response.ShortURL) != sizeUrl {
		t.Errorf("Expected shortened URL length %d; got %d", sizeUrl, len(response.ShortURL))
	}

	// Проверка, что сокращенная ссылка действительно сохранена в базе данных
	originalURL, err := db.Get(response.ShortURL)
	if err != nil {
		t.Fatalf("Failed to get original URL from database: %v", err)
	}

	if originalURL != requestBody.URL {
		t.Errorf("Expected original URL %s; got %s", requestBody.URL, originalURL)
	}
}

func testShortenUrlHandlerPostgres(t *testing.T) {
	db := NewPostgresDatabase()
	defer db.Close()

	str := "aaaaaaaaaa"
	handler := SendUrlHandler(db, &str)

	// Создание тестового запроса
	requestBody := SendRequest{
		URL: "https://www.example.com",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/shorten", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Создание тестового ResponseWriter для записи ответа
	recorder := httptest.NewRecorder()

	// Выполнение запроса
	handler.ServeHTTP(recorder, req)

	// Проверка статуса ответа
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	// Проверка тела ответа
	var response SendResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Проверка корректности сокращенной ссылки
	if len(response.ShortURL) != sizeUrl {
		t.Errorf("Expected shortened URL length %d; got %d", sizeUrl, len(response.ShortURL))
	}

	// Проверка, что сокращенная ссылка действительно сохранена в базе данных
	originalURL, err := db.Get(response.ShortURL)
	if err != nil {
		t.Fatalf("Failed to get original URL from database: %v", err)
	}

	if originalURL != requestBody.URL {
		t.Errorf("Expected original URL %s; got %s", requestBody.URL, originalURL)
	}
}
