package models

import (
	"context"
	"database/sql"
	"time"
)

type Probe struct {
	Id          int
	Target_id   int
	Status_code int
	Latency_ms  int
	Err_msg     string
	Timestamp   time.Time
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

func InsertProbe(db *sql.DB, p *Probe) (*Probe, error) {
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO probe (target_id, status_code, latency_ms, error, timestamp) VALUES (?,?,?,?,?) RETURNING probe.id`,
		p.Target_id, p.Status_code, p.Latency_ms, p.Err_msg, p.Timestamp).
		Scan(&p.Id)
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
		if err := rows.Scan(&p.Id, &p.Target_id, &p.Status_code, &p.Latency_ms, &p.Err_msg, &p.Timestamp); err != nil {
			return nil, err
		}
		probes = append(probes, p)
	}
	return probes, rows.Err()
}
