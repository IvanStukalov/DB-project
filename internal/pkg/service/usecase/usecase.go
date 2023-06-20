package usecase

import (
	"github.com/IvanStukalov/DB_project/internal/pkg/service"
)

type UseCase struct {
	repo service.Repository
}

func NewRepoUsecase(repo service.Repository) service.UseCase {
	return &UseCase{repo: repo}
}
