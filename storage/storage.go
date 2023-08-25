package storage

import (
	"context"
	"telegram-bot/models"
)

type StorageI interface {
	CloseDb()
	Order() OrderRepoI
}

type OrderRepoI interface {
	Create(context.Context, *models.OrderCreate) (string, error)
	GetByID(context.Context, *models.OrderPrimaryKey) (*models.Order, error)
	Update(context.Context, *models.OrderUpdate) (int64, error)
	Delete(context.Context, *models.OrderPrimaryKey) error
}
