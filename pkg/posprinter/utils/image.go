package utils

import (
	"image"
	"image/color"

	"pos-daemon.adcon.dev/pkg/posprinter/imaging"
)

// PrintImage representa una imagen preparada para impresión
type PrintImage struct {
	// La imagen original de Go
	Source image.Image

	// Dimensiones efectivas para impresión
	Width  int
	Height int

	// Metadatos opcionales
	DPI       int
	Threshold uint8

	// Datos pre-procesados opcionales
	MonochromeData []byte

	// Imagen procesada con dithering (si se aplicó)
	ProcessedImage image.Image
	DitherMode     imaging.DitherMode
}

// NewPrintImage crea una nueva imagen para impresión
func NewPrintImage(img image.Image) *PrintImage {
	bounds := img.Bounds()
	return &PrintImage{
		Source:     img,
		Width:      bounds.Dx(),
		Height:     bounds.Dy(),
		DPI:        203,
		Threshold:  128,
		DitherMode: imaging.DitherNone,
	}
}

// ApplyDithering aplica un algoritmo de dithering a la imagen
func (p *PrintImage) ApplyDithering(mode imaging.DitherMode) error {
	processed, err := imaging.ProcessImageWithDithering(p.Source, mode, p.Threshold)
	if err != nil {
		return err
	}

	p.ProcessedImage = processed
	p.DitherMode = mode

	// Invalidar datos monocromáticos anteriores
	p.MonochromeData = nil

	return nil
}

// GetEffectiveImage devuelve la imagen a usar (procesada o original)
func (p *PrintImage) GetEffectiveImage() image.Image {
	if p.ProcessedImage != nil {
		return p.ProcessedImage
	}
	return p.Source
}

// GetPixel obtiene el valor de un pixel como blanco (false) o negro (true)
func (p *PrintImage) GetPixel(x, y int) bool {
	// Si tenemos datos monocromáticos, usarlos
	if p.MonochromeData != nil {
		byteIndex := (y*p.Width + x) / 8
		bitIndex := uint(7 - (x % 8))
		return p.MonochromeData[byteIndex]&(1<<bitIndex) != 0
	}

	// Usar la imagen efectiva (procesada o original)
	img := p.GetEffectiveImage()

	// Si la imagen ya es en escala de grises (resultado de dithering)
	if grayImg, ok := img.(*image.Gray); ok {
		return grayImg.GrayAt(x, y).Y < p.Threshold
	}

	// Convertir a escala de grises
	gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
	return gray.Y < p.Threshold
}

// ToMonochrome convierte la imagen a datos monocromáticos
func (p *PrintImage) ToMonochrome() []byte {
	if p.MonochromeData != nil {
		return p.MonochromeData
	}

	// Calcular bytes necesarios
	bytesPerRow := (p.Width + 7) / 8
	data := make([]byte, bytesPerRow*p.Height)

	// Usar la imagen efectiva
	img := p.GetEffectiveImage()

	// Convertir pixel por pixel
	for y := 0; y < p.Height; y++ {
		for x := 0; x < p.Width; x++ {
			// Para imágenes ya procesadas con dithering
			if grayImg, ok := img.(*image.Gray); ok {
				if grayImg.GrayAt(x, y).Y < p.Threshold {
					byteIndex := y*bytesPerRow + x/8
					bitIndex := uint(7 - (x % 8))
					data[byteIndex] |= 1 << bitIndex
				}
			} else {
				// Para imágenes a color
				if p.GetPixel(x, y) {
					byteIndex := y*bytesPerRow + x/8
					bitIndex := uint(7 - (x % 8))
					data[byteIndex] |= 1 << bitIndex
				}
			}
		}
	}

	p.MonochromeData = data
	return data
}
