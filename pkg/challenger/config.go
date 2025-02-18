package challenger

type Config struct {
	Difficulty int `env:"CHALLENGE_DIFFICULTY,required"`
}
