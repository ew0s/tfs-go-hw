package main

import (
	"errors"
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

var (
	ErrSetupEnviroment            = errors.New("setup enviroment error")
	ErrFilePathFlagEmpty          = errors.New("file-path flag value is empty")
	ErrEnviromentVarEmpty         = errors.New("enviroment variable value is empty")
	ErrStdinInvalidArgumentsCount = errors.New("invalid count of argument in STDIN")
	ErrJSONFilePathIsEmpty        = errors.New("json file path is empty")
)

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
	return fmt.Errorf("%s: %w", ErrSetupEnviroment, err)
}

func (e *appEnviroment) setFlags() (filePath, error) {
	fileFlag := flag.String(FilePathFlagName, "", FilePathFlagUsage)
	flag.Parse()
	if *fileFlag == "" {
		return "", ErrFilePathFlagEmpty
	}
	return *fileFlag, nil
}

func (e *appEnviroment) setEnvVariables() (filePath, error) {
	var value string
	var ok bool
	if value, ok = os.LookupEnv(FilePathEnvVariableName); !ok {
		return "", ErrEnviromentVarEmpty
	}
	return value, nil
}

func (e *appEnviroment) setStdinVariables() (filePath, error) {
	if flag.NArg() != 1 {
		return "", ErrStdinInvalidArgumentsCount
	}
	return flag.Arg(0), nil
}

func (e appEnviroment) getJSONFilePath() (string, error) {
	if e.pathToBillingFile == "" {
		return "", ErrJSONFilePathIsEmpty
	}
	return e.pathToBillingFile, nil
}

func getAppEnvInstance() *appEnviroment {
	if appEnvSingleInstance == nil {
		appEnvSingleInstance = &appEnviroment{}
	}
	return appEnvSingleInstance
}
