package storage

import (
	"github.com/jmoiron/sqlx"
)

type store struct {
	db    *sqlx.DB
	order *orderRepo
	user  *userRepo
}

func New(db *sqlx.DB) (*store, error) {

	return &store{
		db: db,
	}, nil
}

func (s *store) CloseDb() {
	s.db.Close()
}

func (s *store) Order() *orderRepo {
	if s.order == nil {
		s.order = NewOrderRepo(s.db)
	}

	return s.order
}

func (s *store) User() *userRepo {
	if s.order == nil {
		s.user = NewUserRepo(s.db)
	}

	return s.user
}
