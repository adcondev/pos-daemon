package service

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"

	"pos-daemon.adcon.dev/internal/models"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/imaging"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
)

// TicketConstructor maneja la construcción e impresión de tickets
type TicketConstructor struct {
	output   io.Writer
	printer  PrinterInterface // Usa la interfaz definida en este paquete
	template *models.TicketTemplateData
	ticket   *models.TicketData
}

// NewTicketConstructor crea un nuevo constructor de tickets
func NewTicketConstructor(output io.Writer, printer PrinterInterface) *TicketConstructor {
	return &TicketConstructor{
		output:  output,
		printer: printer,
	}
}

// LoadTemplateFromJSON carga la plantilla desde JSON
func (tc *TicketConstructor) LoadTemplateFromJSON(data []byte) error {
	tc.template = &models.TicketTemplateData{}
	return json.Unmarshal(data, tc.template)
}

// LoadTicketFromJSON carga los datos del ticket desde JSON
func (tc *TicketConstructor) LoadTicketFromJSON(data []byte) error {
	tc.ticket = &models.TicketData{}
	return json.Unmarshal(data, tc.ticket)
}

// PrintTicket imprime el ticket completo
func (tc *TicketConstructor) PrintTicket() error {
	if tc.template == nil || tc.ticket == nil {
		return fmt.Errorf("template or ticket data not loaded")
	}

	// Inicializar impresora
	if err := tc.printer.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize printer: %w", err)
	}

	// Imprimir logo si está definido
	if err := tc.printLogo(); err != nil {
		log.Printf("Warning: failed to print logo: %v", err)
		// No fallar el ticket completo por el logo
	}

	// Imprimir encabezado
	if err := tc.printHeader(); err != nil {
		return err
	}

	// Imprimir detalles
	if err := tc.printDetails(); err != nil {
		return err
	}

	// Imprimir footer
	if err := tc.printFooter(); err != nil {
		return err
	}

	// Cortar papel
	return tc.printer.Cut(escpos.CutFeed, 3)
}

// printLogo imprime el logo si existe
func (tc *TicketConstructor) printLogo() error {
	// Verificar si el template tiene configuración de logo
	// Por ahora, usar una ruta hardcodeada o de configuración
	// TODO: Agregar campo Logo a TicketTemplateData si no existe
	logoPath := "./assets/logo.png" // Ruta por defecto

	// Verificar si el archivo existe
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		return nil // No hay logo, continuar sin error
	}

	// Cargar imagen
	file, err := os.Open(logoPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing logo file: %v", err)
		}
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Centrar imagen
	if err := tc.printer.SetJustification(escpos.Center); err != nil {
		return err
	}

	// Imprimir imagen con dithering
	if err := tc.printer.PrintImageWithDithering(img, command.DensitySingle, imaging.DitherFloydSteinberg); err != nil {
		return err
	}

	// Restaurar alineación
	return tc.printer.SetJustification(escpos.Left)
}

// printHeader imprime el encabezado del ticket
func (tc *TicketConstructor) printHeader() error {
	// Centrar texto
	if err := tc.printer.SetJustification(escpos.Center); err != nil {
		return err
	}

	// Texto en negrita
	if err := tc.printer.SetEmphasis(true); err != nil {
		return err
	}

	// Imprimir nombre de la tienda
	if err := tc.printer.TextLn(tc.ticket.SucursalNombreComercial); err != nil {
		return err
	}

	// Desactivar negrita
	if err := tc.printer.SetEmphasis(false); err != nil {
		return err
	}

	// Imprimir dirección
	if err := tc.printer.TextLn(tc.ticket.SucursalCalle + " " + tc.ticket.SucursalNumero); err != nil {
		return err
	}

	if err := tc.printer.TextLn(tc.ticket.SucursalColonia + ", " + tc.ticket.SucursalMunicipio); err != nil {
		return err
	}

	if err := tc.printer.TextLn(tc.ticket.SucursalEstado + " CP: " + tc.ticket.SucursalCP); err != nil {
		return err
	}

	// RFC
	if err := tc.printer.TextLn("RFC: " + tc.ticket.SucursalRFC); err != nil {
		return err
	}

	// Línea separadora
	if err := tc.printer.Feed(1); err != nil {
		return err
	}

	if err := tc.printer.TextLn(strings.Repeat("-", 42)); err != nil {
		return err
	}

	// Restaurar alineación izquierda
	return tc.printer.SetJustification(escpos.Left)
}

// printDetails imprime los detalles del ticket
func (tc *TicketConstructor) printDetails() error {
	// Información del ticket
	if err := tc.printer.TextLn("Folio: " + tc.ticket.Folio); err != nil {
		return err
	}

	if err := tc.printer.TextLn("Fecha: " + tc.ticket.FechaSistema); err != nil {
		return err
	}

	if err := tc.printer.TextLn("Vendedor: " + tc.ticket.Vendedor); err != nil {
		return err
	}

	// Cliente (si existe)
	if tc.ticket.ClienteNombre != "" && tc.ticket.ClienteNombre != "PUBLICO EN GENERAL" {
		if err := tc.printer.Feed(1); err != nil {
			return err
		}
		if err := tc.printer.TextLn("Cliente: " + tc.ticket.ClienteNombre); err != nil {
			return err
		}
		if tc.ticket.ClienteRFC != "" {
			if err := tc.printer.TextLn("RFC: " + tc.ticket.ClienteRFC); err != nil {
				return err
			}
		}
	}

	// Línea separadora
	if err := tc.printer.Feed(1); err != nil {
		return err
	}
	if err := tc.printer.TextLn(strings.Repeat("-", 42)); err != nil {
		return err
	}

	// Encabezados de conceptos
	if err := tc.printer.TextLn("CANT  DESCRIPCION              TOTAL"); err != nil {
		return err
	}
	if err := tc.printer.TextLn(strings.Repeat("-", 42)); err != nil {
		return err
	}

	// Imprimir conceptos
	for _, concepto := range tc.ticket.Conceptos {
		// Formatear línea del concepto
		line := fmt.Sprintf("%-5.0f %-23s %7.2f",
			concepto.Cantidad,
			truncateString(concepto.Descripcion, 23),
			concepto.Total,
		)
		if err := tc.printer.TextLn(line); err != nil {
			return err
		}

		// Si la descripción es muy larga, imprimir en segunda línea
		if len(concepto.Descripcion) > 23 {
			if err := tc.printer.TextLn("      " + concepto.Descripcion[23:]); err != nil {
				return err
			}
		}

		// Precio unitario si es más de una pieza
		if concepto.Cantidad > 1 {
			unitPrice := fmt.Sprintf("      @ $%.2f c/u", concepto.PrecioVenta)
			if err := tc.printer.TextLn(unitPrice); err != nil {
				return err
			}
		}
	}

	return nil
}

// printFooter imprime el pie del ticket
func (tc *TicketConstructor) printFooter() error {
	// Línea separadora
	if err := tc.printer.TextLn(strings.Repeat("-", 42)); err != nil {
		return err
	}

	// Totales alineados a la derecha
	if err := tc.printer.SetJustification(escpos.Right); err != nil {
		return err
	}

	// Descuento (si existe)
	if tc.ticket.Descuento > 0 {
		desc := fmt.Sprintf("Descuento: $%.2f", tc.ticket.Descuento)
		if err := tc.printer.TextLn(desc); err != nil {
			return err
		}
	}

	// Total en negrita
	if err := tc.printer.SetEmphasis(true); err != nil {
		return err
	}
	total := fmt.Sprintf("TOTAL: $%.2f", tc.ticket.Total)
	if err := tc.printer.TextLn(total); err != nil {
		return err
	}
	if err := tc.printer.SetEmphasis(false); err != nil {
		return err
	}

	// Pago y cambio
	if tc.ticket.Pagado > 0 {
		pago := fmt.Sprintf("Pagado: $%.2f", tc.ticket.Pagado)
		if err := tc.printer.TextLn(pago); err != nil {
			return err
		}

		if tc.ticket.Cambio > 0 {
			cambio := fmt.Sprintf("Cambio: $%.2f", tc.ticket.Cambio)
			if err := tc.printer.TextLn(cambio); err != nil {
				return err
			}
		}
	}

	// QR de autofactura (si existe)
	if tc.ticket.AutofacturaLink != "" {
		if err := tc.printer.Feed(2); err != nil {
			return err
		}

		// Centrar
		if err := tc.printer.SetJustification(escpos.Center); err != nil {
			return err
		}

		if err := tc.printer.TextLn("AUTOFACTURA"); err != nil {
			return err
		}

		// Generar y imprimir QR
		qr, err := qrcode.New(tc.ticket.AutofacturaLink, qrcode.Medium)
		if err != nil {
			log.Printf("Error generating QR: %v", err)
		} else {
			qrImage := qr.Image(256)
			if err := tc.printer.PrintImage(qrImage); err != nil {
				log.Printf("Error printing QR: %v", err)
			}
		}
	}

	// Mensaje final
	if err := tc.printer.Feed(2); err != nil {
		return err
	}
	if err := tc.printer.SetJustification(escpos.Center); err != nil {
		return err
	}
	if err := tc.printer.TextLn("¡GRACIAS POR SU COMPRA!"); err != nil {
		return err
	}

	return tc.printer.Feed(2)
}

// truncateString trunca una cadena a la longitud especificada
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
