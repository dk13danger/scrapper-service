package main

import (
	"flag"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/dk13danger/scrapper-service/config"
	"github.com/dk13danger/scrapper-service/server"
	"github.com/dk13danger/scrapper-service/service"
)

var cfgFile = flag.String("config", "cfg/config.yml", "path to config (default: cfg/config.yml)")

func main() {
	flag.Parse()
	cfg := config.MustInit(*cfgFile)

	logger := logrus.New()
	if os.Getenv("DEBUG_MODE") == "true" {
		logger.Level = logrus.DebugLevel
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	srv := service.NewService(logger)
	web := server.NewServer(srv, logger, &cfg.Server)
	web.Run()
}
