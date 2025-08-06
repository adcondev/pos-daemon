package escpos

import (
	"log"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol"
	"strings"
)

// ESCPOSProtocol implementa Protocol para ESC/POS
type ESCPOSProtocol struct {
	// TODO: Mover aquí las propiedades que necesites del ESCPrinter actual
	// Por ejemplo:
	characterTable int
	capabilities   map[string]bool
	paperWidth     int // Ancho de papel en mm
	dpi            int // Dots per inch (DPI) de la impresora
	// NO incluir el conector aquí, eso va en la impresora
}

// NewESCPOSProtocol crea una nueva instancia del protocolo ESC/POS
func NewESCPOSProtocol() protocol.Protocol {
	p := &ESCPOSProtocol{}
	return p
}

// === Implementación de la interfaz Protocol ===

// Initialize genera el comando de inicialización ESC/POS
func (p *ESCPOSProtocol) Initialize() []byte {
	// ESC @ - Reset printer
	return []byte{ESC, '@'}
}

// Close genera comandos de cierre (si los hay)
func (p *ESCPOSProtocol) Close() []byte {
	// ESC/POS no tiene un comando específico de cierre
	// pero podrías incluir un reset o feed final
	return []byte{}
}

// SetJustification convierte el tipo genérico al específico de ESC/POS
func (p *ESCPOSProtocol) SetJustification(justification command.Alignment) []byte {
	// Mapear el tipo genérico al valor ESC/POS
	var escposValue byte
	switch justification {
	case command.AlignLeft:
		escposValue = 0 // ESC/POS: 0 = left
	case command.AlignCenter:
		escposValue = 1 // ESC/POS: 1 = center
	case command.AlignRight:
		escposValue = 2 // ESC/POS: 2 = right
	default:
		escposValue = 0 // Default to left
	}

	// ESC a n
	return []byte{ESC, 'a', escposValue}
}

// SetFont mapea fuentes genéricas a ESC/POS
func (p *ESCPOSProtocol) SetFont(font command.Font) []byte {
	var fontValue byte
	switch font {
	case command.FontA:
		fontValue = 0
	case command.FontB:
		fontValue = 1
	case command.FontC:
		fontValue = 2
	default:
		fontValue = 0
	}

	// ESC M n
	return []byte{ESC, 'M', fontValue}
}

// SetEmphasis activa/desactiva negrita
func (p *ESCPOSProtocol) SetEmphasis(on bool) []byte {
	val := byte(0)
	if on {
		val = 1
	}
	// ESC E n
	return []byte{ESC, 'E', val}
}

// SetDoubleStrike activa/desactiva doble golpe
func (p *ESCPOSProtocol) SetDoubleStrike(on bool) []byte {
	val := byte(0)
	if on {
		val = 1
	}
	// ESC G n
	return []byte{ESC, 'G', val}
}

// SetUnderline configura el subrayado
func (p *ESCPOSProtocol) SetUnderline(underline command.UnderlineMode) []byte {
	var val byte
	switch underline {
	case command.UnderlineNone:
		val = 0
	case command.UnderlineSingle:
		val = 1
	case command.UnderlineDouble:
		val = 2
	default:
		val = 0
	}
	// ESC - n
	return []byte{ESC, '-', val}
}

// Text convierte texto a bytes con encoding apropiado
func (p *ESCPOSProtocol) Text(str string) []byte {
	cmd := strings.ReplaceAll(str, "\n", string(LF))
	return []byte(cmd)
}

// TextLn agrega un salto de línea al final
func (p *ESCPOSProtocol) TextLn(str string) []byte {
	text := p.Text(str)
	// Agregar LF al final
	return append(text, LF)
}

// TextRaw envía bytes sin procesar
func (p *ESCPOSProtocol) TextRaw(str string) []byte {
	return []byte(str)
}

// Cut genera comando de corte
func (p *ESCPOSProtocol) Cut(mode command.CutMode, lines int) []byte {
	// TODO: Implementar validación de lines

	cmd := []byte{GS, 'V'}

	switch mode {
	case command.CutFeed:
		cmd = append(cmd, byte(66), byte(lines)) // o 48 ('0') según el modelo
	case command.Cut:
		cmd = append(cmd, byte(49)) // o 49 ('1') según el modelo
	default:
		cmd = append(cmd, 0)
	}

	return cmd
}

// Feed genera comando de alimentación de papel
func (p *ESCPOSProtocol) Feed(lines int) []byte {
	// TODO: Validar que lines esté en rango válido
	if lines <= 0 {
		return []byte{}
	}

	// ESC d n - Print and feed n lines
	return []byte{ESC, 'd', byte(lines)}
}

// TODO: Implementar el resto de métodos de la interfaz
// Por ahora, implementaciones stub para compilar:

func (p *ESCPOSProtocol) SetTextSize(widthMultiplier, heightMultiplier int) []byte {
	// TODO: Implementar usando GS ! n
	// Hint: n = (widthMultiplier-1)<<4 | (heightMultiplier-1)
	return []byte{}
}

func (p *ESCPOSProtocol) SetLineSpacing(height *int) []byte {
	// TODO: Si height es nil, usar ESC 2 (default)
	// Si no, usar ESC 3 n
	return []byte{}
}

func (p *ESCPOSProtocol) SetPrintLeftMargin(margin int) []byte {
	// TODO: Implementar usando GS L nL nH
	return []byte{}
}

func (p *ESCPOSProtocol) SetPrintWidth(width int) []byte {
	// TODO: Implementar usando GS W nL nH
	return []byte{}
}

func (p *ESCPOSProtocol) SelectCharacterTable(table int) []byte {
	// Validar que table esté en un rango válido
	if (table < 0 || table > 21) || (table > 10 && table < 16) {
		log.Printf("advertencia: tabla de caracteres %d fuera de rango, usando 0 por defecto", table)
		table = 0 // Default a 0 si está fuera de rango
	}
	// ESC t n - Select character code table
	return []byte{ESC, 't', byte(table)}
}

func (p *ESCPOSProtocol) GetCharacterTable() int {
	return p.characterTable
}

func (p *ESCPOSProtocol) SetBarcodeHeight(height int) []byte {
	// TODO: Implementar GS h n
	return []byte{}
}

func (p *ESCPOSProtocol) SetBarcodeWidth(width int) []byte {
	// TODO: Implementar GS w n
	return []byte{}
}

func (p *ESCPOSProtocol) SetBarcodeTextPosition(position command.BarcodeTextPosition) []byte {
	// TODO: Mapear position a valores ESC/POS y usar GS H n
	return []byte{}
}

func (p *ESCPOSProtocol) Barcode(content string, barType command.BarcodeType) ([]byte, error) {
	// TODO: Esta es la más compleja, necesitas:
	// 1. Mapear barType genérico a tipo ESC/POS
	// 2. Validar content según el tipo
	// 3. Generar comando según p.capabilities["barcode_b"]
	return []byte{}, nil
}

func (p *ESCPOSProtocol) FeedReverse(lines int) []byte {
	// TODO: Implementar ESC e n
	return []byte{}
}

func (p *ESCPOSProtocol) FeedForm() []byte {
	// TODO: Implementar FF
	return []byte{}
}

func (p *ESCPOSProtocol) Pulse(pin int, onMS int, offMS int) []byte {
	// TODO: Implementar ESC p m t1 t2
	return []byte{}
}

func (p *ESCPOSProtocol) Release() []byte {
	// TODO: Implementar si es necesario
	return []byte{}
}

func (p *ESCPOSProtocol) Name() string {
	return "ESC/POS"
}

func (p *ESCPOSProtocol) HasCapability(cap string) bool {
	return p.capabilities[cap]
}

// HasNativeImageSupport indica si este protocolo soporta imágenes nativas
func (p *ESCPOSProtocol) HasNativeImageSupport() bool {
	return true // ESC/POS soporta imágenes de forma nativa
}

// GetMaxImageWidth devuelve el ancho máximo de imagen que soporta la impresora
func (p *ESCPOSProtocol) GetMaxImageWidth() int {
	// Cálculo basado en el ancho del papel y resolución
	// Formula: (ancho_papel_mm / 25.4) * dpi
	if p.paperWidth > 0 && p.dpi > 0 {
		return int((float64(p.paperWidth) / 25.4) * float64(p.dpi))
	}

	// Valores predeterminados si no hay configuración
	if p.paperWidth >= 80 {
		return 576 // Para papel de 80mm a 203dpi
// PrintQR implementa el comando ESC Z para imprimir códigos QR
func (p *ESCPOSProtocol) PrintQR(
	data string,
	model command.QRModel,
	moduleSize command.QRModuleSize,
	ecLevel command.QRErrorCorrection,
) ([][]byte, error) {
	// Validación de modelo
	if model < command.Model1 || model > command.Model2 {
		return nil, fmt.Errorf("modelo de QR inválida(0-1): %d", model)
	}

	// Comando para seleccionar tamaño del módulo
	mdl, err := p.SelectQRModel(model)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar modelo de QR: %w", err)
	}

	// Comando para seleccionar tamaño del módulo
	mdlSz, err := p.SelectQRSize(moduleSize)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar tamaño de módulo de QR: %w", err)
	}

	// Obtener el byte correspondiente al nivel de corrección
	ec, err := p.SelectQRErrorCorrection(ecLevel)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar nivel de corrección de QR: %w", err)
	}

	// Almacenamiento de datos para QR
	ct, err := p.SetQRData(data)
	if err != nil {
		return nil, fmt.Errorf("error al preparar datos de QR: %w", err)
	}

	// Comando para imprimir QR
	prnt, err := p.PrintQRData()
	if err != nil {
		return nil, fmt.Errorf("error al generar comando de impresión de QR: %w", err)
	}

	cmdLines := [][]byte{mdl, mdlSz, ec, ct, prnt}
	if len(cmdLines) == 0 {
		return nil, fmt.Errorf("no se generaron comandos para imprimir QR")
	}

	return cmdLines, nil
}

var modelMap = map[command.QRModel]byte{
	command.Model1: '1', // Modelo 1
	command.Model2: '2', // Modelo 2
}

func (p *ESCPOSProtocol) SelectQRModel(model command.QRModel) ([]byte, error) {
	// Validación de modelo
	if model < command.Model1 || model > command.Model2 {
		return nil, fmt.Errorf("modelo de QR inválida(0-1): %d", model)
	}

	pL, pH, err := utils.LengthLowHigh(4)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('A')
	n1 := modelMap[model]
	n2 := byte(0) // Siempre 0, reservado

	cmd := make([]byte, 0, 9)
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, n1, n2)

	return cmd, nil
}

func (p *ESCPOSProtocol) SelectQRSize(moduleSize command.QRModuleSize) ([]byte, error) {
	// Validar tamaño del módulo
	if moduleSize < command.MinType || moduleSize > command.MaxType {
		return nil, fmt.Errorf("tamaño de módulo QR inválido(1-16): %d", moduleSize)
	}

	pL, pH, err := utils.LengthLowHigh(3)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('C')
	n := byte(moduleSize)

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, n)

	return cmd, nil
}

// Mapear los niveles de corrección de errores en QR a sus valores ESCPOS
var ecMap = map[command.QRErrorCorrection]byte{
	command.ECLow:     '0', // 7% de corrección
	command.ECMedium:  '1', // 15% de corrección
	command.ECHigh:    '2', // 25% de corrección
	command.ECHighest: '3', // 30% de corrección
}

func (p *ESCPOSProtocol) SelectQRErrorCorrection(level command.QRErrorCorrection) ([]byte, error) {
	// Validar nivel de corrección
	ec, ok := ecMap[level]
	if !ok {
		return nil, fmt.Errorf("nivel de corrección de QR inválido(0-3): %d", level)
	}

	pL, pH, err := utils.LengthLowHigh(3)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('E')

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, ec)

	return cmd, nil
}

func (p *ESCPOSProtocol) SetQRData(data string) ([]byte, error) {
	// Validar longitud de datos
	if len(data) == 0 || len(data) > 7089 {
		return nil, fmt.Errorf("longitud de datos de QR inválida (1-7089): %d", len(data))
	}

	pL, pH, err := utils.LengthLowHigh(len(data) + 3)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('P')
	m := byte('0') // Siempre 0, reservado

	cmd := make([]byte, 0, 7+len(data))
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, m)
	cmd = append(cmd, data...)

	return cmd, nil
}

func (p *ESCPOSProtocol) PrintQRData() ([]byte, error) {
	// Comando para imprimir QR
	pL, pH, err := utils.LengthLowHigh(3)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('Q')
	m := byte('0') // Siempre 0 para impresion estandard

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, m)

	return cmd, nil
}

// Implementar los métodos restantes de la interfaz Protocol según sea necesario
