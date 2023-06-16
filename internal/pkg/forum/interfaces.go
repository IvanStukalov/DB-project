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
}

type Repository interface {
	GetUser(ctx context.Context, name string) (models.User, error)
	CheckUserEmailOrNicknameUniq(usersS models.User) ([]models.User, error)
	CheckUserEmailUniq(usersS models.User) ([]models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
}
