package escpos

import (
	"image"
)

// === Adaptadores para mantener compatibilidad con el código existente ===

// EscposImage mantiene la estructura antigua para compatibilidad
type EscposImage struct {
	img       image.Image
	threshold uint8
}

// NewEscposImage crea una imagen compatible con código antiguo
func NewEscposImage(img image.Image, threshold uint8) *EscposImage {
	return &EscposImage{
		img:       img,
		threshold: threshold,
	}
}
