package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
	Review string `json:"review,omitempty"`
}

var books = []Book{
	{ID: 1, Title: "Go in Action", Author: "William Kennedy", Rating: 5},
	{ID: 2, Title: "Clean Code", Author: "Robert C. Martin", Rating: 5},
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/", bookByIDHandler)

	log.Println("📘 Book API running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("⚠️ Failed to write health check response: %v", err)
	}
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(books); err != nil {
			log.Printf("⚠️ Failed to encode books: %v", err)
		}
	case http.MethodPost:
		var b Book
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		b.ID = getNextBookID()
		books = append(books, b)
		go autoRemoveBook(b.ID, 2*time.Minute)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(b); err != nil {
			log.Printf("⚠️ Failed to encode created book: %v", err)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func bookByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		for _, b := range books {
			if b.ID == id {
				if err := json.NewEncoder(w).Encode(b); err != nil {
					log.Printf("⚠️ Failed to encode book by ID: %v", err)
				}
				return
			}
		}
		http.NotFound(w, r)
	case http.MethodPut:
		var updated Book
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i, b := range books {
			if b.ID == id {
				updated.ID = id
				books[i] = updated
				if err := json.NewEncoder(w).Encode(updated); err != nil {
					log.Printf("⚠️ Failed to encode updated book: %v", err)
				}
				return
			}
		}
		http.NotFound(w, r)
	case http.MethodDelete:
		for i, b := range books {
			if b.ID == id {
				books = append(books[:i], books[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.NotFound(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getNextBookID() int {
	maxID := 0
	for _, b := range books {
		if b.ID > maxID {
			maxID = b.ID
		}
	}
	return maxID + 1
}

func autoRemoveBook(id int, delay time.Duration) {
	time.Sleep(delay)
	for i, b := range books {
		if b.ID == id {
			log.Printf("⏳ Auto-removing book ID %d\n", id)
			books = append(books[:i], books[i+1:]...)
			return
		}
	}
}
