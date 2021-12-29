package bx24sync

import (
	"fmt"
	"os"
	"testing"
)

func Test_LoadEnv_Success(t *testing.T) {

	rmEnv, err := createTestEnvValues()

	if err != nil {
		t.Fatal(err)
	}

	if err := LoadEnv("env"); err != nil {
		t.Error(err)
	}

	kHost := "172.16.10.10"
	if v, ok := os.LookupEnv("KAFKA_HOST"); ok {
		if v != kHost {
			t.Errorf("Expected value is %s", kHost)
		}
	}

	kPort := "13454"
	if v, ok := os.LookupEnv("KAFKA_PORT"); ok {
		if v != kPort {
			t.Errorf("Expected value is %s", kHost)
		}
	}

	kTopic := "changes"
	if v, ok := os.LookupEnv("KAFKA_TOPIC"); ok {
		if v != kTopic {
			t.Errorf("Expected value is %s", kTopic)
		}
	}

	hHost := "localhost"
	if v, ok := os.LookupEnv("HTTP_HOST"); ok {
		if v != hHost {
			t.Errorf("Expected value is %s", hHost)
		}
	}

	hPort := "8002"
	if v, ok := os.LookupEnv("HTTP_PORT"); ok {
		if v != hPort {
			t.Errorf("Expected value is %s", hPort)
		}
	}

	hCheckInput := "true"
	if v, ok := os.LookupEnv("HTTP_CHECK_INPUT"); ok {
		if v != hCheckInput {
			t.Errorf("Expected value is %s", hCheckInput)
		}
	}

	if err := rmEnv(); err != nil {
		t.Fatal(err)
	}
}




func createTestEnvValues() (rm func() error, err error) {

	file := ".env"

	env := make(map[string]string)
	env["KAFKA_HOST"] = "172.16.10.10"
	env["KAFKA_PORT"] = "13454"
	env["KAFKA_TOPIC"] = "changes"
	env["HTTP_HOST"] = "localhost"
	env["HTTP_PORT"] = "8002"
	env["HTTP_CHECK_INPUT"] = "true"

	if f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666); err == nil {
		defer f.Close()

		for k, v := range env {
			fmt.Fprintf(f, "%s=%s\n", k, v)
		}
		rm := func() error {
			return os.Remove(file)
		}
		return rm, err
	} else {
		return nil, err
	}
}
