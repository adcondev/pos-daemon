package escpos

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"

	// "io"
	"math"
	"regexp"
	"strings"
)

// PrintConnector define la interfaz para la conexión física con la impresora.
// Debes implementar esta interfaz para tu método de conexión (USB, TCP, Serial, etc.).
type PrintConnector interface {
	// Write envía bytes a la impresora.
	Write([]byte) (int, error)
	// Close finaliza la conexión con la impresora.
	Close() error
}

// EscposImage TODO es un placeholder para la representación de una imagen.
// La implementación real para cargar y convertir imágenes (ToRasterFormat, ToColumnFormat)
// debe ser proporcionada. Esto implica manipulación de píxeles y formatos específicos de ESC/POS.
type EscposImage struct {
	// Datos de la imagen (por ejemplo, un objeto image.Image)
	pixels [][]byte // Ejemplo muy simplificado, cada byte representa 8 píxeles verticalmente
	width  int      // Ancho en píxeles
	height int      // Alto en píxeles
}

// NewEscposImageFromBytes TODO es un constructor placeholder.
// La implementación real debería cargar una imagen (PNG, JPEG, etc.) y convertirla a un formato interno adecuado.
func NewEscposImageFromBytes(data []byte) (*EscposImage, error) {
	// Esta es solo una simulación. La carga y el procesamiento de imágenes son complejos.
	// Un implementador real usaría paquetes como "image", "image/png", etc.
	// y lógica de dithering para convertir a 1 bit por píxel.
	return nil, errors.New("image loading and processing not implemented")
}

func (img *EscposImage) GetWidth() int {
	// Placeholder
	return img.width
}
func (img *EscposImage) GetHeight() int {
	// Placeholder
	return img.height
}
func (img *EscposImage) GetWidthBytes() int {
	// Placeholder: Ancho en bytes para formato raster (ancho en píxeles / 8, redondeado hacia arriba)
	return (img.width + 7) / 8
}

// ToRasterFormat TODO convierte la imagen a formato raster ESC/POS (GS v 0).
// Esta es una función placeholder.
func (img *EscposImage) ToRasterFormat() ([]byte, error) {
	// La implementación real requiere dithering y empaquetamiento de bits.
	return nil, errors.New("image raster format conversion not implemented")
}

// ToColumnFormat TODO convierte la imagen a formato de columna ESC/POS (ESC *).
// highDensity indica si se debe usar el modo de 24 puntos verticales.
// Esta es una función placeholder.
func (img *EscposImage) ToColumnFormat(highDensity bool) ([][]byte, error) {
	// La implementación real requiere dithering y empaquetamiento de bits por columna.
	return nil, errors.New("image column format conversion not implemented")
}

// CapabilityProfile describe las capacidades de una impresora ESC/POS específica.
type CapabilityProfile struct {
	SupportsBarcodeB     bool // Soporta el formato de comando GS k m L data (65-73)
	SupportsPdf417Code   bool
	SupportsQrCode       bool
	CodePages            map[int]string // Mapa de índice de codepage a nombre/descripción
	SupportsStarCommands bool           // Indica si soporta comandos específicos de Star (como ESC GS t)
	// Agrega otras capacidades según sea necesario (ej: anchos de papel, fuentes, etc.)
}

// LoadProfile TODO carga un CapabilityProfile predefinido o desde una fuente externa.
// Esta es una función placeholder.
func LoadProfile(name string) (*CapabilityProfile, error) {
	// En una biblioteca real, esto cargaría perfiles desde archivos o datos incrustados.
	// Para este port, devolvemos un perfil dummy "default".
	if name == "default" {
		return &CapabilityProfile{
			SupportsBarcodeB:   true, // Asumimos que el perfil por defecto soporta el formato moderno de código de barras
			SupportsPdf417Code: true, // Asumimos que el perfil por defecto soporta códigos 2D
			SupportsQrCode:     true,
			CodePages: map[int]string{ // Ejemplo de algunas codepages comunes
				0: "CP437",
				1: "CP850",
				2: "CP852",
				3: "CP858",
				4: "CP860",
				5: "CP863",
				6: "CP865",
				// Agrega más codepages soportadas
				16:  "WPC1252",   // Windows 1252
				254: "Shift_JIS", // Ejemplo asiático
				255: "GBK",       // Ejemplo chino (usado en textChinese)
			},
			SupportsStarCommands: false, // Asumimos que el perfil por defecto no soporta comandos Star
		}, nil
	}
	return nil, fmt.Errorf("perfil de capacidad desconocido: %s", name)
}

// Constantes ESC/POS y parámetros.
const (
	NUL byte = 0x00 // Null
	LF  byte = 0x0a // Line Feed
	ESC byte = 0x1b // Escape
	FS  byte = 0x1c // Field Separator / Group Separator
	FF  byte = 0x0c // Form Feed
	GS  byte = 0x1d // Group Separator

	// Códigos de barras
	BARCODE_UPCA    int = 65 // 'A'
	BARCODE_UPCE    int = 66 // 'B'
	BARCODE_JAN13   int = 67 // 'C' (EAN13)
	BARCODE_JAN8    int = 68 // 'D' (EAN8)
	BARCODE_CODE39  int = 69 // 'E'
	BARCODE_ITF     int = 70 // 'F'
	BARCODE_CODABAR int = 71 // 'G'
	BARCODE_CODE93  int = 72 // 'H'
	BARCODE_CODE128 int = 73 // 'I'

	// Posición del texto del código de barras
	BARCODE_TEXT_NONE  int = 0
	BARCODE_TEXT_ABOVE int = 1
	BARCODE_TEXT_BELOW int = 2

	// Color (para impresoras con múltiples colores)
	COLOR_1 int = 0 // Color 1 (generalmente negro)
	COLOR_2 int = 1 // Color 2 (generalmente rojo)

	// Modo de corte de papel
	CUT_FULL    int = 65 // 'A'
	CUT_PARTIAL int = 66 // 'B'

	// Fuentes
	FONT_A int = 0
	FONT_B int = 1
	FONT_C int = 2

	// Tamaño de imagen (para comandos Bit Image)
	IMG_DEFAULT       int = 0
	IMG_DOUBLE_WIDTH  int = 1
	IMG_DOUBLE_HEIGHT int = 2

	// Justificación del texto
	JUSTIFY_LEFT   int = 0
	JUSTIFY_CENTER int = 1
	JUSTIFY_RIGHT  int = 2

	// Modo de impresión (combinación de bits para ESC !)
	MODE_FONT_A        int = 0   // Bit 0 OFF for Font A
	MODE_FONT_B        int = 1   // Bit 0 ON for Font B
	MODE_EMPHASIZED    int = 8   // Bit 3 ON (Negrita)
	MODE_DOUBLE_HEIGHT int = 16  // Bit 4 ON (Doble Altura)
	MODE_DOUBLE_WIDTH  int = 32  // Bit 5 ON (Doble Ancho)
	MODE_UNDERLINE     int = 128 // Bit 7 ON (Subrayado)

	// Opciones de PDF417
	PDF417_STANDARD  int = 0
	PDF417_TRUNCATED int = 1

	// Niveles de corrección de error QR (aproximados)
	QR_ECLEVEL_L int = 0 // 7%
	QR_ECLEVEL_M int = 1 // 15%
	QR_ECLEVEL_Q int = 2 // 25%
	QR_ECLEVEL_H int = 3 // 30%

	// Modelos de QR
	QR_MODEL_1 int = 1
	QR_MODEL_2 int = 2
	QR_MICRO   int = 3

	// Tipos de estado (para comandos de estado, no implementados como métodos públicos en PHP)
	// Se incluyen por completitud de las constantes PHP
	STATUS_PRINTER       int = 1 // GS I 1 (Estado de la impresora)
	STATUS_OFFLINE_CAUSE int = 2 // GS I 2 (Causa de estar offline)
	STATUS_ERROR_CAUSE   int = 3 // GS I 3 (Causa del error)
	STATUS_PAPER_ROLL    int = 4 // GS I 4 (Estado del rollo de papel)
	STATUS_INK_A         int = 7 // GS I 7 (Estado de la tinta/cinta A)
	STATUS_INK_B         int = 6 // GS I 6 (Estado de la tinta/cinta B)
	STATUS_PEELER        int = 8 // GS I 8 (Estado del peeler - para etiquetas)

	// Tipo de subrayado
	UNDERLINE_NONE   int = 0
	UNDERLINE_SINGLE int = 1
	UNDERLINE_DOUBLE int = 2
)

// Printer representa una impresora térmica ESC/POS.
type Printer struct {
	connector      PrintConnector
	profile        *CapabilityProfile
	characterTable int // La tabla de caracteres (codepage) actualmente seleccionada.
}

// *** FUNCIÓN PARA CODIFICAR A CP858 ***
func toCP858(s string) []byte {
	// Obtener el codificador para CP858
	encoder := charmap.CodePage858.NewEncoder()
	// Convertir la string (UTF-8) a bytes codificados en CP858
	encoded, err := encoder.Bytes([]byte(s))
	if err != nil {
		// En caso de error (ej. carácter no representable en CP858),
		// podrías loguear el error, o intentar un fallback.
		// Aquí, por simplicidad, devolvemos la string original (UTF-8),
		// aunque esto no solucionaría el problema del acento si falla la codificación.
		// Una mejor práctica sería reemplazar el carácter desconocido.
		fmt.Printf("Advertencia: No se pudo codificar string a CP858: %v (original: %q)\n", err, s)
		return []byte(s) // Fallback (probablemente no imprimirá bien el carácter problemático)
	}
	return encoded
}

// NewPrinter crea una nueva instancia de Printer.
// Requiere un PrintConnector y opcionalmente un CapabilityProfile.
// Si el perfil es nil, carga el perfil por defecto.
func NewPrinter(connector PrintConnector, profile *CapabilityProfile) (*Printer, error) {
	if connector == nil {
		return nil, errors.New("PrintConnector no puede ser nil")
	}
	if profile == nil {
		// Cargar perfil por defecto si no se proporciona ninguno
		defaultProfile, err := LoadProfile("default")
		if err != nil {
			return nil, fmt.Errorf("falló al cargar el perfil de capacidad por defecto: %w", err)
		}
		profile = defaultProfile
	}

	p := &Printer{
		connector:      connector,
		profile:        profile,
		characterTable: 0, // Tabla de caracteres por defecto
	}

	// Inicializar la impresora
	if err := p.Initialize(); err != nil {
		// Si la inicialización falla, consideramos que la creación de la impresora falla.
		return nil, fmt.Errorf("falló al inicializar la impresora: %w", err)
	}

	return p, nil
}

// --- Métodos Públicos (Espejo de la clase PHP) ---

// Initialize restablece la impresora a su configuración por defecto (ESC @).
func (p *Printer) Initialize() error {
	_, err := p.connector.Write([]byte{ESC, '@'})
	if err == nil {
		p.characterTable = 0 // Resetear el seguimiento de la tabla de caracteres
	}
	return err
}

// Text envía una cadena de texto a la impresora.
// Maneja los saltos de línea '\n' convirtiéndolos a LF.
func (p *Printer) Text(str string) error {
	// Reemplazar los saltos de línea de Go/PHP ('\n') con el carácter LF ESC/POS (0x0a)
	bytesToSend := strings.ReplaceAll(str, "\n", string(LF))
	_, err := p.connector.Write(toCP858(bytesToSend))
	return err
}

// TextRaw envía una cadena de texto (o bytes) a la impresora sin procesar.
func (p *Printer) TextRaw(str string) error {
	_, err := p.connector.Write([]byte(str))
	return err
}

// TextChinese TODO envía texto en chino.
// Esta es una implementación placeholder ya que la conversión de codificación
// (UTF-8 a GBK) es compleja y requiere librerías externas en Go.
// Los comandos de activación/desactivación de modo chino (FS & / FS .) se incluyen.
func (p *Printer) TextChinese(str string) error {
	// Activar modo de caracteres chinos (FS &)
	cmd := []byte{FS, '&'}

	// --- Placeholder: Conversión de UTF-8 a GBK ---
	// En una implementación real, usarías un paquete como golang.org/x/text/encoding/chinese
	// gbkEncoder := chinese.GBK.NewEncoder()
	// gbkBytes, err := gbkEncoder.Bytes([]byte(str))
	// if err != nil {
	//     return fmt.Errorf("falló al codificar texto chino a GBK: %w", err)
	// }
	// cmd = append(cmd, gbkBytes...)
	// --- Fin Placeholder ---

	// Para demostración, enviar los bytes UTF-8 directamente (probablemente imprimirá basura si la impresora no está configurada para UTF-8)
	cmd = append(cmd, []byte(str)...)

	// Desactivar modo de caracteres chinos (FS .)
	cmd = append(cmd, FS, '.')

	_, err := p.connector.Write(cmd)
	return err
}

// Cut corta el papel.
// mode puede ser CUT_FULL o CUT_PARTIAL. lines es el número de líneas para alimentar antes de cortar (0-255).
func (p *Printer) Cut(mode int, lines int) error {
	// PHP usa chr(mode) donde mode es 65 ('A') o 66 ('B').
	// El comando estándar es GS V m [n], donde m es 0,1,48,49 (full/partial)
	// o m es 65,66 ('A'/'B') con un parámetro n adicional para líneas de avance.
	// Replicamos el comportamiento de PHP usando 'A' o 'B' y el parámetro lines.
	if err := validateInteger(mode, CUT_FULL, CUT_PARTIAL, "Cut", "modo"); err != nil {
		return fmt.Errorf("Cut: %w", err)
	} // 65 ('A') o 66 ('B')
	if err := validateInteger(lines, 0, 255, "Cut", "líneas"); err != nil {
		return fmt.Errorf("Cut: %w", err)
	}

	cmd := []byte{GS, 'V', byte(mode), byte(lines)} // GS V 'A'/'B' n
	_, err := p.connector.Write(cmd)
	return err
}

// Feed avanza el papel el número especificado de líneas.
func (p *Printer) Feed(lines int) error {
	if err := validateInteger(lines, 1, 255, "Feed", "líneas"); err != nil {
		return fmt.Errorf("Feed: %w", err)
	}
	if lines <= 1 {
		// Usar solo LF para una línea es un poco más rápido a veces
		_, err := p.connector.Write([]byte{LF})
		return err
	}
	// ESC d n - Imprime los datos del búfer y alimenta n líneas
	cmd := []byte{ESC, 'd', byte(lines)}
	_, err := p.connector.Write(cmd)
	return err
}

// FeedReverse retrocede el papel el número especificado de líneas.
func (p *Printer) FeedReverse(lines int) error {
	if err := validateInteger(lines, 1, 255, "FeedReverse", "líneas"); err != nil {
		return fmt.Errorf("FeedReverse: %w", err)
	}
	// ESC e n - Alimenta el papel hacia atrás n líneas
	cmd := []byte{ESC, 'e', byte(lines)}
	_, err := p.connector.Write(cmd)
	return err
}

// FeedForm alimenta el papel hasta el principio del siguiente formulario (poco común en impresoras de recibos).
func (p *Printer) FeedForm() error {
	// FF - Form Feed
	_, err := p.connector.Write([]byte{FF})
	return err
}

// Release envía un comando (ESC q) que PHP describe como "liberar la impresora del estado de espera".
// Este comando NO es ESC/POS estándar y es probable que sea específico del fabricante (como Star).
func (p *Printer) Release() error {
	// Advertencia: ESC q es probablemente específico del fabricante y no estándar ESC/POS.
	_, err := p.connector.Write([]byte{ESC, 'q'}) // PHP envía ESC seguido del byte 113 ('q')
	return err
}

// SetJustification establece la alineación del texto (izquierda, centro, derecha).
func (p *Printer) SetJustification(justification int) error {
	if err := validateInteger(justification, JUSTIFY_LEFT, JUSTIFY_RIGHT, "SetJustification", "justificación"); err != nil {
		return fmt.Errorf("SetJustification: %w", err)
	}
	// ESC a n - n=0: izquierda, 1: centro, 2: derecha
	cmd := []byte{ESC, 'a', byte(justification)}
	_, err := p.connector.Write(cmd)
	return err
}

// SetFont establece la fuente (A, B o C).
func (p *Printer) SetFont(font int) error {
	if err := validateInteger(font, FONT_A, FONT_C, "SetFont", "fuente"); err != nil {
		return fmt.Errorf("SetFont: %w", err)
	}
	// ESC M n - n=0: Fuente A, 1: Fuente B, 2: Fuente C
	cmd := []byte{ESC, 'M', byte(font)}
	_, err := p.connector.Write(cmd)
	return err
}

// SetEmphasis habilita o deshabilita el modo enfatizado (negrita).
func (p *Printer) SetEmphasis(on bool) error {
	// ESC E n - n=1: habilitar, n=0: deshabilitar
	val := byte(0)
	if on {
		val = 1
	}
	cmd := []byte{ESC, 'E', val}
	_, err := p.connector.Write(cmd)
	return err
}

// SetDoubleStrike habilita o deshabilita el modo doble golpeo.
func (p *Printer) SetDoubleStrike(on bool) error {
	// ESC G n - n=1: habilitar, n=0: deshabilitar
	val := byte(0)
	if on {
		val = 1
	}
	cmd := []byte{ESC, 'G', val}
	_, err := p.connector.Write(cmd)
	return err
}

// SetUnderline establece el modo de subrayado (ninguno, simple, doble).
// Puede aceptar 0 (none), 1 (single), 2 (double).
func (p *Printer) SetUnderline(underline int) error {
	// La clase PHP también acepta booleanos y los convierte.
	// En Go, la validación de tipo estática nos da la garantía, así que solo validamos el rango entero.
	if err := validateInteger(underline, UNDERLINE_NONE, UNDERLINE_DOUBLE, "SetUnderline", "subrayado"); err != nil {
		return fmt.Errorf("SetUnderline: %w", err)
	}
	// ESC - n - n=0: ninguno, 1: simple, 2: doble
	cmd := []byte{ESC, '-', byte(underline)}
	_, err := p.connector.Write(cmd)
	return err
}

// SetColor establece el color de impresión (para impresoras con múltiples colores).
// color puede ser COLOR_1 (negro) o COLOR_2 (rojo).
func (p *Printer) SetColor(color int) error {
	if err := validateInteger(color, COLOR_1, COLOR_2, "SetColor", "color"); err != nil {
		return fmt.Errorf("SetColor: %w", err)
	}
	// ESC r n - n=0: Color 1, 1: Color 2
	cmd := []byte{ESC, 'r', byte(color)}
	_, err := p.connector.Write(cmd)
	return err
}

// SetReverseColors habilita o deshabilita la impresión en colores inversos.
func (p *Printer) SetReverseColors(on bool) error {
	// GS B n - n=1: habilitar, 0: deshabilitar
	val := byte(0)
	if on {
		val = 1
	}
	cmd := []byte{GS, 'B', val}
	_, err := p.connector.Write(cmd)
	return err
}

// SetTextSize establece el tamaño del texto usando multiplicadores de ancho y alto (1-8).
func (p *Printer) SetTextSize(widthMultiplier, heightMultiplier int) error {
	if err := validateInteger(widthMultiplier, 1, 8, "SetTextSize", "multiplicador de ancho"); err != nil {
		return fmt.Errorf("SetTextSize: %w", err)
	}
	if err := validateInteger(heightMultiplier, 1, 8, "SetTextSize", "multiplicador de alto"); err != nil {
		return fmt.Errorf("SetTextSize: %w", err)
	}
	// GS ! n - n es una combinación de bits de los multiplicadores (ancho-1) * 16 + (alto-1)
	c := byte(((widthMultiplier - 1) << 4) | (heightMultiplier - 1))
	cmd := []byte{GS, '!', c}
	_, err := p.connector.Write(cmd)
	return err
}

// SetLineSpacing establece el espaciado entre líneas.
// Si height es nil, restablece al espaciado por defecto (ESC 2).
// Si height no es nil, establece el espaciado a height/180 o height/203 pulgadas (ESC 3 n).
func (p *Printer) SetLineSpacing(height *int) error {
	if height == nil {
		// ESC 2 - Restablecer espaciado de línea por defecto
		_, err := p.connector.Write([]byte{ESC, '2'})
		return err
	}
	if err := validateInteger(*height, 1, 255, "SetLineSpacing", "altura"); err != nil {
		return fmt.Errorf("SetLineSpacing: %w", err)
	}
	// ESC 3 n - Establecer espaciado de línea a n
	cmd := []byte{ESC, '3', byte(*height)}
	_, err := p.connector.Write(cmd)
	return err
}

// SetPrintLeftMargin establece el margen izquierdo de impresión en puntos.
func (p *Printer) SetPrintLeftMargin(margin int) error {
	if err := validateInteger(margin, 0, 65535, "SetPrintLeftMargin", "margen"); err != nil {
		return fmt.Errorf("SetPrintLeftMargin: %w", err)
	}
	// GS L nL nH - Establece el margen izquierdo a nL + nH * 256 puntos
	marginBytes, err := intLowHigh(margin, 2) // 2 bytes (nL nH)
	if err != nil {
		return fmt.Errorf("SetPrintLeftMargin: falló al formatear bytes del margen: %w", err)
	}
	cmd := []byte{GS, 'L'}
	cmd = append(cmd, marginBytes...)
	_, err = p.connector.Write(cmd)
	return err
}

// SetPrintWidth establece el ancho del área de impresión en puntos.
func (p *Printer) SetPrintWidth(width int) error {
	if err := validateInteger(width, 1, 65535, "SetPrintWidth", "ancho"); err != nil {
		return fmt.Errorf("SetPrintWidth: %w", err)
	}
	// GS W nL nH - Establece el ancho del área de impresión a nL + nH * 256 puntos
	widthBytes, err := intLowHigh(width, 2) // 2 bytes (nL nH)
	if err != nil {
		return fmt.Errorf("SetPrintWidth: falló al formatear bytes del ancho: %w", err)
	}
	cmd := []byte{GS, 'W'}
	cmd = append(cmd, widthBytes...)
	_, err = p.connector.Write(cmd)
	return err
}

// SetPrintBuffer no se porta directamente ya que el manejo del texto se simplificó.
// La funcionalidad de `PrintBuffer` (manejo de \n y escritura raw)
// está cubierta por `Text` y `TextRaw`.

// SetBarcodeHeight establece la altura del código de barras en puntos.
func (p *Printer) SetBarcodeHeight(height int) error {
	if err := validateInteger(height, 1, 255, "SetBarcodeHeight", "altura"); err != nil {
		return fmt.Errorf("SetBarcodeHeight: %w", err)
	}
	// GS h n - Establece la altura del código de barras a n puntos
	cmd := []byte{GS, 'h', byte(height)}
	_, err := p.connector.Write(cmd)
	return err
}

// SetBarcodeWidth establece el ancho de los módulos del código de barras.
func (p *Printer) SetBarcodeWidth(width int) error {
	if err := validateInteger(width, 1, 255, "SetBarcodeWidth", "ancho"); err != nil {
		return fmt.Errorf("SetBarcodeWidth: %w", err)
	}
	// GS w n - Establece el ancho horizontal de los módulos a n (normalmente 2 o 3)
	cmd := []byte{GS, 'w', byte(width)}
	_, err := p.connector.Write(cmd)
	return err
}

// SetBarcodeTextPosition establece la posición del texto legible del código de barras.
func (p *Printer) SetBarcodeTextPosition(position int) error {
	if err := validateInteger(position, BARCODE_TEXT_NONE, BARCODE_TEXT_BELOW, "SetBarcodeTextPosition", "posición"); err != nil {
		return fmt.Errorf("SetBarcodeTextPosition: %w", err)
	} // 0: ninguno, 1: arriba, 2: abajo, 3: ambos (no siempre soportado) - PHP valida 0-3
	cmd := []byte{GS, 'H', byte(position)}
	_, err := p.connector.Write(cmd)
	return err
}

// Barcode imprime un código de barras.
// content es la cadena de datos del código de barras.
// barType es el tipo de código de barras (BARCODE_UPCA, BARCODE_CODE39, etc.).
func (p *Printer) Barcode(content string, barType int) error {
	if err := validateInteger(barType, BARCODE_UPCA, BARCODE_CODE128, "Barcode", "tipo de código de barras"); err != nil {
		return fmt.Errorf("Barcode: %w", err)
	}
	contentLen := len(content)

	// --- Validación de contenido basada en el tipo de código de barras (traducir regex y longitud) ---
	var validationErr error
	switch barType {
	case BARCODE_UPCA:
		validationErr = validateStringRegex(content, `^[0-9]{11,12}$`, "contenido UPCA")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 11, 12, "Barcode", "longitud contenido UPCA")
		}
	case BARCODE_UPCE:
		validationErr = validateStringRegex(content, `^([0-9]{6,8}|[0-9]{11,12})$`, "contenido UPCE")
		if validationErr == nil {
			validationErr = validateIntegerMulti(contentLen, [][]int{{6, 8}, {11, 12}}, "Barcode", "longitud contenido UPCE")
		}
	case BARCODE_JAN13:
		validationErr = validateStringRegex(content, `^[0-9]{12,13}$`, "contenido JAN13")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 12, 13, "Barcode", "longitud contenido JAN13")
		}
	case BARCODE_JAN8:
		validationErr = validateStringRegex(content, `^[0-9]{7,8}$`, "contenido JAN8")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 7, 8, "Barcode", "longitud contenido JAN8")
		}
	case BARCODE_CODE39:
		// PHP regex: `^([0-9A-Z $%+\-./]+|\*[0-9A-Z $%+\-./]+\*)$`
		// Requiere un * al principio y al final, o no.
		validationErr = validateStringRegex(content, `^([0-9A-Z $%+\-./]+|\*[0-9A-Z $%+\-./]+\*)$`, "contenido CODE39")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE39")
		}
	case BARCODE_ITF:
		// PHP regex: `^([0-9]{2})+$` - requiere solo dígitos y longitud par.
		validationErr = validateStringRegex(content, `^([0-9]{2})+$`, "contenido ITF")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 2, 255, "Barcode", "longitud contenido ITF")
		}
		// Validar longitud par
		if validationErr == nil && contentLen%2 != 0 {
			validationErr = errors.New("la longitud del contenido ITF debe ser par")
		}
	case BARCODE_CODABAR:
		// PHP regex: `^[A-Da-d][0-9$%+\-./:]+[A-Da-d]$` - inicia/termina con A-D, medio con dígitos/símbolos.
		validationErr = validateStringRegex(content, `^[A-Da-d][0-9$%+\-./:]+[A-Da-d]$`, "contenido Codabar")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 1, 255, "Barcode", "longitud contenido Codabar")
		}
	case BARCODE_CODE93:
		// PHP regex: `^[\x00-\x7F]+$` - solo caracteres ASCII.
		validationErr = validateStringRegex(content, `^[\x00-\x7F]+$`, "contenido CODE93")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE93")
		}
	case BARCODE_CODE128:
		// PHP regex: `^\{[A-C][\\x00-\\x7F]+$` - espera que el contenido empiece con {A, {B o {C y luego ASCII.
		// Esto es un poco inusual, ya que normalmente el usuario no proporciona los códigos de inicio/función de Code128.
		// Replicamos la validación de PHP.
		validationErr = validateStringRegex(content, `^\{[A-C][\x00-\x7F]+$`, "contenido CODE128")
		if validationErr == nil {
			validationErr = validateInteger(contentLen, 1, 255, "Barcode", "longitud contenido CODE128")
		}
		if validationErr == nil && contentLen < 2 { // Necesita al menos '{' y un carácter de tipo
			validationErr = errors.New("el contenido Code128 debe tener al menos 2 caracteres ({A, {B, {C...)")
		}
	}
	if validationErr != nil {
		return fmt.Errorf("Barcode: contenido '%s' inválido para el tipo %d: %w", content, barType, validationErr)
	}
	// --- Fin Validación ---

	// Lógica de envío:
	// PHP usa el comando GS k m d1...dk (m=0-6) si getSupportsBarcodeB() es false.
	// PHP usa el comando GS k m L d1...dL (m=65-73) si getSupportsBarcodeB() es true.
	// 'SupportsBarcodeB' en PHP parece referirse al soporte del formato de comando más nuevo (con byte de longitud L).

	cmd := []byte{GS, 'k'}
	if !p.profile.SupportsBarcodeB {
		// Usar el formato de comando antiguo: GS k m data NUL (m = 0-6)
		// Validar que el tipo solicitado esté en el rango 65-71 (correspondiente a m 0-6)
		if barType < BARCODE_UPCA || barType > BARCODE_CODABAR {
			return fmt.Errorf("Barcode: el perfil de impresora no soporta el tipo de código de barras %d con el formato de comando antiguo", barType)
		}
		cmd = append(cmd, byte(barType-65)) // Tipo de 0 a 6
		cmd = append(cmd, []byte(content)...)
		cmd = append(cmd, NUL) // Terminador NUL
	} else {
		// Usar el formato de comando nuevo: GS k m L data (m = 65-73)
		cmd = append(cmd, byte(barType))      // Tipo de 65 a 73
		cmd = append(cmd, byte(contentLen))   // Byte de longitud L
		cmd = append(cmd, []byte(content)...) // Datos
	}

	_, err := p.connector.Write(cmd)
	return err
}

// QrCode imprime un código QR.
func (p *Printer) QrCode(content string, ecLevel int, size int, model int) error {
	if content == "" {
		return nil
	} // No hacer nada si el contenido está vacío, como PHP
	if !p.profile.SupportsQrCode {
		return errors.New("los códigos QR no están soportados en este perfil de impresora")
	}

	if err := validateString(content, "QrCode", "contenido"); err != nil {
		return fmt.Errorf("QrCode: %w", err)
	}
	if err := validateInteger(ecLevel, QR_ECLEVEL_L, QR_ECLEVEL_H, "QrCode", "nivel EC"); err != nil {
		return fmt.Errorf("QrCode: %w", err)
	} // 0-3
	if err := validateInteger(size, 1, 16, "QrCode", "tamaño"); err != nil {
		return fmt.Errorf("QrCode: %w", err)
	} // Tamaño del punto (1-16)
	if err := validateInteger(model, QR_MODEL_1, QR_MICRO, "QrCode", "modelo"); err != nil {
		return fmt.Errorf("QrCode: %w", err)
	} // 1, 2, 3

	cn := byte('1') // Código de función 1A (para QR Code) para GS ( k

	// 1. Establecer modelo: GS ( k pL pH cn 41 n1 n2
	//    cn=1, fn=65 ('A'), n1=49(modelo 1), 50(modelo 2), 51(micro), n2=0
	//    PHP envía chr(48 + model) para n1. Replicamos.
	if err := p.wrapperSend2dCodeData(byte(65), cn, []byte{byte(48 + model), 0}, 0); err != nil {
		return fmt.Errorf("QrCode: falló al establecer el modelo: %w", err)
	}

	// 2. Establecer tamaño del módulo: GS ( k pL pH cn 67 n
	//    cn=1, fn=67 ('C'), n=tamaño (1-16)
	if err := p.wrapperSend2dCodeData(byte(67), cn, []byte{byte(size)}, 0); err != nil {
		return fmt.Errorf("QrCode: falló al establecer el tamaño: %w", err)
	}

	// 3. Establecer nivel EC: GS ( k pL pH cn 69 n
	//    cn=1, fn=69 ('E'), n=nivel EC (0-3)
	//    PHP envía chr(48 + ecLevel). Replicamos.
	if err := p.wrapperSend2dCodeData(byte(69), cn, []byte{byte(48 + ecLevel)}, 0); err != nil {
		return fmt.Errorf("QrCode: falló al establecer el nivel EC: %w", err)
	}

	// 4. Almacenar datos: GS ( k pL pH cn 80 m d1...dk
	//    cn=1, fn=80 ('P'), m='0' (modo de procesamiento), d1...dk=contenido
	if err := p.wrapperSend2dCodeData(byte(80), cn, []byte(content), byte('0')); err != nil {
		return fmt.Errorf("QrCode: falló al almacenar los datos: %w", err)
	}

	// 5. Imprimir símbolo: GS ( k pL pH cn 81 m
	//    cn=1, fn=81 ('Q'), m='0' (modo de impresión)
	if err := p.wrapperSend2dCodeData(byte(81), cn, []byte{}, byte('0')); err != nil { // Sin datos después de '0'
		return fmt.Errorf("QrCode: falló al imprimir el símbolo: %w", err)
	}

	return nil
}

// Pdf417Code imprime un código PDF417.
func (p *Printer) Pdf417Code(content string, width, heightMultiplier, dataColumnCount int, ec float64, options int) error {
	if content == "" {
		return nil
	} // No hacer nada si el contenido está vacío, como PHP
	if !p.profile.SupportsPdf417Code {
		return errors.New("los códigos PDF417 no están soportados en este perfil de impresora")
	}

	if err := validateString(content, "Pdf417Code", "contenido"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	}
	if err := validateInteger(width, 2, 8, "Pdf417Code", "ancho"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // Ancho del módulo (2-8 puntos)
	if err := validateInteger(heightMultiplier, 2, 8, "Pdf417Code", "multiplicador de alto"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // Multiplicador de alto de fila (2-8)
	if err := validateInteger(dataColumnCount, 0, 30, "Pdf417Code", "contador columnas de datos"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // 0 = automático
	if err := validateFloat(ec, 0.01, 4.00, "Pdf417Code", "nivel EC"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // Nivel EC como flotante (ej: 0.10 para 10%)
	if err := validateInteger(options, PDF417_STANDARD, PDF417_TRUNCATED, "Pdf417Code", "opciones"); err != nil {
		return fmt.Errorf("Pdf417Code: %w", err)
	} // 0: estándar, 1: truncado

	cn := byte('0') // Código de función 1A (para PDF417) para GS ( k

	// 1. Establecer opciones (estándar/truncado): GS ( k pL pH cn 70 n
	//    cn=0, fn=70 ('F'), n=opciones (0 o 1)
	if err := p.wrapperSend2dCodeData(byte(70), cn, []byte{byte(options)}, 0); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer las opciones: %w", err)
	}

	// 2. Establecer contador de columnas de datos: GS ( k pL pH cn 65 n
	//    cn=0, fn=65 ('A'), n=contador columnas de datos (0-30)
	if err := p.wrapperSend2dCodeData(byte(65), cn, []byte{byte(dataColumnCount)}, 0); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer el contador de columnas: %w", err)
	}

	// 3. Establecer ancho del módulo: GS ( k pL pH cn 67 n
	//    cn=0, fn=67 ('C'), n=ancho (2-8)
	if err := p.wrapperSend2dCodeData(byte(67), cn, []byte{byte(width)}, 0); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer el ancho: %w", err)
	}

	// 4. Establecer multiplicador de alto de fila: GS ( k pL pH cn 68 n
	//    cn=0, fn=68 ('D'), n=multiplicador de alto (2-8)
	if err := p.wrapperSend2dCodeData(byte(68), cn, []byte{byte(heightMultiplier)}, 0); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer el multiplicador de alto: %w", err)
	}

	// 5. Establecer nivel EC: GS ( k pL pH cn 69 n m
	//    cn=0, fn=69 ('E'), n=ceil(ec * 10), m='1' (modo)
	//    PHP calcula ec_int = ceil(floatval(ec) * 10) y lo envía como byte, con m='1'. Replicamos.
	ecInt := byte(math.Ceil(ec * 10))
	if err := p.wrapperSend2dCodeData(byte(69), cn, []byte{ecInt}, byte('1')); err != nil {
		return fmt.Errorf("Pdf417Code: falló al establecer el nivel EC: %w", err)
	}

	// 6. Almacenar datos: GS ( k pL pH cn 80 m d1...dk
	//    cn=0, fn=80 ('P'), m='0' (modo), d1...dk=contenido
	if err := p.wrapperSend2dCodeData(byte(80), cn, []byte(content), byte('0')); err != nil {
		return fmt.Errorf("Pdf417Code: falló al almacenar los datos: %w", err)
	}

	// 7. Imprimir símbolo: GS ( k pL pH cn 81 m
	//    cn=0, fn=81 ('Q'), m='0' (modo)
	if err := p.wrapperSend2dCodeData(byte(81), cn, []byte{}, byte('0')); err != nil { // Sin datos después de '0'
		return fmt.Errorf("Pdf417Code: falló al imprimir el símbolo: %w", err)
	}

	return nil
}

// Pulse envía un pulso a un pin del conector del cajón portamonedas para abrirlo.
func (p *Printer) Pulse(pin int, onMS, offMS int) error {
	if err := validateInteger(pin, 0, 1, "Pulse", "pin"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Pin 0 o 1
	if err := validateInteger(onMS, 1, 511, "Pulse", "onMS"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Tiempo ON en ms (1-511)
	if err := validateInteger(offMS, 1, 511, "Pulse", "offMS"); err != nil {
		return fmt.Errorf("Pulse: %w", err)
	} // Tiempo OFF en ms (1-511) - a menudo ignorado por la impresora para el segundo pulso

	// Comando: ESC p m t1 t2
	// m: pin del cajón (0 o 1). PHP usa pin + 48 ('0' o '1'). Replicamos.
	// t1: Tiempo ON (t1 * 2 ms). PHP envía on_ms / 2. Replicamos.
	// t2: Tiempo OFF (t2 * 2 ms). PHP envía off_ms / 2. Replicamos.
	cmd := []byte{ESC, 'p', byte(pin + 48), byte(onMS / 2), byte(offMS / 2)}
	_, err := p.connector.Write(cmd)
	return err
}

// SelectCharacterTable selecciona la tabla de caracteres (codepage) a utilizar.
func (p *Printer) SelectCharacterTable(table int) error {
	if err := validateInteger(table, 0, 255, "SelectCharacterTable", "tabla"); err != nil {
		return fmt.Errorf("SelectCharacterTable: %w", err)
	}

	// Verificar si el perfil soporta esta tabla
	if _, ok := p.profile.CodePages[table]; !ok {
		return fmt.Errorf("SelectCharacterTable: la tabla de caracteres %d no está permitida por el perfil de esta impresora", table)
	}

	// La clase PHP elige entre ESC t n (estándar) y ESC GS t n (posiblemente Star)
	// basándose en SupportsStarCommands. Implementamos esta lógica.
	var cmd []byte
	if p.profile.SupportsStarCommands {
		// Este comando es probable que sea específico de Star Micronics.
		cmd = []byte{ESC, GS, 't', byte(table)}
	} else {
		// Comando ESC/POS estándar para seleccionar tabla de caracteres.
		cmd = []byte{ESC, 't', byte(table)}
	}

	_, err := p.connector.Write(cmd)
	if err == nil {
		p.characterTable = table // Actualizar el estado interno si la escritura fue exitosa
	}
	return err
}

// BitImage imprime una imagen utilizando el comando de imagen de bits (GS v 0).
// Requiere que la imagen sea convertible a formato raster de 1 bit.
func (p *Printer) BitImage(img *EscposImage, size int) error {
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
func (p *Printer) BitImageColumnFormat(img *EscposImage, size int) error {
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
func (p *Printer) Graphics(img *EscposImage, size int) error {
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

// Close finaliza la conexión con la impresora.
func (p *Printer) Close() error {
	return p.connector.Close()
}

// GetCharacterTable devuelve la tabla de caracteres (codepage) actualmente seleccionada.
func (p *Printer) GetCharacterTable() int {
	return p.characterTable
}

// GetPrintConnector devuelve el conector que está utilizando la impresora.
func (p *Printer) GetPrintConnector() PrintConnector {
	return p.connector
}

// GetPrinterCapabilityProfile devuelve el perfil de capacidad de la impresora.
func (p *Printer) GetPrinterCapabilityProfile() *CapabilityProfile {
	return p.profile
}

// --- Métodos de Ayuda Internos (Funciones/Métodos Privados en Go) ---

// wrapperSend2dCodeData envía una parte de un comando de código 2D (GS ( k ...).
// fn y cn son bytes de función y código.
// data son los bytes de datos.
// m es un byte de modo opcional ('0' o '1' para algunas funciones).
func (p *Printer) wrapperSend2dCodeData(fn, cn byte, data []byte, m byte) error {
	// Formato del comando: GS ( k pL pH cn fn [m] d1...dk
	// pL pH: longitud del payload que sigue (cn + fn + [m] + data)
	// cn: código del símbolo (0 para PDF417, 1 para QR)
	// fn: código de función (ej: 65='A' para configurar, 80='P' para almacenar datos, 81='Q' para imprimir)
	// m: byte de modo opcional
	// d1...dk: datos específicos de la función

	payloadLen := 2 // cn (1 byte) + fn (1 byte)
	if m != 0 {     // Si m se proporciona (distinto de cero)
		payloadLen += 1 // + m (1 byte)
	}
	payloadLen += len(data) // + longitud de los datos

	// pL pH es la longitud total del payload en formato low-high (2 bytes)
	headerBytes, err := intLowHigh(payloadLen, 2)
	if err != nil {
		return fmt.Errorf("wrapperSend2dCodeData: falló al crear la cabecera de longitud: %w", err)
	}

	// Construir el comando completo
	var cmd bytes.Buffer
	cmd.Write([]byte{GS, '(', 'k'}) // Prefijo
	cmd.Write(headerBytes)          // pL pH
	cmd.WriteByte(cn)               // cn
	cmd.WriteByte(fn)               // fn
	if m != 0 {
		cmd.WriteByte(m) // [m] opcional
	}
	cmd.Write(data) // d1...dk

	_, err = p.connector.Write(cmd.Bytes())
	return err
}

// wrapperSendGraphicsData envía una parte de un comando gráfico (GS ( L ...).
// m y fn son bytes de modo y función.
// data son los bytes de datos.
func (p *Printer) wrapperSendGraphicsData(m, fn byte, data []byte) error {
	// Formato del comando: GS ( L pL pH m fn [data]
	// pL pH: longitud del payload que sigue (m + fn + data)
	// m: byte de modo ('0' para este conjunto de comandos gráficos)
	// fn: código de función (ej: 'p' para enviar datos, '2' para imprimir)
	// data: datos gráficos

	payloadLen := 2 + len(data) // m (1 byte) + fn (1 byte) + longitud de los datos

	// pL pH es la longitud total del payload en formato low-high (2 bytes)
	headerBytes, err := intLowHigh(payloadLen, 2)
	if err != nil {
		return fmt.Errorf("wrapperSendGraphicsData: falló al crear la cabecera de longitud: %w", err)
	}

	// Construir el comando completo
	var cmd bytes.Buffer
	cmd.Write([]byte{GS, '(', 'L'}) // Prefijo
	cmd.Write(headerBytes)          // pL pH
	cmd.WriteByte(m)                // m
	cmd.WriteByte(fn)               // fn
	cmd.Write(data)                 // [data]

	_, err = p.connector.Write(cmd.Bytes())
	return err
}

// dataHeader formatea enteros de entrada en bytes (bajo/alto o byte único).
// Se utiliza para formatear dimensiones en comandos de imagen.
// long=true: formatar como 2 bytes (nL nH)
// long=false: formatar como 1 byte
func dataHeader(inputs []int, long bool) ([]byte, error) {
	var outp bytes.Buffer
	for _, input := range inputs {
		if long {
			// Formato de 2 bytes (nL nH) - rango 0 a 65535
			bytes, err := intLowHigh(input, 2)
			if err != nil {
				return nil, fmt.Errorf("dataHeader: falló al formatear el entero %d como 2 bytes: %w", input, err)
			}
			outp.Write(bytes)
		} else {
			// Formato de 1 byte - rango 0 a 255
			if input < 0 || input > 255 {
				return nil, fmt.Errorf("dataHeader: el entero %d está fuera del rango para un byte único (0-255)", input)
			}
			outp.WriteByte(byte(input))
		}
	}
	return outp.Bytes(), nil
}

// intLowHigh convierte un entero en un slice de bytes en orden bajo-alto (Little Endian).
// input es el entero a convertir.
// length es el número de bytes deseado (1 a 4).
func intLowHigh(input, length int) ([]byte, error) {
	if length < 1 || length > 4 { // PHP limita a 1-4, nos ceñimos a eso.
		return nil, fmt.Errorf("intLowHigh: la longitud debe estar entre 1 y 4, se recibió %d", length)
	}

	// El rango máximo para `length` bytes es 2^(length*8) - 1.
	// PHP usa (256 << (length*8)) - 1. Para length=1, (256 << 8) - 1 = 2^8 - 1 = 255.
	// Para length=2, (256 << 16) - 1 = 2^16 - 1 = 65535.
	// Para length=4, (256 << 32) - 1 - esto desborda int en PHP.
	// Usemos uint32 para la comparación para manejar hasta 4 bytes correctamente.
	var maxInput uint32
	if length == 4 {
		maxInput = math.MaxUint32 // 2^32 - 1
	} else {
		maxInput = (uint32(1) << uint(length*8)) - 1
	}

	if input < 0 || uint32(input) > maxInput {
		return nil, fmt.Errorf("intLowHigh: la entrada %d está fuera del rango para %d bytes (0-%d)", input, length, maxInput)
	}

	outp := make([]byte, length)
	// Usar encoding/binary para asegurar el orden Little Endian
	// Convertimos el int a uint32 para usar PutUint32
	binary.LittleEndian.PutUint32(outp, uint32(input))

	// Si la longitud es menor a 4, solo tomamos los primeros `length` bytes
	return outp[:length], nil
}

// --- Funciones de Ayuda para Validación ---
// Estas funciones validan los argumentos y devuelven un error si son inválidos.

// validateBoolean es en gran parte redundante en Go debido al tipado estático.
func validateBoolean(test bool, source string) error {
	// En Go, un bool siempre es true o false. No hay necesidad de verificar esto.
	// La función se mantiene por completitud del port, pero siempre devuelve nil.
	return nil
}

func validateFloat(test float64, min, max float64, source, argument string) error {
	if test < min || test > max {
		return fmt.Errorf("el argumento '%s' (%f) dado a %s debe estar en el rango %f a %f", argument, test, source, min, max)
	}
	return nil
}

func validateInteger(test, min, max int, source, argument string) error {
	return validateIntegerMulti(test, [][]int{{min, max}}, source, argument)
}

func validateIntegerMulti(test int, ranges [][]int, source, argument string) error {
	match := false
	for _, r := range ranges {
		if len(r) != 2 {
			// Esto indica un error interno en cómo se llama a esta función de validación
			return fmt.Errorf("error interno: validateIntegerMulti recibió un rango inválido %v", r)
		}
		if test >= r[0] && test <= r[1] {
			match = true
			break
		}
	}

	if !match {
		// Construir el mensaje de rango similar a PHP
		rangeStrs := make([]string, len(ranges))
		for i, r := range ranges {
			rangeStrs[i] = fmt.Sprintf("%d-%d", r[0], r[1])
		}
		rangeStr := strings.Join(rangeStrs, ", ")
		if len(ranges) > 1 {
			// Reemplazar la última coma con " o " si hay más de un rango
			lastCommaIndex := strings.LastIndex(rangeStr, ", ")
			if lastCommaIndex != -1 {
				rangeStr = rangeStr[:lastCommaIndex+2] + "o " + rangeStr[lastCommaIndex+2:]
			}
		}

		return fmt.Errorf("el argumento '%s' (%d) dado a %s debe estar en el rango %s", argument, test, source, rangeStr)
	}
	return nil
}

// validateString es en gran parte redundante en Go debido al tipado estático.
// El chequeo de PHP sobre objetos con __toString no aplica directamente en Go.
func validateString(test string, source, argument string) error {
	// En Go, el tipado estático ya asegura que es una cadena si el argumento es string.
	// La función se mantiene por completitud del port, pero siempre devuelve nil.
	return nil
}

// Cache para expresiones regulares compiladas
var regexCache = make(map[string]*regexp.Regexp)

func validateStringRegex(test string, regexPattern string, argument string) error {
	// Compilar la regex si no está en caché
	re, ok := regexCache[regexPattern]
	if !ok {
		var err error
		re, err = regexp.Compile(regexPattern)
		if err != nil {
			// Error interno: la regex proporcionada no es válida
			return fmt.Errorf("error interno: falló al compilar la regex '%s': %w", regexPattern, err)
		}
		regexCache[regexPattern] = re
	}

	if !re.MatchString(test) {
		// El mensaje de error de PHP incluía el nombre de la función fuente,
		// pero aquí el argumento 'argument' ya describe qué valor es.
		return fmt.Errorf("el argumento '%s' ('%s') es inválido. Debe coincidir con la regex '%s'", argument, test, regexPattern)
	}
	return nil
}

// intPtr es una función de ayuda para obtener un puntero a un int.
// Útil para métodos con parámetros opcionales *int (como SetLineSpacing).
func intPtr(i int) *int {
	return &i
}
