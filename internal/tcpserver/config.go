package tcpserver

type Config struct {
	Port int `env:"SERVER_PORT,required"`
}
