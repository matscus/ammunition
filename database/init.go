package database

import (
	"database/sql"
	"errors"
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
	DB, err = sql.Open("postgres", "host="+os.Getenv("POSTGRESHOST")+" port="+os.Getenv("POSTGRESPORT")+" user="+os.Getenv("POSTGRESUSER")+" password="+os.Getenv("POSTGRESPASSWORD")+" dbname="+os.Getenv("AMMUNITIONBD")+" sslmode=disable")
	if err != nil {
		log.Println("Database not avalible: ", err)
	}
	err = firstStartInit()
	if err != nil {
		log.Println("First Start Init error: ", err)
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

func firstStartInit() (err error) {
	_, err = DB.Exec("CREATE SCHEMA IF NOT EXISTS system AUTHORIZATION " + os.Getenv("POSTGRESUSER"))
	if err != nil {
		return errors.New("Create scheme error: " + err.Error())
	}
	_, err = DB.Exec("create table if not exists system.tDatapools (id SERIAL PRIMARY key,project varchar(50), scriptname varchar(50),CONSTRAINT progectScripts UNIQUE (project, scriptname))")
	if err != nil {
		return errors.New("Create table system.tDatapools error : " + err.Error())
	}
	return nil
}
