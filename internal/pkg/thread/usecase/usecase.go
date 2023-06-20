package usecase

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/post"
	"github.com/IvanStukalov/DB_project/internal/pkg/thread"
	"github.com/IvanStukalov/DB_project/internal/pkg/user"
)

type UseCase struct {
	repo  thread.Repository
	pRepo post.Repository
	uRepo user.Repository
}

func NewRepoUsecase(repo thread.Repository, pRepo post.Repository, uRepo user.Repository) thread.UseCase {
	return &UseCase{repo: repo, pRepo: pRepo, uRepo: uRepo}
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

	if len(posts) == 0 {
		return []models.Post{}, nil
	}

	createdPosts, err := u.repo.CreatePosts(ctx, foundThread.ID, foundForum, posts)
	if err != nil {
		return createdPosts, err
	}
	return createdPosts, nil
}

func (u *UseCase) CreateVote(ctx context.Context, slugOrId string, vote models.Vote) (models.Thread, error) {
	_, err := u.uRepo.GetUser(ctx, vote.Nickname)
	if err == models.NotFound {
		return models.Thread{}, models.NotFound
	}

	foundThread, err := u.repo.GetThread(ctx, slugOrId)
	if err != nil {
		return models.Thread{}, models.NotFound
	}

	err = u.repo.CreateVote(ctx, foundThread.ID, vote)
	if err == models.Conflict {
		errUpdate := u.repo.ChangeVote(ctx, foundThread.ID, vote)
		if errUpdate != nil {
			return foundThread, models.InternalError
		}
	}
	if err == models.InternalError {
		return foundThread, models.InternalError
	}

	foundThread, err = u.repo.GetThread(ctx, slugOrId)
	if err != nil {
		return foundThread, models.NotFound
	}
	return foundThread, nil
}

func (u *UseCase) GetPosts(ctx context.Context, slugOrId string, sort string, limit string, since string, desc string) ([]models.Post, error) {
	foundThread, err := u.repo.GetThread(ctx, slugOrId)
	if err != nil {
		return []models.Post{}, models.NotFound
	}

	switch sort {
	case "flat":
		return u.repo.GetPostsFlat(ctx, foundThread.ID, limit, since, desc)
	case "tree":
		return u.repo.GetPostsTree(ctx, foundThread.ID, limit, since, desc)
	case "parent_tree":
		return u.repo.GetPostsParentTree(ctx, foundThread.ID, limit, since, desc)
	default:
		return u.repo.GetPostsFlat(ctx, foundThread.ID, limit, since, desc)
	}
}
