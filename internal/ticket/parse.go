package ticket

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func JSONFileToBytes(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

func BytesToObj(b []byte) (*Ticket, error) {
	var w Wrapper

	if err := json.Unmarshal(b, &w); err != nil {
		return nil, err
	}
	return &w.Data, nil
}

func (t *Ticket) ToBytes() ([]byte, error) {
	return json.Marshal(Wrapper{Data: *t})
}
