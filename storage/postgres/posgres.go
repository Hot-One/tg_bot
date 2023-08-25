package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"telegram-bot/config"
	"telegram-bot/storage"
)

type store struct {
	db    *pgxpool.Pool
	order *OrderRepo
}

func NewConnectionPostgres(cfg *config.Config) (storage.StorageI, error) {

	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}
	connect.MaxConns = cfg.PostgresMaxConnection

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), connect)
	if err != nil {
		return nil, err
	}

	return &store{
		db: pgxpool,
	}, nil
}

func (s *store) CloseDb() {
	s.db.Close()
}

func (s *store) Order() storage.OrderRepoI {

	if s.order == nil {
		s.order = NewOrderRepo(s.db)
	}

	return s.order
}
