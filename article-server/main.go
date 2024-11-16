package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

// Article represents a blog article with an ID, Heading, Description, and Tags
type Article struct {
	ID          int      `json:"id"`
	Heading     string   `json:"heading"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// in-memory storage for articles
var articles = make(map[int]Article)
var mu sync.Mutex // mutex for safe concurrent access

// Helper function to log HTTP request details
func logRequest(r *http.Request) {
	log.Printf("Received %s request for %s from %s\n", r.Method, r.URL, r.RemoteAddr)
}

// Create a new article
func createArticle(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	var article Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	mu.Lock()
	article.ID = rand.Intn(1000) // generate a random ID for the article
	articles[article.ID] = article
	mu.Unlock()

	log.Printf("Article created with ID %d\n", article.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

// Retrieve a single article by ID
func getArticle(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Invalid ID:", err)
		return
	}

	mu.Lock()
	article, exists := articles[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Article not found", http.StatusNotFound)
		log.Printf("Article with ID %d not found\n", id)
		return
	}

	log.Printf("Article retrieved with ID %d\n", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// Retrieve all articles
func getAllArticles(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	mu.Lock()
	defer mu.Unlock()

	var allArticles []Article
	for _, article := range articles {
		allArticles = append(allArticles, article)
	}

	log.Printf("All articles retrieved, count: %d\n", len(allArticles))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allArticles)
}

// Update an existing article by ID
func updateArticle(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Invalid ID:", err)
		return
	}

	var updatedArticle Article
	if err := json.NewDecoder(r.Body).Decode(&updatedArticle); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	mu.Lock()
	if _, exists := articles[id]; !exists {
		mu.Unlock()
		http.Error(w, "Article not found", http.StatusNotFound)
		log.Printf("Article with ID %d not found for update\n", id)
		return
	}
	updatedArticle.ID = id
	articles[id] = updatedArticle
	mu.Unlock()

	log.Printf("Article updated with ID %d\n", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedArticle)
}

// Delete an article by ID
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Invalid ID:", err)
		return
	}

	mu.Lock()
	if _, exists := articles[id]; !exists {
		mu.Unlock()
		http.Error(w, "Article not found", http.StatusNotFound)
		log.Printf("Article with ID %d not found for deletion\n", id)
		return
	}
	delete(articles, id)
	mu.Unlock()

	log.Printf("Article deleted with ID %d\n", id)

	w.WriteHeader(http.StatusNoContent)
}

// Seed some test articles
func seedArticles() {
	articles[1] = Article{ID: 1, Heading: "Go Concurrency", Description: "Learn about goroutines and channels in Go.", Tags: []string{"Go", "Concurrency"}}
	articles[2] = Article{ID: 2, Heading: "REST API Design", Description: "Best practices for designing RESTful APIs.", Tags: []string{"API", "REST"}}
	articles[3] = Article{ID: 3, Heading: "Microservices Architecture", Description: "An introduction to microservices.", Tags: []string{"Microservices", "Architecture"}}
	log.Println("Seeded initial articles")
}

func main() {
	// Seed articles
	seedArticles()

	http.HandleFunc("/api/v1/article", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createArticle(w, r)
		case http.MethodGet:
			if r.URL.Query().Has("id") {
				getArticle(w, r)
			} else {
				getAllArticles(w, r)
			}
		case http.MethodPut:
			updateArticle(w, r)
		case http.MethodDelete:
			deleteArticle(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		response := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Welcome</title>
			</head>
			<body>
				<h1>Welcome to the Article Server</h1>
				<p>Use the <code>/api/v1/article</code> endpoint to interact with the server.</p>
			</body>
			</html>
		`
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)
	})

	fmt.Println("Server is running on port 9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
