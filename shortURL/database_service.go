package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func ReadConfig() (config string) {
	DbDriver := os.Getenv("DB_driver")
	DbUsername := os.Getenv("DB_username")
	DbPassword := os.Getenv("DB_password")
	DbHost := os.Getenv("DB_host")
	DbPort := os.Getenv("DB_port")
	DbName := os.Getenv("DB_name")
	DbSslmode := os.Getenv("DB_sslmode")

	config = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		DbDriver, DbUsername, DbPassword, DbHost, DbPort, DbName, DbSslmode)
	return
}
func NewDB(dbSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dbSourceName)
	if err != nil {
		return nil, err
	}
	// check connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected!")
	return &DB{db}, nil
}
