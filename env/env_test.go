package env

import (
	"testing"
	"os"
)

func TestEnv(t *testing.T) {
	t.Run("Getenv should return an env var if it exists", func(t *testing.T) {
		varName := "TEST_VAR1"
		engine := "POSTGRES"
		if Getenv(varName, engine) != "" {
			t.Errorf("%s_%s should not exist", engine, varName)
		}
		os.Setenv("POSTGRES_TEST_VAR1", "abc")
		if v := Getenv(varName, engine); v != "abc" {
			t.Errorf("%s should be abc, was %s", varName, v)
		}
		os.Unsetenv("POSTGRES_TEST_VAR1")
	})
	t.Run("Getenv should use <VAR>_KEY as a pointer to the name of the actual variable", func(t *testing.T) {
		varName := "TEST_VAR2"
		engine := "POSTGRES"
		os.Setenv("POSTGRES_TEST_VAR2_KEY", "POSTGRES_TEST_VAR3")
		if v := "POSTGRES_TEST_VAR2"; Getenv(v, engine) != "" {
			t.Errorf("%s should not exist", v)
		}
		os.Setenv("POSTGRES_TEST_VAR3", "abc")
		if Getenv(varName, engine) != "abc" {
			t.Errorf("%s should be abc", varName)
		}
		os.Unsetenv("POSTGRES_TEST_VAR2_KEY")
		os.Unsetenv("POSTGRES_TEST_VAR3")
	})
	t.Run("GetenvWithDefault should return a default value if the value doesn't exist and otherwise return the value", func(t *testing.T) {
		varName := "TEST_VAR4"
		engine := "POSTGRES"
		os.Setenv("POSTGRES_TEST_VAR4", "abc")
		if v := GetenvWithDefault(varName, "xyz", engine); v != "abc" {
			t.Errorf("%s should be abc, was %s", varName, v)
		}
		if v := GetenvWithDefault("NON_EXISTING", "abc", engine); v != "abc" {
			t.Errorf("%s should be abc", "NON_EXISTING")
		}
		os.Unsetenv("POSTGRES_TEST_VAR4")
	})
}

