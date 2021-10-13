package main

import (
	"flag"
	"fmt"
	"os"
)

type filePath = string
type appEnviromentGetter func() (filePath, error)
type getterFunctionOrder int

type appEnviroment struct {
	getters           map[getterFunctionOrder]appEnviromentGetter
	pathToBillingFile filePath
}

const (
	FilePathFlagUsage       = "tell application that json file passed as flag value"
	FilePathFlagName        = "file-path"
	FilePathEnvVariableName = "FILE_PATH"
)

const (
	setFlag getterFunctionOrder = iota
	setEnvVariable
	setStdinVariable
)

var appEnvSingleInstance *appEnviroment

func (e *appEnviroment) init() {
	e.getters = map[getterFunctionOrder]appEnviromentGetter{
		setFlag:          e.setFlags,
		setEnvVariable:   e.setEnvVariables,
		setStdinVariable: e.setStdinVariables}
}

func (e *appEnviroment) setupEnviroment() error {
	e.init()
	var err error
	var path string
	for setterOrder := 0; setterOrder < len(e.getters); setterOrder++ {
		setter := e.getters[getterFunctionOrder(setterOrder)]
		if path, err = setter(); err == nil {
			e.pathToBillingFile = path
			return nil
		}
	}
	return fmt.Errorf("setup eviroment error happened: %w", err)
}

func (e *appEnviroment) setFlags() (filePath, error) {
	fileFlag := flag.String(FilePathFlagName, "", FilePathFlagUsage)
	flag.Parse()
	if *fileFlag == "" {
		return "", fmt.Errorf("file-path flag value is empty")
	}
	return *fileFlag, nil
}

func (e *appEnviroment) setEnvVariables() (filePath, error) {
	var value string
	var ok bool
	if value, ok = os.LookupEnv(FilePathEnvVariableName); !ok {
		return "", fmt.Errorf("enviroment variable value is empty")
	}
	return value, nil
}

func (e *appEnviroment) setStdinVariables() (filePath, error) {
	if flag.NArg() != 1 {
		return "", fmt.Errorf("invalid count of argument in STDIN")
	}
	return flag.Arg(0), nil
}

func (e appEnviroment) getJSONFilePath() (string, error) {
	if e.pathToBillingFile == "" {
		return "", fmt.Errorf("json file path is empty")
	}
	return e.pathToBillingFile, nil
}

func getAppEnvInstance() *appEnviroment {
	if appEnvSingleInstance == nil {
		appEnvSingleInstance = &appEnviroment{}
	}
	return appEnvSingleInstance
}
