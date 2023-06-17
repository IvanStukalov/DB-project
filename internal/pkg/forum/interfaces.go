package forum

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
)

type UseCase interface {
	GetUser(ctx context.Context, user models.User) (models.User, error)
	CreateUser(ctx context.Context, user models.User) ([]models.User, error)
	UpdateUser(ctx context.Context, user models.User) ([]models.User, error)
	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	GetForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error)
	GetThread(ctx context.Context, slugOrId string) (models.Thread, error)
	GetThreadByForumSlug(ctx context.Context, slug string, limit string, since string, desc string) ([]models.Thread, error)
	CreatePosts(ctx context.Context, slugOrId string, posts []models.Post) ([]models.Post, error)
}

type Repository interface {
	GetUser(ctx context.Context, name string) (models.User, error)
	CheckUserEmailOrNicknameUniq(usersS models.User) ([]models.User, error)
	CheckUserEmailUniq(usersS models.User) ([]models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	GetForum(ctx context.Context, slug string) (models.Forum, error)
	GetForumByThread(ctx context.Context, id int) (string, error)
	CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error)
	GetThread(ctx context.Context, slugOrId string) (models.Thread, error)
	GetThreadByForumSlug(ctx context.Context, slug string, limit string, since string, desc string) ([]models.Thread, error)
	CreatePosts(ctx context.Context, thread int, forum string, posts []models.Post) ([]models.Post, error)
}
