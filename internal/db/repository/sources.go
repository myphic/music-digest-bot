package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type SourcesRepository interface {
	GetByName(ctx *context.Context, name string) (SourceModel, error)
}

type SourcesRepositoryImpl struct {
	conn *pgx.Conn
}

func NewSourcesRepository(conn *pgx.Conn) *SourcesRepositoryImpl {
	return &SourcesRepositoryImpl{conn: conn}
}

func (r *SourcesRepositoryImpl) GetByName(ctx context.Context, name string) (SourceModel, error) {
	var Source SourceModel
	err := r.conn.QueryRow(ctx, "SELECT * FROM sources WHERE name = $1", name).Scan(Source)

	if err != nil {
		return SourceModel{}, err
	}

	return Source, nil
}
