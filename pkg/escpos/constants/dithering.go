package constants

type (
	DitherMode int
)

const (
	None       DitherMode = 0 // Sin dithering
	FloydStein DitherMode = 1 // Dithering Floyd-Steinberg
	Ordered    DitherMode = 2 // Dithering ordenado (matriz 4x4)

	// Threshold

	DefaultThreshold = 128
)
