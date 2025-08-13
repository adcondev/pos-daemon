package escpos

import "strings"

// TODO: Comandos fundamentales y de uso común

// Text convierte texto a bytes con encoding apropiado
func (p *Commands) Text(str string) []byte {
	cmd := strings.ReplaceAll(str, "\n", string(LF))
	return []byte(cmd)
}

// TextLn agrega un salto de línea al final
func (p *Commands) TextLn(str string) []byte {
	text := p.Text(str)
	// Agregar LF al final
	return append(text, LF)
}

// TextRaw envía bytes sin procesar
func (p *Commands) TextRaw(str string) []byte {
	return []byte(str)
}
