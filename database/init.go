package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func init() {
	println(os.Getenv("AMMUNITIONBD"))
	DB, err = sql.Open("postgres", "user="+os.Getenv("POSTGRESUSER")+" password="+os.Getenv("POSTGRESPASSWORD")+" dbname="+os.Getenv("AMMUNITIONBD")+" sslmode=disable")
	if err != nil {
		log.Println("Database not avalible: ", err)
	}
	go func() {
		for {
			err := DB.Ping()
			if err != nil {
				log.Println("Database not avalible")
			}
			time.Sleep(5 * time.Second)
		}
	}()
}
