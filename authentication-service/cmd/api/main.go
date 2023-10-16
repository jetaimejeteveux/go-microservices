package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var count = 0

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	fmt.Println("starting authentication service")
	conn := connecToDB()
	if conn == nil {
		log.Panic("cant connect to DB")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connecToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println(" postgres not ready yet ... ")
			count++
		} else {
			log.Println(" connected to postgres ! ")
			return conn
		}
		if count > 10 {
			log.Println(err)
			return nil
		}
		log.Println("backing of for two secs")
		time.Sleep(2 * time.Second)
		continue

	}

}
