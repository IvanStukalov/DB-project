package delivery

import (
	"github.com/IvanStukalov/DB_project/internal/pkg/post"
)

type Handler struct {
	uc post.UseCase
}

func NewPostHandler(PostUseCase post.UseCase) *Handler {
	return &Handler{uc: PostUseCase}
}
