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

// ntLowHigh convierte un entero en un slice de bytes en orden bajo-alto (Little Endian).
// input es el entero a convertir.
// length es el número de bytes deseado (1 a 4).
func IntLowHigh(input, length int) ([]byte, error) {
	if length < 1 || length > 4 { // PHP limita a 1-4, nos ceñimos a eso.
		return nil, fmt.Errorf("IntLowHigh: la longitud debe estar entre 1 y 4, se recibió %d", length)
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
		return nil, fmt.Errorf("IntLowHigh: la entrada %d está fuera del rango para %d bytes (0-%d)", input, length, maxInput)
	}

	buf := make([]byte, Uint32Size)
	// Usar encoding/binary para asegurar el orden Little Endian
	// Convertimos el int a uint32 para usar PutUint32
	binary.LittleEndian.PutUint32(buf, uint32(input))

	// Si la longitud es menor a 4, solo tomamos los primeros `length` bytes
	return buf[:length], nil
}
