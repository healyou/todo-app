package main

import (
	"fmt"
	"log"
	"os"
	"todo/src/controllers"
	"todo/src/di"
	"todo/src/environment"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func main() {
	/* Грузим переменые окружения */
	environment.GetEnvVariables()
	/* Логирование */
	environment.InitLogFileForPROD()
	/* Инициализируем одиночек */
	di.GetInstance()

	// minioExample()
	var router = controllers.SetupRouter()
	err := router.Run(":8222")
	if err != nil {
		err = errors.Wrap(err, "ошибка запуска приложения")
		log.Println(fmt.Printf("%+v", err))
		os.Exit(1)
	}
}
