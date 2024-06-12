package pwrotate

import (
	"log"
	"strconv"
	"os"
	"strings"
	"github.com/emillamm/pwrotate/env"
)

func Run() {
	engine := strings.ToUpper(os.Getenv("PW_ROTATE_ENGINE"))
	if engine == "" {
		log.Fatal("PW_ROTATE_ENGINE cannot be empty")
	}

	seedPassword := env.Getenv("SEED_PASSWORD", engine)
	if seedPassword == "" {
		log.Println(engine + "_SEED_PASSWORD is empty, skipping password rotation")
		return
	}

	user := env.GetenvOrFatal("USER", engine)
	password := env.GetenvOrFatal("PASSWORD", engine)
	host := env.GetenvWithDefault("HOST", "localhost", engine)
	portStr := env.GetenvWithDefault("PORT", "5432", engine)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid %s_PORT %s", engine, portStr)
	}
	if port == 0 {
		log.Fatalf("%s_PORT cannot be 0", engine)
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

