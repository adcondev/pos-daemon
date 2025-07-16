package service

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	tckt "pos-daemon.adcon.dev/internal/ticket"
	tmpt "pos-daemon.adcon.dev/internal/ticket_template"
	esc "pos-daemon.adcon.dev/pkg/escpos"
	cons "pos-daemon.adcon.dev/pkg/escpos/constants"

	"strconv"
)

const (
	LEN_CANT      int = 4
	LEN_DESC      int = 22
	LEN_PRECIO    int = 10
	LEN_TOTAL     int = 11
	LEN_DECIMALES int = 2
)

// TicketConstructor handles the construction and printing of tickets
type TicketConstructor struct {
	template tmpt.TicketTemplate
	ticket   tckt.Ticket
	writer   io.Writer
	printer  *esc.Printer
}

// NewTicketConstructor creates a new ticket constructor with the specified writer
func NewTicketConstructor(writer io.Writer, printer *esc.Printer) *TicketConstructor {
	return &TicketConstructor{
		writer:  writer,
		printer: printer,
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
		return fmt.Errorf("ticket printer: template or ticket data not loaded")
	}

	// Configurar justificación y estilo
	if err := tc.printer.SetJustification(cons.JUSTIFY_CENTER); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	// Tipo de fuente
	if err := tc.printer.SetFont(cons.FONT_A); err != nil {
		log.Printf("Error al establecer fuente: %v", err)
	}

	err := tc.printer.SetPrintWidth(int(tc.template.Data.TicketWidth))
	if err != nil {
		return fmt.Errorf("ticket printer: error al establecer ancho de impresión: %w", err)
	}
	err = tc.printer.SetPrintLeftMargin(0)
	if err != nil {
		return fmt.Errorf("ticket printer: error al establecer ancho de impresió izquierdo: %w", err)
	}
	if err := tc.printer.TextLn("Ticket Width: " + strconv.Itoa(int(tc.template.Data.TicketWidth))); err != nil {
		log.Printf("ticket_printer: error al imprimir texto: %v", err)
	}
	tc.printHeader()
	tc.printCustomerInfo()
	tc.printTicketInfo()

	// Configurar justificación y estilo
	if err := tc.printer.SetJustification(cons.JUSTIFY_LEFT); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err := tc.printer.SetEmphasis(false); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	// Tipo de fuente
	if err := tc.printer.SetFont(cons.FONT_B); err != nil {
		log.Printf("Error al establecer fuente: %v", err)
	}

	taxes := tc.printItems()
	tc.printTaxes(taxes)
	tc.printPaymentInfo()

	// Configurar justificación y estilo
	if err := tc.printer.SetJustification(cons.JUSTIFY_CENTER); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	// Tipo de fuente
	if err := tc.printer.SetFont(cons.FONT_A); err != nil {
		log.Printf("Error al establecer fuente: %v", err)
	}

	tc.printFooter()

	// Alimentar papel al final
	if err = tc.printer.Feed(2); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}

	// Cortar papel
	if err = tc.printer.Cut(esc.CUT_FULL, 0); err != nil { // CUT_FULL o CUT_PARTIAL
		log.Printf("Error al cortar papel: %v", err)
	}

	log.Println("Todos los comandos de impresión han sido enviados a la cola de Windows.")

	return nil
}

// printHeader prints the store information in the header
func (tc *TicketConstructor) printHeader() {
	tmpl := tc.template.Data
	datosTckt := tc.ticket.Data

	// Print custom header if available
	if tmpl.CambiarCabecera != "" {
		if err := tc.printer.TextLn("Cabecera: " + tmpl.CambiarCabecera); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print logo placeholder if configured
	if tmpl.VerLogotipo {
		logoPath := "./img/perro.jpeg"
		logoFile, err := os.Open(logoPath)
		if err != nil {
			log.Fatalf("datosTckt printer: error abriendo archivo de logo (%s): %v", logoPath, err)
		}
		defer func(logoFile *os.File) {
			err := logoFile.Close()
			if err != nil {
				log.Printf("datosTckt printer: error al cerrar archivo de logo")
			}
		}(logoFile)

		// Decodificar según el formato real
		imgLogo, format, err := image.Decode(logoFile)
		if err != nil {
			log.Fatalf("datosTckt printer: error decodificando imagen de logo (%s): %v", logoPath, err)
		}
		log.Printf("datosTckt printer: logo cargado desde %s (formato %s)", logoPath, format)

		// Imprimir la imagen con dithering de Floyd-Steinberg
		if err := tc.printer.ImageWithDithering(imgLogo, cons.IMG_DEFAULT, cons.FloydStein, cons.DefaultPrintSize); err != nil {
			log.Printf("datosTckt printer: error al imprimir logo con dithering: %v", err)
		}

		if err := tc.printer.Feed(1); err != nil {
			log.Printf("ticket_printer: error al alimentar papel después de imprimir logo: %v", err)
		}
	}

	// Print store name
	if tmpl.VerNombre {
		if err := tc.printer.TextLn("Sucursal Nombre: " + datosTckt.SucursalNombre); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print commercial name
	if tmpl.VerNombreC {
		if err := tc.printer.TextLn("Sucursal Nombre Comercial: " + datosTckt.SucursalNombreComercial); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print RFC
	if tmpl.VerRFC {
		if err := tc.printer.TextLn("Sucursal RFC: " + datosTckt.SucursalRFC); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print address
	if dom := ""; tmpl.VerDom {
		dom = fmt.Sprintf("%s %s", datosTckt.SucursalCalle, datosTckt.SucursalNumero)
		if datosTckt.SucursalNumeroInt != "" {
			dom = dom + fmt.Sprintf(" Int. %s", datosTckt.SucursalNumeroInt)
		}
		dom = dom + "\n"
		dom = dom + fmt.Sprintf("Col. %s\n", datosTckt.SucursalColonia)
		dom = dom + fmt.Sprintf("%s, %s, %s\n", datosTckt.SucursalLocalidad, datosTckt.SucursalEstado, datosTckt.SucursalPais)
		dom = dom + fmt.Sprintf("C.P. %s", datosTckt.SucursalCP)

		if err := tc.printer.TextLn("Sucursal Domicilio: " + dom); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print tax regime
	if reg := ""; tmpl.VerRegimen {
		reg = fmt.Sprintf("Sucursal Régimen Fiscal: %s - %s", datosTckt.SucursalRegimenClave, datosTckt.SucursalRegimen)
		if err := tc.printer.TextLn(reg); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print email
	if tmpl.VerEmail {
		if err := tc.printer.TextLn("Sucursal Email: " + datosTckt.SucursalEmails); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print phone if configured
	if tmpl.VerTelefono && datosTckt.SucursalTelefono != "" {
		if err := tc.printer.TextLn("Sucursal Telefono: " + datosTckt.SucursalTelefono); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}
	if err := tc.printer.Feed(1); err != nil {
		log.Printf("ticket_printer: error al imprimir texto: %v", err)
	}
}

// printCustomerInfo prints the customer information
func (tc *TicketConstructor) printCustomerInfo() {
	if tc.template.Data.VerNombreCliente {
		if err := tc.printer.TextLn("Cliente Nombre: " + tc.ticket.Data.ClienteNombre); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.TextLn("Cliente RFC: " + tc.ticket.Data.ClienteRFC); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.Feed(1); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}
}

// printTicketInfo prints folio, date, and store info
func (tc *TicketConstructor) printTicketInfo() {
	tmpl := tc.template.Data
	tcktData := tc.ticket.Data

	if serieFolio := ""; tmpl.VerFolio {
		serieFolio = fmt.Sprintf("Serie Folio: %s %s", tcktData.Serie, tcktData.Folio)
		if err := tc.printer.TextLn(serieFolio); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	if fecha := ""; tmpl.VerFecha {
		fecha = fmt.Sprintf("Fecha: %s", tcktData.FechaSistema)
		if err := tc.printer.TextLn(fecha); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	if sucTienda := ""; tmpl.VerTienda {
		sucTienda = fmt.Sprintf("Sucursal Tienda: %s\n", tcktData.SucursalTienda)
		if err := tc.printer.TextLn(sucTienda); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	if err := tc.printer.TextLn("Vendedor: " + tcktData.Vendedor); err != nil {
		log.Printf("ticket_printer: error al imprimir texto: %v", err)
	}
	if err := tc.printer.Feed(1); err != nil {
		log.Printf("ticket_printer: error al cortarpapel: %v", err)
	}
}

// printItems prints the purchased items
func (tc *TicketConstructor) printItems() map[string]float64 {
	tmpl := tc.template.Data

	cant := ""
	if tmpl.VerCantProductos {
		cant = tckt.PadCenter("CANT", LEN_CANT, ' ')
	}
	producto := tckt.PadCenter("PRODUCTO", LEN_DESC, ' ')
	precio := ""
	if tmpl.VerPrecioU {
		precio = tckt.PadCenter("PRECIO/U", LEN_PRECIO, ' ')
	}
	subtotal := tckt.PadCenter("SUBTOTAL", LEN_TOTAL, ' ')

	if err := tc.printer.TextLn(cant + producto + precio + subtotal); err != nil {
		log.Printf("Error al imprimir artículo 1: %v", err)
	}

	const IVA_TRAS int = 0
	const IEPS_TRAS int = 1
	const IVA_RET int = 2
	const ISR_RET = 3

	var subtotal_sum float64
	var ivaTrasladado_sum float64
	var iepsTrasladado_sum float64
	var ivaRetenido_sum float64
	var isrRetenido_sum float64

	// Print each conc
	conceptoRow := ""
	for _, conc := range tc.ticket.Data.Conceptos {
		cant = ""
		if tmpl.VerCantProductos {
			cant = tckt.PadCenter(strconv.FormatFloat(conc.Cantidad, 'f', -1, 64), LEN_CANT, ' ')
		}
		producto = tckt.PadCenter(tckt.Substr(conc.Descripcion, LEN_DESC-2), LEN_DESC, ' ')
		precio := ""
		if tmpl.VerPrecioU {
			precio = tckt.PadCenter("$"+tckt.FormatFloat(conc.PrecioVenta, LEN_DECIMALES), LEN_PRECIO, ' ')
		}
		subtotal = tckt.PadCenter("$"+tckt.FormatFloat(conc.Total, LEN_DECIMALES), LEN_TOTAL, ' ')

		subtotal_sum = subtotal_sum + conc.Total

		if len(conc.Impuestos) > 0 {
			ivaTrasladado_sum = ivaTrasladado_sum + conc.Impuestos[IVA_TRAS].Importe
		}
		if len(conc.Impuestos) > 1 {
			iepsTrasladado_sum = iepsTrasladado_sum + conc.Impuestos[IEPS_TRAS].Importe
			ivaRetenido_sum = ivaRetenido_sum + conc.Impuestos[IVA_RET].Importe
			isrRetenido_sum = isrRetenido_sum + conc.Impuestos[ISR_RET].Importe
		}

		conceptoRow = cant + producto + precio + subtotal
		if err := tc.printer.TextLn(conceptoRow); err != nil {
			log.Printf("Error al imprimir artículo 2: %v", err)
		}
	}

	// Configurar justificación y estilo
	if err := tc.printer.SetJustification(cons.JUSTIFY_RIGHT); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	// Tipo de fuente
	if err := tc.printer.SetFont(cons.FONT_B); err != nil {
		log.Printf("Error al establecer fuente: %v", err)
	}

	if err := tc.printer.TextLn("Subtotal: $" + tckt.FormatFloat(subtotal_sum, LEN_DECIMALES)); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}

	return map[string]float64{
		"ivaTrasladado":  ivaTrasladado_sum,
		"iepsTrasladado": iepsTrasladado_sum,
		"ivaRetenido":    ivaRetenido_sum,
		"isrRetenido":    isrRetenido_sum,
	}

}

// printTaxes prints tax information if configured
func (tc *TicketConstructor) printTaxes(taxes map[string]float64) {
	if tc.template.Data.VerImpuestos || tc.template.Data.VerImpuestosTotal || tc.template.Data.IncluyeImpuestos {
		if err := tc.printer.TextLn("IVA Trasladado: $" + tckt.FormatFloat(taxes["ivaTrasladado"], LEN_DECIMALES)); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.TextLn("IEPS Trasladado: $" + tckt.FormatFloat(taxes["iepsTrasladado"], LEN_DECIMALES)); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.TextLn("IVA Retenido: $" + tckt.FormatFloat(taxes["ivaRetenido"], LEN_DECIMALES)); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.TextLn("ISR Retenido: $" + tckt.FormatFloat(taxes["isrRetenido"], LEN_DECIMALES)); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
	}
}

// printPaymentInfo prints payment details
func (tc *TicketConstructor) printPaymentInfo() {
	tcktData := tc.ticket.Data.DocumentosPago[0]

	if err := tc.printer.TextLn("Total Field: $" + tckt.FormatFloat(tcktData.Total, LEN_DECIMALES)); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.TextLn("Efectivo: $" + tckt.FormatFloat(tcktData.FormasPago[0].Cantidad, LEN_DECIMALES)); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.TextLn("Cambio: $" + tckt.FormatFloat(tcktData.Cambio, LEN_DECIMALES) + "\n"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
}

// printFooter prints the footer information
func (tc *TicketConstructor) printFooter() {
	tmpl := tc.template.Data

	if tmpl.VerLeyenda && tmpl.CambiarReclamacion != "" {
		if err := tc.printer.TextLn(tmpl.CambiarReclamacion); err != nil {
			log.Printf("Error al imprimir: %v", err)
		}
		if err := tc.printer.Feed(1); err != nil {
			log.Printf("Error al imprimir: %v", err)
		}
	}

	// Print footer message
	if err := tc.printer.TextLn(tmpl.CambiarPie); err != nil {
		log.Printf("Error al imprimir: %v", err)
	}
}
