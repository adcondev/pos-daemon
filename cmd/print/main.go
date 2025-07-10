// Comando print para impresión de tickets usando impresoras ESC/POS en Windows.
//
// Este programa demuestra el uso de la biblioteca ESC/POS con el conector
// de Windows para imprimir un ticket de venta de ejemplo. Utiliza la API
// del spooler de Windows para comunicarse con impresoras térmicas instaladas.
//
// Uso:
//   go run cmd/print/main.go -printer "NombreImpresora" [-debug]
//
// Flags:
//   -printer: Nombre exacto de la impresora instalada en Windows (requerido)
//   -debug:   Habilita logging detallado para depuración (opcional)
//
// Ejemplo:
//   go run cmd/print/main.go -printer "EC-PM-80250" -debug
//
// El programa imprime un ticket de venta de ejemplo que incluye:
//   - Encabezado del establecimiento centrado y en negrita
//   - Información del vendedor y cliente
//   - Lista de productos vendidos con precios
//   - Total de la venta en negrita y alineado a la derecha
//   - Mensaje de agradecimiento centrado
//   - Corte completo del papel al final
//
// Notas importantes:
//   - La impresora debe estar instalada y configurada para aceptar datos RAW
//   - El nombre debe coincidir exactamente con el mostrado en Windows
//   - Se requieren permisos para acceder al spooler de impresión
package main

import (
	"log"
	"os"

	"pos-daemon.adcon.dev/internal/config"
	"pos-daemon.adcon.dev/internal/platform/windows"
	"pos-daemon.adcon.dev/pkg/escpos"
)

func main() {
	cfg := config.ParseFlags()
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("Intentando conectar a la impresora de Windows: %s", cfg.Printer)

	// Crear una instancia del PrintConnector para Windows
	connector, err := windows.NewWindowsPrintConnector(cfg.Printer)
	if err != nil {
		log.Fatalf("Error fatal al crear el conector de Windows para '%s': %v", cfg.Printer, err)
	}

	// Asegurar cierre del conector al finalizar
	defer func() {
		log.Println("Cerrando el conector de la impresora.")
		if closeErr := connector.Close(); closeErr != nil {
			log.Printf("Error al cerrar el conector: %v", closeErr)
		}
	}()
	log.Println("Conector de Windows (API Spooler) creado exitosamente.")

	// Crear una instancia de la impresora ESC/POS
	log.Println("Creando instancia de Printer.")
	printer, err := escpos.NewPrinter(connector, nil) // NewPrinter llama a Initialize() internamente
	if err != nil {
		log.Fatalf("Error fatal al crear e inicializar la impresora: %v", err)
	}
	log.Println("Instancia de Printer creada e inicializada.")

	// Imprimir ticket de venta de ejemplo
	log.Println("Enviando comandos de impresión ESC/POS a la cola de Windows...")

	// Configurar justificación y estilo
	if err = printer.SetJustification(escpos.JUSTIFY_CENTER); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err = printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}

	// Imprimir encabezado
	if err = printer.Text("Mini Súper El Centro\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Vendedor: Juan Pérez\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}
	if err = printer.Text("Cliente: Público en general\n"); err != nil {
		log.Printf("Error al imprimir texto: %v", err)
	}

	// Restablecer estilo e imprimir separador
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
		log.Printf("Error al imprimir mensaje final: %v", err)
	}

	// Alimentar papel al final
	if err = printer.Feed(4); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}

	// Cortar papel
	if err = printer.Cut(escpos.CUT_FULL, 0); err != nil {
		log.Printf("Error al cortar papel: %v", err)
	}

	log.Println("Todos los comandos de impresión han sido enviados a la cola de Windows.")
}
