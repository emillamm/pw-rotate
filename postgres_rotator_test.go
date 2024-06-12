package pwrotate

import (
	"testing"
	"fmt"
	"math/rand"
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
	"github.com/emillamm/pwrotate/env"
)

func TestPostgresRotator(t *testing.T) {

	engine := "POSTGRES"
	user := env.GetenvWithDefault("USER", "postgres", engine)
	password := env.GetenvWithDefault("PASSWORD", "postgres", engine)
	host := env.GetenvWithDefault("HOST", "localhost", engine)
	portStr := env.GetenvWithDefault("PORT", "5432", engine)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Errorf("invalid PORT %s", portStr)
		return
	}

	tmpPassword := "tmpPassword"
	newPassword := "tmpPassword2"

	// Set up connection
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=disable", user, password, host, port)
	db, err := openConnection(connStr)
	if err != nil {
		t.Errorf("failed to open connection %s", err)
		return
	}
	defer db.Close()

	// Create ephemeral test user
	ephemeralPostgresUser(t, db, tmpPassword, func(tmpUser string) {
		rotator := PostgresRotator{Host : host, Port: port}
		// run tests
		RotatorTest(t, rotator, tmpUser, tmpPassword, newPassword)
	})
}

func openConnection(connStr string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connStr)
	if db != nil {
		err = db.Ping()
	}
	return
}

func ephemeralPostgresUser(
	t testing.TB,
	parentSession *sql.DB,
	password string,
	block func(tmpUser string),
) {
	t.Helper()
	var err error

	user := randomString(5)

	createRoleQ := fmt.Sprintf("create role %s with login password '%s';", user, password)
	if _, err = parentSession.Exec(createRoleQ); err != nil {
		t.Errorf("failed to create role %s: %s", user, err)
		return
	}

	createDbQ := fmt.Sprintf("create database %s;", user)
	if _, err = parentSession.Exec(createDbQ); err != nil {
		t.Errorf("failed to create database %s: %s", user, err)
		return
	}

	defer func() {
		dropRoleQ := fmt.Sprintf("drop role %s;", user)
		_, err = parentSession.Exec(dropRoleQ)
		if err != nil {
			t.Errorf("failed to drop role database %s: %s", user, err)
		}

		dropDbQ := fmt.Sprintf("drop database %s;", user)
		_, err = parentSession.Exec(dropDbQ)
		if err != nil {
			t.Errorf("failed to drop database %s: %s", user, err)
		}
	}()

	block(user)
}

func randomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

