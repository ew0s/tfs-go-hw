package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func setEnviroment() (*appEnviroment, error) {
	enviroment := getAppEnvInstance()
	err := enviroment.setupEnviroment()
	if err != nil {
		return nil, fmt.Errorf("setEnviroment: %w", err)
	}
	return enviroment, nil
}

func getPath(enviroment *appEnviroment) (string, error) {
	var path string
	var err error
	if path, err = enviroment.getJSONFilePath(); err != nil {
		return "", fmt.Errorf("getPath: %w", err)
	}
	return path, nil
}

func getBillingEntries(path string) (billings, error) {
	reader := jsonReader{}
	entries, err := reader.readBillings(path)
	if err != nil {
		return nil, fmt.Errorf("getBillingEntries: %w", err)
	}
	return entries, nil
}

func main() {
	enviroment, err := setEnviroment()
	if err != nil {
		log.Println(err)
	}
	path, err := getPath(enviroment)
	if err != nil {
		log.Println(err)
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
