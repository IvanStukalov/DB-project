package usecase

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/thread"
)

type UseCase struct {
	repo thread.Repository
}

func NewRepoUsecase(repo thread.Repository) thread.UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) UpdateThread(ctx context.Context, slugOrId string, thread models.Thread) (models.Thread, error) {
	updatedThread, err := u.repo.UpdateThread(ctx, slugOrId, thread)
	if err != nil {
		return updatedThread, err
	}
	return updatedThread, nil
}

func (u *UseCase) GetThread(ctx context.Context, slugOrId string) (models.Thread, error) {
	return u.repo.GetThread(ctx, slugOrId)
}

func (u *UseCase) CreatePosts(ctx context.Context, slugOrId string, posts []models.Post) ([]models.Post, error) {
	foundThread, err := u.repo.GetThread(ctx, slugOrId)
	if err != nil {
		return posts, models.NotFound
	}

	foundForum, err := u.repo.GetForumByThread(ctx, foundThread.ID)
	if err != nil {
		return posts, models.NotFound
	}

	createdPosts, err := u.repo.CreatePosts(ctx, foundThread.ID, foundForum, posts)
	if err != nil {
		return createdPosts, err
	}
	return createdPosts, nil
}

func (u *UseCase) CreateVote(ctx context.Context, slugOrId string, vote models.Vote) (models.Thread, error) {
	thread, err := u.repo.GetThread(ctx, slugOrId)
	if err != nil {
		return models.Thread{}, models.NotFound
	}

	err = u.repo.CreateVote(ctx, thread.ID, vote)
	if err == models.Conflict {
		errUpdate := u.repo.ChangeVote(ctx, thread.ID, vote)
		if errUpdate != nil {
			return thread, models.InternalError
		}
	}
	if err == models.InternalError {
		return thread, models.InternalError
	}

	thread, err = u.repo.GetThread(ctx, slugOrId)
	if err != nil {
		return thread, models.NotFound
	}
	return thread, nil
}
