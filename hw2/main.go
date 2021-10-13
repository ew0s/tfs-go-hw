package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

var (
	ErrGetPath           = errors.New("get path")
	ErrSetEnviroment     = errors.New("set enviroment")
	ErrGetBillingEntries = errors.New("get billing entries")
)

func setEnviroment() (*appEnviroment, error) {
	enviroment := getAppEnvInstance()
	err := enviroment.setupEnviroment()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSetEnviroment, err)
	}
	return enviroment, nil
}

func getPath(enviroment *appEnviroment) (string, error) {
	var path string
	var err error
	if path, err = enviroment.getJSONFilePath(); err != nil {
		return "", fmt.Errorf("%s: %w", ErrGetPath, err)
	}
	return path, nil
}

func getBillingEntries(path string) (billings, error) {
	reader := jsonReader{}
	entries, err := reader.readBillings(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrGetBillingEntries, err)
	}
	return entries, nil
}

func main() {
	enviroment, err := setEnviroment()
	if err != nil {
		log.Error(err)
		return
	}
	path, err := getPath(enviroment)
	if err != nil {
		log.Error(err)
		return
	}
	billingEntries, err := getBillingEntries(path)
	if err != nil {
		log.Println(err)
	}
	statistic := calculateCompaniesStatistic(billingEntries)
	file, err := json.MarshalIndent(statistic, "", "\t")
	if err != nil {
		log.Println("marshalJson: ", err)
	}
	err = ioutil.WriteFile("out.json", file, 0644)
	if err != nil {
		log.Println("writeFile: ", err)
	}
}
