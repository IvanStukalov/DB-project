package handler

import (
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/forum"
	"github.com/IvanStukalov/DB_project/internal/utils"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
)

type Handler struct {
	uc forum.UseCase
}

func NewForumHandler(ForumUseCase forum.UseCase) *Handler {
	return &Handler{uc: ForumUseCase}
}

func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	newForum := models.Forum{}
	err := easyjson.UnmarshalFromReader(r.Body, &newForum)

	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	finalForum, err := h.uc.CreateForum(r.Context(), newForum)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "error"})
		return
	}
	if err == models.Conflict {
		utils.Response(w, http.StatusConflict, finalForum)
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusCreated, finalForum)
	return
}

func (h *Handler) GetForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid slug"})
		return
	}

	newForum := models.Forum{}
	newForum.Slug = slug

	finalForum, err := h.uc.GetForum(r.Context(), newForum)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "can`t find forum " + slug})
		return
	}
	utils.Response(w, http.StatusOK, finalForum)
	return
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	forumParam, found := vars["slug"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid slug"})
		return
	}

	newThread := models.Thread{}
	err := easyjson.UnmarshalFromReader(r.Body, &newThread)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	newThread.Forum = forumParam

	finalThread, err := h.uc.CreateThread(r.Context(), newThread)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "error"})
		return
	}
	if err == models.Conflict {
		utils.Response(w, http.StatusConflict, finalThread)
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusCreated, finalThread)
	return
}

func (h *Handler) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid slug"})
		return
	}

	queryParams := r.URL.Query()
	limit := queryParams.Get("limit")
	desc := queryParams.Get("desc")
	since := queryParams.Get("since")

	finalThreads, err := h.uc.GetThreadByForumSlug(r.Context(), slug, limit, since, desc)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "can`t find threads " + slug})
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, finalThreads)
	return
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid slug"})
		return
	}

	queryParams := r.URL.Query()
	limit := queryParams.Get("limit")
	desc := queryParams.Get("desc")
	since := queryParams.Get("since")

	users, err := h.uc.GetUsers(r.Context(), slug, limit, since, desc)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "can`t find user from forum " + slug})
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, users)
	return
}
