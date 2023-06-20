package repo

import (
	"context"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/user"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) user.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) GetUser(ctx context.Context, name string) (models.User, error) {
	var userM models.User
	const SelectUserByNickname = `SELECT nickname, fullname, about, email 
																FROM users 
																WHERE nickname=$1 
																LIMIT 1;`

	row := r.Conn.QueryRow(ctx, SelectUserByNickname, name)
	err := row.Scan(&userM.NickName, &userM.FullName, &userM.About, &userM.Email)
	if err != nil {
		return models.User{}, models.NotFound
	}
	return userM, nil
}

func (r *repoPostgres) IsEmailOrNicknameUniq(ctx context.Context, usersS models.User) ([]models.User, error) {
	const SelectUserByEmailOrNickname = `SELECT nickname, fullname, about, email 
																			 FROM users 
																			 WHERE nickname=$1 OR email=$2 
																			 LIMIT 2;`

	rows, err := r.Conn.Query(ctx, SelectUserByEmailOrNickname, usersS.NickName, usersS.Email)
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

func (r *repoPostgres) IsEmailUniq(ctx context.Context, usersS models.User) (models.User, error) {
	var userM models.User
	const SelectUserByEmail = `SELECT nickname, fullname, about, email 
														 FROM users 
														 WHERE email=$1`

	row := r.Conn.QueryRow(ctx, SelectUserByEmail, usersS.Email)
	err := row.Scan(&userM.NickName, &userM.FullName, &userM.About, &userM.Email)
	if err != nil {
		return models.User{}, models.NotFound
	}
	return userM, models.Conflict
}

func (r *repoPostgres) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	const createUser = `INSERT INTO users (Nickname, FullName, About, Email) 
											VALUES ($1, $2, $3, $4);`

	_, err := r.Conn.Exec(ctx, createUser, user.NickName, user.FullName, user.About, user.Email)
	if err != nil {
		return models.User{}, models.InternalError
	}
	return user, nil
}

func (r *repoPostgres) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	const updateUser = `UPDATE users 
											SET FullName=coalesce(nullif($2, ''), FullName), About=coalesce(nullif($3, ''), About), 
											    Email=coalesce(nullif($4, ''), Email) 
											WHERE Nickname=$1 
											RETURNING *;`

	row := r.Conn.QueryRow(ctx, updateUser, user.NickName, user.FullName, user.About, user.Email)
	updatedUser := models.User{}
	err := row.Scan(&updatedUser.NickName, &updatedUser.FullName, &updatedUser.About, &updatedUser.Email)
	if err != nil {
		return updatedUser, err
	}
	return updatedUser, nil
}
