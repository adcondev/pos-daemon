package escpos

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
)

// La implementación real para cargar y convertir imágenes (ToRasterFormat, ToColumnFormat)
// debe ser proporcionada. Esto implica manipulación de píxeles y formatos específicos de ESC/POS.
type EscposImage struct {
	img              image.Image
	threshold        uint8
	width            int
	height           int
	rasterData       []byte
	columnFormatHigh [][]byte
	columnFormatLow  [][]byte
}

func NewEscposImage(img image.Image, threshold uint8) *EscposImage {
	bounds := img.Bounds()
	return &EscposImage{
		img:       img,
		threshold: threshold,
		width:     bounds.Dx(),
		height:    bounds.Dy(),
	}
}

func NewEscposImageFromFile(filename string, threshold uint8) (*EscposImage, error) {
	file, err := openFile(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("error cerrando el archivo: %v\n", err)
		}
	}()

	return NewEscposImageFromReader(file, threshold)
}

func NewEscposImageFromReader(reader io.Reader, threshold uint8) (*EscposImage, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decodificando la image: %w", err)
	}

	return NewEscposImage(img, threshold), nil
}

func NewEscposImageFromBytes(data []byte, threshold uint8) (*EscposImage, error) {
	return NewEscposImageFromReader(bytes.NewReader(data), threshold)
}

func (ei *EscposImage) GetWidth() int {
	return ei.width
}

func (ei *EscposImage) GetHeight() int {
	return ei.height
}

func (ei *EscposImage) GetWidthBytes() int {
	return (ei.width + 7) / 8
}

func (ei *EscposImage) ToRasterFormat() ([]byte, error) {
	if ei.rasterData == nil {
		if err := ei.processRasterData(); err != nil {
			return nil, err
		}
	}
	return ei.rasterData, nil
}

func (ei *EscposImage) ToColumnFormat(highDensity bool) ([][]byte, error) {
	if highDensity {
		if ei.columnFormatHigh == nil {
			if err := ei.processColumnData(true); err != nil {
				return nil, err
			}
		}
		return ei.columnFormatHigh, nil
	}

	if ei.columnFormatLow == nil {
		if err := ei.processColumnData(false); err != nil {
			return nil, err
		}
	}
	return ei.columnFormatLow, nil
}

func (ei *EscposImage) processRasterData() error {
	if ei.img == nil {
		return errors.New("imagen no inicializada")
	}

	// Cada línea de bytes representa 8 píxeles verticales
	widthBytes := ei.GetWidthBytes()
	result := make([]byte, widthBytes*ei.height)

	for y := 0; y < ei.height; y++ {
		for x := 0; x < ei.width; x++ {
			// Determinar si el pixel es negro u oscuro basado en el threshold
			if ei.isBlack(x, y) {
				bytePos := (y * widthBytes) + (x / 8)
				bitPos := 7 - (x % 8) // El bit más significativo es el píxel izquierdo
				result[bytePos] |= 1 << bitPos
			}
		}
	}

	ei.rasterData = result
	return nil
}

func (ei *EscposImage) processColumnData(highDensity bool) error {
	if ei.img == nil {
		return errors.New("imagen no inicializada")
	}

	// Calcular cuántas filas necesitamos
	dotsPerColumn := 24
	if !highDensity {
		dotsPerColumn = 8
	}

	rowCount := (ei.height + dotsPerColumn - 1) / dotsPerColumn
	result := make([][]byte, rowCount)

	// Procesar cada fila
	for row := 0; row < rowCount; row++ {
		startY := row * dotsPerColumn
		endY := startY + dotsPerColumn
		if endY > ei.height {
			endY = ei.height
		}

		rowData := make([]byte, ei.width*3) // 3 bytes por columna en modo 24 dots
		if !highDensity {
			rowData = make([]byte, ei.width) // 1 byte por columna en modo 8 dots
		}

		// Procesar cada columna (píxel horizontal)
		for x := 0; x < ei.width; x++ {
			if highDensity {
				// 24 píxeles verticales = 3 bytes por columna
				for i := 0; i < 3; i++ {
					for b := 0; b < 8; b++ {
						y := startY + (i * 8) + b
						if y < endY && ei.isBlack(x, y) {
							rowData[x*3+i] |= 1 << (7 - b)
						}
					}
				}
			} else {
				// 8 píxeles verticales = 1 byte por columna
				for b := 0; b < 8; b++ {
					y := startY + b
					if y < endY && ei.isBlack(x, y) {
						rowData[x] |= 1 << (7 - b)
					}
				}
			}
		}

		result[row] = rowData
	}

	if highDensity {
		ei.columnFormatHigh = result
	} else {
		ei.columnFormatLow = result
	}

	return nil
}

func (ei *EscposImage) isBlack(x, y int) bool {
	if x < 0 || y < 0 || x >= ei.width || y >= ei.height {
		return false
	}

	// Obtener el color del píxel
	c := ei.img.At(x, y)

	// Convertir a escala de grises (luminancia)
	r, g, b, _ := c.RGBA()
	// Los valores están en el rango 0-65535, por lo que necesitamos convertirlos a 0-255
	gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)

	// Determinar si es negro basado en el threshold
	return gray <= ei.threshold
}

// Función de ayuda para abrir archivos
func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo %s: %w", filename, err)
	}
	return file, nil
}
