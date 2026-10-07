package unit

import (
	"database/sql"
)

type Unit struct {
	ID            int    `json:"id"`
	VendorID      int    `json:"vendor_id"`
	VendorName    string `json:"vendor_name"`
	PlateNumber   string `json:"plate_number"`
	UnitTypeID    int    `json:"unit_type_id"`
	UnitTypeName  string `json:"unit_type_name"`
	CreatedBy     string `json:"created_by"`
	CreatedByName string `json:"created_by_name"`
	UpdatedBy     string `json:"updated_by"`
	UpdatedByName string `json:"updated_by_name"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const selectColumns = `
	u.id,
	u.vendor_id,
	v.name,

	u.plate_number,

	ut.id,
	ut.name,

	creator.id,
	creator.name,

	updater.id,
	updater.name
`

const fromClause = `
	FROM units u

	INNER JOIN vendors v
		ON v.id = u.vendor_id

	INNER JOIN unit_types ut
		ON ut.id = u.unit_type_id

	INNER JOIN users creator
		ON creator.id = u.created_by

	INNER JOIN users updater
		ON updater.id = u.updated_by
`

func scanUnit(scanner interface {
	Scan(dest ...any) error
}) (*Unit, error) {
	var unit Unit

	err := scanner.Scan(
		&unit.ID,

		&unit.VendorID,
		&unit.VendorName,

		&unit.PlateNumber,

		&unit.UnitTypeID,
		&unit.UnitTypeName,

		&unit.CreatedBy,
		&unit.CreatedByName,

		&unit.UpdatedBy,
		&unit.UpdatedByName,
	)

	if err != nil {
		return nil, err
	}

	return &unit, nil
}

func (r *Repository) Create(unit Unit) error {
	_, err := r.db.Exec(`
		INSERT INTO units (
			vendor_id,
			plate_number,
			unit_type_id,
			created_by,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		unit.VendorID,
		unit.PlateNumber,
		unit.UnitTypeID,
		unit.CreatedBy,
		unit.UpdatedBy,
	)

	return err
}

func (r *Repository) GetByID(id int) (*Unit, error) {
	row := r.db.QueryRow(`
		SELECT
			`+selectColumns+`
		`+fromClause+`
		WHERE u.id = ?
	`, id)

	return scanUnit(row)
}

func (r *Repository) GetAll() ([]Unit, error) {
	rows, err := r.db.Query(`
		SELECT
			` + selectColumns + `
		` + fromClause + `
		ORDER BY u.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	units := make([]Unit, 0)

	for rows.Next() {
		unit, err := scanUnit(rows)
		if err != nil {
			return nil, err
		}

		units = append(units, *unit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return units, nil
}

func (r *Repository) Update(unit Unit) error {
	_, err := r.db.Exec(`
		UPDATE units
		SET
			vendor_id = ?,
			plate_number = ?,
			unit_type_id = ?,
			updated_by = ?
		WHERE id = ?
	`,
		unit.VendorID,
		unit.PlateNumber,
		unit.UnitTypeID,
		unit.UpdatedBy,
		unit.ID,
	)

	return err
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM units
		WHERE id = ?
	`, id)

	return err
}
