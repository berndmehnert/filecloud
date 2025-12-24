package main

import (
	"log"
	"net/http"
	"os"

	"example.com/filecloud/internal/database"
	"example.com/filecloud/internal/handler"
	"example.com/filecloud/internal/middleware"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

const (
	databaseFile = "./files.db"
	storageDir   = "./storage"
)

func main() {
	if err := os.MkdirAll(storageDir, 0o755); err != nil {
		log.Fatalf("creating storage dir: %v", err)
	}

	db, err := database.New(databaseFile)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.CorsMiddleware)

	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleUpload(w, r, db)
	})
	r.Get("/files/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleDownload(w, r, db)
	})
	secret := []byte("replace-with-secure-random-secret")
	r.Get("/api/files", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleListFiles(w, r, db, secret)
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
