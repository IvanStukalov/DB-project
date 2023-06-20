package delivery

import (
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/service"
	"github.com/IvanStukalov/DB_project/internal/utils"
	"net/http"
)

type Handler struct {
	uc service.UseCase
}

func NewServiceHandler(ServiceUseCase service.UseCase) *Handler {
	return &Handler{uc: ServiceUseCase}
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	status, err := h.uc.Status(r.Context())

	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, status)
		return
	}

	utils.Response(w, http.StatusOK, status)
	return
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	h.uc.Clear(r.Context())
	utils.Response(w, http.StatusOK, nil)
	return
}
