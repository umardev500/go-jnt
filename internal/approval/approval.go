package approval

import (
	"database/sql"
	"errors"
	"time"
)

type Service struct {
	db *sql.DB
}

var ErrDaytimeOnly = errors.New("operation is only available during daytime")

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

// TodayCreated returns the number of manifests created today.
func (s *Service) TodayCreated() (int, error) {
	var hour int

	err := s.db.QueryRow(`
		SELECT CAST(strftime('%H', 'now', '+7 hours') AS INTEGER)
	`).Scan(&hour)
	if err != nil {
		return 0, err
	}

	if hour < 6 || hour >= 18 {
		return 0, ErrDaytimeOnly
	}

	var total int

	err = s.db.QueryRow(`
		SELECT COALESCE(
			(
				SELECT total_created
				FROM daily_manifest_stats
				WHERE stat_date = DATE('now', '+7 hours')
			),
			0
		)
	`).Scan(&total)

	return total, err
}

// IncrementTodayCreated inserts or increments today's manifest count.
func (s *Service) IncrementTodayCreated() error {
	_, err := s.db.Exec(`
		INSERT INTO daily_manifest_stats (
			stat_date,
			total_created
		)
		VALUES (
			DATE('now', '+7 hours'),
			1
		)
		ON CONFLICT(stat_date)
		DO UPDATE SET
			total_created = total_created + 1,
			updated_at = CURRENT_TIMESTAMP
	`)

	return err
}
