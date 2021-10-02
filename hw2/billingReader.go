package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type billings []billing

type jsonReader struct {
	data []byte
}

func (r *jsonReader) readFile(pathToFile string) error {
	f, err := os.Open(pathToFile)
	if err != nil {
		return fmt.Errorf("unable to open file")
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("unbale to read file")
	}
	r.data = data
	err = f.Close()
	if err != nil {
		return fmt.Errorf("readFile: %w", err)
	}
	return nil
}

func (r jsonReader) readBillings(billingsFilePath string) (billings, error) {
	err := r.readFile(billingsFilePath)
	if err != nil {
		return nil, fmt.Errorf("read billings file error hapanned: %w", err)
	}
	var entries billings
	err = json.Unmarshal(r.data, &entries)
	if err != nil {
		log.Println("unmarshal billings error happened:", err)
	}
	var notSkippedEntries billings
	for _, entry := range entries {
		if err := entry.validate(); err == nil {
			notSkippedEntries = append(notSkippedEntries, entry)
		}
	}
	return notSkippedEntries, nil
}
