package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/skip2/go-qrcode"
	tckt "pos-daemon.adcon.dev/internal/models"
	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/image"
	"pos-daemon.adcon.dev/pkg/posprinter/types"
)

// TODO: Hacer reconfigurable los tamaños de acuerdo al ticket
const (
	LenCant      int = 4
	LenDesc      int = 18
	LenPrecio    int = 9
	LenTotal     int = 9
	LenDecimales int = 2
	MaxRowChars  int = 40
)

// TicketConstructor handles the construction and printing of tickets
type TicketConstructor struct {
	template tckt.NewTicketTemplate
	ticket   tckt.NewTicket
	writer   io.Writer
	printer  *posprinter.GenericPrinter
}

// NewTicketConstructor creates a new ticket constructor with the specified writer
func NewTicketConstructor(writer io.Writer, printer *posprinter.GenericPrinter) *TicketConstructor {
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
	if err := tc.printer.SetJustification(types.AlignCenter); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}

	// Tipo de fuente
	if err := tc.printer.SetFont(types.FontA); err != nil {
		log.Printf("Error al establecer fuente: %v", err)
	}

	// Usar perfil para establecer ancho y márgenes
	profile := tc.printer.GetProfile()
	// Actualizar perfil con el ancho del ticket si es necesario
	if int(tc.template.Data.TicketWidth) != int(profile.PaperWidth) {
		profile.PaperWidth = float64(tc.template.Data.TicketWidth)
		profile.DotsPerLine = int(float64(profile.DPI) * profile.PaperWidth / 25.4) // convertir mm a puntos
		tc.printer.SetProfile(profile)
	}

	tc.printHeader()
	tc.printCustomerInfo()
	tc.printTicketInfo()
	if err := tc.printer.Feed(1); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}
	taxes := tc.printItems()
	tc.printTaxes(taxes)
	tc.printPaymentInfo()
	if err := tc.printer.Feed(1); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}
	tc.printQr()

	tc.printFooter()

	// Alimentar papel al final
	if err := tc.printer.Feed(2); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}

	// Cortar papel
	if err := tc.printer.Cut(types.CutFeed, 3); err != nil {
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
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(tmpl.CambiarCabecera); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
	}

	// Print logo placeholder if configured
	if tmpl.VerLogotipo {
		logoPath := "./img/perro.jpeg"
		img, err := image.LoadImage(logoPath)
		if err != nil {
			log.Printf("ticket_printer: error cargando imagen del logo: %v", err)
		} else {
			if err := tc.printer.Feed(1); err != nil {
				log.Printf("ticket_printer: error al alimentar papel después de imprimir cabecera: %v", err)
			}

			// Imprimir la imagen con dithering usando las nuevas opciones
			opts := posprinter.PrintImageOptions{
				Density:    types.DensitySingle,
				DitherMode: image.DitherFloydSteinberg,
				Threshold:  128,
				Width:      tmpl.LogoWidth * 2,
			}
			if err := tc.printer.PrintImageWithOptions(img, opts); err != nil {
				log.Printf("ticket_printer: error al imprimir logo con dithering: %v", err)
			}

			if err := tc.printer.Feed(1); err != nil {
				log.Printf("ticket_printer: error al alimentar papel después de imprimir cabecera: %v", err)
			}
		}
	}

	// Print store name
	if tmpl.VerNombre && datosTckt.SucursalNombre != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn("Matriz\n" + datosTckt.SucursalNombre); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Feed(1); err != nil {
			log.Printf("ticket_printer: error al alimentar papel después de imprimir logo: %v", err)
		}
	}

	// Print commercial name
	if tmpl.VerNombreC && datosTckt.SucursalNombreComercial != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("Nombre Comercial: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}

		// Si es mas grande fijar al doble de tamaño - en GenericPrinter esto tendría que implementarse
		// a través del protocolo, omitimos por ahora
		if tmpl.RazonSocialSize > 10 {
			if err := tc.printer.SetFont(types.FontB); err != nil {
				log.Printf("Error al establecer fuente: %v", err)
			}
			// GenericPrinter no tiene SetTextSize directamente - podría ser necesario agregar
			// esa funcionalidad a través del protocolo
		}

		if err := tc.printer.TextLn(datosTckt.SucursalNombreComercial); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}

		// Restaurar el tamaño de texto si se cambió
		if tmpl.RazonSocialSize > 10 {
			if err := tc.printer.SetFont(types.FontA); err != nil {
				log.Printf("Error al establecer fuente: %v", err)
			}
		}

		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
	}

	// Print RFC
	if tmpl.VerRFC && datosTckt.SucursalRFC != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("RFC: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(datosTckt.SucursalRFC); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print tax regime
	if tmpl.VerRegimen && datosTckt.SucursalRegimen != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("Régimen Fiscal: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(datosTckt.SucursalRegimen); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print email
	if tmpl.VerEmail && datosTckt.SucursalEmails != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("Email: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(datosTckt.SucursalEmails); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	// Print address
	if dom := ""; tmpl.VerDom && datosTckt.SucursalCalle != "" && datosTckt.SucursalNumero != "" && datosTckt.SucursalColonia != "" {
		dom = fmt.Sprintf("%s %s,", datosTckt.SucursalCalle, datosTckt.SucursalNumero)
		if datosTckt.SucursalNumeroInt != "" {
			dom = dom + fmt.Sprintf(" Int. %s,", datosTckt.SucursalNumeroInt)
		}
		dom = dom + fmt.Sprintf(" Col. %s,", datosTckt.SucursalColonia)
		dom = dom + fmt.Sprintf(" %s, %s, %s, ", datosTckt.SucursalLocalidad, datosTckt.SucursalEstado, datosTckt.SucursalPais)
		dom = dom + fmt.Sprintf(" C.P. %s", datosTckt.SucursalCP)
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("Domicilio: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(dom); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}
}

// printCustomerInfo prints the customer information
func (tc *TicketConstructor) printCustomerInfo() {
	if tc.template.Data.VerNombreCliente && tc.ticket.Data.ClienteNombre != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("Cliente: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(tc.ticket.Data.ClienteNombre); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}
}

// printTicketInfo prints folio, date, and store info
func (tc *TicketConstructor) printTicketInfo() {
	tmpl := tc.template.Data
	tcktData := tc.ticket.Data

	if tmpl.VerFolio && tcktData.Folio != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("Folio: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(tcktData.Folio); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	if tmpl.VerFecha && tcktData.FechaSistema != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("Fecha: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(tcktData.FechaSistema); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}

	if tmpl.VerTienda && tcktData.SucursalTienda != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.Text("Tienda: "); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(tcktData.SucursalTienda); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
	}
}

// printItems prints the purchased items
func (tc *TicketConstructor) printItems() map[string]float64 {
	const IvaTras int = 0
	const IepsTras int = 1
	const IvaRet int = 2
	const IsrRet = 3

	var subtotalSum float64
	var ivatrasladadoSum float64
	var iepstrasladadoSum float64
	var ivaretenidoSum float64
	var isrretenidoSum float64

	tmpl := tc.template.Data

	precioCol := ""
	productoCol := PadCenter("PRODUCTO", LenDesc+LenPrecio, ' ')
	if tmpl.VerPrecioU {
		precioCol = PadCenter("PRECIO/U", LenPrecio, ' ')
		productoCol = PadCenter("PRODUCTO", LenDesc, ' ')
	}
	cantCol := ""
	subtotalCol := PadLeft("SUBTOTAL", LenTotal+LenCant, ' ')
	if tmpl.VerCantProductos {
		cantCol = PadCenter("CANT", LenCant, ' ')
		subtotalCol = PadLeft("SUBTOTAL", LenTotal, ' ')
	}

	// Configurar justificación y estilo
	if err := tc.printer.SetJustification(types.AlignRight); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	columnas := cantCol + productoCol + precioCol + subtotalCol
	if err := tc.printer.TextLn(columnas); err != nil {
		log.Printf("Error al imprimir artículo 1: %v", err)
	}
	if err := tc.printer.SetEmphasis(false); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	if len(columnas) != MaxRowChars {
		log.Printf("Advertencia: la fila del concepto excede o es menor al máximo de caracteres: %d / %d): %s", len(columnas), MaxRowChars, "|"+columnas+"|")
	}

	// Print each concept
	conceptoRow := ""
	for _, conc := range tc.ticket.Data.Conceptos {
		cant := ""
		subtotal := PadLeft("$"+FormatFloat(conc.Total, LenDecimales), LenTotal+LenCant, ' ')
		if tmpl.VerCantProductos {
			cant = PadCenter(strconv.FormatFloat(conc.Cantidad, 'f', -1, 64), LenCant, ' ')
			subtotal = PadLeft("$"+FormatFloat(conc.Total, LenDecimales), LenTotal, ' ')
		}
		precio := ""
		seriesStr := ""
		if tmpl.VerSeries && len(conc.Series) > 0 {
			seriesStr = ", " + strings.Join(conc.Series, ", ")
		}
		productos := SplitString(conc.Descripcion+", "+seriesStr, LenDesc+LenPrecio-2)
		productos[0] = PadCenter(productos[0], LenDesc+LenPrecio, ' ')
		if tmpl.VerPrecioU {
			precio = PadCenter("$"+FormatFloat(conc.PrecioVenta, LenDecimales), LenPrecio, ' ')
			productos = SplitString(conc.Descripcion+seriesStr, LenDesc-2)
			productos[0] = PadCenter(productos[0], LenDesc, ' ')
		}

		conceptoRow = cant + productos[0] + precio + subtotal
		if len(conceptoRow) != MaxRowChars {
			log.Printf("Advertencia: la fila del concepto excede o es menor al máximo de caracteres: %d / %d): %s", len(conceptoRow), MaxRowChars, "|"+conceptoRow+"|")
		}
		if err := tc.printer.TextLn(conceptoRow); err != nil {
			log.Printf("Error al imprimir fila 1 de artículo 1: %v", err)
		}

		if len(productos) > 1 {
			for _, prod := range productos[1:] {
				cant = PadCenter("", LenCant, ' ')
				producto := PadCenter(prod, LenDesc+LenPrecio, ' ')
				if tmpl.VerPrecioU {
					producto = PadCenter(prod, LenDesc, ' ')
				}
				precio = PadCenter("", LenPrecio, ' ')
				subtotal = PadLeft("", LenTotal, ' ')

				conceptoRow = cant + producto + precio + subtotal
				if len(conceptoRow) != MaxRowChars {
					log.Printf("Advertencia: la fila del concepto excede o es menor al máximo de caracteres: %d / %d): %s", len(conceptoRow), MaxRowChars, "|"+conceptoRow+"|")
				}
				if err := tc.printer.TextLn(conceptoRow); err != nil {
					log.Printf("Error al imprimir artículo 2: %v", err)
				}
			}
		}

		subtotalSum = subtotalSum + conc.Total

		if len(conc.Impuestos) > 0 {
			ivatrasladadoSum = ivatrasladadoSum + conc.Impuestos[IvaTras].Importe
		}
		if len(conc.Impuestos) > 1 {
			iepstrasladadoSum = iepstrasladadoSum + conc.Impuestos[IepsTras].Importe
			ivaretenidoSum = ivaretenidoSum + conc.Impuestos[IvaRet].Importe
			isrretenidoSum = isrretenidoSum + conc.Impuestos[IsrRet].Importe
		}
	}

	// Configurar justificación y estilo
	if err := tc.printer.SetJustification(types.AlignRight); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}

	if err := tc.printer.Text("Subtotal: $"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	if err := tc.printer.TextLn(FormatFloat(subtotalSum, LenDecimales)); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.SetEmphasis(false); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}

	return map[string]float64{
		"ivaTrasladado":  ivatrasladadoSum,
		"iepsTrasladado": iepstrasladadoSum,
		"ivaRetenido":    ivaretenidoSum,
		"isrRetenido":    isrretenidoSum,
	}
}

// printTaxes prints tax information if configured
func (tc *TicketConstructor) printTaxes(taxes map[string]float64) {
	if (tc.template.Data.VerImpuestos || tc.template.Data.VerImpuestosTotal) && tc.template.Data.IncluyeImpuestos {
		if err := tc.printer.Text("IVA Trasladado: $"); err != nil {
			log.Printf("Error al imprimir: %v", err)
		}
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(FormatFloat(taxes["ivaTrasladado"], LenDecimales)); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}

		if err := tc.printer.Text("IVA Retenido: $"); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(FormatFloat(taxes["ivaRetenido"], LenDecimales)); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}

		if err := tc.printer.Text("IEPS Trasladado: $"); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(FormatFloat(taxes["iepsTrasladado"], LenDecimales)); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}

		if err := tc.printer.Text("ISR Retenido: $"); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn(FormatFloat(taxes["isrRetenido"], LenDecimales)); err != nil {
			log.Printf("Error al imprimir suma: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
	}
}

// printPaymentInfo prints payment details
func (tc *TicketConstructor) printPaymentInfo() {
	tcktData := tc.ticket.Data.DocumentosPago[0]

	if err := tc.printer.Text("Total: $"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	if err := tc.printer.TextLn(FormatFloat(tcktData.Total, LenDecimales)); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.SetEmphasis(false); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}

	if err := tc.printer.Text("Efectivo: $"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	if err := tc.printer.TextLn(FormatFloat(tcktData.FormasPago[0].Cantidad, LenDecimales)); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.SetEmphasis(false); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}

	if err := tc.printer.Text("Cambio: $"); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	if err := tc.printer.TextLn(FormatFloat(tcktData.Cambio, LenDecimales)); err != nil {
		log.Printf("Error al imprimir suma: %v", err)
	}
	if err := tc.printer.SetEmphasis(false); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
}

// printFooter prints the footer information
func (tc *TicketConstructor) printFooter() {
	tmpl := tc.template.Data

	if err := tc.printer.SetJustification(types.AlignCenter); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}

	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}

	// Texto grande para "PAGADO"
	// Nota: GenericPrinter no tiene SetTextSize, usamos énfasis y fuentes para simular
	if err := tc.printer.SetFont(types.FontB); err != nil { // Fuente más grande
		log.Printf("Error al establecer fuente: %v", err)
	}
	if err := tc.printer.TextLn("PAGADO"); err != nil {
		log.Printf("Error al imprimir: %v", err)
	}
	if err := tc.printer.SetFont(types.FontA); err != nil { // Restaurar fuente normal
		log.Printf("Error al establecer fuente: %v", err)
	}

	if err := tc.printer.SetEmphasis(false); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}

	if err := tc.printer.Feed(1); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}

	cantSum := 0.0
	for _, cant := range tc.ticket.Data.Conceptos {
		cantSum += cant.Cantidad
	}
	if err := tc.printer.TextLn(fmt.Sprintf("Cantidad de Productos: %s", strconv.FormatFloat(cantSum, 'f', -1, 64))); err != nil {
		log.Printf("Error al imprimir: %v", err)
	}

	if tmpl.VerLeyenda && tmpl.CambiarReclamacion != "" {
		if err := tc.printer.TextLn(tmpl.CambiarReclamacion); err != nil {
			log.Printf("Error al imprimir: %v", err)
		}
	}

	// Print phone if configured
	if tmpl.VerTelefono && tc.ticket.Data.SucursalTelefono != "" {
		if err := tc.printer.SetEmphasis(true); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
		if err := tc.printer.TextLn("Teléfono: " + tc.ticket.Data.SucursalTelefono); err != nil {
			log.Printf("ticket_printer: error al imprimir texto: %v", err)
		}
		if err := tc.printer.SetEmphasis(false); err != nil {
			log.Printf("Error al establecer énfasis: %v", err)
		}
	}
	// Print footer message
	if err := tc.printer.SetEmphasis(true); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
	if err := tc.printer.TextLn(tmpl.CambiarPie); err != nil {
		log.Printf("Error al imprimir: %v", err)
	}
	if err := tc.printer.SetEmphasis(false); err != nil {
		log.Printf("Error al establecer énfasis: %v", err)
	}
}

func (tc *TicketConstructor) printQr() {
	if err := tc.printer.SetJustification(types.AlignCenter); err != nil {
		log.Printf("Error al establecer justificación: %v", err)
	}
	if err := tc.printer.TextLn(tc.ticket.Data.AutofacturaLink); err != nil {
		log.Printf("Error al imprimir: %v", err)
	}

	// Generar el código QR
	qr, err := qrcode.New(tc.ticket.Data.AutofacturaLinkQr, qrcode.Medium)
	if err != nil {
		log.Printf("Error generando QR: %v", err)
		return
	}

	// Obtener la imagen del QR
	var size = 256
	qrImage := qr.Image(size)

	// Imprimir usando la nueva API de imágenes
	if err = tc.printer.PrintImage(qrImage); err != nil {
		log.Printf("Error al imprimir QR: %v", err)
	}
}
