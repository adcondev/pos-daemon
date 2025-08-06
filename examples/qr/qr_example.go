package main

import (
	"log"

	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/profile"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
)

func main() {
	// === Crear conector ===
	// Seleccionar la impresora según tu configuración
	// printerName := "80mm EC-PM-80250"
	printerName := "58mm PT-210"

	// === Crear Perfil de impresora ===
	// Puedes definir un perfil si necesitas configuraciones específicas
	// prof := profile.CreateProfile80mm()
	prof := profile.CreatePt210() // Usar perfil de 58mm

	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Error al crear conector: %v", err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error al cerrar conector de impresora: %v", err)
		}
	}(conn)

	// === Crear protocolo ===
	// Aquí es donde seleccionas qué protocolo usar (ESC/POS, ZPL, etc.)
	proto := escpos.NewESCPOSProtocol()

	// === Crear impresora genérica ===
	// === Inicializar impresora ===
	printer, err := posprinter.NewGenericPrinter(proto, conn, prof)
	if err != nil {
		log.Fatalf("Error al crear impresora: %v", err)
	}
	defer func(printer *posprinter.GenericPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar impresora: %v", err)
		}
	}(printer)

	// === Imprimir título ===
	if err := printer.SetFont(command.FontA); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetJustification(command.AlignCenter); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(true); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("PRUEBA DE QR á é í ó ú ñ"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(false); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	// === Imprimir QR Code ===
	if err := printer.PrintQR(
		"https://github.com/AdConDev/pos-daemon", // Contenido del QR Code
		command.Model2,                           // Modelo de QR Code (Model1, Model2)
		command.ECHigh,                           // Nivel de corrección de errores (Low, Medium, High, Highest)
		8,                                        // Tamaño del módulo (1-16)
		256,                                      // Tamaño del QR Code (en pixeles, si el protocolo no soporta QR nativo)
	); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Cut(command.CutFeed, 1); err != nil {
		log.Printf("Error al cortar: %v", err)
	}
	log.Println("Impresión de QR completada exitosamente.")
}
