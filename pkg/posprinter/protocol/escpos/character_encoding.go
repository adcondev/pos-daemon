package escpos

import (
	"log"

	"pos-daemon.adcon.dev/pkg/posprinter/encoding"
	"pos-daemon.adcon.dev/pkg/posprinter/types"
)

// TODO: Comandos para manejo de codificación de caracteres
// - Código de página
// - Caracteres internacionales
// - Caracteres especiales

type CodePage byte

const (
	// Tabla de códigos comunes en ESC/POS
	CP437      CodePage = iota // CP437 U.S.A. / Standard Europe
	Katakana                   // Katakana (JIS X 0201)
	CP850                      // CP850 Multilingual
	CP860                      // CP860 Portuguese
	CP863                      // CP863 Canadian French
	CP865                      // CP865 Nordic
	WestEurope                 // WestEurope (ISO-8859-1)
	Greek                      // Greek (ISO-8859-7)
	Hebrew                     // Hebrew (ISO-8859-8)
	CP755                      // CP755 East Europe (not directly supported)
	Iran                       // Iran (CP720 Arabic)
)

const (
	WCP1252 CodePage = iota + 16 // WCP1252 Windows-1252
	CP866                        // CP866 Cyrillic #2
	CP852                        // CP852 Latin2
	CP858                        // CP858 Multilingual + Euro
	IranII                       // IranII (CP864)
	Latvian                      // Latvian (Windows-1257)
)

func (cp CodePage) IsValid() bool {
	return cp <= Latvian || (cp >= WCP1252 && cp <= Latvian)
}

func (p *Commands) SelectCharacterTable(table types.CharacterSet) []byte {
	charTable := CodePage(encoding.Registry[table].EscPos)
	// Validar que table esté en un rango válido
	if !charTable.IsValid() {
		// Log de advertencia si está fuera de rango
		log.Printf("advertencia: tabla de caracteres %d fuera de rango, usando 0 por defecto", table)
		charTable = 0 // Default a 0 si está fuera de rango
	}
	// ESC t n - Select character code table
	cmd := []byte{ESC, 't', byte(charTable)}

	return cmd
}

// CancelKanjiMode cancela el modo de caracteres Kanji.
//
// Formato:
//
//	ASCII: FS .
//	Hex:   1C 2E
//	Decimal: 28 46
//
// Descripción:
//
//	Deshabilita el modo de caracteres Kanji en la impresora.
//
// Referencia:
//
//	FS &, FS C
func (p *Commands) CancelKanjiMode() []byte {
	return []byte{FS, '.'}
}

func (p *Commands) SelectKanjiMode() []byte {
	return []byte{FS, '&'}
}
