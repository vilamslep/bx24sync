package bx24sync

import (
	"os"
	"testing"
)

func Test_LoadEnv_Success(t *testing.T) {

	if err := LoadEnv(""); err != nil {
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
}
