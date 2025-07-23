// main.go
package main

import (
	// "fmt"
	"log"
	"os"
	"pos-daemon.adcon.dev/internal/models"
	"pos-daemon.adcon.dev/internal/service"
	"pos-daemon.adcon.dev/pkg/escpos"
	cons "pos-daemon.adcon.dev/pkg/escpos/protocol"
	"strconv"

	conn "pos-daemon.adcon.dev/pkg/escpos/connector"
)

func main() {
	jsonBytes, err := models.JSONFileToBytes("./internal/api/rest/config.json")
	if err != nil {
		log.Printf("Error al leer archivo JSON de local_config: %v", err)
		return
	}

	dataConfig, err := models.BytesToConfig(jsonBytes)
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
	connector, err := conn.NewWindowsPrintConnector(dataConfig.Printer)
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

	// --- 2. Crear una instancia de la clase ESCPrinter ---
	// Pasamos el conector y nil para usar el CapabilityProfile por defecto.
	log.Println("Creando instancia de ESCPrinter.")
	printer, err := escpos.NewPrinter(connector, nil) // NewPrinter llama a Initialize() internamente
	if err != nil {
		// El constructor de ESCPrinter llama a Initialize(), que hace un primer Write().
		// Si Initialize falla, puede ser un problema de conexión o que el primer Write no funcionó.
		log.Fatalf("Error fatal al crear e inicializar la impresora: %v", err)
	}
	log.Println("Instancia de ESCPrinter creada e inicializada.")

	// IMPORTANTE: También es buena práctica usar defer en ESCPrinter.Close()
	// Aunque Connector.Close() también cerrará el handle, ESCPrinter.Close()
	// se asegura de que el búfer de impresión esté vacío (si se hubiera usado)
	// y de que el method finalize() del conector se llame (en nuestra simple
	// implementación de conn.Close(), esto es lo mismo).
	// Dejaremos solo el defer conn.Close() por simplicidad ya que ESCPrinter.Close()
	// simplemente llama a conn.Close() en este port.

	// --- 3. Usar los métodos de la clase ESCPrinter para enviar comandos ---
	log.Println("Enviando comandos de impresión ESC/POS a la cola de Windows...")

	jsonBytes, err = models.JSONFileToBytes("./internal/api/rest/ticket.json")
	if err != nil {
		log.Printf("Error al leer archivo JSON de tickets: %v", err)
		return
	}

	dataTicket, err := models.BytesToNewTicket(jsonBytes)
	if err != nil {
		log.Printf("Error al deserializar JSON a objeto: %v", err)
		return
	}

	// Configurar justificación y estilo
	if err = printer.SetJustification(cons.Center); err != nil {
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

	heart := "***     ***\n******* *******\n***************\n*************\n*********\n*****\n***\n*\n"
	if err = printer.Text(heart); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}

	// Imprimir texto
	if err = printer.Text("Matriz\n" + dataTicket.SucursalNombre + "\n\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Nombre Comercial: " + dataTicket.SucursalNombreComercial + "\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("RFC: " + dataTicket.SucursalRFC + "\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Regimen Fiscal: " + dataTicket.SucursalRegimen + "\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Email: " + dataTicket.ClienteEmail + "\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Cliente: " + dataTicket.ClienteNombre + "\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Folio: " + dataTicket.Folio + "\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Fecha: " + dataTicket.FechaSistema + "\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Tienda: " + dataTicket.SucursalTienda + "\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}

	// Restablecer estilo y imprimir separador
	if err = printer.SetEmphasis(false); err != nil {
		log.Printf("Error al restablecer énfasis: %v", err)
	}
	if err = printer.Text("------------------------------------------------\n"); err != nil {
		log.Printf("Error al imprimir separador: %v", err)
	}

	// Imprimir detalles de artículos (alineado a la izquierda)
	if err = printer.SetJustification(cons.Left); err != nil {
		log.Printf("Error al establecer justificación izquierda: %v", err)
	}

	const LEN_CANT int = 4
	const LEN_DESC int = 22
	const LEN_PRECIO int = 10
	const LEN_TOTAL int = 11
	const LEN_DECIMALES int = 2

	cant := service.PadCenter("CANT", LEN_CANT, ' ')
	producto := service.PadCenter("PRODUCTO", LEN_DESC, ' ')
	precio := service.PadCenter("PRECIO/U", LEN_PRECIO, ' ')
	subtotal := service.PadCenter("SUBTOTAL", LEN_TOTAL, ' ')

	if err = printer.Text(cant + producto + precio + subtotal + "\n"); err != nil {
		log.Printf("Error al imprimir artículo 1: %v", err)
	}

	const IvaTras int = 0
	const IepsTras int = 1
	const IvaRet int = 2
	const IsrRet = 3

	var subtotalSum float64
	var ivatrasladadoSum float64
	var iepstrasladadoSum float64
	var ivaretenidoSum float64
	var isrretenidoSum float64
	var totalFinal float64

	for _, v := range dataTicket.Conceptos {
		cant = service.PadCenter(strconv.FormatFloat(v.Cantidad, 'f', -1, 64), LEN_CANT, ' ')
		producto = service.PadCenter(service.Substr(v.Descripcion, LEN_DESC-2), LEN_DESC, ' ')
		precio = service.PadCenter("$"+service.FormatFloat(v.PrecioVenta, LEN_DECIMALES), LEN_PRECIO, ' ')
		subtotal = service.PadCenter("$"+service.FormatFloat(v.Total, LEN_DECIMALES), LEN_TOTAL, ' ')

		subtotalSum = subtotalSum + v.Total

		if len(v.Impuestos) > 0 {
			ivatrasladadoSum = ivatrasladadoSum + v.Impuestos[IvaTras].Importe
		}
		if len(v.Impuestos) > 1 {
			iepstrasladadoSum = iepstrasladadoSum + v.Impuestos[IepsTras].Importe
			ivaretenidoSum = ivaretenidoSum + v.Impuestos[IvaRet].Importe
			isrretenidoSum = isrretenidoSum + v.Impuestos[IsrRet].Importe
		}

		if err = printer.Text(cant + producto + precio + subtotal + "\n"); err != nil {
			log.Printf("Error al imprimir artículo 2: %v", err)
		}
	}

	totalFinal = subtotalSum + ivatrasladadoSum + iepstrasladadoSum + ivaretenidoSum + isrretenidoSum

	// Imprimir detalles de artículos (alineado a la izquierda)
	if err = printer.SetJustification(cons.Center); err != nil {
		log.Printf("Error al establecer justificación izquierda: %v", err)
	}
	// Imprimir total (en negrita y alineado a la derecha)
	if err = printer.Text("------------------------------------------------\n"); err != nil {
		log.Printf("Error al imprimir separador: %v", err)
	}
	if err = printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	if err = printer.SetJustification(cons.Right); err != nil {
		log.Printf("Error al establecer justificación derecha: %v", err)
	}

	if err = printer.Text("Subtotal: $" + service.FormatFloat(subtotalSum, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err = printer.Text("IVA Trasladado: $" + service.FormatFloat(ivatrasladadoSum, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err = printer.Text("IEPS Trasladado: $" + service.FormatFloat(iepstrasladadoSum, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err = printer.Text("IVA Retenido: $" + service.FormatFloat(ivaretenidoSum, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err = printer.Text("ISR Retenido: $" + service.FormatFloat(isrretenidoSum, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err = printer.Text("Total Calc: $" + service.FormatFloat(totalFinal, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err = printer.Text("Total Field: $" + service.FormatFloat(dataTicket.DocumentosPago[0].Total, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err = printer.Text("Efectivo: $" + service.FormatFloat(dataTicket.DocumentosPago[0].FormasPago[0].Cantidad, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err = printer.Text("Cambio: $" + service.FormatFloat(dataTicket.DocumentosPago[0].Cambio, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}

	if err = printer.SetEmphasis(false); err != nil {
		log.Printf("Error al restablecer énfasis: %v", err)
	}
	if err = printer.SetJustification(cons.Center); err != nil {
		log.Printf("Error al establecer justificación centro: %v", err)
	}

	if err = printer.Text("PAGADO\n¡Gracias por tu compra!"); err != nil {
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
	// // Supongamos que tienes una imagen cargada en un objeto _2D *escpos.EscposImage
	// // _2D, err := escpos.NewEscposImageFromBytes(imageData) // Implementar esta función
	// // if err == nil {
	// // 	// Puedes usar BitImage, BitImageColumnFormat o Graphics
	// // 	if printErr := printer.Graphics(_2D, escpos.IMG_DEFAULT); printErr != nil {
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
	if err = printer.Cut(cons.CUT_FULL, 0); err != nil { // CUT_FULL o CUT_PARTIAL
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
