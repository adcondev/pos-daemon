package escpos

import (
	"strings"
)

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
