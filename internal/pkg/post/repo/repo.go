package repo

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/post"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) post.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) GetPost(ctx context.Context, id int) (models.Post, error) {
	selectPost := `SELECT Id, Author, Created, Forum, IsEdited, Message, Parent, Thread, Path
								 FROM posts
								 WHERE Id = $1;`

	row := r.Conn.QueryRow(ctx, selectPost, id)
	finalPost := models.Post{}
	err := row.Scan(&finalPost.ID, &finalPost.Author, &finalPost.Created, &finalPost.Forum, &finalPost.IsEdited, &finalPost.Message, &finalPost.Parent, &finalPost.Thread, &finalPost.Path)
	if err != nil {
		return models.Post{}, models.NotFound
	}
	return finalPost, nil
}

func (r *repoPostgres) UpdatePost(ctx context.Context, post models.Post) (models.Post, error) {
	updatePost := `UPDATE posts 
								 SET Message = $1, IsEdited = true 
								 WHERE Id = $2 
								 RETURNING Id, Author, Created, Forum, IsEdited, Message, Parent, Thread;`

	row := r.Conn.QueryRow(ctx, updatePost, post.Message, post.ID)
	finalPost := models.Post{}
	err := row.Scan(&finalPost.ID, &finalPost.Author, &finalPost.Created, &finalPost.Forum, &finalPost.IsEdited, &finalPost.Message, &finalPost.Parent, &finalPost.Thread)
	if err != nil {
		return post, models.NotFound
	}
	return finalPost, nil
}
