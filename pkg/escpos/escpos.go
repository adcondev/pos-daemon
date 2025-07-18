package escpos

import (
	"errors"
	"fmt"

	cmd "pos-daemon.adcon.dev/pkg/escpos/command"
)

// NewPrinter crea una nueva instancia de ESCPrinter.
// Requiere un Connector y opcionalmente un CapabilityProfile.
// Si el perfil es nil, carga el perfil por defecto.
func NewPrinter(connector cmd.Connector, profile *cmd.CapabilityProfile) (*cmd.ESCPrinter, error) {
	if connector == nil {
		return nil, errors.New("connector no puede ser nil")
	}
	if profile == nil {
		// Cargar perfil por defecto si no se proporciona ninguno
		defaultProfile, err := cmd.LoadProfile("default")
		if err != nil {
			return nil, fmt.Errorf("falló al cargar el perfil de capacidad por defecto: %w", err)
		}
		profile = defaultProfile
	}

	p := &cmd.ESCPrinter{
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
