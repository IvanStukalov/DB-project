package repo

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) service.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) Status(ctx context.Context) (models.Status, error) {
	countUsers := `SELECT count(*) FROM users`
	countForums := `SELECT count(*) FROM forum`
	countThreads := `SELECT count(*) FROM threads`
	countPosts := `SELECT count(*) FROM posts`

	status := models.Status{}

	row := r.Conn.QueryRow(ctx, countUsers)
	err := row.Scan(&status.User)
	if err != nil {
		return status, models.NotFound
	}

	row = r.Conn.QueryRow(ctx, countForums)
	err = row.Scan(&status.Forum)
	if err != nil {
		return status, models.NotFound
	}

	row = r.Conn.QueryRow(ctx, countThreads)
	err = row.Scan(&status.Thread)
	if err != nil {
		return status, models.NotFound
	}

	row = r.Conn.QueryRow(ctx, countPosts)
	err = row.Scan(&status.Post)
	if err != nil {
		return status, models.NotFound
	}

	return status, nil
}

func (r *repoPostgres) Clear(ctx context.Context) {
	clearDB := `TRUNCATE TABLE users, forum, threads, posts, votes, users_forum CASCADE;`
	_, _ = r.Conn.Exec(ctx, clearDB)
}
