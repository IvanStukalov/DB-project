package repo

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/forum"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) forum.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error) {
	const createForum = `INSERT INTO forum (Title, "user", Slug, Posts, Threads) 
										   VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Conn.Exec(ctx, createForum, forum.Title, forum.User, forum.Slug, forum.Posts, forum.Threads)
	if err != nil {
		return forum, models.InternalError
	}
	return forum, nil
}

func (r *repoPostgres) GetForum(ctx context.Context, slug string) (models.Forum, error) {
	const selectForumBySlug = `SELECT Title, "user", Slug, Posts, Threads 
														 FROM forum 
														 WHERE $1 = slug;`

	row := r.Conn.QueryRow(ctx, selectForumBySlug, slug)
	finalForum := models.Forum{}
	err := row.Scan(&finalForum.Title, &finalForum.User, &finalForum.Slug, &finalForum.Posts, &finalForum.Threads)
	if err != nil {
		return models.Forum{}, models.NotFound
	}
	return finalForum, nil
}

func (r *repoPostgres) CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error) {
	const createThread = `INSERT INTO threads (Title, Author, Forum, Message, Votes, Slug, Created) 
												VALUES ($1, $2, $3, $4, $5, $6, $7) 
												RETURNING Id;`

	row := r.Conn.QueryRow(ctx, createThread, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Votes, thread.Slug, thread.Created)
	newThread := models.Thread{}
	err := row.Scan(&newThread.ID)
	if err != nil {
		return thread, models.InternalError
	}
	thread.ID = newThread.ID
	return thread, nil
}

func (r *repoPostgres) GetThreadByForumSlug(ctx context.Context, slug string, limit string, since string, desc string) ([]models.Thread, error) {
	var rows pgx.Rows
	var err error
	if since != "" {
		if desc == "true" {
			if limit != "" {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																		 FROM threads 
																		 WHERE $1 = Forum AND Created <= $2 
																		 ORDER BY Created DESC 
																		 LIMIT $3;`

				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, since, limit)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			} else {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																		 FROM threads 
																		 WHERE $1 = Forum AND Created <= $2 
																		 ORDER BY Created DESC;`

				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, since)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			}
		} else {
			if limit != "" {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																		 FROM threads 
																		 WHERE $1 = Forum AND Created >= $2 
																		 ORDER BY Created 
																		 LIMIT $3;`

				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, since, limit)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			} else {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																		 FROM threads 
																		 WHERE $1 = Forum AND Created >= $2 
																		 ORDER BY Created;`

				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, since)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			}
		}
	} else {
		if desc == "true" {
			if limit != "" {

				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																		 FROM threads 
																		 WHERE $1 = Forum 
																		 ORDER BY Created DESC 
																		 LIMIT $2;`

				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, limit)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			} else {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																		 FROM threads 
																		 WHERE $1 = Forum 
																		 ORDER BY Created DESC;`

				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			}
		} else {
			if limit != "" {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																		 FROM threads 
																		 WHERE $1 = Forum 
																		 ORDER BY Created 
																		 LIMIT $2;`

				rows, err = r.Conn.Query(ctx, selectThreadByForum, slug, limit)
				if err != nil {
					return []models.Thread{}, models.InternalError
				}
			} else {
				const selectThreadByForum = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																		 FROM threads 
																		 WHERE $1 = Forum 
																		 ORDER BY Created;`

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
