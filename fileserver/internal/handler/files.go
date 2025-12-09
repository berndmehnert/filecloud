package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/filecloud/model"

	"github.com/jmoiron/sqlx"
)

var ErrInvalidCursor = errors.New("invalid cursor")

func makeCursor(lastKey, direction string, limit int, secret []byte) (string, error) {
	tok := model.CursorToken{
		LastKey:   lastKey,
		Direction: direction,
		Limit:     limit,
		Exp:       time.Now().Add(1 * time.Hour).Unix(),
	}
	b, err := json.Marshal(tok)
	if err != nil {
		return "", err
	}
	sig := hmacSHA256(b, secret)
	return base64URLEncode(b) + "." + base64URLEncode(sig), nil
}

func decodeCursor(s string, secret []byte) (*model.CursorToken, error) {
	parts := strings.SplitN(s, ".", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidCursor
	}
	b, err := base64URLDecode(parts[0])
	if err != nil {
		return nil, ErrInvalidCursor
	}
	sig, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, ErrInvalidCursor
	}
	if !hmac.Equal(sig, hmacSHA256(b, secret)) {
		return nil, ErrInvalidCursor
	}
	var tok model.CursorToken
	if err := json.Unmarshal(b, &tok); err != nil {
		return nil, ErrInvalidCursor
	}
	if tok.Exp != 0 && time.Now().Unix() > tok.Exp {
		return nil, ErrInvalidCursor
	}
	return &tok, nil
}

func hmacSHA256(b, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(b)
	return mac.Sum(nil)
}

func base64URLEncode(b []byte) string {
	return strings.TrimRight(base64.RawURLEncoding.EncodeToString(b), "=")
}

func base64URLDecode(s string) ([]byte, error) {
	// add padding if required
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	return base64.RawURLEncoding.DecodeString(s)
}

// helpers
func parseLimit(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return def
	}
	if n > 500 {
		return 500
	}
	return n
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

func HandleListFiles(w http.ResponseWriter, r *http.Request, db *sqlx.DB, secret []byte) {
	ctx := r.Context()
	q := r.URL.Query().Get("q") // optional filter substring
	limit := parseLimit(r.URL.Query().Get("limit"), 50)
	cursor := r.URL.Query().Get("cursor")

	var lastCreated time.Time
	var lastID int64
	if cursor != "" {
		tok, err := decodeCursor(cursor, secret)
		if err != nil {
			http.Error(w, "invalid cursor", http.StatusBadRequest)
			return
		}
		parts := strings.SplitN(tok.LastKey, "|", 2)
		if len(parts) != 2 {
			http.Error(w, "invalid cursor key", http.StatusBadRequest)
			return
		}
		lt, err := time.Parse(time.RFC3339Nano, parts[0])
		if err != nil {
			http.Error(w, "invalid cursor time", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			http.Error(w, "invalid cursor id", http.StatusBadRequest)
			return
		}
		lastCreated = lt
		lastID = id
	}

	sqlStr := `
SELECT id, filename, size, mime, checksum, path, created_at
FROM files
WHERE (:q = '' OR LOWER(filename) LIKE '%' || LOWER(:q) || '%')`
	args := map[string]interface{}{"q": q, "limit": limit + 1}

	if cursor == "" {
		sqlStr += " ORDER BY created_at DESC, id DESC LIMIT :limit"
	} else {
		sqlStr += ` AND (created_at, id) < (:lastCreated, :lastID)
ORDER BY created_at DESC, id DESC LIMIT :limit`
		args["lastCreated"] = lastCreated
		args["lastID"] = lastID
	}

	stmt, err := db.PrepareNamedContext(ctx, sqlStr)
	if err != nil {
		http.Error(w, "db prepare: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var rows []model.FileMeta
	if err := stmt.SelectContext(ctx, &rows, args); err != nil {
		http.Error(w, "db select: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	var nextCursor string
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		lastKey := last.CreatedAt.UTC().Format(time.RFC3339Nano) + "|" + strconv.FormatInt(last.ID, 10)
		c, err := makeCursor(lastKey, "forward", limit, secret)
		if err == nil {
			nextCursor = c
		}
	}

	resp := map[string]interface{}{
		"items":      rows,
		"nextCursor": nextCursor,
		"hasMore":    hasMore,
	}
	writeJSON(w, resp)
}
