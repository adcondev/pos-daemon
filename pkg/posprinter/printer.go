package posprinter

import (
	"fmt"
	"image"
	"log"

	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/imaging"
	"pos-daemon.adcon.dev/pkg/posprinter/profile"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol"
	"pos-daemon.adcon.dev/pkg/posprinter/utils"
)

// Printer define la interfaz de alto nivel para cualquier impresora
type Printer interface {
	// === Comandos básicos ===
	Initialize() error
	Close() error

	// === Formato de texto ===
	SetJustification(alignment command.Alignment) error
	SetFont(font command.Font) error
	SetEmphasis(on bool) error
	SetDoubleStrike(on bool) error
	SetUnderline(underline command.UnderlineMode) error

	// === Impresión de texto ===
	Text(str string) error
	TextLn(str string) error

	// === Control de papel ===
	Cut(mode command.CutMode, lines int) error
	Feed(lines int) error

	// === Impresión de imágenes ===
	PrintImage(img image.Image) error
	PrintImageFromFile(filename string) error

	// TODO: Agregar más métodos según necesites
}

// GenericPrinter implementa Printer usando Protocol y Connector
type GenericPrinter struct {
	Protocol  protocol.Protocol
	Connector connector.Connector
	Profile   *profile.Profile

	// TODO: Agregar más campos si necesitas:
	// - Estado actual (font, alignment, etc.)
	// - Buffer de comandos
	// - Configuración
}

// NewGenericPrinter crea una nueva impresora genérica
func NewGenericPrinter(proto protocol.Protocol, conn connector.Connector, prof *profile.Profile) (*GenericPrinter, error) {
	if proto == nil {
		return nil, fmt.Errorf("el protocolo no puede ser nil")
	}
	if conn == nil {
		return nil, fmt.Errorf("el conector no puede ser nil")
	}
	if prof == nil {
		return nil, fmt.Errorf("el perfil no puede ser nil")
	}

	printer := &GenericPrinter{
		Protocol:  proto,
		Connector: conn,
		Profile:   prof,
	}

	// Inicializar la impresora
	if err := printer.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize printer: %w", err)
	}

	return printer, nil
}

// GetProfile devuelve el perfil de la impresora
func (p *GenericPrinter) GetProfile() *profile.Profile {
	return p.Profile
}

// SetProfile establece un nuevo perfil
func (p *GenericPrinter) SetProfile(newProfile *profile.Profile) {
	p.Profile = newProfile
}

// === Implementación de la interfaz Printer ===

// Initialize inicializa la impresora
func (p *GenericPrinter) Initialize() error {
	cmd := p.Protocol.Initialize()
	_, err := p.Connector.Write(cmd)
	return err
}

// Close cierra la conexión con la impresora
func (p *GenericPrinter) Close() error {
	// Primero enviar comandos de cierre del protocolo
	if closeCmd := p.Protocol.Close(); len(closeCmd) > 0 {
		_, _ = p.Connector.Write(closeCmd) // Ignorar error, vamos a cerrar de todos modos
	}

	// Luego cerrar el conector
	return p.Connector.Close()
}

// SetJustification establece la alineación del texto
func (p *GenericPrinter) SetJustification(alignment command.Alignment) error {
	cmd := p.Protocol.SetJustification(alignment)
	_, err := p.Connector.Write(cmd)
	// TODO: Si no hay error, guardar el estado actual
	return err
}

// SetFont establece la fuente
func (p *GenericPrinter) SetFont(font command.Font) error {
	cmd := p.Protocol.SetFont(font)
	_, err := p.Connector.Write(cmd)
	return err
}

// SetEmphasis activa/desactiva negrita
func (p *GenericPrinter) SetEmphasis(on bool) error {
	cmd := p.Protocol.SetEmphasis(on)
	_, err := p.Connector.Write(cmd)
	return err
}

// SetDoubleStrike activa/desactiva doble golpe
func (p *GenericPrinter) SetDoubleStrike(on bool) error {
	cmd := p.Protocol.SetDoubleStrike(on)
	_, err := p.Connector.Write(cmd)
	return err
}

// SetUnderline configura el subrayado
func (p *GenericPrinter) SetUnderline(underline command.UnderlineMode) error {
	cmd := p.Protocol.SetUnderline(underline)
	_, err := p.Connector.Write(cmd)
	return err
}

// Text imprime texto
func (p *GenericPrinter) Text(str string) error {
	cmd := p.Protocol.Text(str)
	_, err := p.Connector.Write(cmd)
	return err
}

// TextLn imprime texto con salto de línea
func (p *GenericPrinter) TextLn(str string) error {
	cmd := p.Protocol.TextLn(str)
	_, err := p.Connector.Write(cmd)
	return err
}

// Cut corta el papel
func (p *GenericPrinter) Cut(mode command.CutMode, lines int) error {
	// TODO: Verificar si la impresora tiene cutter con HasCapability
	cmd := p.Protocol.Cut(mode, lines) // 0 lines feed antes del corte
	_, err := p.Connector.Write(cmd)
	return err
}

// Feed alimenta papel
func (p *GenericPrinter) Feed(lines int) error {
	cmd := p.Protocol.Feed(lines)
	_, err := p.Connector.Write(cmd)
	return err
}

// PrintImageOptions contiene opciones para imprimir imágenes
type PrintImageOptions struct {
	Density    command.Density
	DitherMode imaging.DitherMode
	Threshold  uint8
	Width      int
}

// DefaultPrintImageOptions devuelve opciones por defecto
func DefaultPrintImageOptions() PrintImageOptions {
	return PrintImageOptions{
		Density:    command.DensitySingle,
		DitherMode: imaging.DitherNone,
		Threshold:  128,
		Width:      256,
	}
}

// PrintImage imprime una imagen con opciones por defecto
func (p *GenericPrinter) PrintImage(img image.Image) error {
	opts := DefaultPrintImageOptions()
	return p.PrintImageWithOptions(img, opts)
}

// PrintImageWithOptions imprime una imagen con opciones específicas
func (p *GenericPrinter) PrintImageWithOptions(img image.Image, opts PrintImageOptions) error {
	// Verificar soporte de imágenes
	if !p.Protocol.HasNativeImageSupport() {
		return fmt.Errorf("protocol %s does not support images", p.Protocol.Name())
	}

	// Crear PrintImage
	resizedImg := utils.ResizeToWidth(img, opts.Width, p.Profile.DotsPerLine)
	printImg := utils.NewPrintImage(resizedImg, opts.DitherMode)
	printImg.Threshold = opts.Threshold

	// Aplicar dithering si se especificó
	if opts.DitherMode != imaging.DitherNone {
		if err := printImg.ApplyDithering(opts.DitherMode); err != nil {
			return fmt.Errorf("failed to apply dithering: %w", err)
		}
	}

	// Generar comandos
	cmd, err := p.Protocol.PrintImage(printImg, opts.Density)
	if err != nil {
		return fmt.Errorf("failed to generate image commands: %w", err)
	}

	// Enviar a la impresora
	_, err = p.Connector.Write(cmd)
	return err
}

// Implementar PrintImageFromFile en GenericPrinter
func (p *GenericPrinter) PrintImageFromFile(filename string) error {
	file, err := utils.SafeOpen(filename)
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("printer: error al cerrar imagen: %s", cerr)
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}
	return p.PrintImage(img)
}

// TODO: Implementar el resto de métodos que necesites
