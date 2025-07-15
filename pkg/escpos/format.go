package escpos

import (
	"fmt"
	"pos-daemon.adcon.dev/pkg/escpos/imaging"
)

const (
	// Fuentes
	FONT_A int = 0
	FONT_B int = 1
	FONT_C int = 2

	// Justificación del texto
	JUSTIFY_LEFT   int = 0
	JUSTIFY_CENTER int = 1
	JUSTIFY_RIGHT  int = 2

	// Modo de impresión (combinación de bits para ESC !)
	MODE_FONT_A        int = 0   // Bit 0 OFF for Font A
	MODE_FONT_B        int = 1   // Bit 0 ON for Font B
	MODE_EMPHASIZED    int = 8   // Bit 3 ON (Negrita)
	MODE_DOUBLE_HEIGHT int = 16  // Bit 4 ON (Doble Altura)
	MODE_DOUBLE_WIDTH  int = 32  // Bit 5 ON (Doble Ancho)
	MODE_UNDERLINE     int = 128 // Bit 7 ON (Subrayado)

	// Tipo de subrayado
	UNDERLINE_NONE   int = 0
	UNDERLINE_SINGLE int = 1
	UNDERLINE_DOUBLE int = 2
)

// SetJustification establece la alineación del texto (izquierda, centro, derecha).
func (p *Printer) SetJustification(justification int) error {
	if err := validateInteger(justification, JUSTIFY_LEFT, JUSTIFY_RIGHT, "SetJustification", "justificación"); err != nil {
		return fmt.Errorf("SetJustification: %w", err)
	}
	// ESC a n - n=0: izquierda, 1: centro, 2: derecha
	cmd := []byte{ESC, 'a', byte(justification)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetFont establece la fuente (A, B o C).
func (p *Printer) SetFont(font int) error {
	if err := validateInteger(font, FONT_A, FONT_C, "SetFont", "fuente"); err != nil {
		return fmt.Errorf("SetFont: %w", err)
	}
	// ESC M n - n=0: Fuente A, 1: Fuente B, 2: Fuente C
	cmd := []byte{ESC, 'M', byte(font)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetEmphasis habilita o deshabilita el modo enfatizado (negrita).
func (p *Printer) SetEmphasis(on bool) error {
	// ESC E n - n=1: habilitar, n=0: deshabilitar
	val := byte(0)
	if on {
		val = 1
	}
	cmd := []byte{ESC, 'E', val}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetDoubleStrike habilita o deshabilita el modo doble golpeo.
func (p *Printer) SetDoubleStrike(on bool) error {
	// ESC G n - n=1: habilitar, n=0: deshabilitar
	val := byte(0)
	if on {
		val = 1
	}
	cmd := []byte{ESC, 'G', val}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetUnderline establece el modo de subrayado (ninguno, simple, doble).
// Puede aceptar 0 (none), 1 (single), 2 (double).
func (p *Printer) SetUnderline(underline int) error {
	// La clase PHP también acepta booleanos y los convierte.
	// En Go, la validación de tipo estática nos da la garantía, así que solo validamos el rango entero.
	if err := validateInteger(underline, UNDERLINE_NONE, UNDERLINE_DOUBLE, "SetUnderline", "subrayado"); err != nil {
		return fmt.Errorf("SetUnderline: %w", err)
	}
	// ESC - n - n=0: ninguno, 1: simple, 2: doble
	cmd := []byte{ESC, '-', byte(underline)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetTextSize establece el tamaño del texto usando multiplicadores de ancho y alto (1-8).
func (p *Printer) SetTextSize(widthMultiplier, heightMultiplier int) error {
	if err := validateInteger(widthMultiplier, 1, 8, "SetTextSize", "multiplicador de ancho"); err != nil {
		return fmt.Errorf("SetTextSize: %w", err)
	}
	if err := validateInteger(heightMultiplier, 1, 8, "SetTextSize", "multiplicador de alto"); err != nil {
		return fmt.Errorf("SetTextSize: %w", err)
	}
	// GS ! n - n es una combinación de bits de los multiplicadores (ancho-1) * 16 + (alto-1)
	c := byte(((widthMultiplier - 1) << 4) | (heightMultiplier - 1))
	cmd := []byte{GS, '!', c}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetLineSpacing establece el espaciado entre líneas.
// Si height es nil, restablece al espaciado por defecto (ESC 2).
// Si height no es nil, establece el espaciado a height/180 o height/203 pulgadas (ESC 3 n).
func (p *Printer) SetLineSpacing(height *int) error {
	if height == nil {
		// ESC 2 - Restablecer espaciado de línea por defecto
		_, err := p.Connector.Write([]byte{ESC, '2'})
		return err
	}
	if err := validateInteger(*height, 1, 255, "SetLineSpacing", "altura"); err != nil {
		return fmt.Errorf("SetLineSpacing: %w", err)
	}
	// ESC 3 n - Establecer espaciado de línea a n
	cmd := []byte{ESC, '3', byte(*height)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetPrintLeftMargin establece el margen izquierdo de impresión en puntos.
func (p *Printer) SetPrintLeftMargin(margin int) error {
	if err := validateInteger(margin, 0, 65535, "SetPrintLeftMargin", "margen"); err != nil {
		return fmt.Errorf("SetPrintLeftMargin: %w", err)
	}
	// GS L nL nH - Establece el margen izquierdo a nL + nH * 256 puntos
	marginBytes, err := imaging.intLowHigh(margin, 2) // 2 bytes (nL nH)
	if err != nil {
		return fmt.Errorf("SetPrintLeftMargin: falló al formatear bytes del margen: %w", err)
	}
	cmd := []byte{GS, 'L'}
	cmd = append(cmd, marginBytes...)
	_, err = p.Connector.Write(cmd)
	return err
}

// SetPrintWidth establece el ancho del área de impresión en puntos.
func (p *Printer) SetPrintWidth(width int) error {
	if err := validateInteger(width, 1, 65535, "SetPrintWidth", "ancho"); err != nil {
		return fmt.Errorf("SetPrintWidth: %w", err)
	}
	// GS W nL nH - Establece el ancho del área de impresión a nL + nH * 256 puntos
	widthBytes, err := imaging.intLowHigh(width, 2) // 2 bytes (nL nH)
	if err != nil {
		return fmt.Errorf("SetPrintWidth: falló al formatear bytes del ancho: %w", err)
	}
	cmd := []byte{GS, 'W'}
	cmd = append(cmd, widthBytes...)
	_, err = p.Connector.Write(cmd)
	return err
}

// SetPrintBuffer no se porta directamente ya que el manejo del texto se simplificó.
// La funcionalidad de `PrintBuffer` (manejo de \n y escritura raw)
// está cubierta por `Text` y `TextRaw`.
