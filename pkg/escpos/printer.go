package escpos

import (
	"errors"
	"fmt"
	"os"
)

const (
	// Tipos de estado (para comandos de estado, no implementados como métodos públicos en PHP)
	// Se incluyen por completitud de las constantes PHP
	STATUS_PRINTER       int = 1 // GS I 1 (Estado de la impresora)
	STATUS_OFFLINE_CAUSE int = 2 // GS I 2 (Causa de estar offline)
	STATUS_ERROR_CAUSE   int = 3 // GS I 3 (Causa del error)
	STATUS_PAPER_ROLL    int = 4 // GS I 4 (Estado del rollo de papel)
	STATUS_INK_A         int = 7 // GS I 7 (Estado de la tinta/cinta A)
	STATUS_INK_B         int = 6 // GS I 6 (Estado de la tinta/cinta B)
	STATUS_PEELER        int = 8 // GS I 8 (Estado del peeler - para etiquetas)

)

// Connector define la interfaz para la conexión física con la impresora.
// Debes implementar esta interfaz para tu método de conexión (USB, TCP, Serial, etc.).
type Connector interface {
	// Write envía bytes a la impresora.
	Write([]byte) (int, error)
	// Close finaliza la conexión con la impresora.
	Close() error
}

// CapabilityProfile describe las capacidades de una impresora ESC/POS específica.
type CapabilityProfile struct {
	SupportsBarcodeB     bool // Soporta el formato de comando GS k m L data (65-73)
	SupportsPdf417Code   bool
	SupportsQrCode       bool
	CodePages            map[int]string // Mapa de índice de codepage a nombre/descripción
	SupportsStarCommands bool           // Indica si soporta comandos específicos de Star (como ESC GS t)
	// Agrega otras capacidades según sea necesario (ej: anchos de papel, fuentes, etc.)
}

// LoadProfile TODO carga un CapabilityProfile predefinido o desde una fuente externa.
// Esta es una función placeholder.
func LoadProfile(name string) (*CapabilityProfile, error) {
	// En una biblioteca real, esto cargaría perfiles desde archivos o datos incrustados.
	// Para este port, devolvemos un perfil dummy "default".
	if name == "default" {
		return &CapabilityProfile{
			SupportsBarcodeB:   true, // Asumimos que el perfil por defecto soporta el formato moderno de código de barras
			SupportsPdf417Code: true, // Asumimos que el perfil por defecto soporta códigos 2D
			SupportsQrCode:     true,
			CodePages: map[int]string{ // Ejemplo de algunas codepages comunes
				0: "CP437",
				1: "CP850",
				2: "CP852",
				3: "CP858",
				4: "CP860",
				5: "CP863",
				6: "CP865",
				// Agrega más codepages soportadas
				16:  "WPC1252",   // Windows 1252
				254: "Shift_JIS", // Ejemplo asiático
				255: "GBK",       // Ejemplo chino (usado en textChinese)
			},
			SupportsStarCommands: false, // Asumimos que el perfil por defecto no soporta comandos Star
		}, nil
	}
	return nil, fmt.Errorf("perfil de capacidad desconocido: %s", name)
}

// Printer representa una impresora térmica ESC/POS.
type Printer struct {
	Connector      Connector
	profile        *CapabilityProfile
	characterTable int // La tabla de caracteres (codepage) actualmente seleccionada.
}

// NewPrinter crea una nueva instancia de Printer.
// Requiere un Connector y opcionalmente un CapabilityProfile.
// Si el perfil es nil, carga el perfil por defecto.
func NewPrinter(connector Connector, profile *CapabilityProfile) (*Printer, error) {
	if connector == nil {
		return nil, errors.New("Connector no puede ser nil")
	}
	if profile == nil {
		// Cargar perfil por defecto si no se proporciona ninguno
		defaultProfile, err := LoadProfile("default")
		if err != nil {
			return nil, fmt.Errorf("falló al cargar el perfil de capacidad por defecto: %w", err)
		}
		profile = defaultProfile
	}

	p := &Printer{
		Connector:      connector,
		profile:        profile,
		characterTable: 0, // Tabla de caracteres por defecto
	}

	// Inicializar la impresora
	if err := p.Initialize(); err != nil {
		// Si la inicialización falla, consideramos que la creación de la impresora falla.
		return nil, fmt.Errorf("falló al inicializar la impresora: %w", err)
	}

	return p, nil
}

// --- Métodos Públicos (Espejo de la clase PHP) ---

// Initialize restablece la impresora a su configuración por defecto (ESC @).
func (p *Printer) Initialize() error {
	_, err := p.Connector.Write([]byte{ESC, '@'})
	if err == nil {
		p.characterTable = 0 // Resetear el seguimiento de la tabla de caracteres
	}
	return err
}

// Pulse envía un pulso a un pin del conector del cajón portamonedas para abrirlo.
func (p *Printer) Pulse(pin int, onMS, offMS int) error {
	if err := validateInteger(pin, 0, 1, "Pulse", "pin"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Pin 0 o 1
	if err := validateInteger(onMS, 1, 511, "Pulse", "onMS"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Tiempo ON en ms (1-511)
	if err := validateInteger(offMS, 1, 511, "Pulse", "offMS"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Tiempo OFF en ms (1-511) - a menudo ignorado por la impresora para el segundo pulso

	// Comando: ESC p m t1 t2
	// m: pin del cajón (0 o 1). PHP usa pin + 48 ('0' o '1'). Replicamos.
	// t1: Tiempo ON (t1 * 2 ms). PHP envía on_ms / 2. Replicamos.
	// t2: Tiempo OFF (t2 * 2 ms). PHP envía off_ms / 2. Replicamos.
	cmd := []byte{ESC, 'p', byte(pin + 48), byte(onMS / 2), byte(offMS / 2)}
	_, err := p.Connector.Write(cmd)
	return err
}

// Close finaliza la conexión con la impresora.
func (p *Printer) Close() error {
	return p.Connector.Close()
}

// GetPrintConnector devuelve el conector que está utilizando la impresora.
func (p *Printer) GetPrintConnector() Connector {
	return p.Connector
}

// GetPrinterCapabilityProfile devuelve el perfil de capacidad de la impresora.
func (p *Printer) GetPrinterCapabilityProfile() *CapabilityProfile {
	return p.profile
}

// intPtr es una función de ayuda para obtener un puntero a un int.
// Útil para métodos con parámetros opcionales *int (como SetLineSpacing).
func intPtr(i int) *int {
	return &i
}

// Función de ayuda para abrir archivos
func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo %s: %w", filename, err)
	}
	return file, nil
}
