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
	CutFeed CutMode = iota
	Cut
)

// Density define densidades de impresión estándar para imágenes
type Density int

const (
	DensitySingle Density = iota
	DensityDouble
	DensityTriple
	DensityQuadruple
)

type QRModel byte

const (
	Model1 QRModel = iota // Modelo 1 (estándar)
	Model2                // Modelo 2 (recomendado y estándar)
)

type QRErrorCorrection byte

const (
	ECLow     QRErrorCorrection = iota // 7% de corrección
	ECMedium                           // 15% de corrección
	ECHigh                             // 25% de corrección
	ECHighest                          // 30% de corrección
)

type QRModuleSize byte

const (
	MinType QRModuleSize = 1
	MaxType QRModuleSize = 16
)

type CharacterSet int

const (
	CP437      CharacterSet = iota // CP437 U.S.A. / Standard Europe
	Katakana                       // Katakana (JIS X 0201)
	CP850                          // CP850 Multilingual
	CP860                          // CP860 Portuguese
	CP863                          // CP863 Canadian French
	CP865                          // CP865 Nordic
	WestEurope                     // WestEurope (ISO-8859-1)
	Greek                          // Greek (ISO-8859-7)
	Hebrew                         // Hebrew (ISO-8859-8)
	CP755                          // CP755 East Europe (not directly supported)
	Iran                           // Iran (CP720 Arabic)
	WCP1252                        // WCP1252 Windows-1252
	CP866                          // CP866 Cyrillic #2
	CP852                          // CP852 Latin2
	CP858                          // CP858 Multilingual + Euro
	IranII                         // IranII (CP864)
	Latvian                        // Latvian (Windows-1257)
)

// TODO: Agregar más tipos genéricos según necesites
// Por ejemplo:
// - QRCodeSize
// - PrintSpeed
// - CharacterSet
// etc.
