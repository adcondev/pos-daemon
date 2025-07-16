package escpos

import "fmt"

const (
	// Modo de corte de papel
	CUT_FULL    int = 65 // 'A'
	CUT_PARTIAL int = 66 // 'B'
)

// Pulse envía un pulso a un pin del conector del cajón portamonedas para abrirlo.
func (p *Printer) Pulse(pin int, onMS, offMS int) error {
	if err := validateInteger(pin, 0, 1, "Pulse", "pin"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Pin 0 o 1
	if err := validateInteger(onMS, 1, 511, "Pulse", "onMS"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Tiempo ON en ms (1-511)
	if err := validateInteger(offMS, 1, 511, "Pulse", "offMS"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Tiempo OFF en ms (1-511) - a menudo ignorado por la impresora para el segundo pulso

	// Comando: ESC p m t1 t2
	// m: pin del cajón (0 o 1). PHP usa pin + 48 ('0' o '1'). Replicamos.
	// t1: Tiempo ON (t1 * 2 ms). PHP envía on_ms / 2. Replicamos.
	// t2: Tiempo OFF (t2 * 2 ms). PHP envía off_ms / 2. Replicamos.
	cmd := []byte{ESC, 'p', byte(pin + 48), byte(onMS / 2), byte(offMS / 2)}
	_, err := p.Connector.Write(cmd)
	return err
}

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
