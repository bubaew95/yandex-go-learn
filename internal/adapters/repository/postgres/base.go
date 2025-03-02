package postgres

import "database/sql"

func dbConnect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}
