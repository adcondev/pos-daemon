package config

import (
	"flag"
)

// Config contiene la configuración del daemon POS obtenida desde flags
// de línea de comandos.
type Config struct {
	// Printer especifica el nombre de la impresora Windows a utilizar
	Printer string
	// DebugLog habilita el logging detallado para depuración
	DebugLog bool
}

// ParseFlags procesa los flags de línea de comandos y retorna una instancia
// de Config con los valores especificados por el usuario.
//
// Los flags disponibles son:
//   - printer: Nombre de la impresora (Windows)
//   - debug: Habilitar logs de depuración
func ParseFlags() *Config {
	printer := flag.String("printer", "", "Nombre de la impresora (Windows)")
	debug := flag.Bool("debug", false, "Habilitar logs de depuración")
	flag.Parse()

	return &Config{
		Printer:  *printer,
		DebugLog: *debug,
	}
}
