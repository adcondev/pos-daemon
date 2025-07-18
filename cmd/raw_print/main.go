//go:build windows

package main

import (
	"fmt"
	"golang.org/x/text/encoding/charmap" // Ya importado
	"syscall"
	"time"
	"unsafe"
)

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

type docInfo1 struct {
	DocName    *uint16
	OutputFile *uint16
	DataType   *uint16
}

// *** FUNCIÓN PARA CODIFICAR A CP858 ***
func toCP858(s string) []byte {
	// Obtener el codificador para CP858
	encoder := charmap.CodePage858.NewEncoder()
	// Convertir la string (UTF-8) a bytes codificados en CP858
	encoded, err := encoder.Bytes([]byte(s))
	if err != nil {
		// En caso de error (ej. carácter no representable en CP858),
		// podrías loguear el error, o intentar un fallback.
		// Aquí, por simplicidad, devolvemos la string original (UTF-8),
		// aunque esto no solucionaría el problema del acento si falla la codificación.
		// Una mejor práctica sería reemplazar el carácter desconocido.
		fmt.Printf("Advertencia: No se pudo codificar string a CP858: %v (original: %q)\n", err, s)
		return []byte(s) // Fallback (probablemente no imprimirá bien el carácter problemático)
	}
	return encoded
}

// *** La función openPrinter se mantiene igual ***
func openPrinter(printerName string) (uintptr, error) {
	var h syscall.Handle
	namePtr, _ := syscall.UTF16PtrFromString(printerName)
	// El tercer argumento 0 es PRINTER_ACCESS_USE
	r1, _, err := procOpenPrinter.Call(uintptr(unsafe.Pointer(namePtr)), uintptr(unsafe.Pointer(&h)), 0)
	// OpenPrinterW retorna handle si r1 != 0, y error si r1 == 0.
	// La doc de windows dice que devuelve BOOL, es decir 0 o 1.
	// Go syscall.Call retorna r1, r2, err. r1 es el resultado de la llamada.
	// Si r1 es 0, significa FALSE (error en este caso).
	if r1 == 0 {
		return 0, fmt.Errorf("Error al abrir la impresora %s: %v", printerName, err)
	}
	return uintptr(h), nil
}

// *** La función writeRaw se mantiene igual (llama a openPrinter, StartDocPrinter, WritePrinter, etc.) ***
func writeRaw(printerName string, data []byte) error {
	hPrinter, err := openPrinter(printerName)
	if err != nil {
		return fmt.Errorf("error al abrir impresora: %v", err)
	}
	defer procClosePrinter.Call(hPrinter)

	// docInfo1 requiere punteros a strings UTF-16
	docName, _ := syscall.UTF16PtrFromString("TicketVenta")
	dataType, _ := syscall.UTF16PtrFromString("RAW")

	doc := docInfo1{
		DocName:  docName,
		DataType: dataType,
	}

	// StartDocPrinterW espera puntero a DOCINFOW
	r1, _, err := procStartDocPrinter.Call(hPrinter, 1, uintptr(unsafe.Pointer(&doc)))
	if r1 == 0 { // BOOL result, 0 indicates failure
		return fmt.Errorf("error al iniciar documento: %v", err)
	}
	defer func() {
		// Llamar a EndDocPrinter solo si StartDocPrinterW tuvo éxito
		if r1 != 0 {
			procEndDocPrinter.Call(hPrinter)
		}
	}()

	// StartPagePrinter es opcional para RAW, pero es buena práctica
	procStartPagePrinter.Call(hPrinter)
	// No verificamos error en StartPagePrinter ya que a menudo no es crítico para RAW

	var written uint32 // Variable para recibir cuántos bytes se escribieron

	// WritePrinter espera un puntero al buffer de datos, el tamaño y un puntero a la variable 'written'
	// Nota: data[0] es seguro si len(data) > 0. Si data puede estar vacío, necesitas manejarlo.
	// Para un ticket vacío, WritePrinter con len 0 simplemente no escribe nada.
	dataPtr := uintptr(0)
	if len(data) > 0 {
		dataPtr = uintptr(unsafe.Pointer(&data[0]))
	}

	r1, _, err = procWritePrinter.Call(
		hPrinter,
		dataPtr,                           // Puntero a los datos
		uintptr(len(data)),                // Tamaño de los datos
		uintptr(unsafe.Pointer(&written)), // Puntero a la variable que recibirá los bytes escritos
	)
	// WritePrinter retorna BOOL, 0 indicates failure
	if r1 == 0 {
		return fmt.Errorf("error al escribir en impresora: %v (bytes intentados: %d)", err, len(data))
	}

	// EndPagePrinter es opcional para RAW, pero es buena práctica
	procEndPagePrinter.Call(hPrinter)
	// No verificamos error en EndPagePrinter

	// Opcional: verificar si todos los bytes fueron escritos
	if written != uint32(len(data)) {
		fmt.Printf("Advertencia: Solo se escribieron %d de %d bytes\n", written, len(data))
	}

	return nil // Éxito
}

func main() {
	// 🔧 Asegúrate que este nombre coincide exactamente con el nombre de la impresora en Windows
	printerName := "EC-PM-80250"

	now := time.Now().Format("2006-01-02 15:04")

	// Construir comandos ESC/POS
	cmd := []byte{}
	cmd = append(cmd, 0x1B, 0x40) // ESC @ -> Inicializa impresora
	// command = append(command, 0x1B, 0x74, 0x11) // ESC t 17 -> Windows-1252 / Latin-1
	cmd = append(cmd, 0x1B, 0x74, 0x13) // ESC t 19 -> Código de página CP858 (Mantener si la impresora lo soporta bien)

	cmd = append(cmd, 0x1B, 0x61, 0x01) // ESC a 1 -> Centrado
	// *** Usar toCP858 para strings con acentos o caracteres especiales ***
	cmd = append(cmd, toCP858("Mini Súper El Centro\n")...)
	cmd = append(cmd, 0x1B, 0x61, 0x00) // Alinear a la izquierda

	cmd = append(cmd, toCP858("Fecha: "+now+"\n")...)
	cmd = append(cmd, toCP858("Vendedor: Juan Pérez\n")...)
	cmd = append(cmd, toCP858("Cliente: Público en general\n")...) // Ya usaba toLatin1, ahora cambiamos a CP858 si es necesario

	cmd = append(cmd, []byte("--------------------------------\n")...) // ASCII simple, no necesita codificación especial

	// Asumimos que estos textos también podrían necesitar CP858 si contienen € u otros símbolos
	cmd = append(cmd, toCP858("Coca-Cola 600ml   2 x $18.50 = $37.00\n")...)
	cmd = append(cmd, toCP858("Galletas Oreo     1 x $15.00 = $15.00\n")...)
	cmd = append(cmd, toCP858("Pan Bimbo         1 x $30.00 = $30.00\n")...)
	cmd = append(cmd, []byte("--------------------------------\n")...)

	subtotal := 37.0 + 15.0 + 30.0
	// Es mejor formatear el número *antes* de codificarlo
	subtotalStr := fmt.Sprintf("Subtotal:           $%6.2f\n", subtotal)
	iva := subtotal * 0.16
	ivaStr := fmt.Sprintf("IVA 16%%:            $%6.2f\n", iva) // El %% imprime un solo %
	total := subtotal + iva
	totalStr := fmt.Sprintf("TOTAL:              $%6.2f\n", total)

	cmd = append(cmd, toCP858(subtotalStr)...)
	cmd = append(cmd, toCP858(ivaStr)...)
	cmd = append(cmd, toCP858(totalStr)...)
	cmd = append(cmd, []byte("\n")...)

	cmd = append(cmd, 0x1B, 0x61, 0x01)                        // Centrado
	cmd = append(cmd, toCP858("¡Gracias por tu compra!\n")...) // Contiene '¡' y 'ú'
	// command = append(command, []byte("\n\n\n")...)
	cmd = append(cmd, []byte{0x1B, 'd', byte(3)}...)

	cmd = append(cmd, 0x1D, 0x56, 0x00) // GS V 0 -> Corte total

	// Enviar a la impresora
	err := writeRaw(printerName, cmd)
	if err != nil {
		fmt.Println("❌ Error al imprimir:", err)
	} else {
		fmt.Println("✅ Ticket impreso correctamente.")
	}
}
