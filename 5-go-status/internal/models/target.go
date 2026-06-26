package models

import (
	"context"
	"database/sql"
	"time"
)

type Target struct {
	Id           int
	Url          string
	Interval_sec int
	Contact_info string
	Is_active    bool
	Created_at   time.Time
}

func CreateTableTarget(db *sql.DB) error {

	query := `CREATE TABLE IF NOT EXISTS target (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			interval_sec INTEGER NOT NULL,
			contact_info TEXT,
			is_active BOOLEAN DEFAULT true,
			created_at DATETIME NOT NULL
		)`

	_, err := db.ExecContext(context.Background(), query)
	return err

}

func InsertTarget(db *sql.DB, t *Target) (*Target, error) {
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO target (url, interval_sec, contact_info, is_active, created_at) VALUES (?,?,?,?,?) RETURNING target.id`,
		t.Url, t.Interval_sec, t.Contact_info, t.Is_active, t.Created_at).
		Scan(&t.Id)
	return t, err

}

func SelectAllTargets(db *sql.DB) ([]Target, error) {
	var targets []Target
	rows, err := db.QueryContext(
		context.Background(),
		`SELECT * FROM target`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Target
		if err := rows.Scan(&t.Id, &t.Url, &t.Interval_sec, &t.Contact_info, &t.Is_active, &t.Created_at); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, rows.Err()
}
