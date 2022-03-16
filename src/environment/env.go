package environment

import (
	"log"
	"os"
	"sync"
	"todo/src/utils"

	"github.com/joho/godotenv"
)

var env EnvVariables
var onceLoad sync.Once

func GetEnvVariables() EnvVariables {
	onceLoad.Do(func() {
		loadEnvironmentVariables()
	})
	return env
}

func loadEnvironmentVariables() {
	profile := os.Getenv(utils.ProfileEnvName)
	log.Println("используется профиль запуска - " + profile)

	if profile == "DEV" || profile == "TEST" {
		loadDevEnv()
	} else if profile == "PROD" {
		/* already in env */
	} else {
		log.Fatalln("не указан профиль запуска")
	}

	env.MinioEndpoint = getEnvWithCheckExists(utils.MinioEndpointEnvName)
	env.MinioAccessKey = getEnvWithCheckExists(utils.MinioAccessKeyEnvName)
	env.MinioSecretKey = getEnvWithCheckExists(utils.MinioSecretKeyEnvName)
	env.MySqlDriverName = getEnvWithCheckExists(utils.MySqlDriverNameEnvName)
	env.MySqlDataSource = getEnvWithCheckExists(utils.MySqlDataSourceEnvName)
}

func getEnvWithCheckExists(envVarName string) string {
	value := os.Getenv(envVarName)
	if value == "" {
		log.Fatalln("не указан параметр в .env - " + envVarName)
	}
	return value
}

func loadDevEnv() {
	err := godotenv.Load("../../profile_dev.env")
	if err != nil {
		log.Println("ошибка загрузки переменных окружения")
		log.Fatalln(err)
	}
}

type EnvVariables struct {
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MySqlDriverName string
	MySqlDataSource string
}
