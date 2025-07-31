package main

import (
	"log"
	// Importar el conector
	conn "pos-daemon.adcon.dev/pkg/posprinter/connector"

	// Por ahora, usar el adaptador para mantener compatibilidad
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
	// TODO: Cuando completes la refactorización, cambiar a:
	// "pos-daemon.adcon.dev/pkg/posprinter"
	// "pos-daemon.adcon.dev/pkg/posprinter/command"
	// "pos-daemon.adcon.dev/pkg/posprinter/protocol"
)

func main() {
	// === Configuración ===
	// printerName := "58mm GOOJPRT PT-210"
	printerName := "80mm EC-PM-80250"

	// === Crear conector ===
	log.Printf("Intentando conectar a la impresora: %s", printerName)
	connector, err := conn.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Error al crear el conector para '%s': %v", printerName, err)
	}
	defer func(connector *conn.WindowsPrintConnector) {
		err := connector.Close()
		if err != nil {
			log.Printf("Error cerrando conector de impresora: %v", err)
		}
	}(connector)

	// === Crear impresora usando el adaptador ===
	// Usar el adaptador ESC/POS que mantiene compatibilidad
	printer, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		log.Fatalf("Error al crear e inicializar la impresora: %v", err)
	}
	defer func(printer *escpos.ESCPrinterAdapter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error cerrando impresora: %v", err)
		}
	}(printer)

	// === Prueba básica de impresión ===
	log.Println("Enviando comandos de prueba...")

	// Inicializar (ya se hace en NewPrinter, pero por si acaso)
	if err := printer.Initialize(); err != nil {
		log.Printf("Error al inicializar: %v", err)
	}

	// Texto centrado
	if err := printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error al centrar: %v", err)
	}

	// Texto en negrita
	if err := printer.SetEmphasis(true); err != nil {
		log.Printf("Error al activar negrita: %v", err)
	}

	// Imprimir título
	if err := printer.TextLn("PRUEBA DE IMPRESION"); err != nil {
		log.Printf("Error al imprimir título: %v", err)
	}

	// Desactivar negrita
	if err := printer.SetEmphasis(false); err != nil {
		log.Printf("Error al desactivar negrita: %v", err)
	}

	// Alinear a la izquierda
	if err := printer.SetJustification(escpos.Left); err != nil {
		log.Printf("Error al alinear izquierda: %v", err)
	}

	// Línea separadora
	if err := printer.TextLn("================================"); err != nil {
		log.Printf("Error al imprimir línea: %v", err)
	}

	// Contenido
	if err := printer.TextLn("Esta es una prueba básica"); err != nil {
		log.Printf("Error: %v", err)
	}

	if err := printer.TextLn("de la impresora desacoplada respecto a conectores, impresora y protocolo"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Feed y corte
	if err := printer.Feed(1); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}

	if err := printer.Cut(escpos.CUT_FULL, 0); err != nil {
		log.Printf("Error al cortar: %v", err)
	}

	log.Println("Impresión completada!")
}

// TODO: Cuando completes la refactorización, el código se verá así:
/*
func mainRefactored() {
	// Crear conector
	connector, err := connector.NewWindowsPrintConnector("Mi Impresora")
	if err != nil {
		log.Fatal(err)
	}
	defer connector.Close()

	// Crear protocolo
	protocol := escpos.NewESCPOSProtocol()

	// Crear impresora genérica
	printer, err := posprinter.NewGenericPrinter(protocol, connector)
	if err != nil {
		log.Fatal(err)
	}
	defer printer.Close()

	// Usar con tipos genéricos
	printer.SetJustification(command.AlignCenter)
	printer.SetEmphasis(true)
	printer.TextLn("HOLA MUNDO")
	printer.Cut(command.CutFull)
}
*/
