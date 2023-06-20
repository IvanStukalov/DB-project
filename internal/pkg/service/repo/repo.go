package repo

import (
	"github.com/IvanStukalov/DB_project/internal/pkg/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) service.Repository {
	return &repoPostgres{Conn: Conn}
}
