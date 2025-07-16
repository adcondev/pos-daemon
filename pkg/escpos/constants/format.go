package constants

const (
	// Fuentes
	FONT_A int = 0
	FONT_B int = 1
	FONT_C int = 2

	// Justificación del texto
	JUSTIFY_LEFT   int = 0
	JUSTIFY_CENTER int = 1
	JUSTIFY_RIGHT  int = 2

	// Modo de impresión (combinación de bits para ESC !)
	MODE_FONT_A        int = 0   // Bit 0 OFF for Font A
	MODE_FONT_B        int = 1   // Bit 0 ON for Font B
	MODE_EMPHASIZED    int = 8   // Bit 3 ON (Negrita)
	MODE_DOUBLE_HEIGHT int = 16  // Bit 4 ON (Doble Altura)
	MODE_DOUBLE_WIDTH  int = 32  // Bit 5 ON (Doble Ancho)
	MODE_UNDERLINE     int = 128 // Bit 7 ON (Subrayado)

	// Tipo de subrayado
	UNDERLINE_NONE   int = 0
	UNDERLINE_SINGLE int = 1
	UNDERLINE_DOUBLE int = 2
)
