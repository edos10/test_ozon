package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestShortUrlHandlerInMemory(t *testing.T) {
	db := NewInMemoryDatabase()
	db.InitializeCurrentString()
	defer db.Close()

	// Инициализация базы данных

	handler := SendUrlHandler(db)

	requestBody := SendRequest{
		URL: "https://example.com",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Создание временного HTTP-сервера
	server := httptest.NewServer(handler)
	defer server.Close()

	// Отправка запроса к временному серверу
	client := &http.Client{}
	req, err := http.NewRequest("POST", server.URL+"/send", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200; got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	// Чтение тела ответа
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Распаковка JSON-ответа
	var response SendResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Получение оригинального URL из базы данных
	originalURL, err := db.GetUrl(response.ShortURL)
	if err != nil {
		t.Fatalf("Failed to get original URL from database: %v", err)
	}

	// Проверка корректности оригинального URL
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
