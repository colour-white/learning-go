package models

import (
	"context"
	"database/sql"
	"time"
)

type Probe struct {
	id          int
	target_id   int
	status_code int
	latency_ms  int
	err_msg     string
	timestamp   time.Time
}

func CreateTableProbe(db *sql.DB) error {

	query := `CREATE TABLE IF NOT EXISTS probe (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			target_id INTEGER NOT NULL REFERENCES target(id),
			status_code INTEGER NOT NULL,
			latency_ms INTEGER NOT NULL,
			error TEXT,
			timestamp DATETIME NOT NULL
		)`

	_, err := db.ExecContext(context.Background(), query)
	return err

}

func InsertProbe(db *sql.DB, target_id, status_code, latency_ms int, err_msg string) (Probe, error) {
	var p Probe
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO probe (target_id, status_code, latency_ms, err_msg, timestamp) VALUES (?,?,?,?,?) RETURNING *`,
		target_id, status_code, latency_ms, err_msg, time.Now()).
		Scan(&p.id, &p.target_id, &p.status_code, &p.latency_ms, &p.err_msg, &p.timestamp)
	return p, err
}

func SelectAllProbes(db *sql.DB) ([]Probe, error) {
	var probes []Probe
	rows, err := db.QueryContext(
		context.Background(),
		`SELECT * FROM probe`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Probe
		if err := rows.Scan(&p.id, &p.target_id, &p.status_code, &p.latency_ms, &p.err_msg, &p.timestamp); err != nil {
			return nil, err
		}
		probes = append(probes, p)
	}
	return probes, rows.Err()
}
