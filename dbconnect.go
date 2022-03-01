package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func connectToDB() *sql.DB {
	mysqlChan := make(chan *sql.DB, 1)

	go func() {

		sqlHost := os.Getenv("DB_HOST")
		sqlDSN := fmt.Sprintf("mysql:example@tcp(%s:3307)/library", sqlHost)
		cruddb, err := sql.Open("mysql", sqlDSN)
		if err != nil {
			log.Fatal(fmt.Errorf("error connecting to mysql db %+v", err))
		}

		log.Info("Pinging the mySQL")
		for {
			if cruddbErr := cruddb.Ping(); cruddbErr != nil {
				log.Errorf("an error occurred connecting to the mySQL db trying again in 20 seconds: %v\n", cruddbErr)
				time.Sleep(time.Second * 20)
			} else {
				log.Info("connected to mySQL db")
				break
			}
		}

		mysqlChan <- cruddb
	}()

	// wait for DB to be setup
	// for loop / channel is for appyling multiple db connections if needed
	var mysqlDB *sql.DB

	done := false

	for !done {
		log.Info("done is ", done)
		select {
		case db := <-mysqlChan:
			if mysqlDB == nil {
				mysqlDB = db
				if db != nil {
					done = true
				}
			}
			log.Info("done is ", done)
		}
	}

	// instantiate a new mysql connection
	return mysqlDB
}
