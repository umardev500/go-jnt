package approval

import (
	"database/sql"
	"time"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Grant(minutes int) error {
	_, err := s.db.Exec(`
		UPDATE approval
		SET valid_until = datetime('now', ? || ' minutes'),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`, minutes)

	return err
}

func (s *Service) Revoke() error {
	_, err := s.db.Exec(`
		UPDATE approval
		SET valid_until = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`)

	return err
}

func (s *Service) IsValid() bool {
	var validUntil *time.Time

	err := s.db.QueryRow(`
		SELECT valid_until
		FROM approval
		WHERE id = 1
	`).Scan(&validUntil)

	if err != nil || validUntil == nil {
		return false
	}

	return time.Now().Before(*validUntil)
}
