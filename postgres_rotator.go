package main

import (
	"fmt"
	"database/sql"
	"github.com/lib/pq"
)

type PostgresRotator struct {
	Host string
	Port int
}

func (r PostgresRotator) Rotate(user string, oldPw string, newPw string) error {
	db, err := r.newConnection(user, oldPw)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("alter user %s with password '%s';", user, newPw))
	return err
}

func (r PostgresRotator) Ping(user string, pw string) error {
	db, err := r.newConnection(user, pw)
	if db != nil {
		db.Close()
	}
	return err
}

func checkAuth(db *sql.DB) error {
	err := db.Ping()
	if e, ok := err.(*pq.Error); ok && e.Routine == "auth_failed" {
		return ErrNotAuthenticated
	}
	return err
}

func (r PostgresRotator) newConnection(user string, pw string) (db *sql.DB, err error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=disable", user, pw, r.Host, r.Port)
	db, err = sql.Open("postgres", connStr)
	if db != nil {
		if err = checkAuth(db); err != nil {
			db = nil
		}
	}
	return
}

