package utils

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	Uint32Size     = 4 // Tamaño de uint32 en bytes
	DimensionBytes = 2
)

// IntLowHigh convierte un entero en un slice de bytes en orden bajo-alto (Little Endian).
// input es el entero a convertir.
// length es el número de bytes deseado (1 a 4).
func IntLowHigh(input int, length int) ([]byte, error) {
	if length < 1 || length > 4 { // PHP limita a 1-4, nos ceñimos a eso.
		return nil, fmt.Errorf("IntLowHigh: la longitud debe estar entre 1 y 4, se recibió %d", length)
	}

	// Validar que la entrada no sea negativa antes de convertir a uint32
	if input < 0 {
		return nil, fmt.Errorf("IntLowHigh: la entrada %d no puede ser negativa", input)
	}

	// El rango máximo para `length` bytes es 2^(length*8) - 1.
	var maxInput uint32
	if length == 4 {
		maxInput = math.MaxUint32 // 2^32 - 1
	} else {
		// Usar valores precalculados en lugar de shift para evitar conversión insegura
		switch length {
		case 1:
			maxInput = 255 // 2^8 - 1
		case 2:
			maxInput = 65535 // 2^16 - 1
		case 3:
			maxInput = 16777215 // 2^24 - 1
		}
	}

	// Verificar que input esté en el rango correcto usando uint64 para evitar overflow
	if uint64(input) > uint64(maxInput) {
		return nil, fmt.Errorf("IntLowHigh: la entrada %d está fuera del rango para %d bytes (0-%d)", input, length, maxInput)
	}

	// Ahora es seguro convertir input a uint32
	//nolint:gosec // Seguro porque ya verificamos el rango en líneas anteriores
	inputUint32 := uint32(input)

	buf := make([]byte, Uint32Size)
	// Usar encoding/binary para asegurar el orden Little Endian
	binary.LittleEndian.PutUint32(buf, inputUint32)

	// Si la longitud es menor a 4, solo tomamos los primeros `length` bytes
	return buf[:length], nil
}
