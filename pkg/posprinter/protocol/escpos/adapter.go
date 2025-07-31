package escpos

import (
	"fmt"

	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
)

// ESCPrinterAdapter adapta la nueva arquitectura para mantener compatibilidad
type ESCPrinterAdapter struct {
	printer        *posprinter.GenericPrinter
	Connector      connector.Connector
	CharacterTable int
	// Agregar campos para mantener estado de barcode
	barcodeHeight       int
	barcodeWidth        int
	barcodeTextPosition BarcodeTextPos
}

// NewPrinter crea un adaptador compatible con la API anterior
func NewPrinter(conn connector.Connector, profile *CapabilityProfile) (*ESCPrinterAdapter, error) {
	// Crear protocolo ESC/POS
	proto := NewESCPOSProtocol()

	// Crear impresora genérica
	genericPrinter, err := posprinter.NewGenericPrinter(proto, conn)
	if err != nil {
		return nil, err
	}

	// Crear adaptador
	adapter := &ESCPrinterAdapter{
		printer:             genericPrinter,
		Connector:           conn,
		CharacterTable:      0,
		barcodeHeight:       50,        // Valor por defecto
		barcodeWidth:        2,         // Valor por defecto
		barcodeTextPosition: TextBelow, // Valor por defecto
	}

	return adapter, nil
}

// Initialize mantiene compatibilidad con código antiguo
func (p *ESCPrinterAdapter) Initialize() error {
	return p.printer.Initialize()
}

// SetJustification con tipos antiguos
func (p *ESCPrinterAdapter) SetJustification(justification Justify) error {
	// Convertir tipo antiguo a nuevo
	var alignment command.Alignment
	switch justification {
	case Left:
		alignment = command.AlignLeft
	case Center:
		alignment = command.AlignCenter
	case Right:
		alignment = command.AlignRight
	default:
		alignment = command.AlignLeft
	}

	return p.printer.SetJustification(alignment)
}

// SetFont con tipos antiguos
func (p *ESCPrinterAdapter) SetFont(font Font) error {
	// Convertir tipo antiguo a nuevo
	var newFont command.Font
	switch font {
	case A:
		newFont = command.FontA
	case B:
		newFont = command.FontB
	case C:
		newFont = command.FontC
	default:
		newFont = command.FontA
	}

	return p.printer.SetFont(newFont)
}

// Text mantiene la funcionalidad anterior
func (p *ESCPrinterAdapter) Text(str string) error {
	return p.printer.Text(str)
}

// TextLn mantiene compatibilidad
func (p *ESCPrinterAdapter) TextLn(str string) error {
	return p.printer.TextLn(str)
}

// Cut con la API antigua
func (p *ESCPrinterAdapter) Cut(mode int, lines int) error {
	// Ignorar lines por ahora, usar mode para determinar tipo de corte
	var cutMode command.CutMode
	if mode == CUT_PARTIAL {
		cutMode = command.CutPartial
	} else {
		cutMode = command.CutFeed
	}

	// Si lines > 0, alimentar papel antes del corte
	if lines > 0 {
		if err := p.printer.Feed(lines); err != nil {
			return err
		}
	}

	return p.printer.Cut(cutMode, 0)
}

// Feed mantiene compatibilidad
func (p *ESCPrinterAdapter) Feed(lines int) error {
	return p.printer.Feed(lines)
}

// Close cierra la impresora
func (p *ESCPrinterAdapter) Close() error {
	return p.printer.Close()
}

// SetEmphasis
func (p *ESCPrinterAdapter) SetEmphasis(on bool) error {
	return p.printer.SetEmphasis(on)
}

// SetDoubleStrike
func (p *ESCPrinterAdapter) SetDoubleStrike(on bool) error {
	return p.printer.SetDoubleStrike(on)
}

// SetUnderline con tipos antiguos
func (p *ESCPrinterAdapter) SetUnderline(underline UnderlineMode) error {
	// Convertir tipo antiguo a nuevo
	var newUnderline command.UnderlineMode
	switch underline {
	case NoUnderline:
		newUnderline = command.UnderlineNone
	case Single:
		newUnderline = command.UnderlineSingle
	case Double:
		newUnderline = command.UnderlineDouble
	default:
		newUnderline = command.UnderlineNone
	}

	return p.printer.SetUnderline(newUnderline)
}

// SetTextSize - Por ahora solo almacena los valores
// TODO: Implementar cuando el protocolo genérico lo soporte
func (p *ESCPrinterAdapter) SetTextSize(widthMultiplier, heightMultiplier int) error {
	// Por ahora, solo retornar nil para mantener compatibilidad
	// TODO: Implementar cuando el protocolo lo soporte
	return nil
}

// SetBarcodeHeight almacena la altura para uso posterior
func (p *ESCPrinterAdapter) SetBarcodeHeight(height int) error {
	p.barcodeHeight = height
	return nil
}

// SetBarcodeWidth almacena el ancho para uso posterior
func (p *ESCPrinterAdapter) SetBarcodeWidth(width int) error {
	p.barcodeWidth = width
	return nil
}

// SetBarcodeTextPosition almacena la posición del texto
func (p *ESCPrinterAdapter) SetBarcodeTextPosition(position BarcodeTextPos) error {
	p.barcodeTextPosition = position
	return nil
}

// Barcode imprime un código de barras
func (p *ESCPrinterAdapter) Barcode(content string, barType BarcodeType) error {
	// TODO: Implementar cuando el protocolo genérico lo soporte
	// Por ahora, retornar error indicando que no está implementado
	return fmt.Errorf("barcode printing not yet implemented in adapter")
}
