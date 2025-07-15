package main

import (
	"fmt"
	"os"
	"path/filepath"

	srvc "pos-daemon.adcon.dev/internal/service"
)

func main() {
	// Create a new ticket constructor that outputs to stdout
	constructor := srvc.NewTicketConstructor(os.Stdout)

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
