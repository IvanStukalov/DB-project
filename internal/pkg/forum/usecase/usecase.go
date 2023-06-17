package usecase

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/forum"
)

type UseCase struct {
	repo forum.Repository
}

func NewRepoUsecase(repo forum.Repository) forum.UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) GetUser(ctx context.Context, user models.User) (models.User, error) {
	return u.repo.GetUser(ctx, user.NickName)
}

func (u *UseCase) CreateUser(ctx context.Context, user models.User) ([]models.User, error) {
	chUsers, _ := u.repo.CheckUserEmailOrNicknameUniq(user)
	if len(chUsers) > 0 {
		return chUsers, models.Conflict
	}
	usersS := make([]models.User, 0)
	cUser, _ := u.repo.CreateUser(ctx, user)
	usersS = append(usersS, cUser)
	return usersS, nil
}

func (u *UseCase) UpdateUser(ctx context.Context, user models.User) ([]models.User, error) {
	chUsers, err := u.repo.CheckUserEmailUniq(user)
	if err == models.Conflict {
		return chUsers, models.Conflict
	}

	usersS := make([]models.User, 0)
	updatedUser, err := u.repo.UpdateUser(ctx, user)
	usersS = append(usersS, updatedUser)
	if err != nil {
		return nil, models.NotFound
	}
	return usersS, nil
}

func (u *UseCase) CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error) {
	creator, err := u.repo.GetUser(ctx, forum.User)
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
	creator, err := u.repo.GetUser(ctx, thread.Author)
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
		foundThread, err := u.repo.GetThread(ctx, thread.Slug)
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

func (u *UseCase) GetThread(ctx context.Context, slugOrId string) (models.Thread, error) {
	return u.repo.GetThread(ctx, slugOrId)
}

func (u *UseCase) GetThreadByForumSlug(ctx context.Context, slug string, limit string, since string, desc string) ([]models.Thread, error) {
	_, err := u.repo.GetForum(ctx, slug)
	if err != nil {
		return nil, models.NotFound
	}
	return u.repo.GetThreadByForumSlug(ctx, slug, limit, since, desc)
}
