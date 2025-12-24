package model

import "time"

// FileMeta represents a stored file's metadata.
// Use db tags for sqlx mapping.
type FileMeta struct {
	ID        int64     `db:"id" json:"id"`
	Filename  string    `db:"filename" json:"filename"`
	Size      int64     `db:"size" json:"size"`
	Mime      string    `db:"mime" json:"mime"`
	Checksum  string    `db:"checksum" json:"checksum"`
	Path      string    `db:"path" json:"path"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (f *FileMeta) IsImage() bool {
	switch f.Mime {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true
	}
	return false
}
