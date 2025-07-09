package config

import (
	"flag"
)

type Config struct {
	Printer  string
	DebugLog bool
}

func ParseFlags() *Config {
	printer := flag.String("printer", "", "Nombre de la impresora (Windows)")
	debug := flag.Bool("debug", false, "Habilitar logs de depuración")
	flag.Parse()

	return &Config{
		Printer:  *printer,
		DebugLog: *debug,
	}
}
