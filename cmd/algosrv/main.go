package main

import (
	"log"
	"os"
)

var Cleanfuncs *CleanUpHandlerArray
var Logger *log.Logger

func main() {
	Logger = log.New(os.Stdout, "", log.LstdFlags)
	Cleanfuncs = NewCleanUpHandlerArray(Logger)

	defer func() {
		if Cleanfuncs != nil {
			Cleanfuncs.Wait()
		}
	}()
}
