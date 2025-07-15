// main.go
package main

import (
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/bmp"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"

	"pos-daemon.adcon.dev/internal/local_config"
	"pos-daemon.adcon.dev/internal/ticket"
	"pos-daemon.adcon.dev/pkg/escpos"
	"pos-daemon.adcon.dev/pkg/escpos/connectors"
	cons "pos-daemon.adcon.dev/pkg/escpos/constants"
)

func main() {
	jsonBytes, err := local_config.JSONFileToBytes("./internal/api/schema/local_config.json")
	if err != nil {
		log.Printf("Error al leer archivo JSON de local_config: %v", err)
		return
	}

	dataConfig := &local_config.LocalConfig{}

	dataConfig, err = local_config.BytesToConfig(jsonBytes)
	if err != nil {
		log.Printf("Error al deserializar JSON a objeto: %v", err)
		return
	}

	// Configurar el logger según el valor de DebugLog
	if dataConfig.DebugLog {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println("Modo de depuración activado.")
	} else {
		log.SetOutput(os.Stdout)
		log.SetFlags(0) // Sin detalles adicionales
	}

	// --- CONFIGURACIÓN ---
	// Define el nombre EXACTO de la impresora instalada en Windows.
	// Puedes encontrar este nombre en "Panel de control" -> "Dispositivos e impresoras".
	// Click derecho en la impresora -> "Propiedades de impresora" -> Pestaña "General" -> Nombre.
	// Asegúrate de que la impresora esté configurada para aceptar datos RAW.
	// Para impresoras USB, a veces necesitan un driver especial que crea un puerto serial virtual,
	// o el driver nativo ya permite enviar RAW.
	// Este código asume que el driver está configurado correctamente para RAW data.
	// Si estás probando sin una impresora real o con problemas de driver/configuración,
	// puedes seguir usando la implementación anterior del conector que escribe a un archivo.

	log.Printf("Intentando conectar a la impresora de Windows: %s", dataConfig.Printer)

	// --- 1. Crear una instancia del WindowsPrintConnector ---
	// Usamos el WindowsPrintConnector que usa la API de Spooler.
	connector, err := connectors.NewWindowsPrintConnector(dataConfig.Printer)
	if err != nil {
		log.Fatalf("Error fatal al crear el conector de Windows para '%s': %v", dataConfig.Printer, err)
	}

	// IMPORTANTE: Asegurarse de cerrar el conector al finalizar.
	// Esto llamará a EndDocPrinter y ClosePrinter.
	defer func() {
		log.Println("Cerrando el conector de la impresora.")
		if closeErr := connector.Close(); closeErr != nil {
			// No usar log.Fatalf aquí ya que estamos en un defer y el programa ya terminará.
			log.Printf("Error al cerrar el conector: %v", closeErr)
		}
	}()
	log.Println("Conector de Windows (API Spooler) creado exitosamente.")

	// --- 2. Crear una instancia de la clase Printer ---
	// Pasamos el conector y nil para usar el CapabilityProfile por defecto.
	log.Println("Creando instancia de Printer.")
	printer, err := escpos.NewPrinter(connector, nil) // NewPrinter llama a Initialize() internamente
	if err != nil {
		// El constructor de Printer llama a Initialize(), que hace un primer Write().
		// Si Initialize falla, puede ser un problema de conexión o que el primer Write no funcionó.
		log.Fatalf("Error fatal al crear e inicializar la impresora: %v", err)
	}
	log.Println("Instancia de Printer creada e inicializada.")

	// IMPORTANTE: También es buena práctica usar defer en Printer.Close()
	// Aunque Connector.Close() también cerrará el handle, Printer.Close()
	// se asegura de que el búfer de impresión esté vacío (si se hubiera usado)
	// y de que el method finalize() del conector se llame (en nuestra simple
	// implementación de connector.Close(), esto es lo mismo).
	// Dejaremos solo el defer connector.Close() por simplicidad ya que Printer.Close()
	// simplemente llama a connector.Close() en este port.

	// --- 3. Usar los métodos de la clase Printer para enviar comandos ---
	log.Println("Enviando comandos de impresión ESC/POS a la cola de Windows...")

	jsonBytes, err = ticket.JSONFileToBytes("./internal/api/schema/ticket.json")
	if err != nil {
		log.Printf("Error al leer archivo JSON de tickets: %v", err)
		return
	}

	dataTicket := &ticket.TicketData{}

	dataTicket, err = ticket.BytesToTicket(jsonBytes)
	if err != nil {
		log.Printf("Error al deserializar JSON a objeto: %v", err)
		return
	}

	// Configurar justificación y estilo
	if err = printer.SetJustification(escpos.JUSTIFY_CENTER); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err = printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}

	// Tipo de fuente
	if err = printer.SetFont(0); err != nil {
		log.Printf("Error al establecer fuente: %v", err)
	}

	if err = printer.SetTextSize(1, 1); err != nil {
		log.Printf("Error al establcer fuente: %v", err)
	}

	// Imprimir texto
	if err = printer.Text("BARCODE\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}

	if err = printer.Feed(1); err != nil {
		log.Printf("Error Feed: %v", err)
	} // Espacio después del código de barras

	// Código de barras (UPC-A, requiere 11 o 12 dígitos)
	log.Println("Imprimiendo código de barras...")
	if err = printer.SetBarcodeHeight(80); err != nil {
		log.Printf("Error SetBarcodeHeight: %v", err)
	}
	if err = printer.SetBarcodeWidth(3); err != nil {
		log.Printf("Error SetBarcodeWidth: %v", err)
	}
	if err = printer.SetBarcodeTextPosition(cons.TextBelow); err != nil {
		log.Printf("Error SetBarcodeTextPosition: %v", err)
	}
	// Ejemplo UPC-A: 11 o 12 dígitos. "012345678901"
	if err = printer.Barcode("012345678901", cons.UpcA); err != nil {
		log.Printf("Error Barcode: %v", err)
	}
	if err = printer.Feed(1); err != nil {
		log.Printf("Error Feed: %v", err)
	} // Espacio después del código de barras

	// QR

	if err = printer.Text("QR Code: " + dataTicket.AutofacturaLink + "\n"); err != nil {
		log.Printf("Error al imprimir QR Code: %v", err)
	}

	// Generar el código QR en memoria
	// El parámetro 256 define el tamaño en píxeles
	qr, err := qrcode.New(dataTicket.AutofacturaLinkQr, qrcode.Medium)
	if err != nil {
		log.Fatalf("Error generando QR: %v", err)
	}

	// Obtener la imagen del QR
	var size = 256
	qrImage := qr.Image(size)

	// Guardar imagen como PNG
	pngFile, err := os.Create("./img/qr.png")
	if err != nil {
		log.Fatalf("Error creando archivo PNG: %v", err)
	}
	defer func(pngFile *os.File) {
		err := pngFile.Close()
		if err != nil {
			log.Printf("error al cerrar png: %v", err)
		}
	}(pngFile)

	if err = png.Encode(pngFile, qrImage); err != nil {
		log.Fatalf("Error guardando imagen PNG: %v", err)
	}
	log.Println("Se guardó el archivo qr.png con éxito")

	bmpFile, err := os.Create("./img/qr.bmp")
	if err != nil {
		log.Fatalf("Error creando archivo BMP: %v", err)
	}
	defer func(bmpFile *os.File) {
		err := bmpFile.Close()
		if err != nil {
			log.Printf("error al cerrar bmp: %v", err)
		}
	}(bmpFile)

	// Se puede convertir a BMP directamente, ya que Image es de tipo image.Image
	if err = bmp.Encode(bmpFile, qrImage.(image.Image)); err != nil {
		log.Fatalf("Error guardando imagen BMP: %v", err)
	}
	log.Println("Se guardó el archivo qr.bmp con éxito")

	// Crear un objeto escpos.Image desde la imagen generada
	// El valor 128 es el umbral para determinar qué píxeles son negros (0-255)
	escposQR := escpos.NewEscposImage(qrImage, 128)

	// Imprimir usando uno de los métodos disponibles
	// Opción 1: BitImage - básico pero compatible con la mayoría de impresoras
	if err = printer.BitImage(escposQR, escpos.IMG_DEFAULT); err != nil {
		log.Printf("Error al imprimir QR con BitImage: %v", err)
	}

	logoPath := "./img/perro.jpeg"
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		logoPath = "./img/perro.png"
	}
	logoFile, err := os.Open(logoPath)
	if err != nil {
		log.Fatalf("Error abriendo archivo de logo (%s): %v", logoPath, err)
	}
	defer func(logoFile *os.File) {
		err := logoFile.Close()
		if err != nil {
			log.Printf("main: error al cerrar archivo de logo")
		}
	}(logoFile)

	// Decodificar según el formato real
	imgLogo, format, err := image.Decode(logoFile)
	if err != nil {
		log.Fatalf("Error decodificando imagen de logo (%s): %v", logoPath, err)
	}
	log.Printf("Logo cargado desde %s (formato %s)", logoPath, format)

	// Imprimir la imagen con dithering de Floyd-Steinberg
	if err := printer.ImageWithDithering(imgLogo, escpos.IMG_DEFAULT, cons.FloydStein, cons.DefaultPrintSize); err != nil {
		log.Printf("Error al imprimir logo con dithering: %v", err)
	}

	// Alimentar papel al final
	if err = printer.Feed(4); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}

	// Cortar papel
	if err = printer.Cut(escpos.CUT_FULL, 0); err != nil { // CUT_FULL o CUT_PARTIAL
		log.Printf("Error al cortar papel: %v", err)
	}

	// Abrir cajón portamonedas (si está conectado a la impresora y es compatible)
	// log.Println("Enviando pulso para abrir cajón portamonedas...")
	// // Este comando es ESC p 0/1 t1 t2, que debería funcionar con la mayoría de drivers RAW.
	// if err = printer.Pulse(0, 120, 240); err != nil { // Pin 0, 120ms ON, 240ms OFF
	// 	log.Printf("Error al enviar pulso: %v", err)
	// }

	log.Println("Todos los comandos de impresión han sido enviados a la cola de Windows.")

	// El recibo debería aparecer en la impresora física asociada al nombre proporcionado.
}
