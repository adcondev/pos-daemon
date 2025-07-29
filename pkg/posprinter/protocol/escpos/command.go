package escpos

import (
	"errors"
	"fmt"
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

// ESCPrinter representa una impresora térmica ESC/POS.
type ESCPrinter struct {
	Connector      Connector
	Profile        *CapabilityProfile
	CharacterTable int // La tabla de caracteres (codepage) actualmente seleccionada.
}

// Printer defines the interface for any ESC/POS compatible printer
type Printer interface {
	Initialize() error
	Pulse(int, int, int) error
	Close() error
	Status() (Status, error)
}

// Status represents the printer status
type Status struct {
	Online      bool
	PaperStatus PaperStatus
	DrawerOpen  bool
	// Additional status fields
}

type PaperStatus int

const (
	PaperOK PaperStatus = iota
	PaperLow
	PaperOut
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

// --- Métodos Públicos (Espejo de la clase PHP) ---

// Initialize restablece la impresora a su configuración por defecto (ESC @).
func (p *ESCPrinter) Initialize() error {
	_, err := p.Connector.Write([]byte{ESC, '@'})
	if err == nil {
		p.CharacterTable = 0 // Resetear el seguimiento de la tabla de caracteres
	}
	return err
}

// Close finaliza la conexión con la impresora.
func (p *ESCPrinter) Close() error {
	return p.Connector.Close()
}

// GetPrintConnector devuelve el conector que está utilizando la impresora.
func (p *ESCPrinter) GetPrintConnector() Connector {
	return p.Connector
}

// GetPrinterCapabilityProfile devuelve el perfil de capacidad de la impresora.
func (p *ESCPrinter) GetPrinterCapabilityProfile() *CapabilityProfile {
	return p.Profile
}

// NewPrinter crea una nueva instancia de ESCPrinter.
// Requiere un Connector y opcionalmente un CapabilityProfile.
// Si el perfil es nil, carga el perfil por defecto.
func NewPrinter(connector Connector, profile *CapabilityProfile) (*ESCPrinter, error) {
	if connector == nil {
		return nil, errors.New("connector no puede ser nil")
	}
	if profile == nil {
		// Cargar perfil por defecto si no se proporciona ninguno
		defaultProfile, err := LoadProfile("default")
		if err != nil {
			return nil, fmt.Errorf("falló al cargar el perfil de capacidad por defecto: %w", err)
		}
		profile = defaultProfile
	}

	p := &ESCPrinter{
		Connector:      connector,
		Profile:        profile,
		CharacterTable: 0, // Tabla de caracteres por defecto
	}

	// Inicializar la impresora
	if err := p.Initialize(); err != nil {
		// Si la inicialización falla, consideramos que la creación de la impresora falla.
		return nil, fmt.Errorf("falló al inicializar la impresora: %w", err)
	}

	return p, nil
}
