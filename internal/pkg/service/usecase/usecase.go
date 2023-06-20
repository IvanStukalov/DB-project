package usecase

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/service"
)

type UseCase struct {
	repo service.Repository
}

func NewRepoUsecase(repo service.Repository) service.UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Status(ctx context.Context) (models.Status, error) {
	return u.repo.Status(ctx)
}

func (u *UseCase) Clear(ctx context.Context) {
	u.repo.Clear(ctx)
}
