package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vi-la-muerto/bx24-service/http/generator"
	"github.com/vi-la-muerto/bx24-service/scheme"
)

func main() {

	logFile := "app.log"

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.WithFields(log.Fields{
			"package": "main",
			"func":    "main",
			"file":    logFile,
		}).Fatal("Don'n manage open log file for appending")
	}

	log.SetOutput(f)

	fConf := "database.config.json"

	configContent, err := ioutil.ReadFile(fConf)

	if err != nil {
		log.Fatalf("opening conf file: %s\n", err.Error())
	}

	config := scheme.GeneratorConfig{}
	err = json.Unmarshal(configContent, &config)

	if err != nil {
		log.Fatalf("unmarshal conf file: %s\n", err.Error())
	}

	s := generator.NewServer(config)

	s.Run()
}
