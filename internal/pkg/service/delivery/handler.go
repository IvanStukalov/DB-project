package delivery

import (
	"github.com/IvanStukalov/DB_project/internal/pkg/service"
)

type Handler struct {
	uc service.UseCase
}

func NewServiceHandler(ServiceUseCase service.UseCase) *Handler {
	return &Handler{uc: ServiceUseCase}
}

//func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {

//}
