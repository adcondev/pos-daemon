package main

import (
	"image"
	"log"

	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/imaging"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
)

func main() {
	// Crear conector
	conn, err := connector.NewWindowsPrintConnector("Mi Impresora")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Crear protocolo y impresora
	proto := escpos.NewESCPOSProtocol()
	printer, err := posprinter.NewGenericPrinter(proto, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer printer.Close()

	// Cargar imagen
	img := loadImage("logo.png")

	// === Opción 1: Imprimir sin dithering ===
	log.Println("Imprimiendo sin dithering...")
	if err := printer.PrintImage(img, command.DensitySingle); err != nil {
		log.Printf("Error: %v", err)
	}

	// === Opción 2: Imprimir con Floyd-Steinberg ===
	log.Println("Imprimiendo con Floyd-Steinberg...")
	opts := posprinter.PrintImageOptions{
		Density:    command.DensitySingle,
		DitherMode: imaging.DitherFloydSteinberg,
		Threshold:  128,
	}
	if err := printer.PrintImageWithOptions(img, opts); err != nil {
		log.Printf("Error: %v", err)
	}

	// === Opción 3: Imprimir con Atkinson ===
	log.Println("Imprimiendo con Atkinson...")
	opts.DitherMode = imaging.DitherAtkinson
	if err := printer.PrintImageWithOptions(img, opts); err != nil {
		log.Printf("Error: %v", err)
	}
}

func loadImage(filename string) image.Image {
	// TODO: Implementar carga real de imagen
	return nil
}
