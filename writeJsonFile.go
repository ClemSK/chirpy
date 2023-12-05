package main

import (
	"encoding/json"
	"os"
)

func writeJsonFile(filepath string, data map[string]interface{}) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(data)
}
