package main

import (
	"log"
	"strconv"
	"os"
	"strings"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	engine := strings.ToUpper(getEnv("ENGINE"))
	if engine == "" {
		log.Fatal("ENGINE cannot be empty")
	}

	seedPassword := getEnv(engine + "_SEED_PASSWORD")
	if seedPassword == "" {
		log.Println(engine + "_SEED_PASSWORD is empty, skipping password rotation")
		return
	}

	user := getEnv(engine + "_USER")
	if user == "" {
		log.Fatal(engine + "_USER cannot be empty")
	}

	password := getEnv(engine + "_PASSWORD")
	if password == "" {
		log.Fatal(engine + "_PASSWORD cannot be empty")
	}

	host := getEnv(engine + "_HOST")
	if host == "" {
		log.Fatal(engine + "_HOST cannot be empty")
	}

	port, err := strconv.Atoi(getEnv(engine + "_PORT"))
	if err != nil {
		log.Fatalf("invalid port %d", port)
		return
	}
	if port == 0 {
		log.Fatal(engine + "_PORT cannot be 0")
	}

	var rotator Rotator
	switch engine {
	case "POSTGRES":
		rotator = PostgresRotator{
			Host: host,
			Port: port,
		}
	default:
		log.Fatalf("invalid engine %s", engine)
	}

	if err = rotator.Rotate(user, seedPassword, password); err != nil {
		if err == ErrNotAuthenticated {
			log.Println("Skipping password rotation due to authentication failure. This is expected if password was already rotated.")
			return
		}
		log.Fatalf("failed to rotate password: %s", err)
	}
}

func getEnv(name string) string {
	key := os.Getenv(name + "_KEY")
	return os.Getenv(key)
}

