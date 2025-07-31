package escpos

import (
	"fmt"

	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/utils"
)

// ESCImage ahora es más simple, solo guarda referencia a PrintImage
type ESCImage struct {
	printImage *utils.PrintImage

	// Cache de datos procesados
	rasterData []byte
}

// newESCImageFromPrintImage crea una ESCImage desde PrintImage
func newESCImageFromPrintImage(img *utils.PrintImage) (*ESCImage, error) {
	if img == nil {
		return nil, fmt.Errorf("print image cannot be nil")
	}

	if img.Width <= 0 || img.Height <= 0 {
		return nil, fmt.Errorf("invalid image dimensions: %dx%d", img.Width, img.Height)
	}

	return &ESCImage{
		printImage: img,
	}, nil
}

// GetWidth devuelve el ancho en píxeles
func (e *ESCImage) GetWidth() int {
	return e.printImage.Width
}

// GetHeight devuelve el alto en píxeles
func (e *ESCImage) GetHeight() int {
	return e.printImage.Height
}

// GetWidthBytes devuelve el ancho en bytes
func (e *ESCImage) GetWidthBytes() int {
	return (e.printImage.Width + 7) / 8
}

// toRasterFormat convierte la imagen al formato raster de ESC/POS
func (e *ESCImage) toRasterFormat() ([]byte, error) {
	// Si ya tenemos los datos en cache, devolverlos
	if e.rasterData != nil {
		return e.rasterData, nil
	}

	// Obtener datos monocromáticos de la imagen
	// PrintImage se encarga de aplicar dithering si fue configurado
	e.rasterData = e.printImage.ToMonochrome()

	return e.rasterData, nil
}

// PrintImage implementa el méthodo para el protocolo ESC/POS
func (p *ESCPOSProtocol) PrintImage(img *utils.PrintImage, density command.Density) ([]byte, error) {
	// Crear ESCImage
	escImg, err := newESCImageFromPrintImage(img)
	if err != nil {
		return nil, err
	}

	// Obtener datos raster
	rasterData, err := escImg.toRasterFormat()
	if err != nil {
		return nil, err
	}

	// Mapear densidad a modo ESC/POS
	var mode byte
	switch density {
	case command.DensitySingle:
		mode = 0
	case command.DensityDouble:
		mode = 1
	default:
		mode = 0
	}

	// Construir comando GS v 0
	cmd := []byte{GS, 'v', '0', mode}

	// Agregar dimensiones
	widthBytes, err := utils.IntLowHigh(escImg.GetWidthBytes(), 2)
	if err != nil {
		return nil, err
	}
	heightBytes, err := utils.IntLowHigh(escImg.GetHeight(), 2)
	if err != nil {
		return nil, err
	}

	cmd = append(cmd, widthBytes...)
	cmd = append(cmd, heightBytes...)
	cmd = append(cmd, rasterData...)

	return cmd, nil
}
