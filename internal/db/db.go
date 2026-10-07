package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Init() (*sql.DB, error) {
	fmt.Println("Initializing database...")

	database, err := sql.Open("sqlite", "app.db")
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, err
	}

	if _, err := database.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		database.Close()
		return nil, err
	}

	if err := createTables(database); err != nil {
		database.Close()
		return nil, err
	}

	fmt.Println("Database initialized")

	return database, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		-- =====================================================
		-- Approval
		-- =====================================================

		CREATE TABLE IF NOT EXISTS approval (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			valid_until DATETIME,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		INSERT OR IGNORE INTO approval (
			id,
			valid_until
		)
		VALUES (1, NULL);


		-- =====================================================
		-- Users
		-- =====================================================

		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);


		-- =====================================================
		-- Vendors
		-- =====================================================

		CREATE TABLE IF NOT EXISTS vendors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);


		-- =====================================================
		-- Fixed Unit Types
		-- =====================================================

		CREATE TABLE IF NOT EXISTS unit_types (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);

		INSERT OR IGNORE INTO unit_types (id, name)
		VALUES
			(1, 'FUSO'),
			(2, 'TRONTON'),
			(3, 'CDDL');


		-- =====================================================
		-- Units
		-- =====================================================

		CREATE TABLE IF NOT EXISTS units (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			vendor_id INTEGER NOT NULL,
			plate_number TEXT NOT NULL UNIQUE,
			unit_type_id INTEGER NOT NULL,
			created_by TEXT NOT NULL,
			updated_by TEXT NOT NULL,

			FOREIGN KEY (vendor_id)
				REFERENCES vendors(id),

			FOREIGN KEY (unit_type_id)
				REFERENCES unit_types(id),

			FOREIGN KEY (created_by)
				REFERENCES users(id),

			FOREIGN KEY (updated_by)
				REFERENCES users(id)
		);
	`)

	return err
}
