package imaging

// Constantes de densidad de imagen (para compatibilidad)
const (
	ImgDefault = 0
	ImgDouble  = 1
)

// Alias para dithering methods (para compatibilidad)
const (
	FloydStein = DitherFloydSteinberg // Alias para compatibilidad
)

// Constantes de color
const (
	Color1 = 0 // Negro
	Color2 = 1 // Rojo (en impresoras que lo soporten)
)

// Tamaños predefinidos para impresión
const (
	DefaultPrintSize = 384 // Ancho típico para impresoras de 58mm
	LargePrintSize   = 576 // Ancho típico para impresoras de 80mm
)
