package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SourcesRepository interface {
	GetByName(ctx context.Context, name string) (SourceModel, error)
	Sources(ctx context.Context) ([]SourceModel, error)
	Create(ctx context.Context, source SourceModel) (SourceModel, error)
}

type SourcesRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewSourcesRepository(pool *pgxpool.Pool) *SourcesRepositoryImpl {
	return &SourcesRepositoryImpl{pool: pool}
}

func (s *SourcesRepositoryImpl) GetByName(ctx context.Context, name string) (SourceModel, error) {
	var Source SourceModel
	err := s.pool.QueryRow(ctx, "SELECT id, name, meta, created_at, updated_at FROM sources WHERE name = $1", name).Scan(&Source.ID, &Source.Name, &Source.Meta, &Source.CreatedAt, &Source.UpdatedAt)

	if err != nil {
		return SourceModel{}, err
	}

	return Source, nil
}

func (s *SourcesRepositoryImpl) Sources(ctx context.Context) ([]SourceModel, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name FROM sources")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Sources []SourceModel
	for rows.Next() {
		var r SourceModel
		err := rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return nil, err
		}
		Sources = append(Sources, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return Sources, nil
}

func (s *SourcesRepositoryImpl) Create(ctx context.Context, source SourceModel) (SourceModel, error) {
	err := s.pool.QueryRow(ctx, "INSERT INTO sources (name) VALUES ($1) RETURNING id", source.Name).Scan(&source.ID)
	if err != nil {
		return SourceModel{}, err
	}
	return source, nil
}
