package main

import (
	"log"

	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/config"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/logger"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/server"
)

func main() {
	conf, err := config.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	logger, err := logger.SetupLogger(conf.Env, "")
	if err != nil {
		log.Fatalln(err)
	}

	s := server.NewServer(conf, logger)

	if err = s.Run(); err != nil {
		log.Fatalln(err)
	}
}