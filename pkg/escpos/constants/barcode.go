package constants

type (
	Barcode         int
	BarcodeType     int
	BarcodeTextPos  int
	BarcodeTextFont int
	BarcodeWidth    int
)

const (
	// Códigos de barras

	UpcA    BarcodeType = 65 // 'A'
	UpcE    BarcodeType = 66 // 'B'
	Jan13   BarcodeType = 67 // 'C' (EAN13)
	Jan8    BarcodeType = 68 // 'D' (EAN8)
	Code39  BarcodeType = 69 // 'E'
	Itf     BarcodeType = 70 // 'F'
	Codabar BarcodeType = 71 // 'G'
	Code93  BarcodeType = 72 // 'H'
	Code128 BarcodeType = 73 // 'I'

	// Posición del texto del código de barras

	TextNone  BarcodeTextPos = 0
	TextAbove BarcodeTextPos = 1
	TextBelow BarcodeTextPos = 2
	TextBoth  BarcodeTextPos = 3

	// Fuente de texto del codigo de barras

	TextFontA BarcodeTextFont = 0 // Fuente A
	TextFontB BarcodeTextFont = 1 // Fuente B

	// Ancho del código de barras

	WidthNarrow     BarcodeWidth = 2 // Ancho normal
	WidthMedium     BarcodeWidth = 3 // Ancho medio
	WidthWide       BarcodeWidth = 4 // Ancho ancho
	WidthExtraWide  BarcodeWidth = 5 // Ancho extra ancho
	WidthDoubleWide BarcodeWidth = 6 // Ancho doble ancho
)
