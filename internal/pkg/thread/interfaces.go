package thread

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
)

type UseCase interface {
	UpdateThread(ctx context.Context, slugOrId string, thread models.Thread) (models.Thread, error)
	GetThread(ctx context.Context, slugOrId string) (models.Thread, error)
	CreatePosts(ctx context.Context, slugOrId string, posts []models.Post) ([]models.Post, error)
	CreateVote(ctx context.Context, slugOrId string, vote models.Vote) (models.Thread, error)
	GetPosts(ctx context.Context, slugOrId string, sort string, limit string, since string, desc string) ([]models.Post, error)
}

type Repository interface {
	GetForumByThread(ctx context.Context, id int) (string, error)
	UpdateThread(ctx context.Context, slugOrId string, thread models.Thread) (models.Thread, error)
	GetThread(ctx context.Context, slugOrId string) (models.Thread, error)
	CreatePosts(ctx context.Context, thread int, forum string, posts []models.Post) ([]models.Post, error)
	CreateVote(ctx context.Context, thread int, vote models.Vote) error
	ChangeVote(ctx context.Context, thread int, vote models.Vote) error
	GetPostsFlat(ctx context.Context, thread int, limit string, since string, desc string) ([]models.Post, error)
	GetPostsTree(ctx context.Context, thread int, limit string, since string, desc string) ([]models.Post, error)
	GetPostsParentTree(ctx context.Context, thread int, limit string, since string, desc string) ([]models.Post, error)
}
