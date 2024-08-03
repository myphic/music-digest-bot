package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type DigestRepository interface {
	GetByName(ctx *context.Context, name string) (DigestModel, error)
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
