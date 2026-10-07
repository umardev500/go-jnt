package user

import (
	"database/sql"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (
			id,
			name
		)
		VALUES (?, ?)
	`,
		user.ID,
		user.Name,
	)

	return err
}

func (r *Repository) GetByID(id string) (*User, error) {
	var user User

	err := r.db.QueryRow(`
		SELECT
			id,
			name
		FROM users
		WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.Name,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetAll() ([]User, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			name
		FROM users
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var user User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Update(user User) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET name = ?
		WHERE id = ?
	`,
		user.Name,
		user.ID,
	)

	return err
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(`
		DELETE FROM users
		WHERE id = ?
	`, id)

	return err
}
