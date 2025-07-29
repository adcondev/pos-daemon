package posprinter

// Capability define las capacidades de una impresora
type Capability struct {
	// Características de la impresora
	Model            string
	Vendor           string
	MaxColumns       int
	SupportsQRNative bool
	SupportsCutter   bool
	MaxDPI           int
	CharSets         []int
	// Otras características...
}

// DefaultCapability crea un perfil de capacidades predeterminado
func DefaultCapability() *Capability {
	return &Capability{
		Model:            "Generic",
		Vendor:           "Generic",
		MaxColumns:       48,
		SupportsQRNative: false,
		SupportsCutter:   true,
		MaxDPI:           203,
		CharSets:         []int{0}, // Código de página 0 (CP437 US) por defecto
	}
}

// LoadCapabilityFromJSON carga capacidades desde un archivo JSON
// Esto se puede implementar más adelante cuando sea necesario
func LoadCapabilityFromJSON(jsonPath string) (*Capability, error) {
	// Implementación básica para cargar desde JSON
	// Por ahora devuelve el perfil por defecto
	return DefaultCapability(), nil
}
