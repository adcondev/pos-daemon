package ticket

import (
	"encoding/json"
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

func BytesToTicket(b []byte) (*TicketData, error) {
	var w Ticket

	if err := json.Unmarshal(b, &w); err != nil {
		return nil, err
	}
	return &w.Data, nil
}

func (t *TicketData) ToBytes() ([]byte, error) {
	return json.Marshal(Ticket{Data: *t})
}
