package usecase

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/forum"
	"github.com/IvanStukalov/DB_project/internal/pkg/thread"
	"github.com/IvanStukalov/DB_project/internal/pkg/user"
)

type UseCase struct {
	repo  forum.Repository
	uRepo user.Repository
	tRepo thread.Repository
}

func NewRepoUsecase(repo forum.Repository, uRepo user.Repository, tRepo thread.Repository) forum.UseCase {
	return &UseCase{repo: repo, uRepo: uRepo, tRepo: tRepo}
}

func (u *UseCase) CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error) {
	creator, err := u.uRepo.GetUser(ctx, forum.User)
	if err != nil {
		return forum, models.NotFound
	}
	forum.User = creator.NickName

	createdForum, err := u.repo.CreateForum(ctx, forum)
	if err != nil {
		foundForum, foundError := u.repo.GetForum(ctx, forum.Slug)
		if foundError == nil {
			return foundForum, models.Conflict
		}
		return createdForum, models.InternalError
	}
	return createdForum, nil
}

func (u *UseCase) GetForum(ctx context.Context, forum models.Forum) (models.Forum, error) {
	return u.repo.GetForum(ctx, forum.Slug)
}

func (u *UseCase) CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error) {
	creator, err := u.uRepo.GetUser(ctx, thread.Author)
	if err != nil {
		return thread, models.NotFound
	}
	thread.Author = creator.NickName

	foundForum, err := u.repo.GetForum(ctx, thread.Forum)
	if err != nil {
		return thread, models.NotFound
	}
	thread.Forum = foundForum.Slug

	if thread.Slug != "" {
		foundThread, err := u.tRepo.GetThread(ctx, thread.Slug)
		if err == nil {
			return foundThread, models.Conflict
		}
	}

	createdThread, err := u.repo.CreateThread(ctx, thread)
	if err != nil {
		return createdThread, err
	}
	return createdThread, nil
}

func (u *UseCase) GetThreadByForumSlug(ctx context.Context, slug string, limit string, since string, desc string) ([]models.Thread, error) {
	_, err := u.repo.GetForum(ctx, slug)
	if err != nil {
		return nil, models.NotFound
	}
	return u.repo.GetThreadByForumSlug(ctx, slug, limit, since, desc)
}

func (u *UseCase) GetUsers(ctx context.Context, slug string, limit string, since string, desc string) ([]models.User, error) {
	_, err := u.repo.GetForum(ctx, slug)
	if err == models.NotFound {
		return []models.User{}, err
	}

	return u.repo.GetUsers(ctx, slug, limit, since, desc)
}
