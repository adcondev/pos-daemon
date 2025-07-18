package imaging

const (
	// Tamaño de imagen (para comandos Bit Image)

	ImgDefault      int = 0
	ImgDoubleWidth  int = 1
	ImgDoubleHeight int = 2

	// Color (para impresoras con múltiples colores)

	Color1 int = 0 // Color 1 (generalmente negro)
	Color2 int = 1 // Color 2 (generalmente rojo)

	Uint32Size = 4

	// Tamaño recomendado para QR y otras imágenes en tickets

	DefaultPrintSize = 256
	MaxPrintSize     = 576
)
