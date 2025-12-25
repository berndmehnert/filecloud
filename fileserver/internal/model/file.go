package model

import "time"

type FileMeta struct {
	ID              int64     `db:"id" json:"id"`
	Filename        string    `db:"filename" json:"filename"`
	Size            int64     `db:"size" json:"size"`
	Mime            string    `db:"mime" json:"mime"`
	Checksum        string    `db:"checksum" json:"checksum"`
	Path            string    `db:"path" json:"-"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	ThumbnailPath   *string   `db:"thumbnail_path" json:"-"`
	ThumbnailStatus string    `db:"thumbnail_status" json:"thumbnail_status"`
}

func (f *FileMeta) IsImage() bool {
	switch f.Mime {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true
	}
	return false
}
