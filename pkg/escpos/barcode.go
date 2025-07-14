package escpos

import (
	"errors"
	"fmt"
)

const (
	// Códigos de barras
	BARCODE_UPCA    int = 65 // 'A'
	BARCODE_UPCE    int = 66 // 'B'
	BARCODE_JAN13   int = 67 // 'C' (EAN13)
	BARCODE_JAN8    int = 68 // 'D' (EAN8)
	BARCODE_CODE39  int = 69 // 'E'
	BARCODE_ITF     int = 70 // 'F'
	BARCODE_CODABAR int = 71 // 'G'
	BARCODE_CODE93  int = 72 // 'H'
	BARCODE_CODE128 int = 73 // 'I'

	// Posición del texto del código de barras
	BARCODE_TEXT_NONE  int = 0
	BARCODE_TEXT_ABOVE int = 1
	BARCODE_TEXT_BELOW int = 2
)

// SetBarcodeHeight establece la altura del código de barras en puntos.
func (p *Printer) SetBarcodeHeight(height int) error {
	if err := validateInteger(height, 1, 255, "SetBarcodeHeight", "altura"); err != nil {
		return fmt.Errorf("SetBarcodeHeight: %w", err)
	}
	// GS h n - Establece la altura del código de barras a n puntos
	cmd := []byte{GS, 'h', byte(height)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetBarcodeWidth establece el ancho de los módulos del código de barras.
func (p *Printer) SetBarcodeWidth(width int) error {
	if err := validateInteger(width, 1, 255, "SetBarcodeWidth", "ancho"); err != nil {
		return fmt.Errorf("SetBarcodeWidth: %w", err)
	}
	// GS w n - Establece el ancho horizontal de los módulos a n (normalmente 2 o 3)
	cmd := []byte{GS, 'w', byte(width)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetBarcodeTextPosition establece la posición del texto legible del código de barras.
func (p *Printer) SetBarcodeTextPosition(position int) error {
	if err := validateInteger(position, BARCODE_TEXT_NONE, BARCODE_TEXT_BELOW, "SetBarcodeTextPosition", "posición"); err != nil {
		return fmt.Errorf("SetBarcodeTextPosition: %w", err)
	} // 0: ninguno, 1: arriba, 2: abajo, 3: ambos (no siempre soportado) - PHP valida 0-3
	cmd := []byte{GS, 'H', byte(position)}
	_, err := p.Connector.Write(cmd)
	return err
}

// Barcode imprime un código de barras.
// content es la cadena de datos del código de barras.
// barType es el tipo de código de barras (BARCODE_UPCA, BARCODE_CODE39, etc.).
func (p *Printer) Barcode(content string, barType int) error {
	if err := validateInteger(barType, BARCODE_UPCA, BARCODE_CODE128, "Barcode", "tipo de código de barras"); err != nil {
		return fmt.Errorf("Barcode: %w", err)
	}
	contentLen := len(content)

	// --- Validación de contenido basada en el tipo de código de barras (traducir regex y longitud) ---
	var validationErr error
	switch barType {
	case BARCODE_UPCA:
		validationErr = validateStringRegex(content, `^[0-9]{11,12}$`, "contenido UPCA")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 11, 12, "Barcode", "longitud contenido UPCA")
		}
	case BARCODE_UPCE:
		validationErr = validateStringRegex(content, `^([0-9]{6,8}|[0-9]{11,12})$`, "contenido UPCE")
		if validationErr == nil {
			validationErr = validateIntegerMulti(contentLen, [][]int{{6, 8}, {11, 12}}, "Barcode", "longitud contenido UPCE")
		}
	case BARCODE_JAN13:
		validationErr = validateStringRegex(content, `^[0-9]{12,13}$`, "contenido JAN13")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 12, 13, "Barcode", "longitud contenido JAN13")
		}
	case BARCODE_JAN8:
		validationErr = validateStringRegex(content, `^[0-9]{7,8}$`, "contenido JAN8")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 7, 8, "Barcode", "longitud contenido JAN8")
		}
	case BARCODE_CODE39:
		// PHP regex: `^([0-9A-Z $%+\-./]+|\*[0-9A-Z $%+\-./]+\*)$`
		// Requiere un * al principio y al final, o no.
		validationErr = validateStringRegex(content, `^([0-9A-Z $%+\-./]+|\*[0-9A-Z $%+\-./]+\*)$`, "contenido CODE39")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE39")
		}
	case BARCODE_ITF:
		// PHP regex: `^([0-9]{2})+$` - requiere solo dígitos y longitud par.
		validationErr = validateStringRegex(content, `^([0-9]{2})+$`, "contenido ITF")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 2, 255, "Barcode", "longitud contenido ITF")
		}
		// Validar longitud par
		if validationErr == nil && contentLen%2 != 0 {
			validationErr = errors.New("la longitud del contenido ITF debe ser par")
		}
	case BARCODE_CODABAR:
		// PHP regex: `^[A-Da-d][0-9$%+\-./:]+[A-Da-d]$` - inicia/termina con A-D, medio con dígitos/símbolos.
		validationErr = validateStringRegex(content, `^[A-Da-d][0-9$%+\-./:]+[A-Da-d]$`, "contenido Codabar")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 1, 255, "Barcode", "longitud contenido Codabar")
		}
	case BARCODE_CODE93:
		// PHP regex: `^[\x00-\x7F]+$` - solo caracteres ASCII.
		validationErr = validateStringRegex(content, `^[\x00-\x7F]+$`, "contenido CODE93")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE93")
		}
	case BARCODE_CODE128:
		// PHP regex: `^\{[A-C][\\x00-\\x7F]+$` - espera que el contenido empiece con {A, {B o {C y luego ASCII.
		// Esto es un poco inusual, ya que normalmente el usuario no proporciona los códigos de inicio/función de Code128.
		// Replicamos la validación de PHP.
		validationErr = validateStringRegex(content, `^\{[A-C][\x00-\x7F]+$`, "contenido CODE128")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE128")
		}
		if validationErr == nil && contentLen < 2 { // Necesita al menos '{' y un carácter de tipo
			validationErr = errors.New("el contenido Code128 debe tener al menos 2 caracteres ({A, {B, {C...)")
		}
	}
	if validationErr != nil {
		return fmt.Errorf("Barcode: contenido '%s' inválido para el tipo %d: %w", content, barType, validationErr)
	}
	// --- Fin Validación ---

	// Lógica de envío:
	// PHP usa el comando GS k m d1...dk (m=0-6) si getSupportsBarcodeB() es false.
	// PHP usa el comando GS k m L d1...dL (m=65-73) si getSupportsBarcodeB() es true.
	// 'SupportsBarcodeB' en PHP parece referirse al soporte del formato de comando más nuevo (con byte de longitud L).

	cmd := []byte{GS, 'k'}
	if !p.profile.SupportsBarcodeB {
		// Usar el formato de comando antiguo: GS k m data NUL (m = 0-6)
		// Validar que el tipo solicitado esté en el rango 65-71 (correspondiente a m 0-6)
		if barType < BARCODE_UPCA || barType > BARCODE_CODABAR {
			return fmt.Errorf("Barcode: el perfil de impresora no soporta el tipo de código de barras %d con el formato de comando antiguo", barType)
		}
		cmd = append(cmd, byte(barType-65)) // Tipo de 0 a 6
		cmd = append(cmd, []byte(content)...)
		cmd = append(cmd, NUL) // Terminador NUL
	} else {
		// Usar el formato de comando nuevo: GS k m L data (m = 65-73)
		cmd = append(cmd, byte(barType))      // Tipo de 65 a 73
		cmd = append(cmd, byte(contentLen))   // Byte de longitud L
		cmd = append(cmd, []byte(content)...) // Datos
	}

	_, err := p.Connector.Write(cmd)
	return err
}
