package model

import "time"

type ShortUrl struct {
	URL           string     `json:"url,omitempty"`
	ShortUrl      string     `json:"short_url,omitempty"`
	TotalRedirect string     `json:"total_redirect,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	Updated_at    *time.Time `json:"updated_at,omitempty"`
}
