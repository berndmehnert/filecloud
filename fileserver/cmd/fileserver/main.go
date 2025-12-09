package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"example.com/filecloud/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
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

	db, err := sqlx.Connect("sqlite", databaseFile)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := migrate(db.DB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	r := chi.NewRouter()
	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		handler.HandleUpload(w, r, db)
	})
	r.Get("/files/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		handler.HandleDownload(w, r, db)
	})
	secret := []byte("replace-with-secure-random-secret")
	r.Get("/api/files", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		handler.HandleListFiles(w, r, db, secret)
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func migrate(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  filename TEXT NOT NULL,
  size INTEGER NOT NULL,
  mime TEXT,
  checksum TEXT,
  path TEXT NOT NULL,
  created_at DATETIME NOT NULL
);
`
	_, err := db.Exec(schema)
	return err
}
