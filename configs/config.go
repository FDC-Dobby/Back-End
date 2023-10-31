package configs

import (
	"fmt"
	"os"
	"strconv"
)

type MainConfig struct {
	JWT_SECRET string
	JWT_EXPIRE int
}

func getEnvStr(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		fmt.Printf("Environment variable %s is not set, Keep Going with Default Value '%s' \n", key, defaultValue)
		return defaultValue
	}
	return
}

func getEnvInt(key string, defaultValue int) (intValue int) {
	value := getEnvStr(key, strconv.Itoa(defaultValue))
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return
}

func getAllEnv() MainConfig {
	return MainConfig{
		JWT_SECRET: getEnvStr("JWT_SECRET", "d495ce948d89f228cf4e"),
		JWT_EXPIRE: getEnvInt("AUTH_JWT_EXPIRE", 3600),
	}
}

var Config = getAllEnv()
