package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	fp "path/filepath"
	"strings"
)

// JSONFileToBytes lee un archivo JSON y devuelve su contenido como bytes
func JSONFileToBytes(filepath string) ([]byte, error) {
	// Validar que el archivo esté dentro de un directorio permitido
	allowedDir := "./internal/api/rest" // directorio permitido
	absPath, err := fp.Abs(filepath)
	if err != nil {
		return nil, err
	}
	absAllowedDir, err := fp.Abs(allowedDir)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(absPath, absAllowedDir) {
		return nil, fmt.Errorf("acceso denegado al archivo fuera del directorio permitido")
	}

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

// BytesToConfig convierte bytes a una estructura LocalConfigData
func BytesToConfig(b []byte) (*ConfigData, error) {
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config.Data, nil
}

// ToBytes convierte LocalConfigData a bytes
func (t *ConfigData) ToBytes() ([]byte, error) {
	return json.Marshal(Config{Data: *t})
}

// BytesToTicket convierte bytes a una estructura TicketData
func BytesToTicket(b []byte) (*TicketData, error) {
	var ticket Ticket
	if err := json.Unmarshal(b, &ticket); err != nil {
		return nil, err
	}
	return &ticket.Data, nil
}

// ToBytes convierte TicketData a bytes
func (t *TicketData) ToBytes() ([]byte, error) {
	return json.Marshal(Ticket{Data: *t})
}

// BytesToTicketTemplate convierte bytes a una estructura TicketTemplateData
func BytesToTicketTemplate(b []byte) (*TicketTemplateData, error) {
	var template TicketTemplate
	if err := json.Unmarshal(b, &template); err != nil {
		return nil, err
	}
	return &template.Data, nil
}

// ToBytes convierte TicketTemplateData a bytes
func (t *TicketTemplateData) ToBytes() ([]byte, error) {
	return json.Marshal(TicketTemplate{Data: *t})
}

// BytesToNewTicket convierte bytes a una estructura NewTicketData
func BytesToNewTicket(b []byte) (*NewTicketData, error) {
	var ticket NewTicket
	if err := json.Unmarshal(b, &ticket); err != nil {
		return nil, err
	}
	return &ticket.Data, nil
}

// ToBytes convierte NewTicketData a bytes
func (t *NewTicketData) ToBytes() ([]byte, error) {
	return json.Marshal(NewTicket{Data: *t})
}
