package escpos

import "pos-daemon.adcon.dev/pkg/posprinter/types"

// TODO: Comandos para dar formato al texto
// - Doble ancho/altura
// - Rotación de texto
// - Espaciado de caracteres

// SetFont mapea fuentes genéricas a ESC/POS
func (p *Commands) SetFont(font types.Font) []byte {
	var fontValue byte
	switch font {
	case types.FontA:
		fontValue = 0
	case types.FontB:
		fontValue = 1
	case types.FontC:
		fontValue = 2
	default:
		fontValue = 0
	}

	// ESC M n
	return []byte{ESC, 'M', fontValue}
}

// SetEmphasis activa/desactiva negrita
func (p *Commands) SetEmphasis(on bool) []byte {
	val := byte(0)
	if on {
		val = 1
	}
	// ESC E n
	return []byte{ESC, 'E', val}
}

// SetDoubleStrike activa/desactiva doble golpe
func (p *Commands) SetDoubleStrike(on bool) []byte {
	val := byte(0)
	if on {
		val = 1
	}
	// ESC G n
	return []byte{ESC, 'G', val}
}

// SetUnderline configura el subrayado
func (p *Commands) SetUnderline(underline types.UnderlineMode) []byte {
	var val byte
	switch underline {
	case types.UnderlineNone:
		val = 0
	case types.UnderlineSingle:
		val = 1
	case types.UnderlineDouble:
		val = 2
	default:
		val = 0
	}
	// ESC - n
	return []byte{ESC, '-', val}
}

func (p *Commands) SetTextSize(widthMultiplier, heightMultiplier int) []byte {
	// TODO: Implementar usando GS ! n
	// Hint: n = (widthMultiplier-1)<<4 | (heightMultiplier-1)
	return []byte{}
}

func (p *Commands) SetLineSpacing(height *int) []byte {
	// TODO: Si height es nil, usar ESC 2 (default)
	// Si no, usar ESC 3 n
	return []byte{}
}
