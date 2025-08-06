package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"pos-daemon.adcon.dev/internal/models"
	"pos-daemon.adcon.dev/internal/service"
	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/profile"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
)

// TODO: Modificar padding de líneas para que sea configurable
func main() {
	// Cargar configuración
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

	// Configurar el logger
	if dataConfig.DebugLog {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println("Modo de depuración activado.")
	} else {
		log.SetOutput(os.Stdout)
		log.SetFlags(0)
	}

	// 1. Crear conector
	log.Printf("Conectando a impresora: %s", dataConfig.Printer)
	conn, err := connector.NewWindowsPrintConnector(dataConfig.Printer)
	if err != nil {
		log.Fatalf("Error al crear conector: %v", err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error al cerrar conector: %v", err)
		}
	}(conn)

	// 3. Crear perfil de impresora
	// Detectar tipo de impresora por nombre o usar configuración
	prof := &profile.Profile{}
	if strings.Contains(strings.ToLower(dataConfig.Printer), "58mm") {
		prof = profile.CreateProfile58mm()
		log.Println("Usando perfil para impresora de 58mm")
	} else {
		prof = profile.CreateProfile80mm() // Por defecto 80mm
		log.Println("Usando perfil para impresora de 80mm")
	}

	// 2. Crear protocolo ESC/POS
	protocol := escpos.NewESCPOSProtocol()

	// Personalizar el perfil según necesidad
	prof.Model = dataConfig.Printer
	prof.Vendor = "Generic"

	// 4. Crear impresora genérica
	printer, err := posprinter.NewGenericPrinter(protocol, conn, prof)
	if err != nil {
		log.Fatalf("Error al crear impresora: %v", err)
	}

	writer := os.Stdout
	// 5. Crear constructor de tickets con la impresora genérica
	constructor := service.NewTicketConstructor(writer, printer)

	// Cargar template y datos de ticket
	templateData, err := os.ReadFile(filepath.Join("./internal/api/rest/", "new_ticket_template.json"))
	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
		os.Exit(1)
	}

	ticketData, err := os.ReadFile(filepath.Join("./internal/api/rest/", "new_ticket.json"))
	if err != nil {
		fmt.Printf("Error loading ticket data: %v\n", err)
		os.Exit(1)
	}

	// Parsear template y datos
	if err := constructor.LoadTemplateFromJSON(templateData); err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	if err := constructor.LoadTicketFromJSON(ticketData); err != nil {
		fmt.Printf("Error parsing ticket data: %v\n", err)
		os.Exit(1)
	}

	// Imprimir ticket
	if err = constructor.PrintTicket(); err != nil {
		fmt.Printf("Error printing ticket: %v\n", err)
		os.Exit(1)
	}
}
