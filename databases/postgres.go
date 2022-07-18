package databases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aerodinamicat/thisisme02/models"
	_ "github.com/lib/pq"
)

type PostgresImplementation struct {
	DB *sql.DB

	Name string

	DatabaseHost     string
	DatabasePassword string
	DatabasePort     string
	DatabaseSchema   string
	DatabaseUser     string
}

func NewPostgresImplementation(user string, password, host, port, schema string) (*PostgresImplementation, error) {
	name := "postgres"
	url := buildURL(name, user, password, host, port, schema)

	db, err := sql.Open(name, url)
	if err != nil {
		return nil, err
	}

	return &PostgresImplementation{
		DB:               db,
		Name:             name,
		DatabaseHost:     host,
		DatabasePassword: password,
		DatabasePort:     port,
		DatabaseSchema:   schema,
		DatabaseUser:     user,
	}, nil
}
func buildURL(name, user, password, host, port, schema string) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", name, user, password, host, port, schema)
}

func (pgr *PostgresImplementation) CloseDatabaseConnection() error {
	return pgr.DB.Close()
}

func (pgr *PostgresImplementation) InsertUser(ctx context.Context, user *models.User) error {
	querySentence := `
		INSERT INTO users (
			id, email, password, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := pgr.DB.ExecContext(ctx, querySentence,
		user.Id,
		user.Email,
		user.Password,
		time.Now(),
		time.Now(),
	); err != nil {
		return err
	}

	return nil
}

func (pgr *PostgresImplementation) GetUserById(ctx context.Context, id string) (*models.User, error) {
	querySentence := `
		SELECT
			id, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	rows, err := pgr.DB.QueryContext(ctx, querySentence,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user = new(models.User)
	for rows.Next() {
		var createdAt, updatedAt sql.NullTime
		if err = rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			createdAt,
			updatedAt,
		); err != nil {
			return nil, err
		}

		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}
func (pgr *PostgresImplementation) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	querySentence := `
		SELECT
			id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	rows, err := pgr.DB.QueryContext(ctx, querySentence, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(models.User)
	for rows.Next() {
		var createdAt sql.NullTime
		var updatedAt sql.NullTime
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}
func (pgr *PostgresImplementation) ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, *models.PageInfo, error) {
	if pageInfo.TotalPages == 0 || pageInfo.TotalItems == 0 {
		querySentence := `
			SELECT count(*) AS total_items
			FROM users
		`
		rows, err := pgr.DB.QueryContext(ctx, querySentence)
		if err != nil {
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

	querySentence := `
		SELECT
			id, email, password, created_at, updated_at
		FROM users
		ORDER BY $1 LIMIT $2 OFFSET $3
	`
	rows, err := pgr.DB.QueryContext(ctx, querySentence,
		pageInfo.OrderBy,
		pageInfo.Size,
		pageInfo.Token*pageInfo.Size,
	)
	if err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := new(models.User)
		var createdAt, updatedAt sql.NullTime

		if err = rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			createdAt,
			updatedAt,
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

func (pgr *PostgresImplementation) UpdateUser(ctx context.Context, user *models.User) error {
	querySentence := `
		UPDATE users SET
			email = $1, password = $2, updated_at = $3
		WHERE id = $4
	`
	user.UpdatedAt = time.Now()

	if _, err := pgr.DB.ExecContext(ctx, querySentence,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.Id,
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
	if _, err := pgr.DB.ExecContext(ctx, querySentence,
		id,
	); err != nil {
		return err
	}

	return nil
}
