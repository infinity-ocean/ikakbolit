package config

type Log struct {
	LoggingLevel    string `env:"LOGGING_LEVEL"`
	LoggingFilePath string `env:"LOGGING_FILE_PATH"`
}
