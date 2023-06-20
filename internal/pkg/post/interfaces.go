package post

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
)

type UseCase interface {
	GetPost(ctx context.Context, id string, related string) (models.WrappedPost, error)
	UpdatePost(ctx context.Context, post models.Post) (models.Post, error)
}

type Repository interface {
	GetPost(ctx context.Context, id int) (models.Post, error)
	UpdatePost(ctx context.Context, post models.Post) (models.Post, error)
}
