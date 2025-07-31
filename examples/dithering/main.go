package main

import (
	"log"
	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/imaging"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
	"pos-daemon.adcon.dev/pkg/posprinter/utils"
)

func main() {
	// Crear conector
	// conn, err := connector.NewWindowsPrintConnector("58mm GOOJPRT PT-210")
	conn, err := connector.NewWindowsPrintConnector("80mm EC-PM-80250") // Cambia el nombre según tu impresora
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error cerrando conector: %v", err)
		}
	}(conn)

	// Crear protocolo y impresora
	proto := escpos.NewESCPOSProtocol()
	printer, err := posprinter.NewGenericPrinter(proto, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer func(printer *posprinter.GenericPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error cerrando impresora: %v", err)
		}
	}(printer)

	// Cargar imagen
	img := utils.LoadImage("./img/perro.jpeg")

	// === Opción 1: Imprimir sin dithering ===
	// log.Println("Imprimiendo sin dithering...")
	// if err := printer.PrintImage(img, command.DensitySingle); err != nil {
	// 	log.Printf("Error: %v", err)
	//}

	// === Opción 2: Imprimir con Floyd-Steinberg ===
	log.Println("Imprimiendo con Floyd-Steinberg...")
	opts := posprinter.PrintImageOptions{
		Density:    command.DensitySingle,
		DitherMode: imaging.DitherFloydSteinberg,
		Threshold:  128,
		Width:      256,
	}

	// if err := printer.PrintImageWithOptions(img, opts); err != nil {
	// 	log.Printf("Error: %v", err)
	// }

	// === Opción 3: Imprimir con Atkinson ===
	log.Println("Imprimiendo con Atkinson...")
	opts.DitherMode = imaging.DitherAtkinson
	if err = printer.PrintImageWithOptions(img, opts); err != nil {
		log.Printf("Error: %v", err)
	}

	if err = printer.Feed(1); err != nil {
		log.Printf("Error alimentando papel: %v", err)
	}

	err = printer.Text("Fin del test de imágenes")
	if err != nil {
		return
	}

	err = printer.Feed(3)
	if err != nil {
		return
	}
	err = printer.Cut(command.CutFeed, 3)
	if err != nil {
		return
	}
}
