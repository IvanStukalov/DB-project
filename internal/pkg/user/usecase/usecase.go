package usecase

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/user"
)

type UseCase struct {
	repo user.Repository
}

func NewRepoUsecase(repo user.Repository) user.UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) GetUser(ctx context.Context, user models.User) (models.User, error) {
	return u.repo.GetUser(ctx, user.NickName)
}

func (u *UseCase) CreateUser(ctx context.Context, user models.User) ([]models.User, error) {
	foundUsers, _ := u.repo.IsEmailOrNicknameUniq(ctx, user)
	if len(foundUsers) > 0 {
		return foundUsers, models.Conflict
	}
	users := make([]models.User, 0)
	newUser, _ := u.repo.CreateUser(ctx, user)
	users = append(users, newUser)
	return users, nil
}

func (u *UseCase) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	foundUsers, err := u.repo.IsEmailUniq(ctx, user)
	if err == models.Conflict {
		return foundUsers, models.Conflict
	}

	updatedUser, err := u.repo.UpdateUser(ctx, user)
	if err != nil {
		return models.User{}, models.NotFound
	}
	return updatedUser, nil
}
