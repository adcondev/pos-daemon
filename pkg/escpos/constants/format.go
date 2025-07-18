package constants

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
