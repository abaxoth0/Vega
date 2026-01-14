package app

import (
	"fmt"
	"os"
	"runtime"
	"vega_file_discovery/common/config"
	DB "vega_file_discovery/packages/infrastrcuture/database"

	"github.com/abaxoth0/Vega/libs/go/packages/logger"
)

func StartInit() {
	// All init logs will be shown anyway
	if err := logger.Default.NewForwarding(logger.Stdout); err != nil {
		panic(err.Error())
	}
}

func EndInit() {
	if !config.App.ShowLogs {
		if err := logger.Default.NewForwarding(logger.Stdout); err != nil {
			panic(err.Error())
		}
	}
}

func InitDefault() {
	if runtime.GOOS != "linux" {
		fmt.Println("[ CRITICAL ERROR ] OS is not supported. This program can be used only on Linux-based OS.")
		os.Exit(1)
	}

	config.Init()
	logger.SetServiceInstance(config.App.ServiceID)
	logger.SetServiceName("file_discovery")
	logger.Default.Init()
}

func InitConnections() {
	log.Info("Initializng connections...", nil)

	DB.Database.Connect()

	log.Info("Initializng connections: OK", nil)
}
