package main

import (
	"fmt"
	"go.bug.st/serial"
	"io"
	"log"
	"os"
	"path/filepath"
	"pos-daemon.adcon.dev/internal/local_config"
	"pos-daemon.adcon.dev/pkg/escpos"
	conn "pos-daemon.adcon.dev/pkg/escpos/connectors"
	"strings"

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

	// Create connector based on the printer configuration
	var connector io.ReadWriteCloser

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatalf("Error al listar los puertos: %v", err)
	}

	if len(ports) == 0 {
		log.Println("No se encontraron puertos COM disponibles.")
	}

	fmt.Println("Estado de los puertos COM:")
	for _, port := range ports {
		if conn.IsPortInUse(port) {
			fmt.Printf("%s está en uso.\n", port)
		} else {
			fmt.Printf("%s está disponible.\n", port)
		}
	}

	if strings.HasPrefix(dataConfig.Printer, "COM") || strings.HasPrefix(dataConfig.Printer, "/dev/") || strings.HasPrefix(dataConfig.Printer, "/dev/tty") {
		// Serial port connection
		log.Printf("Intentando conectar a la impresora por puerto serial: %s", dataConfig.Printer)
		serialConfig := conn.DefaultSerialConfig()

		// Override default serial settings if provided in config
		if dataConfig.SerialBaudRate > 0 {
			serialConfig.BaudRate = dataConfig.SerialBaudRate
		}

		connector, err = conn.NewSerialConnector(dataConfig.Printer, serialConfig)
		if err != nil {
			log.Fatalf("Error fatal al crear el conector serial para '%s': %v", dataConfig.Printer, err)
		}
		log.Println("Conector serial creado exitosamente.")
	} else {
		// Windows printer connection (fallback)
		log.Printf("Intentando conectar a la impresora de Windows: %s", dataConfig.Printer)
		connector, err = conn.NewWindowsPrintConnector(dataConfig.Printer)
		if err != nil {
			log.Fatalf("Error fatal al crear el conector de Windows para '%s': %v", dataConfig.Printer, err)
		}
		log.Println("Conector de Windows (API Spooler) creado exitosamente.")
	}

	// IMPORTANTE: Asegurarse de cerrar el conector al finalizar.
	defer func() {
		log.Println("Cerrando el conector de la impresora.")
		if closeErr := connector.Close(); closeErr != nil {
			log.Printf("Error al cerrar el conector: %v", closeErr)
		}
	}()

	// --- 2. Crear una instancia de la clase Printer ---
	log.Println("Creando instancia de Printer.")
	printer, err := escpos.NewPrinter(connector, nil) // NewPrinter llama a Initialize() internamente
	if err != nil {
		log.Fatalf("Error fatal al crear e inicializar la impresora: %v", err)
	}
	log.Println("Instancia de Printer creada e inicializada.")

	// --- 3. Usar los métodos de la clase Printer para enviar comandos ---
	log.Println("Enviando comandos de impresión ESC/POS...")
	// Create a new ticket constructor that outputs to stdout
	constructor := srvc.NewTicketConstructor(os.Stdout, printer)

	// Load template data
	templateData, err := os.ReadFile(filepath.Join("./internal/api/schema/", "ticket_template.json"))
	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
		os.Exit(1)
	}

	// Load ticket data
	ticketData, err := os.ReadFile(filepath.Join("./internal/api/schema/", "new_ticket.json"))
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
