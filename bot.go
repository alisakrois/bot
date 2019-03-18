package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
)

func getConfig() (*Config, error) {
	var config Config
	yamlFile, err := os.Open("config.yaml")

	if err != nil {
		return nil, errors.New("Сan't open .yaml file for reading")
	}

	decoder := yaml.NewDecoder(yamlFile)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, errors.New("Сan't decode .yaml file")
	}

	return &config, nil
}

func checkErrors(config *Config) error {
	configValues := reflect.ValueOf(config).Elem()
	сonfigType := configValues.Type()
	for num := 0; num < configValues.NumField(); num++ {
		configField := configValues.Field(num)
		if configField.Interface() == "" || configField.Interface() == 0 {
			return errors.New("Wrong field in config: " + сonfigType.Field(num).Name)
		}
	}

	_, err := renderTemplate(config.EmailTemplate)
	if err != nil {
		return errors.New("Can't execute the template" + err.Error())
	}
	return nil
}

func main() {
	conf, err := getConfig()

	if err != nil {
		log.Fatal("Can't get config for email: ", err)
	}

	err = checkErrors(conf)
	if err != nil {
		log.Fatal("Invalid config, e: ", err.Error())
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}

	mailService := NewMailService(client, conf.Token, conf.FromEmail, conf.URL)
	botController := BotController{mailService, conf}

	r := mux.NewRouter()
	r.HandleFunc("/slack/events", botController.proccessSlack).Methods("POST")

	port := strconv.Itoa(conf.Port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Print("Failed to listen: ", err.Error())
	}
}
