package main

import (
	"encoding/json"
	"errors"
	"github.com/chfanghr/Backend/car"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	cs []car.Client
}

var config Config

func SetupLogger(logFilePath string, useStdio bool) (logger *log.Logger, err error) {
	var logFile, stdio *os.File
	if len(logFilePath) > 0 {
		logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatalln("open log file :", err)
			return nil, err
		}
	}
	if useStdio {
		stdio = os.Stdout
	}
	var w io.Writer = nil
	if len(logFilePath) > 0 {
		if useStdio {
			w = io.MultiWriter(logFile, stdio)
		} else {
			w = logFile
		}
	} else {
		if useStdio {
			w = stdio
		} else {
			f, _ := os.Open("/dev/null")
			w = f
		}
	}
	logger = log.New(w, "", log.LstdFlags)
	return
}

func LoadConfigFile(file string) error {
	type networkConfig struct {
		Type    string `json:"type"`
		Address string `json:"address"`
	}
	type networkConfigs struct {
		cars []networkConfig `json:"cars"`
	}
	var n networkConfigs
	if len(file) > 0 {
		Logger.Println("load from config file")
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		err = json.Unmarshal(buf, &n)
		if err != nil {
			return err
		}
		for _, v := range n.cars {
			config.cs = append(config.cs, car.NewGeneralClient(v.Type, v.Address, ServiceName))
		}
		for _, v := range config.cs {
			err := v.IsServiceAvailable()
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("invalid config file")
	}
}
