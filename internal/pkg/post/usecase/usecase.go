package usecase

import (
	"github.com/IvanStukalov/DB_project/internal/pkg/post"
)

type UseCase struct {
	repo post.Repository
}

func NewRepoUsecase(repo post.Repository) post.UseCase {
	return &UseCase{repo: repo}
}
