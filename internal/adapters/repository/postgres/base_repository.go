package postgres

import "database/sql"

func dbConnect(connStr string) *sql.DB {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	return db
}
