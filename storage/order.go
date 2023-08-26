package storage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"telegram-bot/models"
)

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *orderRepo {
	return &orderRepo{
		db: db,
	}
}

func (r *orderRepo) Create(ctx context.Context, order models.Order) (models.Order, error) {
	query := `
		INSERT INTO orders(
			user_id,
			lat,
			long,
			block,
			floor,
			flat,
			photo,
			weight,
			comment,
			status 
			) VALUES (:user_id, :lat, :long, :block, :floor, :flat, :photo, :weight, :comment, :status)`

	rows, err := r.db.NamedQueryContext(ctx, query, order)
	if err != nil {
		return models.Order{}, err
	}

	var result models.Order
	for rows.Next() {
		if err := rows.StructScan(&result); err != nil {
			return models.Order{}, err
		}
	}

	return result, nil
}

func (r *orderRepo) Get(ctx context.Context, order models.Order) ([]models.Order, error) {
	var (
		filter string
		args   = make(map[string]interface{})
	)

	query := `
		SELECT
			id,
			user_id,
			lat,
			long,
			block,
			floor,
			flat,
			photo,
			weight,
			comment,
			status,
			updated_at,
			created_at
		FROM orders
		WHERE id = $1`

	if order.ID != "" {
		filter += " AND id = :id"
		args["id"] = order.ID
	}
	if order.UserID != "" {
		filter += " AND user_id = :user_id"
		args["user_id"] = order.UserID
	}
	if order.Lat != "" {
		filter += " AND lat = :lat"
		args["lat"] = order.Lat
	}
	if order.Long != "" {
		filter += " AND long = :long"
		args["long"] = order.Long
	}
	if order.Block != "" {
		filter += " AND block = :block"
		args["block"] = order.Block
	}
	if order.Floor != "" {
		filter += " AND floor = :floor"
		args["floor"] = order.Floor
	}
	if order.Flat != "" {
		filter += " AND flat = :flat"
		args["flat"] = order.Flat
	}
	if order.Photo != "" {
		filter += " AND photo = :photo"
		args["photo"] = order.Photo
	}
	if order.Weight != "" {
		filter += " AND weight = :weight"
		args["weight"] = order.Weight
	}
	if order.Comment != "" {
		filter += " AND comment = :comment"
		args["comment"] = order.Comment
	}
	if order.Status != "" {
		filter += " AND status = :status"
		args["status"] = order.Status
	}

	var result []models.Order

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.StructScan(&order); err != nil {
			return nil, err
		}

		result = append(result, order)
	}

	return result, nil
}

func (r *orderRepo) Update(ctx context.Context, order models.Order) (models.Order, error) {
	return models.Order{}, nil
}
