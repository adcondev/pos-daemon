package encoding

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// CharacterSet representa un conjunto de caracteres con su codificación
type CharacterSet struct {
	Code     int               // Código numérico del charset (ej: 0, 2, 3)
	Name     string            // Nombre descriptivo (ej: "CP437", "CP850")
	Desc     string            // Descripción del charset (opcional)
	Encoding encoding.Encoding // Codificación real de golang.org/x/text
}

// Registry contiene todos los character sets disponibles.
// Numeración "típica" (pero no garantizada universalmente)
var Registry = map[int]*CharacterSet{
	0: {
		Code:     0,
		Name:     "CP437",
		Desc:     "Inglés/EE. UU. y símbolos gráficos DOS",
		Encoding: charmap.CodePage437,
	},
	2: {
		Code:     2,
		Name:     "CP850",
		Desc:     "Europa Occidental (multilingüe)",
		Encoding: charmap.CodePage850,
	},
	3: {
		Code:     3,
		Name:     "CP860",
		Desc:     "Portugués (Portugal)",
		Encoding: charmap.CodePage860,
	},
	4: {
		Code:     4,
		Name:     "CP863",
		Desc:     "Francés canadiense",
		Encoding: charmap.CodePage863,
	},
	5: {
		Code:     5,
		Name:     "CP865",
		Desc:     "Nórdico (escandinavo)",
		Encoding: charmap.CodePage865,
	},
	6: {
		Code:     6,
		Name:     "CP850",
		Desc:     "Europa Central y del Este",
		Encoding: charmap.CodePage850,
	},
	16: {
		Code:     16,
		Name:     "WPC1252",
		Desc:     "Windows Europa Occidental",
		Encoding: charmap.Windows1252,
	},
	17: {
		Code:     17,
		Name:     "CP866",
		Desc:     "Cirílico (Ruso MS-DOS)",
		Encoding: charmap.CodePage866,
	},
	18: {
		Code:     18,
		Name:     "CP852",
		Desc:     "Europa Central (Latin-2)",
		Encoding: charmap.CodePage852,
	},
	19: {
		Code:     19,
		Name:     "CP858",
		Encoding: charmap.CodePage858,
	},
	// Agregar más según necesites
	// IMPORTANTE: No existe un estándar universal obligatorio para la numeración de
	// tablas de codificación (code pages) en impresoras térmicas.
}

// GetEncoder devuelve un encoder para el charset especificado
func GetEncoder(charsetCode int) *encoding.Encoder {
	if cs, ok := Registry[charsetCode]; ok {
		return cs.Encoding.NewEncoder()
	}
	// Default a CP437 si no se encuentra
	return charmap.CodePage437.NewEncoder()
}

// EncodeString codifica un string usando el charset especificado
func EncodeString(str string, charsetCode int) ([]byte, error) {
	encoder := GetEncoder(charsetCode)
	return encoder.Bytes([]byte(str))
}
