package app

import(
	"os"
	"strconv"
)


func GetEnvWithFallback(env string, fallback string) string {
	val := os.Getenv(env)
	if len(val) == 0 {
		return fallback
	}
	return val
}

func StringToInt(val string, fallback int) int {
	if res, err := strconv.Atoi(val); err != nil {
		return fallback
	} else {
		return res
	}
}

