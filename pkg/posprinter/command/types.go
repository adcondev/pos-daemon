package command

// Alignment define las alineaciones de texto estándar
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
	AlignJustified // Algunos protocolos podrían soportar esto
)

// Font define los tipos de fuente estándar
type Font int

const (
	FontA Font = iota
	FontB
	FontC
	// Otros tipos genéricos...
)

// UnderlineMode define los modos de subrayado estándar
type UnderlineMode int

const (
	UnderlineNone UnderlineMode = iota
	UnderlineSingle
	UnderlineDouble
)

// BarcodeType define los tipos de código de barras estándar
type BarcodeType int

const (
	BarcodeUPCA BarcodeType = iota
	BarcodeUPCE
	BarcodeEAN13
	BarcodeEAN8
	BarcodeCode39
	BarcodeITF
	BarcodeCodebar
	BarcodeCode93
	BarcodeCode128
	// Otros tipos comunes...
)

// BarcodeTextPosition define posiciones estándar para texto en códigos de barras
type BarcodeTextPosition int

const (
	BarcodeTextNone BarcodeTextPosition = iota
	BarcodeTextAbove
	BarcodeTextBelow
	BarcodeTextBoth
)

// CutMode define modos de corte estándar
type CutMode int

const (
	CutFull CutMode = iota
	CutPartial
)

// Density define densidades de impresión estándar para imágenes
type Density int

const (
	DensitySingle Density = iota
	DensityDouble
	DensityTriple
	DensityQuadruple
)
