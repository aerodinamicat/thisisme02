package databases

import (
	"context"
	"database/sql"

	"github.com/aerodinamicat/thisisme02/models"
)

type PostgresImplementation struct {
	db *sql.DB
}

func NewPostgresImplementation(databaseURL string) (*PostgresImplementation, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	return &PostgresImplementation{db}, nil
}

func (pgr *PostgresImplementation) CloseDatabaseConnection() error {
	return pgr.db.Close()
}

func (pgr *PostgresImplementation) InsertUser(ctx context.Context, user *models.User) error {
	return nil
}
func (pgr *PostgresImplementation) InsertPerson(ctx context.Context, person *models.Person) error {
	return nil
}

func (pgr *PostgresImplementation) GetUserById(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}
func (pgr *PostgresImplementation) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (pgr *PostgresImplementation) GetPersonByUserId(ctx context.Context, userId string) (*models.Person, error) {
	return nil, nil
}
func (pgr *PostgresImplementation) ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, error) {
	return nil, nil
}
func (pgr *PostgresImplementation) ListPersons(ctx context.Context, pageInfo *models.PageInfo) ([]*models.Person, error) {
	return nil, nil
}

func (pgr *PostgresImplementation) UpdateUser(ctx context.Context, user *models.User) error {
	return nil
}
func (pgr *PostgresImplementation) UpdatePerson(ctx context.Context, person *models.Person) error {
	return nil
}

func (pgr *PostgresImplementation) DeleteUser(ctx context.Context, id string) error {
	return nil
}
func (pgr *PostgresImplementation) DeletePerson(ctx context.Context, userId string) error {
	return nil
}
