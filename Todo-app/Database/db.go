package Database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var DBConnection *sqlx.DB

/*
	1) No need of init
	2) Keep the db credential in environment variable and use it using os.Getenv()
	3) create a separate function for open and closing db and also for migration and transaction
*/

func init() {
	//var err error
	dbs := "host=localhost port=5433 user=local password=local dbname=todo sslmode=disable"
	DB, err := sqlx.Open("postgres", dbs)
	if err != nil {
		log.Fatal("error in db connection")
	}
	DBConnection = DB

}

func CloseDatabase() error {
	return DBConnection.Close()
}
