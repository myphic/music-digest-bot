package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type DigestRepository interface {
	GetByID(ctx context.Context, ID int) (DigestModel, error)
	CreateAndGetID(ctx context.Context, digest DigestModel) (int, error)
	AllNotPosted(ctx context.Context) ([]DigestModel, error)
	MarkAsPosted(ctx context.Context, article DigestModel) error
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
	query := "INSERT INTO digest (source_id, digest_id, title, description, genre, published_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	var id int
	err := r.conn.QueryRow(ctx, query, digest.SourceID, digest.DigestID, digest.Title, digest.Description, digest.Genre, digest.PublishedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *DigestRepositoryImpl) AllNotPosted(ctx context.Context) ([]DigestModel, error) {
	var digests []DigestModel

	rows, err := r.conn.Query(ctx, "SELECT id, title, description, genre FROM digest WHERE posted IS FALSE")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r DigestModel
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.Genre)
		if err != nil {
			return nil, err
		}
		digests = append(digests, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return digests, nil
}

func (r *DigestRepositoryImpl) MarkAsPosted(ctx context.Context, digest DigestModel) error {
	query := "UPDATE digest SET posted=TRUE, updated_at=now() WHERE id = $1"

	_, err := r.conn.Exec(ctx, query, digest.ID)
	if err != nil {
		return err
	}
	return nil
}
