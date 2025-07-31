package protocol

import (
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/utils"
)

// Protocol define una interfaz para cualquier protocolo de impresión.
// Esta interfaz devuelve comandos en bytes que el conector enviará a la impresora
type Protocol interface {
	// === Comandos básicos ===
	Initialize() []byte
	Close() []byte

	// === Manipulación de texto ===
	SetJustification(justification command.Alignment) []byte
	SetFont(font command.Font) []byte
	SetEmphasis(on bool) []byte
	SetDoubleStrike(on bool) []byte
	SetUnderline(underline command.UnderlineMode) []byte
	SetTextSize(widthMultiplier int, heightMultiplier int) []byte
	SetLineSpacing(height *int) []byte
	SetPrintLeftMargin(margin int) []byte
	SetPrintWidth(width int) []byte

	// === Manejo de Character Tables/Code Pages ===
	SelectCharacterTable(table int) []byte
	GetCharacterTable() int

	// === Comandos de texto ===
	Text(str string) []byte
	TextLn(str string) []byte
	TextRaw(str string) []byte

	// === Códigos de barras ===
	SetBarcodeHeight(height int) []byte
	SetBarcodeWidth(width int) []byte
	SetBarcodeTextPosition(position command.BarcodeTextPosition) []byte
	Barcode(content string, barType command.BarcodeType) ([]byte, error)

	// === Imágenes ===
	// PrintImage recibe una imagen genérica y la convierte a comandos del protocolo
	PrintImage(img *utils.PrintImage, density command.Density) ([]byte, error)

	// HasNativeImageSupport indica si el protocolo soporta imágenes nativas
	// (algunos protocolos solo soportan ciertos formatos)
	HasNativeImageSupport() bool

	// GetMaxImageWidth devuelve el ancho máximo de imagen soportado
	GetMaxImageWidth() int

	// === Control de papel ===
	Cut(mode command.CutMode, lines int) []byte
	Feed(lines int) []byte
	FeedReverse(lines int) []byte
	FeedForm() []byte

	// === Hardware ===
	Pulse(pin int, onMS int, offMS int) []byte
	Release() []byte

	// === Información del protocolo ===
	Name() string
	HasCapability(cap string) bool

	// TODO: Agregar más métodos según necesites:
	// - PrintQRCode(data string, size int) ([]byte, error)
	// - SetPrintSpeed(speed int) []byte
	// - GetStatus() []byte
	// etc.
}

// ProtocolFactory es una función que crea una instancia de un protocolo
type ProtocolFactory func() Protocol

// TODO: Definir capabilities estándar que todos los protocolos pueden reportar
const (
	CapabilityQRNative   = "qr_native"
	CapabilityCutter     = "cutter"
	CapabilityColorPrint = "color"
	CapabilityBarcodeB   = "barcode_b"
	// Agregar más según necesites
)
