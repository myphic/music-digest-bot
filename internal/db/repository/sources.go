package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

type SourcesRepository interface {
	GetByName(ctx context.Context, name string) (SourceModel, error)
	Sources(ctx context.Context) ([]SourceModel, error)
}

type SourcesRepositoryImpl struct {
	conn *pgx.Conn
}

func NewSourcesRepository(conn *pgx.Conn) *SourcesRepositoryImpl {
	return &SourcesRepositoryImpl{conn: conn}
}

func (s *SourcesRepositoryImpl) GetByName(ctx context.Context, name string) (SourceModel, error) {
	var Source SourceModel
	err := s.conn.QueryRow(ctx, "SELECT * FROM sources WHERE name = $1", name).Scan(Source)

	if err != nil {
		return SourceModel{}, err
	}

	return Source, nil
}

func (s *SourcesRepositoryImpl) Sources(ctx context.Context) ([]SourceModel, error) {
	rows, err := s.conn.Query(ctx, "SELECT id, name FROM sources")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Sources []SourceModel
	for rows.Next() {
		var r SourceModel
		err := rows.Scan(&r.ID, &r.Name)
		if err != nil {
			log.Fatal(err)
		}
		Sources = append(Sources, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return Sources, nil
}
