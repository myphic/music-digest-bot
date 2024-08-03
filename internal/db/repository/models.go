package repository

import "time"

type SourceModel struct {
	ID        int
	Name      string
	Meta      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DigestModel struct {
	ID          int
	SourceID    int
	Title       string
	Description string
	PublishedAt time.Time
	CreatedAt   time.Time
}
