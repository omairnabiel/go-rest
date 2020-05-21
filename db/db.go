package db

import (
	"database/sql"

	_ "github.com/lib/pq" // postgres driver
	"github.com/volatiletech/sqlboiler/boil"
)

func init() {
	db, err := sql.Open("postgres", "host=localhost dbname=postgres user=postgres password=postgres")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	boil.SetDB(db)

}
