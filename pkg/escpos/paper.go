package escpos

import (
	"fmt"
)

const (
	// Modo de corte de papel
	CUT_FULL    int = 65 // 'A'
	CUT_PARTIAL int = 66 // 'B'
)

// Cut corta el papel.
// mode puede ser CUT_FULL o CUT_PARTIAL. lines es el número de líneas para alimentar antes de cortar (0-255).
func (p *Printer) Cut(mode int, lines int) error {
	// PHP usa chr(mode) donde mode es 65 ('A') o 66 ('B').
	// El comando estándar es GS V m [n], donde m es 0,1,48,49 (full/partial)
	// o m es 65,66 ('A'/'B') con un parámetro n adicional para líneas de avance.
	// Replicamos el comportamiento de PHP usando 'A' o 'B' y el parámetro lines.
	if err := validateInteger(mode, CUT_FULL, CUT_PARTIAL, "Cut", "modo"); err != nil {
		return fmt.Errorf("Cut: %w", err)
	} // 65 ('A') o 66 ('B')
	if err := validateInteger(lines, 0, 255, "Cut", "líneas"); err != nil {
		return fmt.Errorf("Cut: %w", err)
	}

	cmd := []byte{GS, 'V', byte(mode), byte(lines)} // GS V 'A'/'B' n
	_, err := p.Connector.Write(cmd)
	return err
}

// Feed avanza el papel el número especificado de líneas.
func (p *Printer) Feed(lines int) error {
	if err := validateInteger(lines, 1, 255, "Feed", "líneas"); err != nil {
		return fmt.Errorf("Feed: %w", err)
	}
	if lines <= 1 {
		// Usar solo LF para una línea es un poco más rápido a veces
		_, err := p.Connector.Write([]byte{LF})
		return err
	}
	// ESC d n - Imprime los datos del búfer y alimenta n líneas
	cmd := []byte{ESC, 'd', byte(lines)}
	_, err := p.Connector.Write(cmd)
	return err
}

// FeedReverse retrocede el papel el número especificado de líneas.
func (p *Printer) FeedReverse(lines int) error {
	if err := validateInteger(lines, 1, 255, "FeedReverse", "líneas"); err != nil {
		return fmt.Errorf("FeedReverse: %w", err)
	}
	// ESC e n - Alimenta el papel hacia atrás n líneas
	cmd := []byte{ESC, 'e', byte(lines)}
	_, err := p.Connector.Write(cmd)
	return err
}

// FeedForm alimenta el papel hasta el principio del siguiente formulario (poco común en impresoras de recibos).
func (p *Printer) FeedForm() error {
	// FF - Form Feed
	_, err := p.Connector.Write([]byte{FF})
	return err
}

// Release envía un comando (ESC q) que PHP describe como "liberar la impresora del estado de espera".
// Este comando NO es ESC/POS estándar y es probable que sea específico del fabricante (como Star).
func (p *Printer) Release() error {
	// Advertencia: ESC q es probablemente específico del fabricante y no estándar ESC/POS.
	_, err := p.Connector.Write([]byte{ESC, 'q'}) // PHP envía ESC seguido del byte 113 ('q')
	return err
}
