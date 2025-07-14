package escpos

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
)

// Constantes para el dithering
const (
	DitherNone       = 0 // Sin dithering
	DitherFloydStein = 1 // Dithering Floyd-Steinberg
	DitherOrdered    = 2 // Dithering ordenado (matriz 4x4)

	// Tamaño recomendado para QR y otras imágenes en tickets
	DefaultPrintSize = 256
	MaxPrintSize     = 576
)

// ImageWithDithering procesa una imagen con dithering y la imprime
// usando el comando de imagen de bits (GS v 0).
// density: modo de densidad (0-3)
// ditherMethod: método de dithering a utilizar
func (p *Printer) ImageWithDithering(img image.Image, density int, ditherMethod int, size int) error {
	if img == nil {
		return fmt.Errorf("ImageWithDithering: la imagen no puede ser nil")
	}

	// Procesar la imagen con dithering
	processedImg, err := ProcessImageWithDithering(img, ditherMethod, size)
	if err != nil {
		return fmt.Errorf("ImageWithDithering: error al procesar la imagen: %w", err)
	}

	// Convertir a formato escpos.Image
	escposImg := NewEscposImage(processedImg, 128)

	// Imprimir usando BitImage
	return p.BitImage(escposImg, density)
}

// ProcessImageWithDithering procesa una imagen con el método de dithering especificado
// Devuelve una imagen en escala de grises o binaria según el método de dithering
func ProcessImageWithDithering(img image.Image, ditherMethod int, size int) (image.Image, error) {
	// Redimensionar a 256x256 si es necesario
	img = ResizeImage(img, size, size)

	// Convertir a escala de grises primero
	grayImg := image.NewGray(img.Bounds())
	draw.Draw(grayImg, grayImg.Bounds(), img, img.Bounds().Min, draw.Src)

	// Aplicar dithering según el método seleccionado
	switch ditherMethod {
	case DitherNone:
		// No aplicar dithering, solo binarizar con un threshold
		return ThresholdImage(grayImg, 128), nil

	case DitherFloydStein:
		return FloydSteinbergDither(grayImg, 128), nil

	case DitherOrdered:
		return OrderedDither(grayImg, 128), nil

	default:
		return nil, fmt.Errorf("método de dithering no soportado: %d", ditherMethod)
	}
}

// FloydSteinbergDither aplica el algoritmo de dithering de Floyd-Steinberg a una imagen en escala de grises
func FloydSteinbergDither(img *image.Gray, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	// Clonar la imagen original para no modificarla
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.SetGray(x, y, img.GrayAt(x, y))
		}
	}

	// Matriz para almacenar los valores en coma flotante durante el procesamiento
	width := bounds.Dx()
	height := bounds.Dy()
	buffer := make([][]float64, height)
	for y := 0; y < height; y++ {
		buffer[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			buffer[y][x] = float64(result.GrayAt(x+bounds.Min.X, y+bounds.Min.Y).Y)
		}
	}

	// Aplicar el algoritmo de Floyd-Steinberg
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := buffer[y][x]
			newPixel := float64(0) // Negro
			if oldPixel > float64(threshold) {
				newPixel = 255.0 // Blanco
			}

			// Establecer el nuevo valor del píxel
			result.SetGray(x+bounds.Min.X, y+bounds.Min.Y, color.Gray{Y: uint8(newPixel)})

			// Calcular el error
			quant_error := oldPixel - newPixel

			// Distribuir el error a los píxeles vecinos
			if x < width-1 {
				buffer[y][x+1] += quant_error * 7.0 / 16.0
			}
			if y < height-1 {
				if x > 0 {
					buffer[y+1][x-1] += quant_error * 3.0 / 16.0
				}
				buffer[y+1][x] += quant_error * 5.0 / 16.0
				if x < width-1 {
					buffer[y+1][x+1] += quant_error * 1.0 / 16.0
				}
			}
		}
	}

	return result
}

// ThresholdImage convierte una imagen en escala de grises a binaria usando un umbral simple
func ThresholdImage(img *image.Gray, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.GrayAt(x, y).Y > threshold {
				result.SetGray(x, y, color.Gray{Y: 255})
			} else {
				result.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return result
}

// OrderedDither aplica dithering ordenado usando una matriz 4x4
func OrderedDither(img *image.Gray, baseThreshold uint8) *image.Gray {
	// Matriz de umbral 4x4 para dithering ordenado
	threshold := [4][4]uint8{
		{0, 128, 32, 160},
		{192, 64, 224, 96},
		{48, 176, 16, 144},
		{240, 112, 208, 80},
	}

	bounds := img.Bounds()
	result := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Obtener el umbral de la matriz para esta posición
			tx := (x - bounds.Min.X) % 4
			ty := (y - bounds.Min.Y) % 4
			t := threshold[ty][tx]

			// Ajustar el umbral base con la matriz
			adjustedThreshold := (baseThreshold * t) / 128

			// Aplicar umbral
			if img.GrayAt(x, y).Y > adjustedThreshold {
				result.SetGray(x, y, color.Gray{Y: 255})
			} else {
				result.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return result
}

// ResizeImage redimensiona una imagen a un tamaño específico
// Implementación simplificada - para una implementación más sofisticada,
// considera usar paquetes como github.com/nfnt/resize
func ResizeImage(img image.Image, width, height int) image.Image {
	// Aquí deberías usar una biblioteca de redimensionado más sofisticada
	// Como simplificación, solo creamos una nueva imagen del tamaño deseado
	// y copiamos la imagen original escalándola

	bounds := img.Bounds()
	if bounds.Dx() == width && bounds.Dy() == height {
		return img // No es necesario redimensionar
	}

	// Crea una nueva imagen del tamaño deseado
	result := image.NewRGBA(image.Rect(0, 0, width, height))

	// Factores de escalado
	scaleX := float64(bounds.Dx()) / float64(width)
	scaleY := float64(bounds.Dy()) / float64(height)

	// Escala básica por aproximación de vecino más cercano
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := bounds.Min.X + int(float64(x)*scaleX)
			srcY := bounds.Min.Y + int(float64(y)*scaleY)
			result.Set(x, y, img.At(srcX, srcY))
		}
	}

	return result
}
