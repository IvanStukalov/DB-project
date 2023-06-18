package delivery

import (
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/user"
	"github.com/IvanStukalov/DB_project/internal/utils"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
)

type Handler struct {
	uc user.UseCase
}

func NewUserHandler(UserUseCase user.UseCase) *Handler {
	return &Handler{uc: UserUseCase}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, found := vars["nickname"]
	if !found {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	userS := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &userS)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	userS.NickName = nickname

	finalUser, err := h.uc.CreateUser(r.Context(), userS)
	if err != nil {
		utils.Response(w, http.StatusConflict, finalUser)
		return
	}
	newU := finalUser[0]
	utils.Response(w, http.StatusCreated, newU) // TODO сразу прокинуть finalUser[0]
	return
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, found := vars["nickname"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid nickname"})
		return
	}

	userS := models.User{}
	userS.NickName = nickname

	finalUser, err := h.uc.GetUser(r.Context(), userS)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "can`t find user " + nickname})
		return
	}
	utils.Response(w, http.StatusOK, finalUser)
	return
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, found := vars["nickname"]
	if !found {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	newUser := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &newUser)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	newUser.NickName = nickname

	finalUser, err := h.uc.UpdateUser(r.Context(), newUser)
	if err == models.Conflict {
		utils.Response(w, http.StatusConflict, models.ErrMsg{Msg: "error"})
		return
	}
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "error"})
		return
	}

	utils.Response(w, http.StatusOK, finalUser[0])
	return
}
