package escpos_test

import (
	"fmt"
	"testing"

	"pos-daemon.adcon.dev/pkg/escpos"
)

// MockConnector implementa PrintConnector para pruebas
type MockConnector struct {
	data    []byte
	closed  bool
	writeErr error
}

func NewMockConnector() *MockConnector {
	return &MockConnector{
		data: make([]byte, 0),
	}
}

func (m *MockConnector) Write(data []byte) (int, error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	m.data = append(m.data, data...)
	return len(data), nil
}

func (m *MockConnector) Close() error {
	m.closed = true
	return nil
}

func (m *MockConnector) GetData() []byte {
	return m.data
}

func ExampleNewPrinter() {
	// Crear un conector mock para el ejemplo
	connector := NewMockConnector()
	
	// Crear una nueva impresora con perfil por defecto
	_, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Verificar que la impresora se inicializó
	data := connector.GetData()
	if len(data) > 0 && data[0] == 0x1b && data[1] == '@' {
		fmt.Println("Impresora inicializada correctamente")
	}
	
	// Output: Impresora inicializada correctamente
}

func ExamplePrinter_Text() {
	connector := NewMockConnector()
	printer, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Limpiar datos de inicialización
	connector.data = make([]byte, 0)
	
	// Enviar texto
	err = printer.Text("Hola Mundo\n")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Verificar que se envió texto (contiene caracteres esperados)
	data := connector.GetData()
	if len(data) > 0 {
		fmt.Println("Texto enviado a la impresora")
	}
	
	// Output: Texto enviado a la impresora
}

func ExamplePrinter_SetJustification() {
	connector := NewMockConnector()
	printer, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Limpiar datos de inicialización
	connector.data = make([]byte, 0)
	
	// Configurar justificación centrada
	err = printer.SetJustification(escpos.JUSTIFY_CENTER)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Verificar comando ESC a 1 (centrado)
	data := connector.GetData()
	if len(data) >= 3 && data[0] == 0x1b && data[1] == 'a' && data[2] == 1 {
		fmt.Println("Justificación centrada configurada")
	}
	
	// Output: Justificación centrada configurada
}

func ExamplePrinter_Cut() {
	connector := NewMockConnector()
	printer, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Limpiar datos de inicialización
	connector.data = make([]byte, 0)
	
	// Corte completo sin líneas de avance
	err = printer.Cut(escpos.CUT_FULL, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Verificar comando GS V A 0 (corte completo)
	data := connector.GetData()
	if len(data) >= 4 && data[0] == 0x1d && data[1] == 'V' && data[2] == 65 && data[3] == 0 {
		fmt.Println("Comando de corte enviado")
	}
	
	// Output: Comando de corte enviado
}

func TestPrinterBasicFunctionality(t *testing.T) {
	connector := NewMockConnector()
	
	// Probar creación de impresora
	printer, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		t.Fatalf("Error creando impresora: %v", err)
	}
	
	// Verificar que se ejecutó la inicialización
	data := connector.GetData()
	if len(data) < 2 || data[0] != 0x1b || data[1] != '@' {
		t.Error("Comando de inicialización no encontrado")
	}
	
	// Limpiar datos para siguientes pruebas
	connector.data = make([]byte, 0)
	
	// Probar envío de texto
	err = printer.Text("Test")
	if err != nil {
		t.Errorf("Error enviando texto: %v", err)
	}
	
	if len(connector.GetData()) == 0 {
		t.Error("No se enviaron datos para el texto")
	}
}

func TestPrinterValidation(t *testing.T) {
	connector := NewMockConnector()
	printer, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		t.Fatalf("Error creando impresora: %v", err)
	}
	
	// Probar validación de parámetros inválidos
	err = printer.Cut(99, 0) // Modo inválido
	if err == nil {
		t.Error("Se esperaba error con modo de corte inválido")
	}
	
	err = printer.Feed(0) // Líneas inválidas
	if err == nil {
		t.Error("Se esperaba error con número de líneas inválido")
	}
	
	err = printer.SetJustification(99) // Justificación inválida
	if err == nil {
		t.Error("Se esperaba error con justificación inválida")
	}
}

func TestPrinterWithNilConnector(t *testing.T) {
	// Probar que no se puede crear impresora con conector nil
	_, err := escpos.NewPrinter(nil, nil)
	if err == nil {
		t.Error("Se esperaba error con conector nil")
	}
}