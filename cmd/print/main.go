// main.go
package main

import (
	// "fmt"
	"log"
	"os"
	// "time" // Descomentar si necesitas pausas
	"pos-print.adcon.dev/internal/config"
	"pos-print.adcon.dev/pkg/connectors/windows" // !!! REEMPLAZA con la ruta real de tu módulo
	"pos-print.adcon.dev/pkg/escpos"             // !!! REEMPLAZA con la ruta real de tu módulo
)

func main() {
	cfg := config.ParseFlags()
	// Configurar el logger para incluir información útil
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

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

	log.Printf("Intentando conectar a la impresora de Windows: %s", cfg.Printer)

	// --- 1. Crear una instancia del PrintConnector ---
	// Usamos el WindowsPrintConnector que usa la API de Spooler.
	connector, err := connectors.NewWindowsPrintConnector(cfg.Printer)
	if err != nil {
		log.Fatalf("Error fatal al crear el conector de Windows para '%s': %v", cfg.Printer, err)
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
	// y de que el método finalize() del conector se llame (en nuestra simple
	// implementación de connector.Close(), esto es lo mismo).
	// Dejaremos solo el defer connector.Close() por simplicidad ya que Printer.Close()
	// simplemente llama a connector.Close() en este port.

	// --- 3. Usar los métodos de la clase Printer para enviar comandos ---
	log.Println("Enviando comandos de impresión ESC/POS a la cola de Windows...")

	// Configurar justificación y estilo
	if err = printer.SetJustification(escpos.JUSTIFY_CENTER); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err = printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}

	// Imprimir texto
	if err = printer.Text("Mini Súper El Centro\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Vendedor: Juan Pérez\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Cliente: Público en general\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}

	// Restablecer estilo y imprimir separador
	if err = printer.SetEmphasis(false); err != nil {
		log.Printf("Error al restablecer énfasis: %v", err)
	}
	if err = printer.Text("--------------------------------\n"); err != nil {
		log.Printf("Error al imprimir separador: %v", err)
	}

	// Imprimir detalles de artículos (alineado a la izquierda)
	if err = printer.SetJustification(escpos.JUSTIFY_LEFT); err != nil {
		log.Printf("Error al establecer justificación izquierda: %v", err)
	}

	if err = printer.Text("Coca-Cola 600ml   2 x $18.50 = $37.00\n"); err != nil {
		log.Printf("Error al imprimir artículo 1: %v", err)
	}
	if err = printer.Text("Galletas Oreo     1 x $15.00 = $15.00\n"); err != nil {
		log.Printf("Error al imprimir artículo 2: %v", err)
	}
	if err = printer.Text("Pan Bimbo         1 x $30.00 = $30.00\n"); err != nil {
		log.Printf("Error al imprimir artículo 3: %v", err)
	}

	// Imprimir total (en negrita y alineado a la derecha)
	if err = printer.Text("--------------------------------\n"); err != nil {
		log.Printf("Error al imprimir separador: %v", err)
	}
	if err = printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	if err = printer.SetJustification(escpos.JUSTIFY_RIGHT); err != nil {
		log.Printf("Error al establecer justificación derecha: %v", err)
	}
	if err = printer.Text("TOTAL: 51.50\n"); err != nil {
		log.Printf("Error al imprimir total: %v", err)
	}
	if err = printer.SetEmphasis(false); err != nil {
		log.Printf("Error al restablecer énfasis: %v", err)
	}
	if err = printer.SetJustification(escpos.JUSTIFY_CENTER); err != nil {
		log.Printf("Error al establecer justificación centro: %v", err)
	}

	if err = printer.Text("¡Gracias por tu compra!"); err != nil {
		log.Printf("Error al imprimir artículo 3: %v", err)
	}

	// --- Ejemplos de otras funcionalidades (descomentar para probar) ---
	// NOTA: La compatibilidad con códigos de barras, QR, imágenes, etc.
	// depende del driver de la impresora y de la configuración de RAW data.
	// Si el driver intenta interpretar los comandos, podría no funcionar.
	// Asegúrate de que el driver permita el paso de comandos ESC/POS crudos.

	// Código de barras (UPC-A, requiere 11 o 12 dígitos)
	// log.Println("Imprimiendo código de barras...")
	// if err = printer.SetBarcodeHeight(80); err != nil { log.Printf("Error SetBarcodeHeight: %v", err) }
	// if err = printer.SetBarcodeWidth(3); err != nil { log.Printf("Error SetBarcodeWidth: %v", err) }
	// if err = printer.SetBarcodeTextPosition(escpos.BARCODE_TEXT_BELOW); err != nil { log.Printf("Error SetBarcodeTextPosition: %v", err) }
	// // Ejemplo UPC-A: 11 o 12 dígitos. "012345678901"
	// if err = printer.Barcode("012345678901", escpos.BARCODE_UPCA); err != nil { log.Printf("Error Barcode: %v", err) }
	// if err = printer.Feed(2); err != nil { log.Printf("Error Feed: %v", err) } // Espacio después del código de barras

	// Código QR
	// log.Println("Imprimiendo código QR...")
	// // Contenido, nivel EC (L, M, Q, H), tamaño (1-16), modelo (1, 2, Micro)
	// if err = printer.QrCode("https://github.com/your-repo", escpos.QR_ECLEVEL_M, 6, escpos.QR_MODEL_2); err != nil { log.Printf("Error QrCode: %v", err) }
	// if err = printer.Feed(2); err != nil { log.Printf("Error Feed: %v", err) } // Espacio después del código QR

	// Impresión de imagen (requiere implementar EscposImage y sus métodos)
	// log.Println("Intentando imprimir imagen...")
	// // Supongamos que tienes una imagen cargada en un objeto img *escpos.EscposImage
	// // img, err := escpos.NewEscposImageFromBytes(imageData) // Implementar esta función
	// // if err == nil {
	// // 	// Puedes usar BitImage, BitImageColumnFormat o Graphics
	// // 	if printErr := printer.Graphics(img, escpos.IMG_DEFAULT); printErr != nil {
	// // 		log.Printf("Error al imprimir imagen: %v", printErr)
	// // 	}
	// // 	if feedErr := printer.Feed(2); feedErr != nil { log.Printf("Error Feed: %v", feedErr) }
	// // } else {
	// // 	log.Printf("Error al cargar la imagen: %v", err)
	// // }

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
