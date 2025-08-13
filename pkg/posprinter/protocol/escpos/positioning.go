package escpos

import "pos-daemon.adcon.dev/pkg/posprinter/types"

// TODO: Comandos para posicionar texto e imágenes
// - Tabulación
// - Posicionamiento absoluto
// - Posicionamiento relativo

func (p *Commands) SetPrintLeftMargin(margin int) []byte {
	// TODO: Implementar usando GS L nL nH
	return []byte{}
}

func (p *Commands) SetPrintWidth(width int) []byte {
	// TODO: Implementar usando GS W nL nH
	return []byte{}
}

// SetJustification convierte el tipo genérico al específico de ESC/POS
func (p *Commands) SetJustification(justification types.Alignment) []byte {
	// Mapear el tipo genérico al valor ESC/POS
	var escposValue byte
	switch justification {
	case types.AlignLeft:
		escposValue = 0 // ESC/POS: 0 = left
	case types.AlignCenter:
		escposValue = 1 // ESC/POS: 1 = center
	case types.AlignRight:
		escposValue = 2 // ESC/POS: 2 = right
	default:
		escposValue = 0 // Default to left
	}

	// ESC a n
	return []byte{ESC, 'a', escposValue}
}
