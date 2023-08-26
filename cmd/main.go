package main

import (
	"context"
	"fmt"

	"telegram-bot/config"
	"telegram-bot/storage"
	"telegram-bot/usecases/bot_initer"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	cfg *config.Config
}

func NewStore(cfg *config.Config) *Store {
	return &Store{
		cfg: cfg,
	}
}

func main() {
	cfg := config.Load()

	conn, err := sqlx.ConnectContext(context.Background(), "postgres", fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))
	if err != nil {
		panic(err)
	}

	conn.SetMaxIdleConns(cfg.PostgresMaxConnection)

	strg, err := storage.New(conn)
	if err != nil {
		panic(err)
	}

	botIninter := bot_initer.New(cfg, strg.User(), strg.Order())
	if err = botIninter.Execute(context.Background()); err != nil {
		panic(err)
	}

}
