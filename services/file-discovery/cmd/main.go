package main

import (
	"fmt"
	"time"
	"vega_file_discovery/common/config"
	DB "vega_file_discovery/packages/infrastrcuture/database"

	"github.com/abaxoth0/Vega/libs/go/packages/logger"
)

var log = logger.NewSource("MAIN", logger.Default)

func main() {
	if err := logger.Default.NewForwarding(logger.Stdout); err != nil {
		panic(err.Error())
	}

	config.Init()
	logger.SetServiceInstance(config.App.ServiceID)
	logger.SetServiceName("file_discovery")
	logger.Default.Init()

	logger.Debug.Store(config.Debug.Enabled)
	logger.Trace.Store(config.App.TraceLogsEnabled)

	go func() {
		if err := logger.Default.Start(config.Debug.Enabled); err != nil {
			panic(err.Error())
		}
	}()
	defer func() {
		if err := logger.Default.Stop(); err != nil {
			log.Error("Failed to stop logger", err.Error(), nil)
		}
	}()

	// Reserve some time for logger to start up
	time.Sleep(time.Millisecond * 50)

	if err := DB.Database.Connect(); err != nil {
		panic(err)
	}
	if err := DB.Database.Disconnect(); err != nil {
		panic(err)
	}

	fmt.Println("DONE")
	x := ""
	fmt.Scan(&x)
}
