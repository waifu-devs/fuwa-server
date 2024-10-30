package database

import "database/sql"

func GetDBFromQueries(q *Queries) *sql.DB {
	if tracedDB, ok := q.db.(*TracedDB); ok {
		return tracedDB.db
	}
	return q.db.(*sql.DB)
}
