package storage

import (
	"database/sql"

	"go-status/internal/models"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDatabase(dbPath string)(*sql.DB, error) {

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	err = models.CreateTableTarget(db)
	if err != nil {
		return nil, err
	}

	err = models.CreateTableProbe(db)
	return db, err

}
