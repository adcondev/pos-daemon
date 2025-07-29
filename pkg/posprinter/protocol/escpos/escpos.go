package escpos

import (
	"pos-daemon.adcon.dev/pkg/posprinter/command"
)

// TODO: Implementar los métodos de la interfaz Printer y Connector
// ESCPOSProtocol implementa Protocol para ESC/POS
type ESCPOSProtocol struct {
	// ...configuración
}

// TODO: Implementar los métodos de la interfaz Printer y Connector y adaptar los comandos de ESC/POS
// SetJustification convierte el tipo genérico al específico de ESC/POS
func (p *ESCPOSProtocol) SetJustification(justification command.Alignment) []byte {
	// Convierte de comando.Alignment a escpos.Justify
	var escposJustify Justify
	switch justification {
	case command.AlignLeft:
		escposJustify = Left
	case command.AlignCenter:
		escposJustify = Center
	case command.AlignRight:
		escposJustify = Right
	default:
		escposJustify = Left // Default
	}

	// Genera el comando con el tipo específico de ESC/POS
	return []byte{ESC, 'a', byte(escposJustify)}
}

// SetBarcodeTextPosition convierte el tipo genérico al específico de ESC/POS
func (p *ESCPOSProtocol) SetBarcodeTextPosition(position command.BarcodeTextPosition) []byte {
	// Conversión similar
	var escposPosition BarcodeTextPos
	// ...conversión según corresponda

	return []byte{GS, 'H', byte(escposPosition)}
}

// Otras conversiones...
