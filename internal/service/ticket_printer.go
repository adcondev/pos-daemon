package service

import (
	"encoding/json"
	"fmt"
	"io"
	tckt "pos-daemon.adcon.dev/internal/ticket"
	tmpt "pos-daemon.adcon.dev/internal/ticket_template"
	"strings"
)

// TicketConstructor handles the construction and printing of tickets
type TicketConstructor struct {
	template tmpt.TicketTemplate
	ticket   tckt.Ticket
	writer   io.Writer
}

// NewTicketConstructor creates a new ticket constructor with the specified writer
func NewTicketConstructor(writer io.Writer) *TicketConstructor {
	return &TicketConstructor{
		writer: writer,
	}
}

func (tc *TicketConstructor) LoadTemplateFromJSON(data []byte) error {
	if err := json.Unmarshal(data, &tc.template); err != nil {
		return fmt.Errorf("failed to parse template JSON: %w", err)
	}
	return nil
}

// LoadTicketFromJSON loads ticket data from JSON data
func (tc *TicketConstructor) LoadTicketFromJSON(data []byte) error {
	if err := json.Unmarshal(data, &tc.ticket); err != nil {
		return fmt.Errorf("failed to parse ticket JSON: %w", err)
	}
	return nil
}

// PrintTicket prints the ticket according to the template configuration
func (tc *TicketConstructor) PrintTicket() error {
	// Check if template and ticket data are loaded
	if tc.template.Data.TicketWidth == 0 || tc.ticket.Data.Identificador == "" {
		return fmt.Errorf("template or ticket data not loaded")
	}

	tc.printHeader()
	tc.printCustomerInfo()
	tc.printTicketInfo()
	tc.printItems()
	tc.printTotals()
	tc.printTaxes()
	tc.printPaymentInfo()
	tc.printFooter()

	return nil
}

// printHeader prints the store information in the header
func (tc *TicketConstructor) printHeader() {
	tmpl := tc.template.Data
	ticket := tc.ticket.Data

	// Print logo placeholder if configured
	if tmpl.VerLogotipo {
		fmt.Fprintln(tc.writer, "[LOGO]")
	}

	// Print store name
	if tmpl.VerNombre {
		fmt.Fprintln(tc.writer, ticket.SucursalNombre)
	}

	// Print commercial name
	if tmpl.VerNombreC {
		fmt.Fprintln(tc.writer, ticket.SucursalNombreComercial)
	}

	// Print RFC
	if tmpl.VerRFC {
		fmt.Fprintf(tc.writer, "RFC: %s\n", ticket.SucursalRFC)
	}

	// Print address
	if tmpl.VerDom {
		fmt.Fprintf(tc.writer, "%s %s", ticket.SucursalCalle, ticket.SucursalNumero)
		if ticket.SucursalNumeroInt != "" {
			fmt.Fprintf(tc.writer, " Int. %s", ticket.SucursalNumeroInt)
		}
		fmt.Fprintln(tc.writer)
		fmt.Fprintf(tc.writer, "Col. %s\n", ticket.SucursalColonia)
		fmt.Fprintf(tc.writer, "%s, %s, %s\n",
			ticket.SucursalLocalidad, ticket.SucursalEstado, ticket.SucursalPais)
		fmt.Fprintf(tc.writer, "C.P. %s\n", ticket.SucursalCP)
	}

	// Print tax regime
	if tmpl.VerRegimen {
		fmt.Fprintf(tc.writer, "Régimen Fiscal: %s - %s\n",
			ticket.SucursalRegimenClave, ticket.SucursalRegimen)
	}

	// Print email
	if tmpl.VerEmail {
		fmt.Fprintf(tc.writer, "Email: %s\n", ticket.SucursalEmails)
	}

	// Print phone if configured
	if tmpl.VerTelefono && ticket.SucursalTelefono != "" {
		fmt.Fprintf(tc.writer, "Tel: %s\n", tc.ticket.Data.SucursalTelefono)
	}

	// Print custom header if available
	if tmpl.CambiarCabecera != "" {
		fmt.Fprintln(tc.writer, tmpl.CambiarCabecera)
	}

	fmt.Fprintln(tc.writer, strings.Repeat("-", int(tmpl.TicketWidth)))
}

// printCustomerInfo prints the customer information
func (tc *TicketConstructor) printCustomerInfo() {
	if tc.template.Data.VerNombreCliente {
		fmt.Fprintf(tc.writer, "Cliente: %s\n", tc.ticket.Data.ClienteNombre)
		fmt.Fprintf(tc.writer, "RFC: %s\n", tc.ticket.Data.ClienteRFC)
		fmt.Fprintln(tc.writer, strings.Repeat("-", int(tc.template.Data.TicketWidth)))
	}
}

// printTicketInfo prints folio, date, and store info
func (tc *TicketConstructor) printTicketInfo() {
	tmpl := tc.template.Data
	ticket := tc.ticket.Data

	if tmpl.VerFolio {
		fmt.Fprintf(tc.writer, "Folio: %s%s\n", ticket.Serie, ticket.Folio)
	}

	if tmpl.VerFecha {
		fmt.Fprintf(tc.writer, "Fecha: %s\n", ticket.FechaSistema)
	}

	if tmpl.VerTienda {
		fmt.Fprintf(tc.writer, "Tienda: %s\n", ticket.SucursalTienda)
	}

	fmt.Fprintf(tc.writer, "Vendedor: %s\n", ticket.Vendedor)
	fmt.Fprintln(tc.writer, strings.Repeat("-", int(tmpl.TicketWidth)))
}

// printItems prints the purchased items
func (tc *TicketConstructor) printItems() {
	tmpl := tc.template.Data

	// Print each item
	for _, item := range tc.ticket.Data.Conceptos {
		fmt.Fprintf(tc.writer, "%.2f x %s\n", item.Cantidad, item.Descripcion)

		if tmpl.VerPrecioU {
			fmt.Fprintf(tc.writer, "Precio: $%.2f\n", item.PrecioVenta)
		}

		fmt.Fprintf(tc.writer, "Total: $%.2f\n", item.Total)
		fmt.Fprintln(tc.writer, strings.Repeat("-", int(tmpl.TicketWidth)/2))
	}

	if tmpl.VerCantProductos {
		// Count total items
		totalItems := len(tc.ticket.Data.Conceptos)
		fmt.Fprintf(tc.writer, "Artículos: %d\n", totalItems)
	}
}

// printTotals prints the subtotal, discount and total
func (tc *TicketConstructor) printTotals() {
	ticket := tc.ticket.Data

	// If there's a discount, show it
	if ticket.Descuento > 0 {
		fmt.Fprintf(tc.writer, "Subtotal: $%.2f\n", calculateSubtotal(ticket))
		fmt.Fprintf(tc.writer, "Descuento: $%.2f\n", ticket.Descuento)
	}

	fmt.Fprintf(tc.writer, "TOTAL: $%.2f\n", ticket.Total)
	fmt.Fprintln(tc.writer, strings.Repeat("=", int(tc.template.Data.TicketWidth)))
}

// printTaxes prints tax information if configured
func (tc *TicketConstructor) printTaxes() {
	if tc.template.Data.VerImpuestos || tc.template.Data.VerImpuestosTotal {
		fmt.Fprintln(tc.writer, "DESGLOSE DE IMPUESTOS:")

		// Simple implementation - print all taxes
		for _, item := range tc.ticket.Data.Conceptos {
			for _, tax := range item.Impuestos {
				fmt.Fprintf(tc.writer, "%s (%s): $%.2f\n",
					tax.Codigo, tax.Tipo, tax.Importe)
			}
		}

		fmt.Fprintln(tc.writer, strings.Repeat("-", int(tc.template.Data.TicketWidth)))
	}
}

// printPaymentInfo prints payment details
func (tc *TicketConstructor) printPaymentInfo() {
	ticket := tc.ticket.Data

	fmt.Fprintln(tc.writer, "FORMA DE PAGO:")
	for _, payment := range ticket.Pagos {
		fmt.Fprintf(tc.writer, "%s: $%.2f\n", payment.FormaPago, payment.Cantidad)
	}

	fmt.Fprintf(tc.writer, "Pagado: $%.2f\n", ticket.Pagado)
	fmt.Fprintf(tc.writer, "Cambio: $%.2f\n", ticket.Cambio)

	fmt.Fprintln(tc.writer, strings.Repeat("-", int(tc.template.Data.TicketWidth)))
}

// printFooter prints the footer information
func (tc *TicketConstructor) printFooter() {
	tmpl := tc.template.Data

	if tmpl.VerLeyenda && tmpl.CambiarReclamacion != "" {
		fmt.Fprintln(tc.writer, tmpl.CambiarReclamacion)
		fmt.Fprintln(tc.writer, strings.Repeat("-", int(tmpl.TicketWidth)))
	}

	// Print footer message
	fmt.Fprintln(tc.writer, tmpl.CambiarPie)
}

func calculateSubtotal(ticket tckt.TicketData) float64 {
	sum := 0.0
	for _, item := range ticket.Conceptos {
		sum += item.Total
	}

	return ticket.Total // Placeholder
}
