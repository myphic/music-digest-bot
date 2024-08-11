package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type DigestRepository interface {
	GetByID(ctx context.Context, ID int) (DigestModel, error)
	CreateAndGetID(ctx context.Context, digest DigestModel) (int, error)
}

type DigestRepositoryImpl struct {
	conn *pgx.Conn
}

func NewDigestRepository(conn *pgx.Conn) *DigestRepositoryImpl {
	return &DigestRepositoryImpl{conn: conn}
}

func (r *DigestRepositoryImpl) GetByID(ctx context.Context, ID int) (DigestModel, error) {
	var Digest DigestModel
	err := r.conn.QueryRow(ctx, "SELECT * FROM digest WHERE id = $1", ID).Scan(Digest)

	if err != nil {
		return DigestModel{}, err
	}

	return Digest, nil
}

func (r *DigestRepositoryImpl) CreateAndGetID(ctx context.Context, digest DigestModel) (int, error) {
	query := "INSERT INTO digest (id, source_id, title, description, published_at, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	var id int
	err := r.conn.QueryRow(ctx, query, digest.ID, digest.SourceID, digest.Title, digest.Description, digest.PublishedAt, digest.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
