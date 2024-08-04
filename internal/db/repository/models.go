package repository

import "time"

type SourceModel struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Meta      string    `db:"meta"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type DigestModel struct {
	ID          int       `db:"id"`
	SourceID    int       `db:"source_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
}
