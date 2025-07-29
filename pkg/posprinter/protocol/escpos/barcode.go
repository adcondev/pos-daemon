package escpos

import (
	"errors"
	"fmt"
)

// SetBarcodeHeight establece la altura del código de barras en puntos.
func (p *ESCPrinter) SetBarcodeHeight(height int) error {
	if err := ValidateInteger(height, 1, 255, "SetBarcodeHeight", "altura"); err != nil {
		return fmt.Errorf("SetBarcodeHeight: %w", err)
	}
	// GS h n - Establece la altura del código de barras a n puntos
	cmd := []byte{GS, 'h', byte(height)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetBarcodeWidth establece el ancho de los módulos del código de barras.
func (p *ESCPrinter) SetBarcodeWidth(width int) error {
	if err := ValidateInteger(width, 1, 255, "SetBarcodeWidth", "ancho"); err != nil {
		return fmt.Errorf("SetBarcodeWidth: %w", err)
	}
	// GS w n - Establece el ancho horizontal de los módulos a n (normalmente 2 o 3)
	cmd := []byte{GS, 'w', byte(width)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetBarcodeTextPosition establece la posición del texto legible del código de barras.
func (p *ESCPrinter) SetBarcodeTextPosition(position BarcodeTextPos) error {
	if err := ValidateBarcodeTextPosition(position); err != nil {
		return fmt.Errorf("SetBarcodeTextPosition: %w", err)
	} // 0: ninguno, 1: arriba, 2: abajo, 3: ambos (no siempre soportado) - PHP válida 0-3
	cmd := []byte{GS, 'H', byte(position)}
	_, err := p.Connector.Write(cmd)
	return err
}

// Barcode imprime un código de barras.
// content es la cadena de datos del código de barras.
// barType es el tipo de código de barras (BarcodeUpca, BarcodeCode39, etc.).
func (p *ESCPrinter) Barcode(content string, barType BarcodeType) error {
	if err := ValidateBarcodeType(barType); err != nil {
		return fmt.Errorf("barcode: %w", err)
	}
	contentLen := len(content)

	// --- Validación de contenido basada en el tipo de código de barras (traducir regex y longitud) ---
	var validationErr error
	switch barType {
	case UpcA:
		validationErr = ValidateStringRegex(content, `^[0-9]{11,12}$`, "contenido UPCA")
		if validationErr == nil {
			validationErr = ValidateInteger(contentLen, 11, 12, "Barcode", "longitud contenido UPCA")
		}
	case UpcE:
		validationErr = ValidateStringRegex(content, `^([0-9]{6,8}|[0-9]{11,12})$`, "contenido UPCE")
		if validationErr == nil {
			validationErr = ValidateIntegerMulti(contentLen, [][]int{{6, 8}, {11, 12}}, "Barcode", "longitud contenido UPCE")
		}
	case Jan13:
		validationErr = ValidateStringRegex(content, `^[0-9]{12,13}$`, "contenido JAN13")
		if validationErr == nil {
			validationErr = ValidateInteger(contentLen, 12, 13, "Barcode", "longitud contenido JAN13")
		}
	case Jan8:
		validationErr = ValidateStringRegex(content, `^[0-9]{7,8}$`, "contenido JAN8")
		if validationErr == nil {
			validationErr = ValidateInteger(contentLen, 7, 8, "Barcode", "longitud contenido JAN8")
		}
	case Code39:
		// PHP regex: `^([0-9A-Z $%+\-./]+|\*[0-9A-Z $%+\-./]+\*)$`
		// Requiere un * al principio y al final, o no.
		validationErr = ValidateStringRegex(content, `^([0-9A-Z $%+\-./]+|\*[0-9A-Z $%+\-./]+\*)$`, "contenido CODE39")
		if validationErr == nil {
			validationErr = ValidateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE39")
		}
	case Itf:
		// PHP regex: `^([0-9]{2})+$` - requiere solo dígitos y longitud par.
		validationErr = ValidateStringRegex(content, `^([0-9]{2})+$`, "contenido ITF")
		if validationErr == nil {
			validationErr = ValidateInteger(contentLen, 2, 255, "Barcode", "longitud contenido ITF")
		}
		// Validar longitud par
		if validationErr == nil && contentLen%2 != 0 {
			validationErr = errors.New("la longitud del contenido ITF debe ser par")
		}
	case Codabar:
		// PHP regex: `^[A-Da-d][0-9$%+\-./:]+[A-Da-d]$` - inicia/termina con A-D, medio con dígitos/símbolos.
		validationErr = ValidateStringRegex(content, `^[A-Da-d][0-9$%+\-./:]+[A-Da-d]$`, "contenido Codabar")
		if validationErr == nil {
			validationErr = ValidateInteger(contentLen, 1, 255, "Barcode", "longitud contenido Codabar")
		}
	case Code93:
		// PHP regex: `^[\x00-\x7F]+$` - solo caracteres ASCII.
		validationErr = ValidateStringRegex(content, `^[\x00-\x7F]+$`, "contenido CODE93")
		if validationErr == nil {
			validationErr = ValidateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE93")
		}
	case Code128:
		// PHP regex: `^\{[A-C][\\x00-\\x7F]+$` - espera que el contenido empiece con {A, {B o {C y luego ASCII.
		// Esto es un poco inusual, ya que normalmente el usuario no proporciona los códigos de inicio/función de Code128.
		// Replicamos la validación de PHP.
		validationErr = ValidateStringRegex(content, `^\{[A-C][\x00-\x7F]+$`, "contenido CODE128")
		if validationErr == nil {
			validationErr = ValidateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE128")
		}
		if validationErr == nil && contentLen < 2 { // Necesita al menos '{' y un carácter de tipo
			validationErr = errors.New("el contenido Code128 debe tener al menos 2 caracteres ({A, {B, {C...)")
		}
	}
	if validationErr != nil {
		return fmt.Errorf("barcode: contenido '%s' inválido para el tipo %d: %w", content, barType, validationErr)
	}
	// --- Fin Validación ---

	// Lógica de envío:
	// PHP usa el comando GS k m d1...dk (m=0-6) si getSupportsBarcodeB() es false.
	// PHP usa el comando GS k m L d1...dL (m=65-73) si getSupportsBarcodeB() es true.
	// 'SupportsBarcodeB' en PHP parece referirse al soporte del formato de comando más nuevo (con byte de longitud L).

	cmd := []byte{GS, 'k'}
	if !p.Profile.SupportsBarcodeB {
		// Usar el formato de comando antiguo: GS k m data NUL (m = 0-6)
		// Validar que el tipo solicitado esté en el rango 65-71 (correspondiente a m 0-6)
		if barType < UpcA || barType > Codabar {
			return fmt.Errorf("barcode: el perfil de impresora no soporta el tipo de código de barras %d con el formato de comando antiguo", barType)
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
