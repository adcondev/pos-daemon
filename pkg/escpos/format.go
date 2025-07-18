package escpos

import (
	"fmt"
	"log"
	"math"
	bin "pos-daemon.adcon.dev/pkg/escpos/constants"
	cons "pos-daemon.adcon.dev/pkg/escpos/constants"
)

// SetJustification establece la alineación del texto (izquierda, centro, derecha).
func (p *Printer) SetJustification(justification cons.Justify) error {
	if err := ValidateJustifyMode(justification); err != nil {
		log.Printf("Justificacion no valida")
	}
	// ESC a n - n=0: izquierda, 1: centro, 2: derecha
	cmd := []byte{ESC, 'a', byte(justification)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetFont establece la fuente (A, B o C).
func (p *Printer) SetFont(font cons.Font) error {
	if err := ValidateFont(font); err != nil {
		log.Printf("Fuente no válida: %v", err)
	}
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
func (p *Printer) SetUnderline(underline cons.UnderlineMode) error {
	if err := ValidateUnderline(underline); err != nil {
		log.Printf("Subrayado no valido: %v", err)
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
	marginBytes, err := intLowHigh(margin, 2) // 2 bytes (nL nH)
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

	const dotsPerMM float64 = 560.0 / 80
	dots := int(math.Round(float64(width) * dotsPerMM))

	// GS W nL nH - Establece el ancho del área de impresión a nL + nH * 256 puntos
	widthBytes, err := intLowHigh(dots, bin.DimensionBytes) // 2 bytes (nL nH)
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
