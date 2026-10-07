package vendor

import (
	"database/sql"
)

type Vendor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(vendor Vendor) error {
	_, err := r.db.Exec(`
		INSERT INTO vendors (
			name
		)
		VALUES (?)
	`,
		vendor.Name,
	)

	return err
}

func (r *Repository) GetByID(id int) (*Vendor, error) {
	var vendor Vendor

	err := r.db.QueryRow(`
		SELECT
			id,
			name
		FROM vendors
		WHERE id = ?
	`, id).Scan(
		&vendor.ID,
		&vendor.Name,
	)

	if err != nil {
		return nil, err
	}

	return &vendor, nil
}

func (r *Repository) GetAll() ([]Vendor, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			name
		FROM vendors
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vendors := make([]Vendor, 0)

	for rows.Next() {
		var vendor Vendor

		if err := rows.Scan(
			&vendor.ID,
			&vendor.Name,
		); err != nil {
			return nil, err
		}

		vendors = append(vendors, vendor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
}

func (r *Repository) Update(vendor Vendor) error {
	_, err := r.db.Exec(`
		UPDATE vendors
		SET name = ?
		WHERE id = ?
	`,
		vendor.Name,
		vendor.ID,
	)

	return err
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM vendors
		WHERE id = ?
	`, id)

	return err
}
