package protocol

import (
	"image"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
)

// TODO: Implementar los métodos de la interfaz Protocol para diferentes protocolos de impresión
// Protocol define una interfaz para cualquier protocolo de impresión.
type Protocol interface {
	// Comandos básicos
	Initialize() []byte
	Close() []byte

	// Manipulación de texto
	SetJustification(justification command.Alignment) []byte
	SetFont(font command.Font) []byte
	SetEmphasis(on bool) []byte
	SetDoubleStrike(on bool) []byte
	SetUnderline(underline command.UnderlineMode) []byte
	SetTextSize(widthMultiplier int, heightMultiplier int) []byte
	SetLineSpacing(height *int) []byte
	SetPrintLeftMargin(margin int) []byte
	SetPrintWidth(width int) []byte
	SelectCharacterTable(table int) []byte
	GetCharacterTable() int

	// Texto
	Text(str string) []byte
	TextLn(str string) []byte
	TextRaw(str string) []byte

	// Códigos de barras
	SetBarcodeHeight(height int) []byte
	SetBarcodeWidth(width int) []byte
	SetBarcodeTextPosition(position command.BarcodeTextPosition) []byte
	Barcode(content string, barType command.BarcodeType) ([]byte, error)

	// Imágenes
	PrintImage(img image.Image, density command.Density) ([]byte, error)

	// Control de papel
	Cut(mode command.CutMode, lines int) []byte
	Feed(lines int) []byte
	FeedReverse(lines int) []byte
	FeedForm() []byte

	// Hardware
	Pulse(pin int, onMS int, offMS int) []byte
	Release() []byte

	// Otros
	Name() string
	HasCapability(cap string) bool
}
