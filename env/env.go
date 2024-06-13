package env

import (
	"os"
	"log"
	"fmt"
)

func Getenv(name string, engine string) string {
	key := os.Getenv(fmt.Sprintf("%s_%s_KEY", engine, name))
	if key == "" {
		key = fmt.Sprintf("%s_%s", engine, name)
	}
	return os.Getenv(key)
}

func GetenvOrFatal(name string, engine string) string {
	v := Getenv(name, engine)
	if v == "" {
		log.Fatalf("%s_%s cannot be empty", engine, name)
	}
	return v
}

func GetenvWithDefault(name string, defaultValue string, engine string) string {
	value := Getenv(name, engine)
	if value == "" {
		value = defaultValue
	}
	return value
}

