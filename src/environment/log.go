package environment

import (
	"io"
	"log"
	"os"
	"todo/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

/* для профиля PROD устанавливаем файл для логирования */
func InitLogFileForPROD() {
	if (os.Getenv(utils.ProfileEnvName) == "PROD") {
		log.Println("PROD профиль, открываем лог файл.")
		logFilePathWithName := os.Getenv(utils.LOG_FILE_PATH_WITH_NAME)

		if (logFilePathWithName != "") {
			logFile, err := os.OpenFile(logFilePathWithName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				log.Println(errors.Wrap(err, "ошибка открытия лог файла"))
				log.Println("не получилось открыть файл логов - логируем только в консоль")
			} else {
				log.Println("открыт лог файл - " + logFilePathWithName)
				multiWriter := io.MultiWriter(logFile, os.Stdout)

				/* Логируем в файл, всё, что идёт в консоль и весь лог gin */
				gin.DefaultWriter = multiWriter
				gin.DefaultErrorWriter = multiWriter
				log.SetOutput(multiWriter)
				log.SetFlags(log.Lshortfile | log.LstdFlags)
			}
		} else {
			log.Println("не указан параметр '" + utils.LOG_FILE_PATH_WITH_NAME + "' - логируем в консоль")
		}
	} else {
		log.Println("используется профиль запуска, отличный от PROD. логируем в консоль")
	}
}