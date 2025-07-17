package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	healthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	if body := rr.Body.String(); body != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", body)
	}
}

func TestBooksHandler_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rr := httptest.NewRecorder()

	booksHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Go in Action") {
		t.Errorf("Expected response body to contain 'Go in Action'")
	}
}

func TestBooksHandler_Post(t *testing.T) {
	newBook := Book{Title: "DevOps 101", Author: "Alice", Rating: 4}
	b, _ := json.Marshal(newBook)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	booksHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	var result Book
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if result.Title != newBook.Title {
		t.Errorf("Expected book with title '%s', got '%s'", newBook.Title, result.Title)
	}
}

func TestBookByIDHandler_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	rr := httptest.NewRecorder()

	r := req.Clone(req.Context())
	r.URL.Path = "/books/1"

	bookByIDHandler(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestBookByIDHandler_Put(t *testing.T) {
	update := Book{Title: "Updated Book", Author: "Bob", Rating: 5}
	b, _ := json.Marshal(update)
	id := 1
	req := httptest.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(id), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := req.Clone(req.Context())
	r.URL.Path = "/books/" + strconv.Itoa(id)

	bookByIDHandler(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var result Book
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if result.Title != update.Title {
		t.Errorf("Expected updated title '%s', got '%s'", update.Title, result.Title)
	}
}

func TestBookByIDHandler_Delete(t *testing.T) {
	id := 2
	req := httptest.NewRequest(http.MethodDelete, "/books/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()

	r := req.Clone(req.Context())
	r.URL.Path = "/books/" + strconv.Itoa(id)

	bookByIDHandler(rr, r)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rr.Code)
	}
}
