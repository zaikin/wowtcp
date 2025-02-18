package main

import (
	"flag"
	"time"
	"wowtcp/pkg/client"
	"wowtcp/pkg/logger"
)

func main() {
	host := flag.String("host", "localhost", "Host to connect to")
	port := flag.String("port", "8080", "Port to connect to")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	loggerCfg := logger.Config{Console: true, Caller: true, Level: "info"}
	if *debug {
		loggerCfg.Level = "debug"
	}

	log := logger.NewLogger(&loggerCfg)

	cli := client.NewClient(*host, *port)
	if err := cli.Connect(); err != nil {
		//nolint:gocritic
		log.Fatal().Err(err).Msg("Failed to connect to server")
	}
	defer cli.Close()

	for {
		quote, err := cli.RequestQuote()
		if err != nil {
			//nolint:gocritic
			log.Fatal().Err(err).Msg("Failed to request quote")
		}
		log.Info().Str("quote", quote).Msg("Received quote")
		//nolint:mnd
		time.Sleep(5 * time.Second)
	}
}
