package approval

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	fmt.Println("Initializing database...")

	DB, err = sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	createTable()

	fmt.Println("Database initialized")
}

func createTable() {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS approval (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		valid_until DATETIME,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	INSERT OR IGNORE INTO approval (id, valid_until)
	VALUES (1, NULL);
	`)

	if err != nil {
		log.Fatal(err)
	}
}

func Grant(minutes int) error {

	_, err := DB.Exec(`
	UPDATE approval
	SET valid_until = datetime('now', ? || ' minutes'),
	    updated_at = CURRENT_TIMESTAMP
	WHERE id = 1
	`, minutes)

	return err
}

func Revoke() error {

	_, err := DB.Exec(`
	UPDATE approval
	SET valid_until = NULL,
	    updated_at = CURRENT_TIMESTAMP
	WHERE id = 1
	`)

	return err
}

func IsValid() bool {

	var validUntil *time.Time

	err := DB.QueryRow(`
		SELECT valid_until
		FROM approval
		WHERE id = 1
	`).Scan(&validUntil)

	if err != nil || validUntil == nil {
		return false
	}

	return time.Now().Before(*validUntil)
}
