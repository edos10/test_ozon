package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SendRequest struct {
	URL string `json:"url"`
}

type SendResponse struct {
	ShortURL string `json:"short_url"`
}

type GetRequest struct {
	ShortURL string `json:"short_url"`
}

type GetResponse struct {
	URL string `json:"url"`
}

func SendUrlHandler(db Database, currentURL *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var req SendRequest
		err := decoder.Decode(&req)

		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		shortURL := NextUrlString(*currentURL)

		// обрезаем справа от ссылки /, чтобы например youtube.com и youtube.com/ были одинаковы

		req.URL = strings.TrimRight(req.URL, "/")
		err = db.Save(shortURL, req.URL)
		if len(req.URL) == 0 {
			http.Error(w, "Failed to parse link or error in json", http.StatusBadRequest)
			return
		}
		if err != nil {
			errStr := fmt.Sprintf("Failed to save URL: %s", err)
			http.Error(w, errStr, http.StatusInternalServerError)
			return
		}

		*currentURL = shortURL
		resp := SendResponse{
			ShortURL: shortURL,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func GetUrlHandler(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var req GetRequest
		err := decoder.Decode(&req)

		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		originalURL, err := db.Get(req.ShortURL)

		if err != nil {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		resp := GetResponse{
			URL: originalURL,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
