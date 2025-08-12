package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
)

// ========== Database & Models ==========

var db *pg.DB

type Album struct {
	ID     string  `json:"id" pg:"id"`
	Title  string  `json:"title" pg:"title"`
	Artist string  `json:"artist" pg:"artist"`
	Price  float64 `json:"price" pg:"price"`
}

func (Album) TableName() string {
	return "albums"
}

// ========== Database Connection ==========

func connectDB() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Get database configuration from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost" // default to localhost if not specified
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432" // default PostgreSQL port
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		log.Fatal("DB_USER environment variable must be set")
	}

	password := os.Getenv("DB_PASSWORD")
	// Password can be empty if your database doesn't require one

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		log.Fatal("DB_NAME environment variable must be set")
	}

	// Construct the address from host and port
	addr := host + ":" + port

	opts := &pg.Options{
		Addr:     addr,
		User:     user,
		Password: password,
		Database: dbname,
		OnConnect: func(ctx context.Context, conn *pg.Conn) error {
			log.Println("Connected to PostgreSQL!")
			return nil
		},
	}

	db = pg.Connect(opts)

	var exists bool
	_, err := db.QueryOne(pg.Scan(&exists), `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'albums')`)
	if err != nil || !exists {
		log.Fatal("Albums table doesn't exist or can't be accessed")
	}

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println(" Database connected successfully")
}

// ========== HTTP Handlers ==========

func albumsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAlbums(w, r)
	case http.MethodPost:
		postAlbum(w, r)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func albumByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Trim "/albums/" from URL path for ID
	id := strings.TrimPrefix(r.URL.Path, "/albums/")
	if id == "" {
		sendError(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getAlbumByID(w, r, id)
	case http.MethodDelete:
		deleteAlbumByID(w, id)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ========== CRUD Operations ==========

func getAlbums(w http.ResponseWriter, r *http.Request) {
	var albums []Album
	if err := db.Model(&albums).Select(); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, http.StatusOK, albums)
}

func getAlbumByID(w http.ResponseWriter, r *http.Request, id string) {
	var album Album
	err := db.Model(&album).Where("id = ?", id).Select()

	switch err {
	case nil:
		sendJSON(w, http.StatusOK, album)
	case pg.ErrNoRows:
		sendError(w, "album not found", http.StatusNotFound)
	default:
		sendError(w, err.Error(), http.StatusInternalServerError)
	}
}

func postAlbum(w http.ResponseWriter, r *http.Request) {
	var newAlbum Album
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&newAlbum); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := db.Model(&newAlbum).Insert(); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, http.StatusCreated, newAlbum)
}

func deleteAlbumByID(w http.ResponseWriter, id string) {
	res, err := db.Model(&Album{ID: id}).WherePK().Delete()
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.RowsAffected() == 0 {
		sendError(w, "album not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// ========== Helper Functions ==========

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Error response consistently in JSON with "error" key
	resp := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode JSON error response: %v", err)
	}
}

// ========== Main Function ==========

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connectDB()
	defer db.Close()

	http.HandleFunc("/albums", albumsHandler)
	http.HandleFunc("/albums/", albumByIDHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
