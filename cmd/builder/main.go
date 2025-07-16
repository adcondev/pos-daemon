package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pos-daemon.adcon.dev/internal/local_config"
	"pos-daemon.adcon.dev/pkg/escpos"
	"pos-daemon.adcon.dev/pkg/escpos/connectors"

	srvc "pos-daemon.adcon.dev/internal/service"
)

func main() {
	jsonBytes, err := local_config.JSONFileToBytes("./internal/api/schema/local_config.json")
	if err != nil {
		log.Printf("Error al leer archivo JSON de local_config: %v", err)
		return
	}

	dataConfig := &local_config.LocalConfigData{}

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
	log.Println("Creando instancia de Printer.")
	printer, err := escpos.NewPrinter(connector, nil) // NewPrinter llama a Initialize() internamente
	if err != nil {
		log.Fatalf("Error fatal al crear e inicializar la impresora: %v", err)
	}
	log.Println("Instancia de Printer creada e inicializada.")

	// --- 3. Usar los métodos de la clase Printer para enviar comandos ---
	log.Println("Enviando comandos de impresión ESC/POS a la cola de Windows...")
	// Create a new ticket constructor that outputs to stdout
	constructor := srvc.NewTicketConstructor(os.Stdout, printer)

	// Load template data
	templateData, err := os.ReadFile(filepath.Join("./internal/api/schema/", "ticket_template.json"))
	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
		os.Exit(1)
	}

	// Load ticket data
	ticketData, err := os.ReadFile(filepath.Join("./internal/api/schema/", "ticket.json"))
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

	// Print the ticket
	if err := constructor.PrintTicket(); err != nil {
		fmt.Printf("Error printing ticket: %v\n", err)
		os.Exit(1)
	}
}
