//go:build windows

// Comando raw_print para impresión directa usando la API de Windows.
//
// Este programa demuestra la impresión RAW directa a impresoras Windows
// sin usar las abstracciones de la biblioteca ESC/POS. Utiliza directamente
// las funciones de la API del spooler de Windows para máximo control.
//
// Características:
//   - Acceso directo a la API winspool.drv de Windows
//   - Codificación de texto a CP858 para caracteres latinos
//   - Construcción manual de comandos ESC/POS
//   - Impresión de ticket con cálculos de IVA
//   - Manejo de errores detallado
//
// Este comando es específico para Windows y requiere permisos para
// acceder al spooler de impresión. Es útil para casos donde se necesita
// control total sobre los comandos enviados a la impresora.
//
// Uso:
//   go run cmd/raw_print/main.go
//
// Nota: El nombre de la impresora está hardcodeado como "EC-PM-80250"
// en la variable printerName dentro de main().
package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/text/encoding/charmap"
)

// Variables globales para acceder a las funciones del DLL winspool.drv
var (
	winspool             = syscall.NewLazyDLL("winspool.drv")
	procOpenPrinter      = winspool.NewProc("OpenPrinterW")
	procClosePrinter     = winspool.NewProc("ClosePrinter")
	procStartDocPrinter  = winspool.NewProc("StartDocPrinterW")
	procEndDocPrinter    = winspool.NewProc("EndDocPrinter")
	procStartPagePrinter = winspool.NewProc("StartPagePrinter")
	procEndPagePrinter   = winspool.NewProc("EndPagePrinter")
	procWritePrinter     = winspool.NewProc("WritePrinter")
)

// docInfo1 representa la estructura DOC_INFO_1 de la API de Windows
// utilizada para especificar información sobre un documento de impresión.
type docInfo1 struct {
	DocName    *uint16 // Nombre del documento
	OutputFile *uint16 // Archivo de salida (nil para imprimir directamente)
	DataType   *uint16 // Tipo de datos ("RAW" para datos sin procesar)
}

// toCP858 convierte una cadena UTF-8 a bytes codificados en CP858.
// Esta codificación es comúnmente soportada por impresoras ESC/POS
// para caracteres latinos incluyendo acentos y símbolos especiales.
//
// En caso de error de codificación, devuelve la cadena original como fallback.
func toCP858(s string) []byte {
	encoder := charmap.CodePage858.NewEncoder()
	encoded, err := encoder.Bytes([]byte(s))
	if err != nil {
		fmt.Printf("Advertencia: No se pudo codificar string a CP858: %v (original: %q)\n", err, s)
		return []byte(s) // Fallback
	}
	return encoded
}

// openPrinter abre una conexión con la impresora especificada.
// Utiliza la función OpenPrinterW de la API de Windows para obtener
// un handle que puede ser usado en operaciones subsequentes.
//
// Parámetros:
//   - printerName: nombre exacto de la impresora instalada en Windows
//
// Retorna el handle de la impresora o un error si no se puede abrir.
func openPrinter(printerName string) (uintptr, error) {
	var h syscall.Handle
	namePtr, _ := syscall.UTF16PtrFromString(printerName)
	r1, _, err := procOpenPrinter.Call(uintptr(unsafe.Pointer(namePtr)), uintptr(unsafe.Pointer(&h)), 0)
	if r1 == 0 {
		return 0, fmt.Errorf("Error al abrir la impresora %s: %v", printerName, err)
	}
	return uintptr(h), nil
}

// writeRaw envía datos RAW a la impresora usando la API de Windows.
// Maneja todo el flujo: abrir impresora, iniciar documento, escribir datos,
// finalizar documento y cerrar impresora.
//
// Parámetros:
//   - printerName: nombre de la impresora destino
//   - data: bytes RAW a enviar (comandos ESC/POS, texto, etc.)
//
// Retorna un error si cualquier paso del proceso falla.
func writeRaw(printerName string, data []byte) error {
	hPrinter, err := openPrinter(printerName)
	if err != nil {
		return fmt.Errorf("error al abrir impresora: %v", err)
	}
	defer procClosePrinter.Call(hPrinter)

	// Preparar información del documento
	docName, _ := syscall.UTF16PtrFromString("TicketVenta")
	dataType, _ := syscall.UTF16PtrFromString("RAW")

	doc := docInfo1{
		DocName:  docName,
		DataType: dataType,
	}

	// Iniciar documento de impresión
	r1, _, err := procStartDocPrinter.Call(hPrinter, 1, uintptr(unsafe.Pointer(&doc)))
	if r1 == 0 {
		return fmt.Errorf("error al iniciar documento: %v", err)
	}
	defer func() {
		if r1 != 0 {
			procEndDocPrinter.Call(hPrinter)
		}
	}()

	// Iniciar página (opcional para RAW pero recomendado)
	procStartPagePrinter.Call(hPrinter)

	var written uint32
	dataPtr := uintptr(0)
	if len(data) > 0 {
		dataPtr = uintptr(unsafe.Pointer(&data[0]))
	}

	// Escribir datos a la impresora
	r1, _, err = procWritePrinter.Call(
		hPrinter,
		dataPtr,
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
	)
	if r1 == 0 {
		return fmt.Errorf("error al escribir en impresora: %v (bytes intentados: %d)", err, len(data))
	}

	// Finalizar página
	procEndPagePrinter.Call(hPrinter)

	// Verificar que se escribieron todos los bytes
	if written != uint32(len(data)) {
		fmt.Printf("Advertencia: Solo se escribieron %d de %d bytes\n", written, len(data))
	}

	return nil
}

// main ejecuta la demostración de impresión RAW directa.
// Construye manualmente un ticket de venta con comandos ESC/POS
// y lo envía a la impresora usando la API de Windows.
//
// El ticket incluye:
//   - Inicialización de impresora y configuración de codepage CP858
//   - Encabezado centrado con nombre del establecimiento
//   - Información de fecha, vendedor y cliente
//   - Lista de productos con precios
//   - Cálculo automático de subtotal, IVA (16%) y total
//   - Mensaje de agradecimiento centrado
//   - Alimentación de papel y corte automático
func main() {
	// NOTA: Cambiar este nombre por el de tu impresora instalada
	printerName := "EC-PM-80250"

	now := time.Now().Format("2006-01-02 15:04")

	// Construir secuencia de comandos ESC/POS manualmente
	cmd := []byte{}
	
	// Inicializar impresora y configurar codepage
	cmd = append(cmd, 0x1B, 0x40)       // ESC @ - Inicializar impresora
	cmd = append(cmd, 0x1B, 0x74, 0x13) // ESC t 19 - Seleccionar CP858

	// Encabezado centrado
	cmd = append(cmd, 0x1B, 0x61, 0x01)                     // ESC a 1 - Centrar
	cmd = append(cmd, toCP858("Mini Súper El Centro\n")...) // Texto con codificación
	cmd = append(cmd, 0x1B, 0x61, 0x00)                     // ESC a 0 - Alinear izquierda

	// Información del ticket
	cmd = append(cmd, toCP858("Fecha: "+now+"\n")...)
	cmd = append(cmd, toCP858("Vendedor: Juan Pérez\n")...)
	cmd = append(cmd, toCP858("Cliente: Público en general\n")...)

	// Separador
	cmd = append(cmd, []byte("--------------------------------\n")...)

	// Productos
	cmd = append(cmd, toCP858("Coca-Cola 600ml   2 x $18.50 = $37.00\n")...)
	cmd = append(cmd, toCP858("Galletas Oreo     1 x $15.00 = $15.00\n")...)
	cmd = append(cmd, toCP858("Pan Bimbo         1 x $30.00 = $30.00\n")...)
	cmd = append(cmd, []byte("--------------------------------\n")...)

	// Cálculos
	subtotal := 37.0 + 15.0 + 30.0
	iva := subtotal * 0.16
	total := subtotal + iva

	// Totales con formato
	subtotalStr := fmt.Sprintf("Subtotal:           $%6.2f\n", subtotal)
	ivaStr := fmt.Sprintf("IVA 16%%:            $%6.2f\n", iva)
	totalStr := fmt.Sprintf("TOTAL:              $%6.2f\n", total)

	cmd = append(cmd, toCP858(subtotalStr)...)
	cmd = append(cmd, toCP858(ivaStr)...)
	cmd = append(cmd, toCP858(totalStr)...)
	cmd = append(cmd, []byte("\n")...)

	// Mensaje final centrado
	cmd = append(cmd, 0x1B, 0x61, 0x01)                         // Centrar
	cmd = append(cmd, toCP858("¡Gracias por tu compra!\n")...)

	// Alimentar papel y cortar
	cmd = append(cmd, []byte{0x1B, 'd', byte(3)}...) // ESC d 3 - Alimentar 3 líneas
	cmd = append(cmd, 0x1D, 0x56, 0x00)              // GS V 0 - Corte completo

	// Enviar comando completo a la impresora
	err := writeRaw(printerName, cmd)
	if err != nil {
		fmt.Println("❌ Error al imprimir:", err)
	} else {
		fmt.Println("✅ Ticket impreso correctamente.")
	}
}
