package handler

import (
	"encoding/json"
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

func (h *Handler) GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid slug or id"})
		return
	}

	finalThread, err := h.uc.GetThread(r.Context(), slug)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "can`t find thread " + slug})
		return
	}
	utils.Response(w, http.StatusOK, finalThread)
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

func (h *Handler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "invalid slug"})
		return
	}

	var newPosts []models.Post
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newPosts)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if len(newPosts) == 0 {
		utils.Response(w, http.StatusCreated, []models.Post{})
		return
	}

	finalPosts, err := h.uc.CreatePosts(r.Context(), slugOrId, newPosts)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, models.ErrMsg{Msg: "not found thread " + slugOrId})
		return
	}
	if err == models.Conflict {
		utils.Response(w, http.StatusConflict, finalPosts)
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusCreated, finalPosts)
	return
}
