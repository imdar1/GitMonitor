package db

import (
	"database/sql"
	"gitmonitor/constants"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type DBConfig struct {
	Driver *sql.DB
	Close  func() error
}

func (d *DBConfig) InitDB() (DBConfig, error) {
	db, err := sql.Open(constants.DB_DRIVER, constants.DB_PATH)
	if err != nil {
		return DBConfig{}, err
	}

	statement, err := db.Prepare(constants.INIT_PROJECT_TABLE)
	if err != nil {
		return DBConfig{}, err
	}
	statement.Exec()

	statement, err = db.Prepare(constants.INIT_BRANCH_TABLE)
	if err != nil {
		return DBConfig{}, err
	}
	statement.Exec()

	statement, err = db.Prepare(constants.INIT_TASK_TABLE)
	if err != nil {
		return DBConfig{}, err
	}
	statement.Exec()

	dbConfig := DBConfig{
		Driver: db,
		Close: func() error {
			return db.Close()
		},
	}
	return dbConfig, nil
}
