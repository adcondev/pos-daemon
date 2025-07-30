package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pos-daemon.adcon.dev/internal/models"
	"pos-daemon.adcon.dev/internal/service"
	conn "pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
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

	// Windows printer connection (fallback)
	log.Printf("Intentando conectar a la impresora de Windows: %s", dataConfig.Printer)
	// Crear conector
	connector, err := conn.NewWindowsPrintConnector(dataConfig.Printer)
	if err != nil {
		log.Fatalf("Error fatal al crear el conector de Windows para '%s': %v", dataConfig.Printer, err)
	}
	defer func(connector *conn.WindowsPrintConnector) {
		err := connector.Close()
		if err != nil {
			log.Printf("Error cerrando conector de impresora: %v", err)
		}
	}(connector)

	// === Opción 1: Usar el adaptador para compatibilidad ===
	escposPrinter, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		log.Fatalf("Error fatal al crear la impresora: %v", err)
	}

	// Crear un adaptador que implemente PrinterInterface
	printerAdapter := service.NewPrinterAdapter(nil, escposPrinter)

	// === Opción 2: Usar la nueva arquitectura directamente ===
	/*
		// Crear protocolo ESC/POS
		protocol := escpos.NewESCPOSProtocol()

		// Crear impresora genérica
		genericPrinter, err := posprinter.NewGenericPrinter(protocol, connector)
		if err != nil {
			log.Fatalf("Error al crear impresora genérica: %v", err)
		}

		// Crear adaptador ESC/POS para compatibilidad con tipos antiguos
		escposAdapter := escpos.NewESCPrinterAdapter(genericPrinter)

		// Crear adaptador para el ticket builder
		printerAdapter := service.NewPrinterAdapter(genericPrinter, escposAdapter)
	*/

	// Crear ticket constructor con el adaptador
	constructor := service.NewTicketConstructor(os.Stdout, printerAdapter)

	// Load template data
	templateData, err := os.ReadFile(filepath.Join("./internal/api/rest/", "new_ticket_template.json"))
	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
		os.Exit(1)
	}

	// Load ticket data
	ticketData, err := os.ReadFile(filepath.Join("./internal/api/rest/", "new_ticket.json"))
	if err != nil {
		fmt.Printf("Error loading ticket data: %v\n", err)
		os.Exit(1)
	}

	// Parse the template
	if err := constructor.LoadTemplateFromJSON(templateData); err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// Parse the ticket data
	if err := constructor.LoadTicketFromJSON(ticketData); err != nil {
		fmt.Printf("Error parsing ticket data: %v\n", err)
		os.Exit(1)
	}

	// Imprimir el ticket
	if err = constructor.PrintTicket(); err != nil {
		log.Printf("Error printing ticket: %v\n", err)
		os.Exit(1)
	}
}
