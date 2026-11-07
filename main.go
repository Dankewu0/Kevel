package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	urlStore = make(map[string]string)
	mu       sync.RWMutex
)

func generateID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL is empty", http.StatusBadRequest)
		return
	}

	id := generateID()

	mu.Lock()
	urlStore[id] = longURL
	mu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Готово!</h1><p>Твоя ссылка: <a href=\"/%s\">localhost:8080/%s</a></p><p><a href=\"/\">Сократить еще одну</a></p>", id, id)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "index.html")
		return
	}

	id := r.URL.Path[1:]

	mu.RLock()
	originalURL, exists := urlStore[id]
	mu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "styles.css")
	})

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
