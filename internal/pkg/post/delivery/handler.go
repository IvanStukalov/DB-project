package delivery

import (
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/post"
	"github.com/IvanStukalov/DB_project/internal/utils"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
)

type Handler struct {
	uc post.UseCase
}

func NewPostHandler(PostUseCase post.UseCase) *Handler {
	return &Handler{uc: PostUseCase}
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid slug"})
		return
	}

	queryParams := r.URL.Query()
	related := queryParams.Get("related")

	foundPost, err := h.uc.GetPost(r.Context(), id, related)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "can`t find user " + id})
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, foundPost)
	return
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid slug"})
		return
	}

	integerId, err := strconv.Atoi(id)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	updatingPost := models.Post{}
	err = easyjson.UnmarshalFromReader(r.Body, &updatingPost)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	updatingPost.ID = integerId

	updatedPost, err := h.uc.UpdatePost(r.Context(), updatingPost)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "post not found"})
		return
	}
	utils.Response(w, http.StatusOK, updatedPost)
	return
}
