package usecase

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/forum"
	"github.com/IvanStukalov/DB_project/internal/pkg/post"
	"github.com/IvanStukalov/DB_project/internal/pkg/thread"
	"github.com/IvanStukalov/DB_project/internal/pkg/user"
	"strconv"
	"strings"
)

type UseCase struct {
	repo  post.Repository
	fRepo forum.Repository
	uRepo user.Repository
	tRepo thread.Repository
}

func NewRepoUsecase(repo post.Repository, fRepo forum.Repository, uRepo user.Repository, tRepo thread.Repository) post.UseCase {
	return &UseCase{repo: repo, fRepo: fRepo, uRepo: uRepo, tRepo: tRepo}
}

func (u *UseCase) GetPost(ctx context.Context, id string, related string) (models.PostFull, error) {
	var foundAuthor models.User
	var foundForum models.Forum
	var foundThread models.Thread

	wrappedPost := models.PostFull{}

	integerId, err := strconv.Atoi(id)
	if err != nil {
		return wrappedPost, models.InternalError
	}

	foundPost, err := u.repo.GetPost(ctx, integerId)
	if err == models.NotFound {
		return wrappedPost, models.NotFound
	}
	wrappedPost.Post = foundPost

	isUser := strings.Contains(related, "user")
	isForum := strings.Contains(related, "forum")
	isThread := strings.Contains(related, "thread")

	if isUser {
		foundAuthor, err = u.uRepo.GetUser(ctx, foundPost.Author)
		if err == models.NotFound {
			return wrappedPost, models.NotFound
		}
		wrappedPost.Author = &foundAuthor
	}

	if isForum {
		foundForum, err = u.fRepo.GetForum(ctx, foundPost.Forum)
		if err == models.NotFound {
			return wrappedPost, models.NotFound
		}
		wrappedPost.Forum = &foundForum
	}

	if isThread {
		slugOrId := strconv.Itoa(foundPost.Thread)
		foundThread, err = u.tRepo.GetThread(ctx, slugOrId)
		if err == models.NotFound {
			return wrappedPost, models.NotFound
		}
		wrappedPost.Thread = &foundThread
	}

	return wrappedPost, nil
}

func (u *UseCase) UpdatePost(ctx context.Context, post models.Post) (models.Post, error) {
	var foundPost models.Post
	var updatedPost models.Post
	var err error

	foundPost, err = u.repo.GetPost(ctx, post.ID)
	if err != nil {
		return models.Post{}, models.NotFound
	}
	if post.Message == "" {
		return foundPost, nil
	}

	if foundPost.Message == post.Message {
		return foundPost, nil
	}

	updatedPost, err = u.repo.UpdatePost(ctx, post)
	if err == models.NotFound {
		return foundPost, models.NotFound
	}

	return updatedPost, nil
}
