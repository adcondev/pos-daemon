package escpos

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

const (
	// Tipos de estado (para comandos de estado, no implementados como métodos públicos en PHP)
	// Se incluyen por completitud de las constantes PHP
	STATUS_PRINTER       int = 1 // GS I 1 (Estado de la impresora)
	STATUS_OFFLINE_CAUSE int = 2 // GS I 2 (Causa de estar offline)
	STATUS_ERROR_CAUSE   int = 3 // GS I 3 (Causa del error)
	STATUS_PAPER_ROLL    int = 4 // GS I 4 (Estado del rollo de papel)
	STATUS_INK_A         int = 7 // GS I 7 (Estado de la tinta/cinta A)
	STATUS_INK_B         int = 6 // GS I 6 (Estado de la tinta/cinta B)
	STATUS_PEELER        int = 8 // GS I 8 (Estado del peeler - para etiquetas)

	// Modo de corte de papel
	CUT_FULL    int = 65 // 'A'
	CUT_PARTIAL int = 66 // 'B'

	// Opciones de PDF417
	PDF417_STANDARD  int = 0
	PDF417_TRUNCATED int = 1
)

// PrintConnector define la interfaz para la conexión física con la impresora.
// Debes implementar esta interfaz para tu método de conexión (USB, TCP, Serial, etc.).
type PrintConnector interface {
	// Write envía bytes a la impresora.
	Write([]byte) (int, error)
	// Close finaliza la conexión con la impresora.
	Close() error
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
			data, err := intLowHigh(input, 2)
			if err != nil {
				return nil, fmt.Errorf("dataHeader: falló al formatear el entero %d como 2 bytes: %w", input, err)
			}
			outp.Write(data)
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

// Image La implementación real para cargar y convertir imágenes (ToRasterFormat, ToColumnFormat)
// debe ser proporcionada. Esto implica manipulación de píxeles y formatos específicos de ESC/POS.
type Image struct {
	img              image.Image
	threshold        uint8
	width            int
	height           int
	rasterData       []byte
	columnFormatHigh [][]byte
	columnFormatLow  [][]byte
}

func NewEscposImage(img image.Image, threshold uint8) *Image {
	bounds := img.Bounds()
	return &Image{
		img:       img,
		threshold: threshold,
		width:     bounds.Dx(),
		height:    bounds.Dy(),
	}
}

func NewEscposImageFromFile(filename string, threshold uint8) (*Image, error) {
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

func NewEscposImageFromReader(reader io.Reader, threshold uint8) (*Image, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decodificando la image: %w", err)
	}

	return NewEscposImage(img, threshold), nil
}

func NewEscposImageFromBytes(data []byte, threshold uint8) (*Image, error) {
	return NewEscposImageFromReader(bytes.NewReader(data), threshold)
}

func (ei *Image) GetWidth() int {
	return ei.width
}

func (ei *Image) GetHeight() int {
	return ei.height
}

func (ei *Image) GetWidthBytes() int {
	return (ei.width + 7) / 8
}

func (ei *Image) ToRasterFormat() ([]byte, error) {
	if ei.rasterData == nil {
		if err := ei.processRasterData(); err != nil {
			return nil, err
		}
	}
	return ei.rasterData, nil
}

func (ei *Image) ToColumnFormat(highDensity bool) ([][]byte, error) {
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

func (ei *Image) processRasterData() error {
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

func (ei *Image) processColumnData(highDensity bool) error {
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

func (ei *Image) isBlack(x, y int) bool {
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
