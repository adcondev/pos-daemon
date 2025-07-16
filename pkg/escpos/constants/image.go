package constants

const (
	// Tamaño de imagen (para comandos Bit Image)
	IMG_DEFAULT       int = 0
	IMG_DOUBLE_WIDTH  int = 1
	IMG_DOUBLE_HEIGHT int = 2
	IMG_QUADRUPLE         = 3

	// Color (para impresoras con múltiples colores)
	COLOR_1 int = 0 // Color 1 (generalmente negro)
	COLOR_2 int = 1 // Color 2 (generalmente rojo)

	Uint32Size = 4

	// Tamaño recomendado para QR y otras imágenes en tickets
	DefaultPrintSize = 256
	MaxPrintSize     = 576
)
