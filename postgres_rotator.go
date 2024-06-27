package pwrotate

import (
	"fmt"
	"database/sql"
)

type PostgresRotator struct {
	Host string
	Port int
}

func (r PostgresRotator) Rotate(user string, oldPw string, newPw string) error {
	// Check if password was already rotated
	if err := r.Ping(user, newPw); err == nil {
		return ErrAlreadyRotated
	}
	// Rotate old to new password
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

// Creating a new connection will ping to see if auth is valid
func (r PostgresRotator) newConnection(user string, pw string) (db *sql.DB, err error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=disable", user, pw, r.Host, r.Port)
	db, err = sql.Open("postgres", connStr)
	if db != nil {
		if err = db.Ping(); err != nil {
			db = nil
		}
	}
	return
}

