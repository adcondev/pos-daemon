package windows

import (
	"errors"
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

// Variables globales para acceder a las funciones del DLL winspool.drv
var (
	winspool            = syscall.NewLazyDLL("winspool.drv")
	procOpenPrinter     = winspool.NewProc("OpenPrinterW")
	procClosePrinter    = winspool.NewProc("ClosePrinter")
	procStartDocPrinter = winspool.NewProc("StartDocPrinterW")
	procEndDocPrinter   = winspool.NewProc("EndDocPrinter")
	procAbortDocPrinter = winspool.NewProc("AbortDocPrinter")
	procWritePrinter    = winspool.NewProc("WritePrinter")
)

// docInfo1 representa la estructura DOC_INFO_1 de la API de Windows.
// Se utiliza para especificar información sobre un documento de impresión.
type docInfo1 struct {
	DocName    *uint16 // Nombre del documento
	OutputFile *uint16 // Archivo de salida (nil para imprimir directamente)
	DataType   *uint16 // Tipo de datos ("RAW" para datos sin procesar)
}

// WindowsPrintConnector implementa la interfaz PrintConnector para
// impresoras Windows utilizando la API del spooler.
type WindowsPrintConnector struct {
	printerName   string         // Nombre de la impresora Windows
	printerHandle syscall.Handle // Handle de la impresora abierta
	jobStarted    bool           // Indica si se ha iniciado un trabajo
	docInfo       *docInfo1      // Información del documento
}

// NewWindowsPrintConnector crea una nueva instancia del conector Windows.
// Abre una conexión con la impresora especificada y prepara la información
// del documento para trabajos de impresión RAW.
//
// Parámetros:
//   - printerName: nombre exacto de la impresora instalada en Windows
//
// Retorna un conector configurado o un error si no se puede abrir la impresora.
func NewWindowsPrintConnector(printerName string) (*WindowsPrintConnector, error) {
	if printerName == "" {
		return nil, errors.New("el nombre de la impresora no puede estar vacío")
	}

	printerNameUTF16, err := syscall.UTF16PtrFromString(printerName)
	if err != nil {
		return nil, fmt.Errorf("error al convertir el nombre de la impresora: %w", err)
	}

	handle, err := openPrinter(printerNameUTF16)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir la impresora '%s': %w", printerName, err)
	}

	docName, _ := syscall.UTF16PtrFromString("ESC/POS Print Job")
	dataType, _ := syscall.UTF16PtrFromString("RAW")

	doc := &docInfo1{
		DocName:    docName,
		OutputFile: nil,
		DataType:   dataType,
	}

	return &WindowsPrintConnector{
		printerName:   printerName,
		printerHandle: handle,
		jobStarted:    false,
		docInfo:       doc,
	}, nil
}

// Write envía datos a la impresora. Implementa la interfaz PrintConnector.
// Inicia automáticamente un trabajo de impresión en la primera escritura.
//
// Parámetros:
//   - data: bytes a enviar a la impresora
//
// Retorna el número de bytes escritos o un error si falla la operación.
func (c *WindowsPrintConnector) Write(data []byte) (int, error) {
	if c.printerHandle == 0 {
		return 0, errors.New("handle de impresora no válido")
	}

	if !c.jobStarted {
		jobID, err := startDocPrinter(c.printerHandle, c.docInfo)
		if err != nil {
			return 0, fmt.Errorf("no se pudo iniciar el trabajo de impresión: %w", err)
		}
		log.Printf("Trabajo de impresión iniciado (ID: %d)", jobID)
		c.jobStarted = true
	}

	bytesWritten, err := writePrinter(c.printerHandle, data)
	if err != nil {
		return int(bytesWritten), fmt.Errorf("falló al escribir en la impresora: %w", err)
	}

	if int(bytesWritten) != len(data) {
		log.Printf("Advertencia: solo se escribieron %d de %d bytes", bytesWritten, len(data))
		return int(bytesWritten), fmt.Errorf("solo se escribieron %d de %d bytes", bytesWritten, len(data))
	}

	return int(bytesWritten), nil
}

// Close cierra la conexión con la impresora y libera recursos.
// Implementa la interfaz PrintConnector.
//
// Finaliza cualquier trabajo de impresión en curso y cierra el handle.
// Retorna un error si no se pueden liberar correctamente los recursos.
func (c *WindowsPrintConnector) Close() error {
	var finalErr error

	if c.jobStarted {
		err := endDocPrinter(c.printerHandle)
		if err != nil {
			log.Printf("Falló EndDocPrinter: %v, intentando AbortDocPrinter...", err)
			if abortErr := abortDocPrinter(c.printerHandle); abortErr != nil {
				log.Printf("Falló AbortDocPrinter: %v", abortErr)
				finalErr = fmt.Errorf("falló EndDoc y AbortDoc: %v", abortErr)
			}
		} else {
			log.Println("Trabajo de impresión finalizado correctamente.")
		}
	}

	if c.printerHandle != 0 {
		if err := closePrinter(c.printerHandle); err != nil {
			log.Printf("Falló ClosePrinter: %v", err)
			if finalErr == nil {
				finalErr = fmt.Errorf("falló ClosePrinter: %w", err)
			}
		}
		c.printerHandle = 0
		c.jobStarted = false
	}

	return finalErr
}

// openPrinter abre una conexión con la impresora especificada.
// Función auxiliar que encapsula la llamada a OpenPrinterW.
func openPrinter(name *uint16) (handle syscall.Handle, err error) {
	var h syscall.Handle
	r1, _, err := procOpenPrinter.Call(
		uintptr(unsafe.Pointer(name)),
		uintptr(unsafe.Pointer(&h)),
		0,
	)
	if r1 == 0 {
		return 0, err
	}
	return h, nil
}

// closePrinter cierra el handle de la impresora.
// Función auxiliar que encapsula la llamada a ClosePrinter.
func closePrinter(handle syscall.Handle) error {
	r1, _, err := procClosePrinter.Call(uintptr(handle))
	if r1 == 0 {
		return err
	}
	return nil
}

// startDocPrinter inicia un nuevo trabajo de impresión.
// Función auxiliar que encapsula la llamada a StartDocPrinterW.
func startDocPrinter(handle syscall.Handle, docInfo *docInfo1) (uint32, error) {
	r1, _, err := procStartDocPrinter.Call(
		uintptr(handle),
		1,
		uintptr(unsafe.Pointer(docInfo)),
	)
	if r1 == 0 {
		return 0, err
	}
	return uint32(r1), nil
}

// endDocPrinter finaliza el trabajo de impresión actual.
// Función auxiliar que encapsula la llamada a EndDocPrinter.
func endDocPrinter(handle syscall.Handle) error {
	r1, _, err := procEndDocPrinter.Call(uintptr(handle))
	if r1 == 0 {
		return err
	}
	return nil
}

// abortDocPrinter cancela el trabajo de impresión actual.
// Función auxiliar que encapsula la llamada a AbortDocPrinter.
func abortDocPrinter(handle syscall.Handle) error {
	r1, _, err := procAbortDocPrinter.Call(uintptr(handle))
	if r1 == 0 {
		return err
	}
	return nil
}

// writePrinter envía datos al spooler de impresión.
// Función auxiliar que encapsula la llamada a WritePrinter.
func writePrinter(handle syscall.Handle, data []byte) (uint32, error) {
	var bytesWritten uint32
	r1, _, err := procWritePrinter.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	if r1 == 0 {
		return 0, err
	}
	return bytesWritten, nil
}
