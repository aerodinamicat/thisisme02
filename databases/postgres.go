package databases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aerodinamicat/thisisme02/models"
	_ "github.com/lib/pq"
)

const (
	DEFAULT_PAGE_SIZE = 5
	DEFAULT_ORDER_BY  = "created_at desc"
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
	//* Construimos la sentencia SQL.
	querySentence := `
		INSERT INTO users (
			id, email, password, created_at, created_by, updated_at, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	//* Ejecutamos la sentencia.
	if _, err := pgr.DB.ExecContext(ctx, querySentence,
		user.Id,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.CreatedBy,
		user.UpdatedAt,
		user.UpdatedBy,
	); err != nil {
		return err
	}

	return nil
}

func (pgr *PostgresImplementation) GetUserById(ctx context.Context, id string) (*models.User, error) {
	//* Construimos la consulta SQL.
	querySentence := `
		SELECT
			id, email, password, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM users
		WHERE id = $1
	`
	//* Ejecutamos la consulta.
	rows, err := pgr.DB.QueryContext(ctx, querySentence,
		id,
	)
	if err != nil {
		return nil, err
	}
	//* Cerramos la consulta al final de éste proceso.
	defer rows.Close()

	//* Obtenemos los resultados de la ejecución.
	var user = new(models.User)
	for rows.Next() {
		var deletedAt sql.NullTime
		var deletedBy sql.NullString
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.CreatedBy,
			&user.UpdatedAt,
			&user.UpdatedBy,
			&deletedAt,
			&deletedBy,
		); err != nil {
			return nil, err
		}

		//* Si los campos 'sql.NullTime' son válidos, es decir que no son nulos, los asignamos a 'user'.
		if deletedAt.Valid {
			user.DeletedAt = deletedAt.Time
		}
		if deletedBy.Valid {
			user.DeletedBy = deletedBy.String
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	//* Devolvemos la información obtenida.
	return user, nil
}
func (pgr *PostgresImplementation) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	//* Contruimos la consulta SQL.
	querySentence := `
		SELECT
			id, email, password, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM users
		WHERE email = $1
	`
	//* Ejecutamos la consulta.
	rows, err := pgr.DB.QueryContext(ctx, querySentence, email)
	if err != nil {
		return nil, err
	}
	//* Cerramos la consulta al final de éste proceso.
	defer rows.Close()

	//* Obtenemos los resultados de la ejecución.
	user := new(models.User)
	for rows.Next() {
		var deletedAt sql.NullTime
		var deletedBy sql.NullString
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.CreatedBy,
			&user.UpdatedAt,
			&user.UpdatedBy,
			&deletedAt,
			&deletedBy,
		); err != nil {
			return nil, err
		}

		//* Si los campos 'sql.NullTime' son válidos, es decir que no son nulos, los asignamos a 'user'.
		if deletedAt.Valid {
			user.DeletedAt = deletedAt.Time
		}
		if deletedBy.Valid {
			user.DeletedBy = deletedBy.String
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	//* Devolvemos la información obtenida.
	return user, nil
}
func (pgr *PostgresImplementation) ListUsers(ctx context.Context, pageInfo *models.PageInfo) ([]*models.User, *models.PageInfo, error) {
	//* Si 'pageInfo' no tiene valores, lo poblamos.
	if pageInfo.TotalPages == 0 || pageInfo.TotalItems == 0 {
		//* Construimos la consulta SQL.
		querySentence := `
			SELECT count(*) AS total_items
			FROM users
		`
		//* Ejecutamos la consulta.
		rows, err := pgr.DB.QueryContext(ctx, querySentence)
		if err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}
		//* Cerramos la consulta al final de éste proceso.
		defer rows.Close()

		//* Obtenemos los resultado de la ejecución.
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

		//* Calculamos el total de páginas en función del número de elementos por cada una.
		pageInfo.TotalPages = pageInfo.TotalItems / pageInfo.Size
		//* Si los elementos no caben en un total de páginas exacto, añadimos una página mas.
		if pageInfo.TotalItems%pageInfo.Size != 0 {
			pageInfo.TotalPages++
		}
		pageInfo.Token = 1
	}

	//* Construimos la consulta SQL.
	querySentence := `
		SELECT
			id, email, password, created_at, updated_at
		FROM users
		ORDER BY $1 LIMIT $2 OFFSET $3
	`
	//* Ejecutamos la consulta.
	rows, err := pgr.DB.QueryContext(ctx, querySentence,
		pageInfo.OrderBy,
		pageInfo.Size,
		(pageInfo.Token-1)*pageInfo.Size,
	)
	if err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}
	//* Cerramos la consulta al final de éste proceso.
	defer rows.Close()

	//* Obtenemos los resultados de la consulta.
	//* Dado que esperamos una lista, usamos un 'array' vacío.
	var users []*models.User
	for rows.Next() {
		//* En cada iteración:
		//* Creamos un usuario vacío y los campos de tiempo adicionales necesarios nuevamente.
		user := new(models.User)
		var createdAt, updatedAt, deletedAt sql.NullTime

		//* Poblamos 'user' y los campos de tiempo.
		if err = rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			createdAt,
			&user.CreatedBy,
			updatedAt,
			&user.UpdatedBy,
			deletedAt,
			&user.DeletedBy,
		); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}

		//* Si los campos 'sql.NullTime' son válidos, es decir que no son nulos, los asignamos a 'user'.
		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			user.DeletedAt = deletedAt.Time
		}

		//* Añadimos al 'array' el nuevo usuario.
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}

	//* Actualizamos 'pageInfo' y la preparamos para la siguiente iteración, si la hubiera.
	pageInfo.Token++

	//* Devolvemos la información obtenida.
	return users, pageInfo, nil
}
func (pgr *PostgresImplementation) UpdateUser(ctx context.Context, user *models.User) error {
	//* Construimos la sentencia SQL.
	querySentence := `
		UPDATE users SET
			email = $1, password = $2, updated_at = $3, updated_by = $4
		WHERE id = $4
	`

	//* Ejecutamos la sentencia.
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
	//* Construímos las sentencia SQL.
	querySentence := `
		DELETE FROM users
		WHERE id = $1
	`
	//* Ejecutamos la sentencia.
	if _, err := pgr.DB.ExecContext(ctx, querySentence,
		id,
	); err != nil {
		return err
	}

	return nil
}

func (pgr *PostgresImplementation) InsertPropertyChangeLog(ctx context.Context, propertyChange *models.PropertyChange) error {
	//* Construímos las sentencia SQL.
	querySentence := `
		INSERT INTO users_properties_changes_history (
			user_id, name, changed_from, changed_to, created_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	//* Ejecutamos la sentencia.
	if _, err := pgr.DB.ExecContext(ctx, querySentence,
		propertyChange.UserId,
		propertyChange.Name,
		propertyChange.From,
		propertyChange.To,
		propertyChange.CreatedAt,
		propertyChange.CreatedBy,
	); err != nil {
		return err
	}

	return nil
}
func (pgr *PostgresImplementation) ListPropertyChangesByUserId(ctx context.Context, id string, pageInfo *models.PageInfo) ([]*models.PropertyChange, *models.PageInfo, error) {
	//* Si 'pageInfo' no tiene valores, lo poblamos.
	pageInfo.OrderBy = DEFAULT_ORDER_BY
	if pageInfo.Size == 0 {
		pageInfo.Size = DEFAULT_PAGE_SIZE
	}

	if pageInfo.TotalPages == 0 || pageInfo.TotalItems == 0 {
		//* Construimos la consulta SQL.
		querySentence := `
			SELECT count(*) AS total_items
			FROM users_properties_changes_history WHERE user_id = $1
		`
		//* Ejecutamos la consulta.
		rows, err := pgr.DB.QueryContext(ctx, querySentence,
			id,
		)
		if err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}
		//* Cerramos la consulta al final de éste proceso.
		defer rows.Close()

		//* Obtenemos los resultado de la ejecución.
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

		//* Calculamos el total de páginas en función del número de elementos por cada una.
		pageInfo.TotalPages = pageInfo.TotalItems / pageInfo.Size
		//* Si los elementos no caben en un total de páginas exacto, añadimos una página mas.
		if pageInfo.TotalItems%pageInfo.Size != 0 {
			pageInfo.TotalPages++
		}
		pageInfo.Token = 1
	}

	//* Construimos la consulta SQL.
	querySentence := fmt.Sprintf(`
		SELECT
			user_id, name, changed_from, changed_to, created_at, created_by
		FROM users_properties_changes_history WHERE user_id = $1
		ORDER BY %s LIMIT $2 OFFSET $3
	`, pageInfo.OrderBy)
	//* Ejecutamos la consulta.
	rows, err := pgr.DB.QueryContext(ctx, querySentence,
		id,
		pageInfo.Size,
		(pageInfo.Token-1)*pageInfo.Size,
	)
	if err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}
	//* Cerramos la consulta al final de éste proceso.
	defer rows.Close()

	//* Obtenemos los resultados de la consulta.
	//* Dado que esperamos una lista, usamos un 'array' vacío.
	var propertyChanges []*models.PropertyChange
	for rows.Next() {
		//* En cada iteración:
		//* Creamos un usuario vacío y los campos de tiempo adicionales necesarios nuevamente.
		propertyChange := new(models.PropertyChange)
		var createdAt sql.NullTime

		//* Poblamos 'user' y los campos de tiempo.
		if err = rows.Scan(
			&propertyChange.UserId,
			&propertyChange.Name,
			&propertyChange.From,
			&propertyChange.To,
			createdAt,
			&propertyChange.CreatedBy,
		); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}

		//* Si los campos 'sql.NullTime' son válidos, es decir que no son nulos, los asignamos a 'user'.
		if createdAt.Valid {
			propertyChange.CreatedAt = createdAt.Time
		}

		//* Añadimos al 'array' el nuevo usuario.
		propertyChanges = append(propertyChanges, propertyChange)
	}
	if err = rows.Err(); err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}

	//* Actualizamos 'pageInfo' y la preparamos para la siguiente iteración, si la hubiera.
	pageInfo.Token++

	//* Devolvemos la información obtenida.
	return propertyChanges, pageInfo, nil
}
func (pgr *PostgresImplementation) ListPropertyChangesByUserIdAndName(ctx context.Context, id, name string, pageInfo *models.PageInfo) ([]*models.PropertyChange, *models.PageInfo, error) {
	//* Si 'pageInfo' no tiene valores, lo poblamos.
	pageInfo.OrderBy = DEFAULT_ORDER_BY
	if pageInfo.Size == 0 {
		pageInfo.Size = DEFAULT_PAGE_SIZE
	}

	if pageInfo.TotalPages == 0 || pageInfo.TotalItems == 0 {
		//* Construimos la consulta SQL.
		querySentence := `
			SELECT count(*) AS total_items
			FROM users_properties_changes_history WHERE user_id = $1 AND name = $2
		`
		//* Ejecutamos la consulta.
		rows, err := pgr.DB.QueryContext(ctx, querySentence,
			id,
			name,
		)
		if err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}
		//* Cerramos la consulta al final de éste proceso.
		defer rows.Close()

		//* Obtenemos los resultado de la ejecución.
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

		//* Calculamos el total de páginas en función del número de elementos por cada una.
		pageInfo.TotalPages = pageInfo.TotalItems / pageInfo.Size
		//* Si los elementos no caben en un total de páginas exacto, añadimos una página mas.
		if pageInfo.TotalItems%pageInfo.Size != 0 {
			pageInfo.TotalPages++
		}
		pageInfo.Token = 1
	}

	//* Construimos la consulta SQL.
	querySentence := fmt.Sprintf(`
		SELECT
			user_id, name, changed_from, changed_to, created_at, created_by
		FROM users_properties_changes_history WHERE user_id = $1 AND name = $2
		ORDER BY %s LIMIT $3 OFFSET $4
	`, pageInfo.OrderBy)
	//* Ejecutamos la consulta.
	rows, err := pgr.DB.QueryContext(ctx, querySentence,
		id,
		name,
		pageInfo.Size,
		(pageInfo.Token-1)*pageInfo.Size,
	)
	if err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}
	//* Cerramos la consulta al final de éste proceso.
	defer rows.Close()

	//* Obtenemos los resultados de la consulta.
	//* Dado que esperamos una lista, usamos un 'array' vacío.
	var propertyChanges []*models.PropertyChange
	for rows.Next() {
		//* En cada iteración:
		//* Creamos un usuario vacío y los campos de tiempo adicionales necesarios nuevamente.
		propertyChange := new(models.PropertyChange)
		var createdAt sql.NullTime

		//* Poblamos 'user' y los campos de tiempo.
		if err = rows.Scan(
			&propertyChange.UserId,
			&propertyChange.Name,
			&propertyChange.From,
			&propertyChange.To,
			&propertyChange.CreatedAt,
			&propertyChange.CreatedBy,
		); err != nil {
			pageInfo.TotalPages = 0
			pageInfo.TotalItems = 0

			return nil, pageInfo, err
		}

		//* Si los campos 'sql.NullTime' son válidos, es decir que no son nulos, los asignamos a 'user'.
		if createdAt.Valid {
			propertyChange.CreatedAt = createdAt.Time
		}

		//* Añadimos al 'array' el nuevo usuario.
		propertyChanges = append(propertyChanges, propertyChange)
	}
	if err = rows.Err(); err != nil {
		pageInfo.TotalPages = 0
		pageInfo.TotalItems = 0

		return nil, pageInfo, err
	}

	//* Devolvemos la información obtenida.
	return propertyChanges, pageInfo, nil
}
