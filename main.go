package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var telegramAPI = "https://api.telegram.org"
var secretPath = "${{ secrets.MYSECRET_WAY }}"

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Проверяем секрет
	if !strings.HasPrefix(path, "/"+secretPath+"/") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Убираем секрет из пути
	cleanPath := strings.TrimPrefix(path, "/"+secretPath)

	targetURL := telegramAPI + cleanPath

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	defer resp.Body.Close()

	// копируем заголовки
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Proxy started on :" + port)
	http.HandleFunc("/", proxyHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
