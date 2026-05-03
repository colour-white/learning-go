package storage

import (
	"database/sql"

	"go-status/internal/models"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDatabase(dbPath string) error {

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	err = models.CreateTableTarget(db)
	if err != nil {
		return err
	}

	err = models.CreateTableProbe(db)
	return err

}
