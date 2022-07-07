package databases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aerodinamicat/thisisme02/models"
)

type PostgresImplementation struct {
	db *sql.DB

	Name string

	DatabaseHost     string
	DatabasePassword string
	DatabasePort     string
	DatabaseSchema   string
	DatabaseUser     string
}

func NewPostgresImplementation(user string, password string, host string, port string, schema string) (*PostgresImplementation, error) {
	var pgr PostgresImplementation

	name := "postgres"
	url := pgr.BuildURL(name, user, password, host, port, schema)

	db, err := sql.Open(name, url)
	if err != nil {
		return nil, err
	}

	pgr.db = db
	pgr.Name = name
	pgr.DatabaseHost = host
	pgr.DatabasePassword = password
	pgr.DatabasePort = port
	pgr.DatabaseSchema = schema
	pgr.DatabaseUser = user

	return &pgr, nil
}
func (pgr *PostgresImplementation) BuildURL(name string, user string, password string, host string, port string, schema string) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", name, user, password, host, port, schema)
}

func (pgr *PostgresImplementation) CloseDatabaseConnection() error {
	return pgr.db.Close()
}

func (pgr *PostgresImplementation) InsertUser(ctx context.Context, u *models.User) error {
	querySentence := `
		INSERT INTO users (
			id, email, password, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := pgr.db.ExecContext(ctx, querySentence,
		u.Id, u.Email, u.Password, time.Now(), time.Now(),
	); err != nil {
		return err
	}

	return nil
}
func (pgr *PostgresImplementation) InsertPerson(ctx context.Context, p *models.Person) error {
	querySentence := `
		INSERT INTO persons (
			first_name, second_name, first_surname, second_surname, gender, birth_date, user_id
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	if _, err := pgr.db.ExecContext(ctx, querySentence,
		p.FirstName, p.SecondName, p.FirstSurname, p.SecondSurname, p.Gender, p.BirthDate, p.UserId,
	); err != nil {
		return err
	}

	return nil
}

func (pgr *PostgresImplementation) GetUserById(ctx context.Context, id string) (*models.User, error) {
	var querySentence string
	var rows *sql.Rows
	var err error

	var user *models.User
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	querySentence = `
		SELECT
			id, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	if rows, err = pgr.db.QueryContext(ctx, querySentence,
		id,
	); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&user.Id, &user.Email, &user.Password, createdAt, updatedAt,
		); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}

	return user, nil
}
func (pgr *PostgresImplementation) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var querySentence string
	var rows *sql.Rows
	var err error

	var user *models.User
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	querySentence = `
		SELECT
			id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	if rows, err = pgr.db.QueryContext(ctx, querySentence, email); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&user.Id, &user.Email, &user.Password, createdAt, updatedAt,
		); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}

	return user, nil
}
func (pgr *PostgresImplementation) GetPersonByUserId(ctx context.Context, userId string) (*models.Person, error) {
	var querySentence string
	var rows *sql.Rows
	var err error

	var person *models.Person
	var birthDate sql.NullTime
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	querySentence = `
		SELECT
			first_name, second_name, first_surname, second_surname, gender, birth_date, created_at, updated_at, user_id
		FROM persons
		WHERE user_id = $1
	`
	if rows, err = pgr.db.QueryContext(ctx, querySentence,
		userId,
	); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&person.FirstName, &person.SecondName, &person.FirstSurname, &person.SecondSurname,
			&person.Gender, birthDate, createdAt, updatedAt, &person.UserId,
		); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if birthDate.Valid {
		person.BirthDate = birthDate.Time
	}
	if createdAt.Valid {
		person.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		person.UpdatedAt = updatedAt.Time
	}

	return person, nil
}
func (pgr *PostgresImplementation) ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, *models.PageInfo, error) {
	var querySentence string
	var rows *sql.Rows
	var err error

	var users []*models.User

	var user *models.User
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	if pageInfo.TotalPages == 0 || pageInfo.TotalItems == 0 {
		querySentence = `
			SELECT count(*) AS total_items
			FROM users
		`
		if rows, err = pgr.db.QueryContext(ctx, querySentence); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}
		defer rows.Close()

		for rows.Next() {
			if err = rows.Scan(
				&pageInfo.TotalItems,
			); err != nil {
				pageInfo.TotalPages = 0
				pageInfo.TotalItems = 0

				return nil, pageInfo, err
			}
		}
		if err = rows.Err(); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}

		pageInfo.TotalPages = pageInfo.TotalItems / pageInfo.Size
		if pageInfo.TotalItems%pageInfo.Size == 0 {
			pageInfo.TotalPages++
		}
	}

	querySentence = `
		SELECT
			id, email, password, created_at, updated_at
		FROM users
		ORDER BY $1 LIMIT $2 OFFSET $3
	`
	if rows, err = pgr.db.QueryContext(ctx, querySentence,
		pageInfo.OrderBy, pageInfo.Size, pageInfo.Token*pageInfo.Size,
	); err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}
	defer rows.Close()

	for rows.Next() {
		user = &models.User{}
		if err = rows.Scan(
			&user.Id, &user.Email, &user.Password, createdAt, updatedAt,
		); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}

		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}

	return users, pageInfo, nil
}
func (pgr *PostgresImplementation) ListPersons(ctx context.Context, pageInfo *models.PageInfo) ([]*models.Person, *models.PageInfo, error) {
	var querySentence string
	var rows *sql.Rows
	var err error

	var persons []*models.Person

	var person *models.Person
	var birthDate sql.NullTime
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	if pageInfo.TotalPages == 0 || pageInfo.TotalItems == 0 {
		querySentence = `
			SELECT count(*) AS total_items
			FROM persons
		`
		if rows, err = pgr.db.QueryContext(ctx, querySentence); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}
		defer rows.Close()

		for rows.Next() {
			if err = rows.Scan(
				&pageInfo.TotalItems,
			); err != nil {
				pageInfo.TotalPages = 0
				pageInfo.TotalItems = 0

				return nil, pageInfo, err
			}
		}
		if err = rows.Err(); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}

		pageInfo.TotalPages = pageInfo.TotalItems / pageInfo.Size
		if pageInfo.TotalItems%pageInfo.Size == 0 {
			pageInfo.TotalPages++
		}
	}

	querySentence = `
		SELECT
			id, email, password, created_at, updated_at
		FROM users
		ORDER BY $1 LIMIT $2 OFFSET $3
	`
	if rows, err = pgr.db.QueryContext(ctx, querySentence,
		pageInfo.OrderBy, pageInfo.Size, pageInfo.Token,
	); err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}
	defer rows.Close()

	for rows.Next() {
		person = &models.Person{}
		if err = rows.Scan(
			&person.FirstName, &person.SecondName, &person.FirstSurname, &person.SecondSurname,
			&person.Gender, birthDate, createdAt, updatedAt, &person.UserId,
		); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}

		if birthDate.Valid {
			person.BirthDate = birthDate.Time
		}

		if createdAt.Valid {
			person.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			person.UpdatedAt = updatedAt.Time
		}

		persons = append(persons, person)
	}
	if err = rows.Err(); err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}

	return persons, pageInfo, nil
}

func (pgr *PostgresImplementation) UpdateUser(ctx context.Context, user *models.User) error {
	querySentence := `
		UPDATE users SET
			email = $1, password = $2, updated_at = $3
		WHERE id = $4
	`
	user.UpdatedAt = time.Now()
	if _, err := pgr.db.ExecContext(ctx, querySentence,
		user.Email, user.Password, user.UpdatedAt, user.Id,
	); err != nil {
		return err
	}

	return nil
}
func (pgr *PostgresImplementation) UpdatePerson(ctx context.Context, person *models.Person) error {
	querySentence := `
		UPDATE persons SET
			first_name = $1, second_name = $2, first_surname = $3, second_surname = $4,
			gender = $5, birth_date = $6, updated_at = $7
		WHERE user_id = $8
	`
	person.UpdatedAt = time.Now()
	if _, err := pgr.db.ExecContext(ctx, querySentence,
		person.FirstName, person.SecondName, person.FirstSurname, person.SecondSurname,
		person.Gender, person.BirthDate, person.UpdatedAt, person.UserId,
	); err != nil {
		return err
	}

	return nil
}

func (pgr *PostgresImplementation) DeleteUser(ctx context.Context, id string) error {
	querySentence := `
		DELETE FROM users
		WHERE id = $1
	`
	if _, err := pgr.db.ExecContext(ctx, querySentence,
		id,
	); err != nil {
		return err
	}

	return nil
}
func (pgr *PostgresImplementation) DeletePerson(ctx context.Context, userId string) error {
	querySentence := `
		DELETE FROM persons
		WHERE user_id = $1
	`
	if _, err := pgr.db.ExecContext(ctx, querySentence,
		userId,
	); err != nil {
		return err
	}

	return nil
}
