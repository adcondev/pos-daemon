package utils

// IntPtr es una función de ayuda para obtener un puntero a un int.
// Útil para métodos con parámetros opcionales *int (como SetLineSpacing).
func IntPtr(i int) *int {
	return &i
}
