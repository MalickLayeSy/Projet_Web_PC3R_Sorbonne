package database

import "database/sql"

import (
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/projet-pc3r")
	if err != nil {
		panic(err.Error())
		return nil
	} else {
		return db
	}
}
