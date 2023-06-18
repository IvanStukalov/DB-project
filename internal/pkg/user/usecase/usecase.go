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
	chUsers, _ := u.repo.CheckUserEmailOrNicknameUniq(user)
	if len(chUsers) > 0 {
		return chUsers, models.Conflict
	}
	usersS := make([]models.User, 0)
	cUser, _ := u.repo.CreateUser(ctx, user)
	usersS = append(usersS, cUser)
	return usersS, nil
}

func (u *UseCase) UpdateUser(ctx context.Context, user models.User) ([]models.User, error) {
	chUsers, err := u.repo.CheckUserEmailUniq(user)
	if err == models.Conflict {
		return chUsers, models.Conflict
	}

	usersS := make([]models.User, 0)
	updatedUser, err := u.repo.UpdateUser(ctx, user)
	usersS = append(usersS, updatedUser)
	if err != nil {
		return nil, models.NotFound
	}
	return usersS, nil
}
