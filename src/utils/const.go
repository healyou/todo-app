package utils

const ProfileEnvName = "APP_PROFILE_ACTIVE"
const MinioEndpointEnvName = "MINIO_ENDPOINT"
const MinioAccessKeyEnvName = "MINIO_ACCESS_KEY"
const MinioSecretKeyEnvName = "MINIO_SECRET_KEY"
const MySqlDriverNameEnvName = "MY_SQL_DRIVER_NAME"
const MySqlDataSourceEnvName = "MY_SQL_DATASOURCE_URL"
/* Код header с содержашимся jwt токеном пользователя */
const JSON_IN_ACCESS_TOKEN_CODE = "X-Access-Token"
/* Код данных пользователя в gin context */
const GIN_CONTEXT_USER_AUTH_DATA = "USER_AUTH_DATA"
/* Только для PROD профиля */
const LOG_FILE_PATH_WITH_NAME = "LOG_FILE_PATH_WITH_NAME"