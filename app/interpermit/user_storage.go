package interpermit

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *storage {
	return &storage{db: db}
}

func (s *storage) SaveUser(ctx context.Context, data GetUserByTokenResponse) error {
	_, err := s.db.Exec(ctx, "insert into users(cif_no, cid) values($1, $2)", data.Data.CifNo, data.Data.CitizenID)
	return err
}
