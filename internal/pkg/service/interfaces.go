package service

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
)

type UseCase interface {
	Status(ctx context.Context) (models.Status, error)
	Clear(ctx context.Context)
}

type Repository interface {
	Status(ctx context.Context) (models.Status, error)
	Clear(ctx context.Context)
}
