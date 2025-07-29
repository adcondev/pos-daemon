package escpos

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"log"
	imaging2 "pos-daemon.adcon.dev/pkg/posprinter/imaging"
	utils2 "pos-daemon.adcon.dev/pkg/posprinter/utils"
)

// Image La implementación real para cargar y convertir imágenes (ToRasterFormat, ToColumnFormat)
// debe ser proporcionada. Esto implica manipulación de píxeles y formatos específicos de ESC/POS.
type ESCImage struct {
	img              image.Image
	threshold        uint8
	width            int
	height           int
	rasterData       []byte
	columnFormatHigh [][]byte
	columnFormatLow  [][]byte
}

func NewEscposImage(img image.Image, threshold uint8) *ESCImage {
	bounds := img.Bounds()
	return &ESCImage{
		img:       img,
		threshold: threshold,
		width:     bounds.Dx(),
		height:    bounds.Dy(),
	}
}

// BitImage imprime una imagen utilizando el comando de imagen de bits (GS v 0).
// Requiere que la imagen sea convertible a formato raster de 1 bit.
// Importante no imprimir sin salto de linea \n o Feed(1), ya que colisiona y genera caracteres no legibles.
func (p *ESCPrinter) BitImage(img *ESCImage, density int) error {
	if img == nil {
		return errors.New("BitImage: la imagen no puede ser nil")
	}
	if err := ValidateInteger(density, imaging2.ImgDefault, imaging2.ImgDoubleHeight|imaging2.ImgDoubleWidth, "BitImage", "tamaño"); err != nil {
		return fmt.Errorf("BitImage: %w", err)
	} // Combinación de IMG_DEFAULT, IMG_DOUBLE_WIDTH, IMG_DOUBLE_HEIGHT

	rasterData, err := img.ToRasterFormat() // Requiere implementación real de EscposImage
	if err != nil {
		return fmt.Errorf("BitImage: falló al obtener los datos raster: %w", err)
	}

	// Cabecera de datos: xL xH yL yH
	// xL xH: ancho en bytes (_2D.GetWidthBytes()) - 2 bytes
	// yL yH: alto en puntos (_2D.GetHeight()) - 2 bytes
	headerBytes, err := dataHeader([]int{img.GetWidthBytes(), img.GetHeight()}, true) // true para 2 bytes por valor
	if err != nil {
		return fmt.Errorf("BitImage: falló al crear la cabecera de datos: %w", err)
	}

	// Comando: GS v 0 m xL xH yL yH d1...dk
	// m es el modo de tamaño (0-3)
	cmdHeader := []byte{GS, 'v', '0', byte(density)}
	cmdHeader = append(cmdHeader, headerBytes...)

	_, err = p.Connector.Write(cmdHeader)
	if err != nil {
		return fmt.Errorf("BitImage: falló al enviar la cabecera del comando: %w", err)
	}

	_, err = p.Connector.Write(rasterData) // Enviar los datos de la imagen
	if err != nil {
		return fmt.Errorf("BitImage: falló al enviar los datos raster: %w", err)
	}

	return nil
}

// TODO: Revisar implementación ya que parece no implementar bien el comando de imagen.
// BitImageColumnFormat imprime una imagen utilizando el comando de modo gráfico (ESC *).
// Este comando imprime por líneas de 8 o 24 puntos verticales.
func (p *ESCPrinter) BitImageColumnFormat(img *ESCImage, size int) error {
	if img == nil {
		return errors.New("BitImageColumnFormat: la imagen no puede ser nil")
	}
	// PHP valida size 0-3. La lógica interna usa los bits 1 y 2.
	if err := ValidateInteger(size, imaging2.ImgDefault, imaging2.ImgDoubleHeight|imaging2.ImgDoubleWidth, "BitImageColumnFormat", "tamaño"); err != nil {
		return fmt.Errorf("BitImageColumnFormat: %w", err)
	}

	// La clase PHP establece el espaciado de línea a 16 (ESC 3 16) antes de imprimir líneas de imagen
	// y lo restablece después. Esto es necesario para que las líneas de imagen no tengan espacio entre ellas.
	if err := p.SetLineSpacing(utils2.IntPtr(16)); err != nil {
		return fmt.Errorf("BitImageColumnFormat: falló al establecer el espaciado de línea: %w", err)
	}
	// Asegurar que el espaciado se restablezca incluso si hay un error.
	defer func(p *ESCPrinter, height *int) {
		err := p.SetLineSpacing(height)
		if err != nil {
			log.Printf("image: error al restablecer espaciado")
		}
	}(p, nil) // nil restablece al espaciado por defecto

	// Lógica de densidad basada en los bits del parámetro size.
	// ESC * m - m define la densidad vertical y horizontal.
	// m=0: 8 puntos verticales, densidad horizontal normal.
	// m=1: 8 puntos verticales, doble densidad horizontal.
	// m=32: 24 puntos verticales, densidad horizontal normal.
	// m=33: 24 puntos verticales, doble densidad horizontal.
	// La lógica de PHP basada en IMG_DOUBLE_HEIGHT (2) e IMG_DOUBLE_WIDTH (1) parece un poco confusa
	// en comparación con la documentación estándar (donde "doble" en IMG_DOUBLE_... suele significar "menos denso" en términos de puntos por pulgada física, resultando en caracteres más grandes).
	// Vamos a interpretar el significado de los bits 1 y 2 de `size` de la manera más estándar:
	// Si IMG_DOUBLE_HEIGHT (bit 1, valor 2) está activado, usa 8 puntos verticales (m sin bit 5/32).
	// Si IMG_DOUBLE_WIDTH (bit 2, valor 1) está activado, usa densidad horizontal normal (m sin bit 0/1).
	// El modo por defecto (IMG_DEFAULT=0) suele ser 24 puntos verticales, doble densidad horizontal (m=33).

	densityCode := 33 // Valor por defecto: 24 puntos verticales, doble densidad horizontal
	if (size & imaging2.ImgDoubleHeight) == imaging2.ImgDoubleHeight {
		densityCode &^= 32 // Desactivar bit 5 (32) -> 8 puntos verticales
	}
	if (size & imaging2.ImgDoubleWidth) == imaging2.ImgDoubleWidth {
		densityCode &^= 1 // Desactivar bit 0 (1) -> densidad horizontal normal
	}

	// Determinar si la conversión de la imagen debe usar alta densidad vertical (24 puntos)
	// basándose en el `densityCode` calculado. Si el bit 5 (32) está activo, sí.
	useHighDensityVerticalForConversion := (densityCode & 32) != 0

	colFormatData, err := img.ToColumnFormat(useHighDensityVerticalForConversion) // Requiere implementación real
	if err != nil {
		return fmt.Errorf("BitImageColumnFormat: falló al obtener los datos en formato de columna: %w", err)
	}

	// Cabecera de datos: nL nH (número de puntos horizontales) - 2 bytes
	headerBytes, err := dataHeader([]int{img.GetWidth()}, true) // true para 2 bytes (ancho en puntos)
	if err != nil {
		return fmt.Errorf("BitImageColumnFormat: falló al crear la cabecera de datos: %w", err)
	}

	for _, lineData := range colFormatData {
		// Comando para cada línea: ESC * m nL nH d1...dk
		cmd := []byte{ESC, '*', byte(densityCode)}
		cmd = append(cmd, headerBytes...)
		cmd = append(cmd, lineData...) // Datos de la línea de la imagen

		_, err := p.Connector.Write(cmd)
		if err != nil {
			return fmt.Errorf("BitImageColumnFormat: falló al enviar la línea de imagen: %w", err)
		}

		// Avanzar papel una línea después de imprimir cada segmento de imagen vertical.
		// PHP hace esto con `feed()`.
		if err := p.Feed(1); err != nil {
			return fmt.Errorf("BitImageColumnFormat: falló al alimentar después de la línea: %w", err)
		}
	}

	// El espaciado de línea se restablece automáticamente debido a `defer`.

	return nil
}

// SetColor establece el color de impresión (para impresoras con múltiples colores).
// color puede ser COLOR_1 (negro) o COLOR_2 (rojo).
func (p *ESCPrinter) SetColor(color int) error {
	if err := ValidateInteger(color, imaging2.Color1, imaging2.Color2, "SetColor", "color"); err != nil {
		return fmt.Errorf("SetColor: %w", err)
	}
	// ESC r n - n=0: Color 1, 1: Color 2
	cmd := []byte{ESC, 'r', byte(color)}
	_, err := p.Connector.Write(cmd)
	return err
}

// SetReverseColors habilita o deshabilita la impresión en colores inversos.
func (p *ESCPrinter) SetReverseColors(on bool) error {
	// GS B n - n=1: habilitar, 0: deshabilitar
	val := byte(0)
	if on {
		val = 1
	}
	cmd := []byte{GS, 'B', val}
	_, err := p.Connector.Write(cmd)
	return err
}

func (ei *ESCImage) GetWidth() int {
	return ei.width
}

func (ei *ESCImage) GetHeight() int {
	return ei.height
}

func (ei *ESCImage) GetWidthBytes() int {
	return (ei.width + 7) / 8
}

func (ei *ESCImage) ToRasterFormat() ([]byte, error) {
	if ei.rasterData == nil {
		if err := ei.processRasterData(); err != nil {
			return nil, err
		}
	}
	return ei.rasterData, nil
}

func (ei *ESCImage) ToColumnFormat(highDensity bool) ([][]byte, error) {
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

func (ei *ESCImage) processRasterData() error {
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

func (ei *ESCImage) processColumnData(highDensity bool) error {
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

func (ei *ESCImage) isBlack(x, y int) bool {
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

// dataHeader formatea enteros de entrada en bytes (bajo/alto o byte único).
// Se utiliza para formatear dimensiones en comandos de imagen.
// long=true: formatar como 2 bytes (nL nH)
// long=false: formatar como 1 byte
func dataHeader(inputs []int, long bool) ([]byte, error) {
	var buf bytes.Buffer
	for _, input := range inputs {
		if long {
			// Formato de 2 bytes (nL nH) - rango 0 a 65535
			data, err := utils2.IntLowHigh(input, utils2.DimensionBytes)
			if err != nil {
				return nil, fmt.Errorf("dataHeader: falló al formatear el entero %d como 2 bytes: %w", input, err)
			}
			buf.Write(data)
		} else {
			// Formato de 1 byte - rango 0 a 255
			if input < 0 || input > 255 {
				return nil, fmt.Errorf("dataHeader: el entero %d está fuera del rango para un byte único (0-255)", input)
			}
			buf.WriteByte(byte(input))
		}
	}
	return buf.Bytes(), nil
}

// ImageWithDithering procesa una imagen con dithering y la imprime
// usando el comando de imagen de bits (GS v 0).
// density: modo de densidad (0-3)
// ditherMethod: dithering a utilizar
func (p *ESCPrinter) ImageWithDithering(img image.Image, density int, ditherMethod imaging2.DitherMode, size int) error {
	if img == nil {
		return fmt.Errorf("ImageWithDithering: la imagen no puede ser nil")
	}

	// Procesar la imagen con dithering
	processedImg, err := imaging2.ProcessImageWithDithering(img, ditherMethod, size)
	if err != nil {
		return fmt.Errorf("ImageWithDithering: error al procesar la imagen: %w", err)
	}

	// Convertir a formato escpos.Image
	escposImg := NewEscposImage(processedImg, 128)

	// Imprimir usando BitImage
	return p.BitImage(escposImg, density)
}
