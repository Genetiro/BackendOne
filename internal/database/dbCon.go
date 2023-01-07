package database

import (
	"database/sql"
)

func DbCon() error {
	db, err := sql.Open("sqlite3", "links.db")
	if err != nil {
		return err
	}

	defer db.Close()
	return nil
}
