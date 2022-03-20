package environment

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"todo/src/utils"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
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
	if (profile != "") {
		log.Println("используется профиль запуска из ENV - " + profile)
	} else {
		/* try load from args */
		argProfile := flag.String(utils.ProfileEnvName, "", "профиль запуска")
		flag.Parse()

		if (*argProfile == "") {
			err := errors.New("не указан профиль запуска")
			log.Println(fmt.Printf("%+v", err))
			os.Exit(1)
		}

		log.Println("используется профиль запуска из args - " + *argProfile)
		profile = *argProfile
	}

	if profile == "DEV" || profile == "TEST" {
		loadDevEnv()
	} else if profile == "PROD" {
		/* other var already in env */
	} else {
		err := errors.New("неизвестный профиль запуска - " + profile)
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
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
		err := errors.New("не указан параметр в .env - " + envVarName)
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
	}
	return value
}

func loadDevEnv() {
	err := godotenv.Load("../../profile_dev.env")
	if err != nil {
		err = errors.Wrap(err, "ошибка загрузки переменных окружения")
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
	}
}

type EnvVariables struct {
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MySqlDriverName string
	MySqlDataSource string
}
