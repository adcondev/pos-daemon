package main

import (
	"log"

	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/image"
	"pos-daemon.adcon.dev/pkg/posprinter/profile"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
	"pos-daemon.adcon.dev/pkg/posprinter/types"
)

func main() {
	// === Crear conector ===
	// Seleccionar la impresora según tu configuración
	// printerName := "80mm EC-PM-80250"
	printerName := "58mm GP-58N"

	// === Crear Perfil de impresora ===
	// Puedes definir un perfil si necesitas configuraciones específicas
	// prof := profile.CreateProfile80mm()
	prof := profile.CreateProfile58mm() // Usar perfil de 58mm

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

	// === Inicializar impresora ===
	if err := printer.Initialize(); err != nil {
		log.Printf("Error al inicializar: %v", err)
	}

	// === Cargar imagen ===
	img, err := image.LoadImage("./img/perro.jpeg")
	if err != nil {
		log.Fatalf("Error al cargar imagen: %v", err)
	}

	// === Imprimir título ===
	if err := printer.SetJustification(types.AlignCenter); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(true); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("PRUEBA DE DITHERING"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(false); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}

	// === Opción 1: Imprimir sin dithering ===
	if err := printer.TextLn("Imagen sin dithering:"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.PrintImage(img); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}

	opts := posprinter.PrintImageOptions{
		Density:    types.DensitySingle,
		DitherMode: image.DitherFloydSteinberg,
		Threshold:  128,
		Width:      256, // 0 = usar ancho original de imagen. La imagen podría salir más ancha que el papel
	}

	// === Opción 2: Imprimir con Floyd-Steinberg ===
	/*
		if err := printer.TextLn("Imagen con Floyd-Steinberg:"); err != nil {
			log.Printf("Error: %v", err)
		}

		if err := printer.PrintImageWithOptions(img, opts); err != nil {
			log.Printf("Error: %v", err)
		}
		if err := printer.Feed(2); err != nil {
			log.Printf("Error: %v", err)
		}
	*/

	// === Opción 3: Imprimir con Atkinson ===
	if err := printer.TextLn("Imagen con Atkinson:"); err != nil {
		log.Printf("Error: %v", err)
	}
	opts.DitherMode = image.DitherAtkinson
	if err := printer.PrintImageWithOptions(img, opts); err != nil {
		log.Printf("Error: %v", err)
	}

	// === Finalizar impresión ===
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("Fin del test de imágenes"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(3); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Cut(types.CutFeed, 3); err != nil {
		log.Printf("Error: %v", err)
	}
}
