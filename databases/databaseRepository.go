package databases

import (
	"context"

	"github.com/aerodinamicat/thisisme02/models"
)

type DatabaseRepository interface {
	CloseDatabaseConnection() error

	//* Standards methods - Create
	InsertUser(ctx context.Context, user *models.User) error

	//* Standards methods - Read
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, *models.PageInfo, error)

	//* Standards methods - Update
	UpdateUser(ctx context.Context, user *models.User) error

	//* Standards methods - Delete
	DeleteUser(ctx context.Context, id string) error
}

var dbrImplementation DatabaseRepository

func SetDatabaseRepository(dbr DatabaseRepository) {
	dbrImplementation = dbr
}

func CloseDatabaseConnection() error {
	return dbrImplementation.CloseDatabaseConnection()
}

func InsertUser(ctx context.Context, user *models.User) error {
	return dbrImplementation.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return dbrImplementation.GetUserById(ctx, id)
}
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return dbrImplementation.GetUserByEmail(ctx, email)
}

func ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, *models.PageInfo, error) {
	return dbrImplementation.ListUsers(ctx, pageInfo)
}

func UpdateUser(ctx context.Context, user *models.User) error {
	return dbrImplementation.UpdateUser(ctx, user)
}

func DeleteUser(ctx context.Context, id string) error {
	return dbrImplementation.DeleteUser(ctx, id)
}
