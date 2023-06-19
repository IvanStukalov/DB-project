package repo

import (
	"github.com/IvanStukalov/DB_project/internal/pkg/post"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) post.Repository {
	return &repoPostgres{Conn: Conn}
}
