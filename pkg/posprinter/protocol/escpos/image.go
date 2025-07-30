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
	rasterData     []byte
	columnDataHigh [][]byte
	columnDataLow  [][]byte
}

// newESCImageFromPrintImage crea una ESCImage desde PrintImage
func newESCImageFromPrintImage(img *utils.PrintImage) (*ESCImage, error) {
	if img == nil {
		return nil, fmt.Errorf("print image cannot be nil")
	}

	if img.Width <= 0 || img.Height <= 0 {
		return nil, fmt.Errorf("invalid image dimensions: %dx%d", img.Width, img.Height)
	}

	// Validar ancho máximo
	const maxWidth = 576 // Típico para 80mm
	if img.Width > maxWidth {
		return nil, fmt.Errorf("image width %d exceeds maximum %d", img.Width, maxWidth)
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
func (e *ESCImage) toRasterFormat(density command.Density) ([]byte, error) {
	// Si ya tenemos los datos en cache, devolverlos
	if e.rasterData != nil {
		return e.rasterData, nil
	}

	// Obtener datos monocromáticos de la imagen
	// PrintImage se encarga de aplicar dithering si fue configurado
	e.rasterData = e.printImage.ToMonochrome()

	return e.rasterData, nil
}

// toColumnFormat convierte la imagen al formato de columna
func (e *ESCImage) toColumnFormat(highDensity bool) ([][]byte, error) {
	// Verificar cache
	if highDensity && e.columnDataHigh != nil {
		return e.columnDataHigh, nil
	}
	if !highDensity && e.columnDataLow != nil {
		return e.columnDataLow, nil
	}

	// Calcular parámetros
	dotsPerColumn := 24
	if !highDensity {
		dotsPerColumn = 8
	}

	rowCount := (e.printImage.Height + dotsPerColumn - 1) / dotsPerColumn
	result := make([][]byte, rowCount)

	// Procesar cada fila
	for row := 0; row < rowCount; row++ {
		startY := row * dotsPerColumn
		endY := startY + dotsPerColumn
		if endY > e.printImage.Height {
			endY = e.printImage.Height
		}

		var rowData []byte
		if highDensity {
			rowData = make([]byte, e.printImage.Width*3) // 3 bytes por columna
		} else {
			rowData = make([]byte, e.printImage.Width) // 1 byte por columna
		}

		// Procesar cada columna
		for x := 0; x < e.printImage.Width; x++ {
			if highDensity {
				// 24 píxeles = 3 bytes
				for i := 0; i < 3; i++ {
					for b := 0; b < 8; b++ {
						y := startY + (i * 8) + b
						if y < endY && e.printImage.GetPixel(x, y) {
							rowData[x*3+i] |= 1 << (7 - b)
						}
					}
				}
			} else {
				// 8 píxeles = 1 byte
				for b := 0; b < 8; b++ {
					y := startY + b
					if y < endY && e.printImage.GetPixel(x, y) {
						rowData[x] |= 1 << (7 - b)
					}
				}
			}
		}

		result[row] = rowData
	}

	// Guardar en cache
	if highDensity {
		e.columnDataHigh = result
	} else {
		e.columnDataLow = result
	}

	return result, nil
}

// PrintImage implementa el método para el protocolo ESC/POS
func (p *ESCPOSProtocol) PrintImage(img *utils.PrintImage, density command.Density) ([]byte, error) {
	// Crear ESCImage
	escImg, err := newESCImageFromPrintImage(img)
	if err != nil {
		return nil, err
	}

	// Obtener datos raster
	rasterData, err := escImg.toRasterFormat(density)
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
