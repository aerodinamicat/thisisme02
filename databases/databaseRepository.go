package databases

import (
	"context"

	"github.com/aerodinamicat/thisisme02/models"
)

type DatabaseRepository interface {
	CloseDatabaseConnection() error

	//* User related methods:
	//* Standard REST API
	//* Create
	InsertUser(ctx context.Context, user *models.User) error
	//* Read
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, *models.PageInfo, error)
	//* Update
	UpdateUser(ctx context.Context, user *models.User) error
	//* Delete
	DeleteUser(ctx context.Context, id string) error

	//* Property changes related methods:
	//* Standard REST API
	//* Create
	InsertPropertyChangeLog(ctx context.Context, propertyChange *models.PropertyChange) error
	//* Read
	ListPropertyChangesByUserId(ctx context.Context, id string, pageInfo *models.PageInfo) ([]*models.PropertyChange, *models.PageInfo, error)
	ListPropertyChangesByUserIdAndName(ctx context.Context, id, name string, pageInfo *models.PageInfo) ([]*models.PropertyChange, *models.PageInfo, error)
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

func InsertPropertyChangeLog(ctx context.Context, propertyChange *models.PropertyChange) error {
	return dbrImplementation.InsertPropertyChangeLog(ctx, propertyChange)
}
func ListPropertyChangesByUserId(ctx context.Context, id string, pageInfo *models.PageInfo) ([]*models.PropertyChange, *models.PageInfo, error) {
	return dbrImplementation.ListPropertyChangesByUserId(ctx, id, pageInfo)
}
func ListPropertyChangesByUserIdAndName(ctx context.Context, id, name string, pageInfo *models.PageInfo) ([]*models.PropertyChange, *models.PageInfo, error) {
	return dbrImplementation.ListPropertyChangesByUserIdAndName(ctx, id, name, pageInfo)
}
