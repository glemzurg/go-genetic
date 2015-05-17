package genetic

import (
	"database/sql"
	_ "github.com/ziutek/mymysql/godrv" // Makes driver present.
	"log"
)

const (
	_DB_NAME      = "genetic" // Mysql database name.
	_DB_USER_NAME = "genetic" // Mysql user name.
	_DB_USER_PASS = "genetic" // Mysql user password.
)

// newDatabaseConnection connects to the local database.
func newDatabaseConnection() *sql.DB {
	var err error

	// Datasource Name Format docs: http://localhost:6060/pkg/github.com/ziutek/mymysql/godrv/
	var dataSourceName string = _DB_NAME + "/" + _DB_USER_NAME + "/" + _DB_USER_PASS
	var db *sql.DB
	if db, err = sql.Open("mymysql", dataSourceName); err != nil {
		log.Panic(err)
	}
	return db
}
