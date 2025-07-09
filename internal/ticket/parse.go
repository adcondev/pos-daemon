package ticket

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// JSONFileToBytes lee un archivo JSON desde el sistema de archivos y 
// retorna su contenido como slice de bytes.
//
// Parámetros:
//   - filepath: ruta al archivo JSON a leer
//
// Retorna los bytes del archivo o un error si no se puede leer.
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

// BytesToJSONFile escribe bytes a un archivo JSON validando primero
// que el contenido sea JSON válido.
//
// Parámetros:
//   - data: bytes a escribir (debe ser JSON válido)
//   - filename: nombre del archivo destino
//
// Retorna un error si el JSON es inválido o no se puede escribir el archivo.
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

// BytesToObj convierte bytes JSON a una estructura Ticket.
// Espera que el JSON tenga la estructura de Wrapper con el ticket
// en el campo "data".
//
// Parámetros:
//   - b: bytes JSON a convertir
//
// Retorna un puntero a Ticket o un error si el JSON es inválido.
func BytesToObj(b []byte) (*Ticket, error) {
	var w Wrapper

	if err := json.Unmarshal(b, &w); err != nil {
		return nil, err
	}
	return &w.Data, nil
}

// ToBytes convierte el Ticket a bytes JSON envolviéndolo en un Wrapper.
// Esto asegura que el JSON tenga la estructura esperada por el sistema.
//
// Retorna los bytes JSON o un error si no se puede serializar.
func (t *Ticket) ToBytes() ([]byte, error) {
	return json.Marshal(Wrapper{Data: *t})
}
