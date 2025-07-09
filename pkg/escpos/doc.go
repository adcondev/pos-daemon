// Package escpos proporciona una biblioteca completa para impresoras térmicas
// ESC/POS (Electronic Script/Point of Sale).
//
// Esta biblioteca permite controlar impresoras térmicas que utilizan comandos
// ESC/POS estándar, incluyendo impresión de texto, códigos de barras, códigos QR,
// imágenes, y control de formato como justificación, fuentes, y corte de papel.
//
// Características principales:
//   - Soporte para múltiples conectores (USB, TCP, Serial, Windows Spooler)
//   - Impresión de texto con codificación CP858 para caracteres latinos
//   - Códigos de barras: UPC-A, UPC-E, EAN13, EAN8, Code39, ITF, Codabar, Code93, Code128
//   - Códigos QR con diferentes niveles de corrección de error
//   - Códigos PDF417 para aplicaciones que requieren más datos
//   - Control de formato: justificación, fuentes, énfasis, subrayado
//   - Corte de papel (total y parcial)
//   - Apertura de cajón portamonedas
//   - Perfiles de capacidad para diferentes modelos de impresora
//
// Ejemplo básico de uso:
//
//	// Crear conector (implementar PrintConnector)
//	connector := &MiConector{...}
//	
//	// Crear impresora con perfil por defecto
//	printer, err := escpos.NewPrinter(connector, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	
//	// Imprimir texto centrado
//	printer.SetJustification(escpos.JUSTIFY_CENTER)
//	printer.Text("Ticket de Venta\n")
//	
//	// Imprimir código de barras
//	printer.Barcode("123456789012", escpos.BARCODE_UPCA)
//	
//	// Cortar papel
//	printer.Cut(escpos.CUT_FULL, 0)
//
// Para usar con Windows, combine con el paquete windows:
//
//	connector, err := windows.NewWindowsPrintConnector("NombreImpresora")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer connector.Close()
//	
//	printer, err := escpos.NewPrinter(connector, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	
//	// Usar impresora...
package escpos