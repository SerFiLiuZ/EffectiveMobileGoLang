package main

import (
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/apiserver"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
)

func main() {
	logger := utils.NewLogger()
	logger.EnableDebug()

	apiserver.LoadEnv(logger)

	config := apiserver.GetConfig()

	logger.Debugf("config: %v", config)

	logger.Infof("Starting server...")

	if err := apiserver.Start(config, logger); err != nil {
		logger.Fatal(err.Error())
	}
}
