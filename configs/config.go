package configs

import (
	"os"
	"strconv"
)

type MainConfig struct {
	JWT_SECRET string
	JWT_EXPIRE int
}

func getEnv_s(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Environment variable " + key + " not found")
}

func getAllEnv() MainConfig {
	var err error
	jwtExpire, err := strconv.Atoi(getEnv_s("JWT_EXPIRE"))
	if err != nil {
		panic("Environment variable JWT_EXPIRE is not a number")
	}

	return MainConfig{
		JWT_SECRET: getEnv_s("JWT_SECRET"),
		JWT_EXPIRE: jwtExpire,
	}
}

var Config = getAllEnv()
