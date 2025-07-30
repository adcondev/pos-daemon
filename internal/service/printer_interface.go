package service

import (
	"fmt"
	"image"

	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/imaging"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
)

// PrinterInterface define la interfaz que el ticket builder necesita
type PrinterInterface interface {
	// Métodos de texto
	Text(str string) error
	TextLn(str string) error

	// Métodos de formato
	SetJustification(justification escpos.Justify) error
	SetEmphasis(on bool) error
	SetUnderline(underline escpos.UnderlineMode) error
	SetTextSize(widthMultiplier, heightMultiplier int) error
	SetFont(font escpos.Font) error

	// Métodos de control
	Feed(lines int) error
	Cut(mode int, lines int) error
	Initialize() error

	// Métodos de imagen
	PrintImage(img image.Image, density command.Density) error
	PrintImageWithDithering(img image.Image, density command.Density, dither imaging.DitherMode) error

	// Métodos de código de barras y QR
	SetBarcodeHeight(height int) error
	SetBarcodeWidth(width int) error
	SetBarcodeTextPosition(position escpos.BarcodeTextPos) error
	Barcode(content string, barType escpos.BarcodeType) error
}

// PrinterAdapter adapta la nueva arquitectura para el ticket builder
type PrinterAdapter struct {
	printer       posprinter.Printer // Cambiar a interfaz
	escposAdapter *escpos.ESCPrinterAdapter
}

// NewPrinterAdapter crea un nuevo adaptador
func NewPrinterAdapter(printer posprinter.Printer, escposAdapter *escpos.ESCPrinterAdapter) *PrinterAdapter {
	return &PrinterAdapter{
		printer:       printer,
		escposAdapter: escposAdapter,
	}
}

// === Implementación de métodos ===

func (p *PrinterAdapter) Text(str string) error {
	return p.printer.Text(str)
}

func (p *PrinterAdapter) TextLn(str string) error {
	return p.printer.TextLn(str)
}

func (p *PrinterAdapter) SetJustification(justification escpos.Justify) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.SetJustification(justification)
	}
	// Convertir y usar printer directamente
	var alignment command.Alignment
	switch justification {
	case escpos.Left:
		alignment = command.AlignLeft
	case escpos.Center:
		alignment = command.AlignCenter
	case escpos.Right:
		alignment = command.AlignRight
	}
	return p.printer.SetJustification(alignment)
}

func (p *PrinterAdapter) SetEmphasis(on bool) error {
	return p.printer.SetEmphasis(on)
}

func (p *PrinterAdapter) SetUnderline(underline escpos.UnderlineMode) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.SetUnderline(underline)
	}
	// Convertir y usar printer directamente
	var newUnderline command.UnderlineMode
	switch underline {
	case escpos.NoUnderline:
		newUnderline = command.UnderlineNone
	case escpos.Single:
		newUnderline = command.UnderlineSingle
	case escpos.Double:
		newUnderline = command.UnderlineDouble
	}
	return p.printer.SetUnderline(newUnderline)
}

func (p *PrinterAdapter) SetTextSize(widthMultiplier, heightMultiplier int) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.SetTextSize(widthMultiplier, heightMultiplier)
	}
	// Por ahora, ignorar si no hay adaptador
	return nil
}

func (p *PrinterAdapter) SetFont(font escpos.Font) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.SetFont(font)
	}
	// Convertir y usar printer directamente
	var newFont command.Font
	switch font {
	case escpos.A:
		newFont = command.FontA
	case escpos.B:
		newFont = command.FontB
	case escpos.C:
		newFont = command.FontC
	}
	return p.printer.SetFont(newFont)
}

func (p *PrinterAdapter) Feed(lines int) error {
	return p.printer.Feed(lines)
}

func (p *PrinterAdapter) Cut(mode int, lines int) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.Cut(mode, lines)
	}
	// Usar printer directamente
	var cutMode command.CutMode
	if mode == escpos.CUT_PARTIAL {
		cutMode = command.CutPartial
	} else {
		cutMode = command.CutFull
	}
	if lines > 0 {
		p.printer.Feed(lines)
	}
	return p.printer.Cut(cutMode)
}

func (p *PrinterAdapter) Initialize() error {
	return p.printer.Initialize()
}

func (p *PrinterAdapter) PrintImage(img image.Image, density command.Density) error {
	return p.printer.PrintImage(img, density)
}

func (p *PrinterAdapter) PrintImageWithDithering(img image.Image, density command.Density, dither imaging.DitherMode) error {
	// Crear opciones para impresión con dithering
	opts := posprinter.PrintImageOptions{
		Density:    density,
		DitherMode: dither,
		Threshold:  128,
	}

	// Usar PrintImageWithOptions si está disponible
	if genPrinter, ok := p.printer.(*posprinter.GenericPrinter); ok {
		return genPrinter.PrintImageWithOptions(img, opts)
	}

	// Fallback: imprimir sin dithering
	return p.printer.PrintImage(img, density)
}

// Métodos de código de barras
func (p *PrinterAdapter) SetBarcodeHeight(height int) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.SetBarcodeHeight(height)
	}
	return fmt.Errorf("barcode not supported without ESC/POS adapter")
}

func (p *PrinterAdapter) SetBarcodeWidth(width int) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.SetBarcodeWidth(width)
	}
	return fmt.Errorf("barcode not supported without ESC/POS adapter")
}

func (p *PrinterAdapter) SetBarcodeTextPosition(position escpos.BarcodeTextPos) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.SetBarcodeTextPosition(position)
	}
	return fmt.Errorf("barcode not supported without ESC/POS adapter")
}

func (p *PrinterAdapter) Barcode(content string, barType escpos.BarcodeType) error {
	if p.escposAdapter != nil {
		return p.escposAdapter.Barcode(content, barType)
	}
	return fmt.Errorf("barcode not supported without ESC/POS adapter")
}
