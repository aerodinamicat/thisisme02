package databases

import (
	"context"

	"github.com/aerodinamicat/thisisme02/models"
)

type DatabaseRepository interface {
	CloseDatabaseConnection() error

	//* Standards methods - Create
	InsertUser(ctx context.Context, user *models.User) error
	InsertPerson(ctx context.Context, person *models.Person) error

	//* Standards methods - Read
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetPersonByUserId(ctx context.Context, userId string) (*models.Person, error)

	ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, error)
	ListPersons(ctx context.Context, pageInfo *models.PageInfo) ([]*models.Person, error)

	//* Standards methods - Update
	UpdateUser(ctx context.Context, user *models.User) error
	UpdatePerson(ctx context.Context, person *models.Person) error

	//* Standards methods - Delete
	DeleteUser(ctx context.Context, id string) error
	DeletePerson(ctx context.Context, userId string) error
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
func InsertPerson(ctx context.Context, person *models.Person) error {
	return dbrImplementation.InsertPerson(ctx, person)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return dbrImplementation.GetUserById(ctx, id)
}
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return dbrImplementation.GetUserByEmail(ctx, email)
}
func GetPersonByUserId(ctx context.Context, userId string) (*models.Person, error) {
	return dbrImplementation.GetPersonByUserId(ctx, userId)
}
func ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, error) {
	return dbrImplementation.ListUsers(ctx, pageInfo)
}
func ListPersons(ctx context.Context, pageInfo *models.PageInfo) ([]*models.Person, error) {
	return dbrImplementation.ListPersons(ctx, pageInfo)
}

func UpdateUser(ctx context.Context, user *models.User) error {
	return dbrImplementation.UpdateUser(ctx, user)
}
func UpdatePerson(ctx context.Context, person *models.Person) error {
	return dbrImplementation.UpdatePerson(ctx, person)
}

func DeleteUser(ctx context.Context, id string) error {
	return dbrImplementation.DeleteUser(ctx, id)
}
func DeletePerson(ctx context.Context, userId string) error {
	return dbrImplementation.DeletePerson(ctx, userId)
}
