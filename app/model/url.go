package model

import "time"

type ShortUrl struct {
	URL           string     `json:"url"`
	ShortUrl      string     `json:"short_url,omitempty"`
	TotalRedirect int64      `json:"total_redirect,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	Updated_at    *time.Time `json:"updated_at,omitempty"`
}
