package escpos

import (
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

// SelectCharacterTable selecciona la tabla de caracteres (codepage) a utilizar.
func (p *Printer) SelectCharacterTable(table int) error {
	if err := validateInteger(table, 0, 255, "SelectCharacterTable", "tabla"); err != nil {
		return fmt.Errorf("SelectCharacterTable: %w", err)
	}

	// Verificar si el perfil soporta esta tabla
	if _, ok := p.profile.CodePages[table]; !ok {
		return fmt.Errorf("SelectCharacterTable: la tabla de caracteres %d no está permitida por el perfil de esta impresora", table)
	}

	// La clase PHP elige entre ESC t n (estándar) y ESC GS t n (posiblemente Star)
	// basándose en SupportsStarCommands. Implementamos esta lógica.
	var cmd []byte
	if p.profile.SupportsStarCommands {
		// Este comando es probable que sea específico de Star Micronics.
		cmd = []byte{ESC, GS, 't', byte(table)}
	} else {
		// Comando ESC/POS estándar para seleccionar tabla de caracteres.
		cmd = []byte{ESC, 't', byte(table)}
	}

	_, err := p.connector.Write(cmd)
	if err == nil {
		p.characterTable = table // Actualizar el estado interno si la escritura fue exitosa
	}
	return err
}

// GetCharacterTable devuelve la tabla de caracteres (codepage) actualmente seleccionada.
func (p *Printer) GetCharacterTable() int {
	return p.characterTable
}

// *** FUNCIÓN PARA CODIFICAR A CP858 ***
func toCP858(s string) []byte {
	// Obtener el codificador para CP858
	encoder := charmap.CodePage858.NewEncoder()
	// Convertir la string (UTF-8) a bytes codificados en CP858
	encoded, err := encoder.Bytes([]byte(s))
	if err != nil {
		// En caso de error (ej. carácter no representable en CP858),
		// podrías loguear el error, o intentar un fallback.
		// Aquí, por simplicidad, devolvemos la string original (UTF-8),
		// aunque esto no solucionaría el problema del acento si falla la codificación.
		// Una mejor práctica sería reemplazar el carácter desconocido.
		fmt.Printf("Advertencia: No se pudo codificar string a CP858: %v (original: %q)\n", err, s)
		return []byte(s) // Fallback (probablemente no imprimirá bien el carácter problemático)
	}
	return encoded
}
