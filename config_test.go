package bx24sync

import (
	"os"
	"strconv"
	"testing"
)

func Test_Socket_String_Success(t *testing.T) {

	sock := Socket{Host: "127.0.0.1", Port: 9000}
	expectedValue := "127.0.0.1:9000"
	result := sock.String()
	if expectedValue != result {
		t.Errorf("Expected %s. Result %s", expectedValue, result)
	}

}

func Test_Endpoint_Url_Success(t *testing.T) {
	endpointWithSSl := Endpoint{
		Socket: Socket{Host: "127.0.0.1", Port: 9000},
		Method: "test",
		SSL:    true,
	}

	endpointWithoutSSl := Endpoint{
		Socket: Socket{Host: "127.0.0.1", Port: 9000},
		Method: "test",
		SSL:    false,
	}

	resultWithSSL := endpointWithSSl.URL()
	resultWithoutSSL := endpointWithoutSSl.URL()

	expectedValueWithSSL := "https://127.0.0.1:9000/test"
	expectedValueWithoutSSL := "http://127.0.0.1:9000/test"

	if expectedValueWithSSL != resultWithSSL {
		t.Errorf("Expected %s. Result %s", expectedValueWithSSL, resultWithSSL)
	}

	if expectedValueWithoutSSL != resultWithoutSSL {
		t.Errorf("Expected %s. Result %s", expectedValueWithoutSSL, resultWithoutSSL)
	}
}

func Test_BasicAuth_GetPair_Success(t *testing.T) {
	auth := BasicAuth{
		User:     "user1",
		Password: "password1",
	}

	expectedPair := "user1:password1"
	result := auth.GetPair(":")

	if expectedPair != result {
		t.Errorf("Expected %s. Result %s", expectedPair, result)
	}
}

func Test_NewRegistrarConfigFromEnv_Success(t *testing.T) {

	httpPort := 9000
	kafkaHost := "kafka"
	kafkaPort := 9090
	kafkaTopic := "testtopic"
	checkInput := true
	correctConfig := RegistrarConfig{
		Http: Socket{
			Port: httpPort,
		},
		ProducerConfig: ProducerConfig{
			Broker: Socket{
				Host: kafkaHost,
				Port: kafkaPort,
			},
			Topic: kafkaTopic,
		},
		CheckInput: checkInput,
	}

	if err := os.Setenv("HTTP_PORT", strconv.Itoa(httpPort)); err != nil {
		t.Fatalf("can't set HTTP_PORT. %v", err)
	}

	if err := os.Setenv("KAFKA_HOST", kafkaHost); err != nil {
		t.Fatalf("can't set KAFKA_HOST. %v", err)
	}

	if err := os.Setenv("KAFKA_PORT", strconv.Itoa(kafkaPort)); err != nil {
		t.Fatalf("can't set KAFKA_PORT. %v", err)
	}

	if err := os.Setenv("KAFKA_TOPIC", kafkaTopic); err != nil {
		t.Fatalf("can't set KAFKA_TOPIC. %v", err)
	}

	if err := os.Setenv("HTTP_CHECK_INPUT", "1"); err != nil {
		t.Fatalf("can't set HTTP_CHECK_INPUT. %v", err)
	}

	config := NewRegistrarConfigFromEnv()
	if config != correctConfig {
		t.Errorf("the config from env isn't equal correct config")
		t.Logf("the config from env %v", config)
		t.Logf("the correct config %v", correctConfig)
	}
}

func Test_NewDataBaseConnectionFromEnv_Success(t *testing.T) {
	dbHost := "127.0.0.1"
	dbPort := 1413
	dbUser := "user"
	dbPassword := "password"
	dbName := "database"
	correctConfig := DataBaseConnection{
		Socket:    Socket{Host: dbHost, Port: dbPort},
		BasicAuth: BasicAuth{User: dbUser, Password: dbPassword},
		Database:  dbName,
	}

	if err := os.Setenv("DB_HOST", dbHost); err != nil {
		t.Fatalf("can't set DB_HOST. %v", err)
	}

	if err := os.Setenv("DB_PORT", strconv.Itoa(dbPort)); err != nil {
		t.Fatalf("can't set DB_PORT. %v", err)
	}

	if err := os.Setenv("DB_USER", dbUser); err != nil {
		t.Fatalf("can't set DB_USER. %v", err)
	}

	if err := os.Setenv("DB_PASSWORD", dbPassword); err != nil {
		t.Fatalf("can't set HTTP_CHECK_INPUT. %v", err)
	}

	if err := os.Setenv("DB_BASE", dbName); err != nil {
		t.Fatalf("can't set DB_BASE. %v", err)
	}

	config := NewDataBaseConnectionFromEnv()
	if config != correctConfig {
		t.Errorf("the config from env isn't equal correct config")
		t.Logf("the config from env %v", config)
		t.Logf("the correct config %v", correctConfig)
	}
}
