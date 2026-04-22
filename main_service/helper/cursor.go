package helper

import (
	"encoding/base64"
	"net/http"
	"strconv"
)

type CursorMeta struct {
	Total      int64  `json:"total"`
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

func EncodeCursor(id int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(id, 10)))
}

func DecodeCursor(s string) (int64, bool) {
	if s == "" {
		return 0, false
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return 0, false
	}
	id, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func ParseCursorPage(r *http.Request) (afterID int64, limit int) {
	cursor := r.URL.Query().Get("cursor")
	afterID, _ = DecodeCursor(cursor)

	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return
}

func NewCursorMeta(limit int, hasMore bool, lastID int64, total int64) CursorMeta {
	meta := CursorMeta{Total: total, Limit: limit, HasMore: hasMore}
	if hasMore {
		meta.NextCursor = EncodeCursor(lastID)
	}
	return meta
}
