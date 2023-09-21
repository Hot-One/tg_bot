package storage

import (
	"context"
	"telegram-bot/models"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	query := `
		INSERT INTO users(
			first_name,
			last_name,
			phone_number,
			telegram_user_id,
			status,
			user_name) 
			VALUES (:first_name, :last_name, :phone_number, :telegram_user_id, :status, :user_name)`

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return models.User{}, err
	}

	var result models.User
	for rows.Next() {
		if err := rows.StructScan(&result); err != nil {
			return models.User{}, err
		}
	}

	return result, nil
}

func (r *userRepo) Get(ctx context.Context, order models.User) ([]models.User, error) {
	var (
		filter string
		args   = make(map[string]interface{})
	)
	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone_number,
			telegram_user_id,
			status,
			user_name,
			created_at,
			updated_at
		FROM users
		WHERE true`

	if order.ID != "" {
		filter += " AND id = :id"
		args["id"] = order.ID
	}
	if order.FirstName != "" {
		filter += " AND first_name = :first_name"
		args["first_name"] = order.FirstName
	}
	if order.LastName != "" {
		filter += " AND last_name = :last_name"
		args["last_name"] = order.LastName
	}
	if order.PhoneNumber != "" {
		filter += " AND phone_number = :phone_number"
		args["phone_number"] = order.PhoneNumber
	}
	if order.TelegramUserID != 0 {
		filter += " AND telegram_user_id = :telegram_user_id"
		args["telegram_user_id"] = order.TelegramUserID
	}
	if order.Status != "" {
		filter += " AND status = :status"
		args["status"] = order.Status
	}
	if order.UserName != "" {
		filter += " AND user_name = :user_name"
		args["user_name"] = order.UserName
	}

	var result []models.User

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, nil
}

func (r *userRepo) Update(ctx context.Context, user models.User) (models.User, error) {
	return models.User{}, nil
}
