package escpos

import (
	"errors"
	"fmt"
)

const (
	// Niveles de corrección de error QR (aproximados)
	QR_ECLEVEL_L int = 0 // 7%
	QR_ECLEVEL_M int = 1 // 15%
	QR_ECLEVEL_Q int = 2 // 25%
	QR_ECLEVEL_H int = 3 // 30%

	// Modelos de QR
	QR_MODEL_1 int = 1
	QR_MODEL_2 int = 2
	QR_MICRO   int = 3
)

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
