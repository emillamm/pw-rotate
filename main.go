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

	engine := strings.ToUpper(os.Getenv("PW_ROTATE_ENGINE"))
	if engine == "" {
		log.Fatal("PW_ROTATE_ENGINE cannot be empty")
	}

	seedPassword := getEnv("SEED_PASSWORD", engine)
	if seedPassword == "" {
		log.Println(engine + "_SEED_PASSWORD is empty, skipping password rotation")
		return
	}

	user := getEnv("USER", engine)
	if user == "" {
		log.Fatal(engine + "_USER cannot be empty")
	}

	password := getEnv("PASSWORD", engine)
	if password == "" {
		log.Fatal(engine + "_PASSWORD cannot be empty")
	}

	host := getEnvWithDefault("HOST", engine)
	if host == "" {
		log.Fatal(engine + "_HOST cannot be empty")
	}

	port, err := strconv.Atoi(getEnvWithDefault("PORT", engine))
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
	log.Printf("%s password successfully rotated for user %s\n", engine, user)
}

func getEnv(name string, engine string) string {
	key := os.Getenv(engine + "_" + name + "_KEY")
	return os.Getenv(key)
}

func getEnvWithDefault(name string, engine string) string {
	value := getEnv(name, engine)
	if value == "" {
		value = os.Getenv("DEFAULT_" + engine + "_" + name)	
	}
	return value
}

