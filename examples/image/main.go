package main

import (
	"image"
	"log"

	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
	// En el futuro: "pos-daemon.adcon.dev/pkg/posprinter/protocol/zpl"
)

func main() {
	// === Crear conector ===
	conn, err := connector.NewWindowsPrintConnector("Mi Impresora")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("dithering: error al cerrar conector")
		}
	}(conn)

	// === Opción 1: Usar protocolo ESC/POS ===
	useESCPOS(conn)

	// === Opción 2: Usar protocolo ZPL (cuando esté implementado) ===
	// useZPL(conn)
}

func useESCPOS(conn connector.Connector) {
	// Crear protocolo ESC/POS
	protocol := escpos.NewESCPOSProtocol()

	// Crear impresora
	printer, err := posprinter.NewGenericPrinter(protocol, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer func(printer *posprinter.GenericPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("dithering: error al cerrar impresora")
		}
	}(printer)

	// Cargar imagen
	img := loadTestImage()

	// Imprimir con densidad normal
	if err := printer.PrintImage(img, command.DensitySingle); err != nil {
		log.Printf("Error imprimiendo imagen: %v", err)
	}
}

// Ejemplo de cómo sería con otro protocolo
// func useZPL(conn connector.Connector) {
// TODO: Cuando implementes ZPL
/*
	protocol := zpl.NewZPLProtocol()
	printer, err := posprinter.NewGenericPrinter(protocol, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer printer.Close()

	img := loadTestImage()

	// ZPL procesará la imagen de manera diferente internamente,
	// pero la API es la misma
	if err := printer.PrintImage(img, command.DensitySingle); err != nil {
		log.Printf("Error: %v", err)
	}
*/
// }

func loadTestImage() image.Image {
	// TODO: Cargar una imagen real
	// Por ahora, crear una imagen de prueba
	return image.NewGray(image.Rect(0, 0, 100, 100))
}
