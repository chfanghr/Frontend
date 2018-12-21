package main

import (
	"flag"
	"log"
	"os"
)

var ServiceName = "backendService"
var CleanUpfuncs *CleanUpHandlerArray
var Logger *log.Logger

func main() {
	configFile := flag.String("configFile", "", "path to config file")
	logFilePath := flag.String("logFile", "", "path to log file")
	closeStdio := flag.Bool("closeStdio", false, "daemon close stdio or not")
	flag.StringVar(&ServiceName, "serviceName", ServiceName, "name of rpc service")
	Logger = log.New(os.Stdout, "", log.LstdFlags)

	Logger, err := SetupLogger(*logFilePath, !*closeStdio)
	if err != nil {
		log.Fatalln("error occur while setting up logger :", err)
	}
	CleanUpfuncs = NewCleanUpHandlerArray(Logger)
	err = LoadConfigFile(*configFile)
	if err != nil {
		log.Fatalln("error occur while loading config file :", err)
	}
	defer func() {
		if CleanUpfuncs != nil {
			CleanUpfuncs.Wait()
		}
	}()
}
