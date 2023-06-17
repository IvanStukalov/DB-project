package repo

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/forum"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) forum.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) GetUser(ctx context.Context, name string) (models.User, error) {
	var userM models.User
	const SelectUserByNickname = "SELECT nickname, fullname, about, email FROM users WHERE nickname=$1 LIMIT 1;"
	row := r.Conn.QueryRow(ctx, SelectUserByNickname, name)
	err := row.Scan(&userM.NickName, &userM.FullName, &userM.About, &userM.Email)
	if err != nil {
		return models.User{}, models.NotFound
	}
	return userM, nil
}

func (r *repoPostgres) CheckUserEmailOrNicknameUniq(usersS models.User) ([]models.User, error) {
	const SelectUserByEmailOrNickname = "SELECT nickname, fullname, about, email FROM users WHERE nickname=$1 OR email=$2 LIMIT 2;"
	rows, err := r.Conn.Query(context.Background(), SelectUserByEmailOrNickname, usersS.NickName, usersS.Email)
	defer rows.Close()
	if err != nil {
		return []models.User{}, models.InternalError
	}
	users := make([]models.User, 0)
	for rows.Next() {
		userOne := models.User{}
		err := rows.Scan(&userOne.NickName, &userOne.FullName, &userOne.About, &userOne.Email)
		if err != nil {
			return []models.User{}, models.InternalError
		}
		users = append(users, userOne)
	}
	return users, nil
}

func (r *repoPostgres) CheckUserEmailUniq(usersS models.User) ([]models.User, error) {
	var userM models.User
	const SelectUserByEmail = "SELECT nickname, fullname, about, email FROM users WHERE email=$1 LIMIT 1;"
	row := r.Conn.QueryRow(context.Background(), SelectUserByEmail, usersS.Email) // TODO ctx
	err := row.Scan(&userM.NickName, &userM.FullName, &userM.About, &userM.Email)
	if err != nil {
		return []models.User{}, models.NotFound
	}
	users := make([]models.User, 0)
	users = append(users, userM)
	return users, models.Conflict
}

func (r *repoPostgres) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	const createUser = "INSERT INTO users (Nickname, FullName, About, Email) VALUES ($1, $2, $3, $4);"
	_, err := r.Conn.Exec(ctx, createUser, user.NickName, user.FullName, user.About, user.Email)
	if err != nil {
		return models.User{}, models.InternalError
	}
	return user, nil
}

func (r *repoPostgres) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	const updateUser = "UPDATE users SET FullName=coalesce(nullif($2, ''), FullName), About=coalesce(nullif($3, ''), About), Email=coalesce(nullif($4, ''), Email) WHERE Nickname=$1 RETURNING *"
	row := r.Conn.QueryRow(ctx, updateUser, user.NickName, user.FullName, user.About, user.Email)
	updatedUser := models.User{}
	err := row.Scan(&updatedUser.NickName, &updatedUser.FullName, &updatedUser.About, &updatedUser.Email)
	if err != nil {
		return updatedUser, err
	}
	return updatedUser, nil
}

func (r *repoPostgres) CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error) {
	const createForum = `INSERT INTO forum (Title, "user", Slug, Posts, Threads) VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Conn.Exec(ctx, createForum, forum.Title, forum.User, forum.Slug, forum.Posts, forum.Threads)
	if err != nil {
		return forum, models.InternalError
	}
	return forum, nil
}

func (r *repoPostgres) GetForum(ctx context.Context, slug string) (models.Forum, error) {
	const selectForumBySlug = `SELECT Title, "user", Slug, Posts, Threads FROM forum WHERE $1 = slug;`
	row := r.Conn.QueryRow(ctx, selectForumBySlug, slug)
	finalForum := models.Forum{}
	err := row.Scan(&finalForum.Title, &finalForum.User, &finalForum.Slug, &finalForum.Posts, &finalForum.Threads)
	if err != nil {
		return models.Forum{}, models.NotFound
	}
	return finalForum, nil
}

func (r *repoPostgres) CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error) {
	const createThread = "INSERT INTO threads (Title, Author, Forum, Message, Votes, Slug, Created) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING Id;"
	row := r.Conn.QueryRow(ctx, createThread, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Votes, thread.Slug, thread.Created)
	newThread := models.Thread{}
	err := row.Scan(&newThread.ID)
	thread.ID = newThread.ID
	if err != nil {
		return thread, models.InternalError
	}
	return thread, nil
}

func (r *repoPostgres) GetThread(ctx context.Context, slugOrId string) (models.Thread, error) {
	var row pgx.Row
	if id, err := strconv.Atoi(slugOrId); err == nil {
		const selectThreadBySlug = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Id;`
		row = r.Conn.QueryRow(ctx, selectThreadBySlug, id)
	} else {
		const selectThreadBySlug = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Slug;`
		row = r.Conn.QueryRow(ctx, selectThreadBySlug, slugOrId)
	}
	finalThread := models.Thread{}
	err := row.Scan(&finalThread.ID, &finalThread.Title, &finalThread.Author, &finalThread.Forum, &finalThread.Message, &finalThread.Votes, &finalThread.Slug, &finalThread.Created)
	if err != nil {
		return finalThread, models.NotFound
	}
	return finalThread, nil
}

func (r *repoPostgres) GetThreadByForumSlug(ctx context.Context, slug string, limit string, since string, desc string) ([]models.Thread, error) {
	var rows pgx.Rows
	var err error
	if since != "" {
		if desc == "true" {
			if limit != "" {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Forum AND Created <= $2 order by Created desc limit $3;`
				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, since, limit)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			} else {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Forum AND Created <= $2 order by Created desc;`
				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, since)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			}
		} else {
			if limit != "" {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Forum AND Created >= $2 order by Created limit $3;`
				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, since, limit)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			} else {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Forum AND Created >= $2 order by Created;`
				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, since)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			}
		}
	} else {
		if desc == "true" {
			if limit != "" {

				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Forum order by Created desc limit $2;`
				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, limit)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			} else {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Forum order by Created desc;`
				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			}
		} else {
			if limit != "" {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Forum order by Created limit $2;`
				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, limit)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			} else {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM threads WHERE $1 = Forum order by Created;`
				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			}
		}
	}
	defer rows.Close()
	threads := make([]models.Thread, 0)
	for rows.Next() {
		threadOne := models.Thread{}
		err = rows.Scan(&threadOne.ID, &threadOne.Title, &threadOne.Author, &threadOne.Forum, &threadOne.Message, &threadOne.Votes, &threadOne.Slug, &threadOne.Created)
		if err != nil {
			return []models.Thread{}, models.InternalError
		}
		threads = append(threads, threadOne)
	}
	return threads, nil
}
