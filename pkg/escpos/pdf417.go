package escpos

import (
	"errors"
	"fmt"
	"math"
)

const (
	// Opciones de PDF417
	PDF417_STANDARD  int = 0
	PDF417_TRUNCATED int = 1
)

// Pdf417Code imprime un código PDF417.
func (p *Printer) Pdf417Code(content string, width, heightMultiplier, dataColumnCount int, ec float64, options int) error {
	if content == "" {
		return nil
	} // No hacer nada si el contenido está vacío, como PHP
	if !p.profile.SupportsPdf417Code {
		return errors.New("los códigos PDF417 no están soportados en este perfil de impresora")
	}

	if err := validateString(content, "Pdf417Code", "contenido"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	}
	if err := validateInteger(width, 2, 8, "Pdf417Code", "ancho"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // Ancho del módulo (2-8 puntos)
	if err := validateInteger(heightMultiplier, 2, 8, "Pdf417Code", "multiplicador de alto"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // Multiplicador de alto de fila (2-8)
	if err := validateInteger(dataColumnCount, 0, 30, "Pdf417Code", "contador columnas de datos"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // 0 = automático
	if err := validateFloat(ec, 0.01, 4.00, "Pdf417Code", "nivel EC"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // Nivel EC como flotante (ej: 0.10 para 10%)
	if err := validateInteger(options, PDF417_STANDARD, PDF417_TRUNCATED, "Pdf417Code", "opciones"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // 0: estándar, 1: truncado

	cn := byte('0') // Código de función 1A (para PDF417) para GS ( k

	// 1. Establecer opciones (estándar/truncado): GS ( k pL pH cn 70 n
	//    cn=0, fn=70 ('F'), n=opciones (0 o 1)
	if err := p.wrapperSend2dCodeData(byte(70), cn, []byte{byte(options)}, 0); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer las opciones: %w", err)
	}

	// 2. Establecer contador de columnas de datos: GS ( k pL pH cn 65 n
	//    cn=0, fn=65 ('A'), n=contador columnas de datos (0-30)
	if err := p.wrapperSend2dCodeData(byte(65), cn, []byte{byte(dataColumnCount)}, 0); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer el contador de columnas: %w", err)
	}

	// 3. Establecer ancho del módulo: GS ( k pL pH cn 67 n
	//    cn=0, fn=67 ('C'), n=ancho (2-8)
	if err := p.wrapperSend2dCodeData(byte(67), cn, []byte{byte(width)}, 0); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer el ancho: %w", err)
	}

	// 4. Establecer multiplicador de alto de fila: GS ( k pL pH cn 68 n
	//    cn=0, fn=68 ('D'), n=multiplicador de alto (2-8)
	if err := p.wrapperSend2dCodeData(byte(68), cn, []byte{byte(heightMultiplier)}, 0); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer el multiplicador de alto: %w", err)
	}

	// 5. Establecer nivel EC: GS ( k pL pH cn 69 n m
	//    cn=0, fn=69 ('E'), n=ceil(ec * 10), m='1' (modo)
	//    PHP calcula ec_int = ceil(floatval(ec) * 10) y lo envía como byte, con m='1'. Replicamos.
	ecInt := byte(math.Ceil(ec * 10))
	if err := p.wrapperSend2dCodeData(byte(69), cn, []byte{ecInt}, byte('1')); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer el nivel EC: %w", err)
	}

	// 6. Almacenar datos: GS ( k pL pH cn 80 m d1...dk
	//    cn=0, fn=80 ('P'), m='0' (modo), d1...dk=contenido
	if err := p.wrapperSend2dCodeData(byte(80), cn, []byte(content), byte('0')); err != nil {
		return fmt.Errorf("Pdf417Code: falló al almacenar los datos: %w", err)
	}

	// 7. Imprimir símbolo: GS ( k pL pH cn 81 m
	//    cn=0, fn=81 ('Q'), m='0' (modo)
	if err := p.wrapperSend2dCodeData(byte(81), cn, []byte{}, byte('0')); err != nil { // Sin datos después de '0'
		return fmt.Errorf("Pdf417Code: falló al imprimir el símbolo: %w", err)
	}

	return nil
}
