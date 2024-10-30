package database

import "database/sql"

func GetDBFromQueries(q *Queries) *sql.DB {
	return q.db.(*sql.DB)
}
