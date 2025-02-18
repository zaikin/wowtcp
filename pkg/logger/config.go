package logger

type Config struct {
	Level   string `env:"LOGGER_LEVEL,required"`
	Caller  bool   `env:"LOGGER_ENABLE_CALLER,required"`
	Console bool   `env:"LOGGER_ENABLE_CONSOLE,required"`
}
