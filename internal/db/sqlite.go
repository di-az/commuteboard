package db

import (
	"commuteboard/internal/domain"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

func NewSQLite(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	return db
}

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS routes (
		id TEXT PRIMARY KEY,
		origin_id TEXT NOT NULL,
		destination_id TEXT NOT NULL,
		distance_meters INTEGER NOT NULL,
		duration_seconds INTEGER NOT NULL,
		recorded_at TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

func InsertCommute(db *sql.DB, c domain.Commute) error {
	query := `
	INSERT OR REPLACE INTO routes
	(id, origin_id, destination_id, distance_meters, duration_seconds, recorded_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(
		query,
		c.ID,
		c.OriginID,
		c.DestinationID,
		c.DistanceMeters,
		int(c.Duration.Seconds()),
		c.RecordedAt.Format(time.RFC3339),
	)

	return err
}

func GetAllCommutes(db *sql.DB) ([]domain.Commute, error) {
	rows, err := db.Query(`SELECT id, origin_id, destination_id, distance_meters, duration_seconds, recorded_at FROM routes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commutes []domain.Commute

	for rows.Next() {
		var c domain.Commute
		var durationSeconds int
		var recordedAt string

		err := rows.Scan(
			&c.ID,
			&c.OriginID,
			&c.DestinationID,
			&c.DistanceMeters,
			&durationSeconds,
			&recordedAt,
		)
		if err != nil {
			return nil, err
		}

		c.Duration = time.Duration(durationSeconds) * time.Second
		c.RecordedAt, _ = time.Parse(time.RFC3339, recordedAt)

		commutes = append(commutes, c)
	}

	return commutes, nil
}
