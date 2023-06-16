package forum

import "github.com/IvanStukalov/DB_project/internal/models"

type UseCase interface {
	GetUser(user models.User) (models.User, error)
	CreateUser(user models.User) ([]models.User, error)
}

type Repository interface {
	GetUser(name string) (models.User, error)
	CheckUserEmailUniq(usersS models.User) ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
}
