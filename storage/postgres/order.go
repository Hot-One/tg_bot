package postgres

import (
	"context"
	"database/sql"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"telegram-bot/models"
	"telegram-bot/pkg/helper"
)

type OrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) Create(ctx context.Context, req *models.OrderCreate) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO orders(id, name, phone, lat, long, address, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.Phone,
		req.Lat,
		req.Long,
		req.Address,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *OrderRepo) GetByID(ctx context.Context, req *models.OrderPrimaryKey) (*models.Order, error) {

	var (
		query string

		id      sql.NullString
		name    sql.NullString
		phone   sql.NullString
		lat     sql.NullString
		long    sql.NullString
		address sql.NullString
	)

	query = `
		SELECT
			id,
			name,
			phone,
			lat,
			long,
			address
		FROM orders
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&phone,
		&lat,
		&long,
		&address,
	)

	if err != nil {
		return nil, err
	}

	return &models.Order{
		Id:      id.String,
		Name:    name.String,
		Phone:   phone.String,
		Lat:     lat.String,
		Long:    long.String,
		Address: address.String,
	}, nil
}

// func (r *OrderRepo) GetList(ctx context.Context, req *models.OrderGetListRequest) (*models.OrderGetListResponse, error) {

// 	var (
// 		resp    = &models.OrderGetListResponse{}
// 		query   string
// 		where   = " WHERE deleted = false"
// 		offset  = " OFFSET 0"
// 		limit   = " LIMIT 10"
// 		ordered = " ORDER BY created_at desc"
// 	)

// 	query = `
// 		SELECT
// 			COUNT(*) OVER(),
// 			id,
// 			name,
// 			address,
// 			created_at,
// 			updated_at,
// 			deleted,
// 			deleted_at
// 		FROM branch
// 	`

// 	if req.Offset > 0 {
// 		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
// 	}

// 	if req.Limit > 0 {
// 		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
// 	}

// 	if req.Search != "" {
// 		where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'`
// 	}

// 	if req.SearchByAddress != "" {
// 		where += ` AND address ILIKE '%' || '` + req.Search + `' || '%'`
// 	}

// 	query += where + ordered + offset + limit
// 	fmt.Println(query)

// 	rows, err := r.db.Query(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var (
// 			id        sql.NullString
// 			name      sql.NullString
// 			address   sql.NullString
// 			createdAt sql.NullString
// 			updatedAt sql.NullString
// 			deleted   sql.NullBool
// 			deletedAt sql.NullString
// 		)

// 		err := rows.Scan(
// 			&resp.Count,
// 			&id,
// 			&name,
// 			&address,
// 			&createdAt,
// 			&updatedAt,
// 			&deleted,
// 			&deletedAt,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		resp.Branches = append(resp.Branches, &models.Order{
// 			Id:        id.String,
// 			Name:      name.String,
// 			Address:   address.String,
// 			CreatedAt: createdAt.String,
// 			UpdatedAt: updatedAt.String,
// 			Deleted:   deleted.Bool,
// 			DeletedAt: deletedAt.String,
// 		})
// 	}

// 	return resp, nil
// }

func (r *OrderRepo) Update(ctx context.Context, req *models.OrderUpdate) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			orders
		SET
			name = :name,
			phone = :phone,
			lat = :lat,
			long = :long,
			address = :address,
			photo = :photo,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":      req.Id,
		"name":    req.Name,
		"phone":   req.Phone,
		"lat":     req.Lat,
		"long":    req.Long,
		"address": req.Address,
		"photo":   req.Photo,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *OrderRepo) Delete(ctx context.Context, req *models.OrderPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM orders WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil

}
