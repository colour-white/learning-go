package models

import (
	"context"
	"database/sql"
	"time"
)

type Target struct {
	id           int
	url          string
	interval_sec int
	contact_info string
	is_active    bool
	created_at   time.Time
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

func InsertTarget(db *sql.DB, url string, interval_sec int, contact_info string, is_active bool) (Target, error) {
	var t Target
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO target (url, interval_sec, contact_info, is_active, created_at) VALUES (?,?,?,?,?) RETURNING *`,
		url, interval_sec, contact_info, is_active, time.Now()).
		Scan(&t.id, &t.url, &t.interval_sec, &t.contact_info, &t.is_active, &t.created_at)
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
		if err := rows.Scan(&t.id, &t.url, &t.interval_sec, &t.contact_info, &t.is_active, &t.created_at); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, rows.Err()
}
