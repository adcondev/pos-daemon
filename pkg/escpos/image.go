package escpos

import (
	"errors"
	"fmt"
)

const (
	// Tamaño de imagen (para comandos Bit Image)
	IMG_DEFAULT       int = 0
	IMG_DOUBLE_WIDTH  int = 1
	IMG_DOUBLE_HEIGHT int = 2

	// Color (para impresoras con múltiples colores)
	COLOR_1 int = 0 // Color 1 (generalmente negro)
	COLOR_2 int = 1 // Color 2 (generalmente rojo)
)

// BitImage imprime una imagen utilizando el comando de imagen de bits (GS v 0).
// Requiere que la imagen sea convertible a formato raster de 1 bit.
func (p *Printer) BitImage(img *Image, size int) error {
	if img == nil {
		return errors.New("BitImage: la imagen no puede ser nil")
	}
	if err := validateInteger(size, IMG_DEFAULT, IMG_DOUBLE_HEIGHT|IMG_DOUBLE_WIDTH, "BitImage", "tamaño"); err != nil {
		return fmt.Errorf("BitImage: %w", err)
	} // Combinación de IMG_DEFAULT, IMG_DOUBLE_WIDTH, IMG_DOUBLE_HEIGHT

	rasterData, err := img.ToRasterFormat() // Requiere implementación real de EscposImage
	if err != nil {
		return fmt.Errorf("BitImage: falló al obtener los datos raster: %w", err)
	}

	// Cabecera de datos: xL xH yL yH
	// xL xH: ancho en bytes (img.GetWidthBytes()) - 2 bytes
	// yL yH: alto en puntos (img.GetHeight()) - 2 bytes
	headerBytes, err := dataHeader([]int{img.GetWidthBytes(), img.GetHeight()}, true) // true para 2 bytes por valor
	if err != nil {
		return fmt.Errorf("BitImage: falló al crear la cabecera de datos: %w", err)
	}

	// Comando: GS v 0 m xL xH yL yH d1...dk
	// m es el modo de tamaño (0-3)
	cmdHeader := []byte{GS, 'v', '0', byte(size)}
	cmdHeader = append(cmdHeader, headerBytes...)

	_, err = p.connector.Write(cmdHeader)
	if err != nil {
		return fmt.Errorf("BitImage: falló al enviar la cabecera del comando: %w", err)
	}

	_, err = p.connector.Write(rasterData) // Enviar los datos de la imagen
	if err != nil {
		return fmt.Errorf("BitImage: falló al enviar los datos raster: %w", err)
	}

	return nil
}

// BitImageColumnFormat imprime una imagen utilizando el comando de modo gráfico (ESC *).
// Este comando imprime por líneas de 8 o 24 puntos verticales.
func (p *Printer) BitImageColumnFormat(img *Image, size int) error {
	if img == nil {
		return errors.New("BitImageColumnFormat: la imagen no puede ser nil")
	}
	// PHP valida size 0-3. La lógica interna usa los bits 1 y 2.
	if err := validateInteger(size, IMG_DEFAULT, IMG_DOUBLE_HEIGHT|IMG_DOUBLE_WIDTH, "BitImageColumnFormat", "tamaño"); err != nil {
		return fmt.Errorf("BitImageColumnFormat: %w", err)
	}

	// La clase PHP establece el espaciado de línea a 16 (ESC 3 16) antes de imprimir líneas de imagen
	// y lo restablece después. Esto es necesario para que las líneas de imagen no tengan espacio entre ellas.
	if err := p.SetLineSpacing(intPtr(16)); err != nil {
		return fmt.Errorf("BitImageColumnFormat: falló al establecer el espaciado de línea: %w", err)
	}
	// Asegurar que el espaciado se restablezca incluso si hay un error.
	defer p.SetLineSpacing(nil) // nil restablece al espaciado por defecto

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
	if (size & IMG_DOUBLE_HEIGHT) == IMG_DOUBLE_HEIGHT {
		densityCode &^= 32 // Desactivar bit 5 (32) -> 8 puntos verticales
	}
	if (size & IMG_DOUBLE_WIDTH) == IMG_DOUBLE_WIDTH {
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

		_, err := p.connector.Write(cmd)
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

// Graphics imprime una imagen utilizando los comandos de gráfico GS ( L.
// Este método es a menudo más robusto para imágenes grandes o de alta calidad.
func (p *Printer) Graphics(img *Image, size int) error {
	if img == nil {
		return errors.New("Graphics: la imagen no puede ser nil")
	}
	if err := validateInteger(size, IMG_DEFAULT, IMG_DOUBLE_HEIGHT|IMG_DOUBLE_WIDTH, "Graphics", "tamaño"); err != nil {
		return fmt.Errorf("Graphics: %w", err)
	} // Combinación de IMG_DEFAULT, IMG_DOUBLE_WIDTH, IMG_DOUBLE_HEIGHT

	rasterData, err := img.ToRasterFormat() // Requiere implementación real
	if err != nil {
		return fmt.Errorf("Graphics: falló al obtener los datos raster: %w", err)
	}

	// Cabecera de imagen: xL xH yL yH (ancho en puntos, alto en puntos) - 2 bytes cada uno
	imgHeaderBytes, err := dataHeader([]int{img.GetWidth(), img.GetHeight()}, true) // true para 2 bytes por valor
	if err != nil {
		return fmt.Errorf("Graphics: falló al crear la cabecera de imagen: %w", err)
	}

	// Construir los datos para el comando 'p' (imprimir datos gráficos definidos por el usuario)
	// Formato: tono xm ym colors imgHeader rasterData
	// tono: '0' (normal)
	// xm: multiplicador horizontal ('1' o '2')
	// ym: multiplicador vertical ('1' o '2')
	// colors: '1' (1 bit por píxel)
	// PHP usa chr(1) o chr(2) para xm/ym. Replicamos.
	xm := byte(1)
	if (size & IMG_DOUBLE_WIDTH) == IMG_DOUBLE_WIDTH {
		xm = 2
	}
	ym := byte(1)
	if (size & IMG_DOUBLE_HEIGHT) == IMG_DOUBLE_HEIGHT {
		ym = 2
	}

	graphicsDataP := []byte{'0', xm, ym, '1'}                // tono, xm, ym, colors
	graphicsDataP = append(graphicsDataP, imgHeaderBytes...) // Cabecera de imagen
	graphicsDataP = append(graphicsDataP, rasterData...)     // Datos raster

	// Enviar comando para definir/imprimir los datos gráficos (fn='p')
	// El wrapper calcula pL pH.
	if err := p.wrapperSendGraphicsData(byte('0'), byte('p'), graphicsDataP); err != nil {
		return fmt.Errorf("Graphics: falló al enviar los datos gráficos (fn 'p'): %w", err)
	}

	// Enviar comando para imprimir el último dato gráfico definido (fn='2')
	// Este comando no tiene datos adicionales después de m y fn.
	if err := p.wrapperSendGraphicsData(byte('0'), byte('2'), []byte{}); err != nil {
		return fmt.Errorf("Graphics: falló al enviar el comando de impresión (fn '2'): %w", err)
	}

	return nil
}
