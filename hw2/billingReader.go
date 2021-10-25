package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type billings []billing

type jsonReader struct {
	data []byte
}

var (
	ErrOpenFile     = errors.New("unable to open file")
	ErrReadFile     = errors.New("unable to read file")
	ErrCloseFile    = errors.New("close file")
	ErrReadBillings = errors.New("read billings file error happened")
)

func (r *jsonReader) readFile(pathToFile string) error {
	f, err := os.Open(pathToFile)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrOpenFile, err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrReadFile, err)
	}
	r.data = data
	if err := f.Close(); err != nil {
		return fmt.Errorf("%s: %w", ErrCloseFile, err)
	}
	return nil
}

func (r jsonReader) readBillings(billingsFilePath string) (billings, error) {
	err := r.readFile(billingsFilePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrReadBillings, err)
	}
	var entries []billingRaw
	if err := json.Unmarshal(r.data, &entries); err != nil {
		return billings{}, fmt.Errorf("%s: %w", ErrReadBillings, err)
	}

	var billings billings
	for _, entry := range entries {
		unwrappedBillingRaw := entry.toUnwrappedBillingRaw()
		validatedBilling, err := unwrappedBillingRaw.validate()
		if err != nil {
			log.Warn(fmt.Errorf("skip billing: %w", err))
			continue
		}
		billings = append(billings, validatedBilling)
	}

	return billings, nil
}
