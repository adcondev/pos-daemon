package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func JSONFileToBytes(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error al cerrar JSON: %v", err)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func BytesToJSONFile(data []byte, filename string) error {
	var jsonCheck interface{}
	if err := json.Unmarshal(data, &jsonCheck); err != nil {
		return fmt.Errorf("Invalid JSON data: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("Failed to write to file: %w", err)
	}

	return nil
}

func BytesToConfig(b []byte) (*Config, error) {
	var w Wrapper

	if err := json.Unmarshal(b, &w); err != nil {
		return nil, err
	}
	return &w.Data, nil
}

func (t *Config) ToBytes() ([]byte, error) {
	return json.Marshal(Wrapper{Data: *t})
}
