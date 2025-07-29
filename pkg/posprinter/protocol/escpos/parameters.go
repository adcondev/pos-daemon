package escpos

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

type (
	Font          int
	Justify       int
	PrintMode     int
	UnderlineMode int
)

const (
	// Fuentes

	A Font = 0
	B Font = 1
	C Font = 2

	// Justificación del texto

	Left   Justify = 0
	Center Justify = 1
	Right  Justify = 2

	// Modo de impresión (combinación de bits para ESC !)

	FontA        PrintMode = 0x0  // Bit 0 OFF for Font A
	FontB        PrintMode = 0x1  // Bit 0 ON for Font B
	Emphasized   PrintMode = 0x8  // Bit 3 ON (Negrita)
	DoubleHeight PrintMode = 0x10 // Bit 4 ON (Doble Altura)
	DoubleWidth  PrintMode = 0x20 // Bit 5 ON (Doble Ancho)
	Underline    PrintMode = 0x80 // Bit 7 ON (Subrayado)

	// Tipo de subrayado
	NoUnderline UnderlineMode = 0
	Single      UnderlineMode = 1
	Double      UnderlineMode = 2
)

const (
	// Modo de corte de papel
	CUT_FULL    int = 65 // 'A'
	CUT_PARTIAL int = 66 // 'B'
)
