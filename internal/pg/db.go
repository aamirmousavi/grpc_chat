package pg

import (
	"database/sql"

	_ "github.com/lib/pq"

	"os"
)

var Hand = connect("host=127.0.0.1 port=5432 user=postgres password=admin dbname=postgres sslmode=disable")

type Postgres struct {
	Db *sql.DB
}

func connect(uri string) *Postgres {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		println("pg connection Error: ", err.Error())
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		println("unable to send ping to db Error: ", err.Error())
	}
	return &Postgres{
		Db: db,
	}
}
