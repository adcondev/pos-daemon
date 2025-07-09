package ticket_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"pos-daemon.adcon.dev/internal/ticket"
)

func ExampleTicket_ToBytes() {
	// Crear un ticket de ejemplo
	t := &ticket.Ticket{
		Identificador: "123456",
		Vendedor:      "Juan Pérez",
		Folio:         "001",
		Serie:         "A",
		FechaSistema:  "10/01/2024 15:30:00",
		TipoOperacion: "NOTA_VENTA",
		Total:         100.50,
		Conceptos: []ticket.Concepto{
			{
				Clave:       "PROD001",
				Descripcion: "Producto de prueba",
				Cantidad:    1.0,
				PrecioVenta: 100.50,
				Total:       100.50,
			},
		},
	}

	// Convertir a bytes JSON
	bytes, err := t.ToBytes()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Validar que es JSON válido
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		fmt.Printf("Error validando JSON: %v\n", err)
		return
	}

	fmt.Println("Ticket serializado correctamente a JSON")
	// Output: Ticket serializado correctamente a JSON
}

func ExampleBytesToObj() {
	// JSON de ejemplo con estructura Wrapper
	jsonData := `{
		"data": {
			"identificador": "123456",
			"vendedor": "Juan Pérez",
			"folio": "001",
			"serie": "A",
			"fecha_sistema": "10/01/2024 15:30:00",
			"tipo_operacion": "NOTA_VENTA",
			"total": "100.50",
			"conceptos": [
				{
					"clave": "PROD001",
					"descripcion": "Producto de prueba",
					"cantidad": "1.0",
					"precio_venta": "100.50",
					"total": "100.50"
				}
			]
		}
	}`

	// Convertir a objeto Ticket
	t, err := ticket.BytesToObj([]byte(jsonData))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Ticket ID: %s, Vendedor: %s\n", t.Identificador, t.Vendedor)
	// Output: Ticket ID: 123456, Vendedor: Juan Pérez
}

func ExampleJSONFileToBytes() {
	// Crear archivo temporal de prueba
	tempFile, err := os.CreateTemp("", "ticket_test_*.json")
	if err != nil {
		fmt.Printf("Error creando archivo temporal: %v\n", err)
		return
	}
	defer os.Remove(tempFile.Name())

	// Escribir JSON de prueba
	testJSON := `{"data": {"identificador": "123", "vendedor": "Test"}}`
	if _, err := tempFile.WriteString(testJSON); err != nil {
		fmt.Printf("Error escribiendo archivo: %v\n", err)
		return
	}
	tempFile.Close()

	// Leer archivo con JSONFileToBytes
	bytes, err := ticket.JSONFileToBytes(tempFile.Name())
	if err != nil {
		fmt.Printf("Error leyendo archivo: %v\n", err)
		return
	}

	// Validar que es JSON válido
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		fmt.Printf("Error validando JSON: %v\n", err)
		return
	}

	fmt.Println("Archivo JSON leído correctamente")
	// Output: Archivo JSON leído correctamente
}

func TestTicketSerialization(t *testing.T) {
	// Crear ticket de prueba
	original := &ticket.Ticket{
		Identificador: "TEST123",
		Vendedor:      "Test User",
		Folio:         "001",
		Serie:         "A",
		FechaSistema:  time.Now().Format("02/01/2006 15:04:05"),
		TipoOperacion: "NOTA_VENTA",
		Total:         150.75,
		Descuento:     10.25,
		Conceptos: []ticket.Concepto{
			{
				Clave:       "PROD001",
				Descripcion: "Producto de prueba",
				Cantidad:    2.0,
				PrecioVenta: 80.50,
				Total:       161.00,
			},
		},
	}

	// Serializar a bytes
	bytes, err := original.ToBytes()
	if err != nil {
		t.Fatalf("Error serializando ticket: %v", err)
	}

	// Deserializar de vuelta
	restored, err := ticket.BytesToObj(bytes)
	if err != nil {
		t.Fatalf("Error deserializando ticket: %v", err)
	}

	// Verificar que los datos son los mismos
	if restored.Identificador != original.Identificador {
		t.Errorf("Identificador no coincide: got %s, want %s", restored.Identificador, original.Identificador)
	}

	if restored.Vendedor != original.Vendedor {
		t.Errorf("Vendedor no coincide: got %s, want %s", restored.Vendedor, original.Vendedor)
	}

	if restored.Total != original.Total {
		t.Errorf("Total no coincide: got %f, want %f", restored.Total, original.Total)
	}

	if len(restored.Conceptos) != len(original.Conceptos) {
		t.Errorf("Número de conceptos no coincide: got %d, want %d", len(restored.Conceptos), len(original.Conceptos))
	}
}

func TestBytesToJSONFile(t *testing.T) {
	// Crear archivo temporal
	tempFile, err := os.CreateTemp("", "ticket_output_*.json")
	if err != nil {
		t.Fatalf("Error creando archivo temporal: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// JSON válido
	validJSON := `{"data": {"identificador": "123", "vendedor": "Test"}}`
	err = ticket.BytesToJSONFile([]byte(validJSON), tempFile.Name())
	if err != nil {
		t.Errorf("Error con JSON válido: %v", err)
	}

	// JSON inválido
	invalidJSON := `{"data": {"identificador": "123", "vendedor": "Test"`
	err = ticket.BytesToJSONFile([]byte(invalidJSON), tempFile.Name())
	if err == nil {
		t.Error("Se esperaba error con JSON inválido, pero no ocurrió")
	}
}