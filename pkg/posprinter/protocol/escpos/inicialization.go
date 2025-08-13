package escpos

// TODO: Comandos para inicialización y configuración básica de la impresora
// - Configuración de página de códigos
// - Configuración regional
// - Reinicio de impresora
// - Selección de modo

// Initialize genera el comando de inicialización ESC/POS
func (p *Commands) Initialize() []byte {
	// ESC @ - Reset printer
	return []byte{ESC, '@'}
}

// TODO: Comando compuesto para el final Feed(1) y Cut(1)

// Close genera comandos de cierre (si los hay)
func (p *Commands) Close() []byte {
	// ESC/POS no tiene un comando específico de cierre
	// pero podrías incluir un reset o feed final o ambos
	return []byte{}
}
