package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

/**
* Create the connection pool
 */
func OpenPostgres() (*sql.DB, error) {

	poolSize, _ := strconv.Atoi(os.Getenv("DB_POOLSIZE"))
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), port, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", psqlconn)

	db.SetMaxOpenConns(poolSize)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return db, err
}

// You can stop the worker by closing the quit channel: close(quit)
func StartKeepAlive() {

	heartbeat, _ := strconv.Atoi(os.Getenv("DB_HEARTBEAT"))

	ticker := time.NewTicker(time.Duration(heartbeat) * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				//log.Printf("Ping...\n")

				if !DB_Ping() {
					fmt.Println("Error: Database unavailable")
					_connection = nil
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func ConnectAndKeepAlive() (*sql.DB, error) {

	log.Printf("Connecting to database %s:%s...\n", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	conn, err := OpenPostgres()

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	err = conn.Ping()

	if err == nil {

		_connection = conn

		log.Printf("Starting keep alive function...\n")
		StartKeepAlive()

	} else {
		fmt.Printf("Error: %s\n", err.Error())
		_connection = nil
	}

	return _connection, err
}
