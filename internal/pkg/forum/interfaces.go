package forum

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
)

type UseCase interface {
	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	GetForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error)
	GetThreadByForumSlug(ctx context.Context, slug string, limit string, since string, desc string) ([]models.Thread, error)
}

type Repository interface {
	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	GetForum(ctx context.Context, slug string) (models.Forum, error)
	CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error)
	GetThreadByForumSlug(ctx context.Context, slug string, limit string, since string, desc string) ([]models.Thread, error)
}
