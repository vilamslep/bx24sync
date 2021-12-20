package bx24sync

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func LoadEnv(path string) (err error) {

	if path == "" {
		path = ".env"
	}

	if r, err := os.OpenFile(path, os.O_RDONLY, 0666); err == nil {
		return load(r)
	} else {
		return err
	}
}

func load(f *os.File) (err error) {

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := scanner.Text()

		envData := strings.Split(text,"=")
		
		if len(envData) == 0 {
			continue
		}

		var key, value string 

		if len(envData) > 0 {
			key = envData[0]
		}

		if key == "" { continue }

		if len(envData) > 1 {
			value = envData[1]
		}

		if err := os.Setenv(key,value); err != nil {
			return err
		}
	}
	return err
}

func GetEnvAsInt(key string, defaultValue int) int {
	value := getEnv(key,"")
	if value, err := strconv.Atoi(value); err == nil {
		return value
	}
	return defaultValue
}

func GetEnvAsBool(key string, defaultValue bool) bool {
	value := getEnv(key, "")
	if value, err := strconv.ParseBool(value); err == nil {
		return value
	}

	return defaultValue
}

func GetEnvAsString(key string, defaultValue string) string {
	return getEnv(key, defaultValue)
}

func GetEnvAsStringSlice(key string, sep string, defaultValue []string) ([]string) {

	value := getEnv(key, "")

	if value == "" { return defaultValue }

	return strings.Split(value, sep)
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
