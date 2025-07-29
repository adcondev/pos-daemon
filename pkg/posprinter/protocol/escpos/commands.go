package escpos

// Constantes ESC/POS y parámetros.
const (
	// HT representa el comando Tabulación Horizontal.
	//
	// Nombre:
	//   Tabulación horizontal
	//
	// Formato:
	//   ASCII: HT
	//   Hex: 09
	//   Decimal: 9
	//
	// Descripción:
	//   Mueve la posición de impresión a la siguiente posición de tabulación horizontal.
	//
	// Detalles:
	//   - Este comando se ignora a menos que se haya configurado la siguiente posición de tabulación horizontal.
	//   - Si la siguiente posición de tabulación horizontal excede el área de impresión, la impresora establece la posición de impresión en [Ancho del área de impresión + 1].
	//   - Las posiciones de tabulación horizontal se configuran con ESC D.
	//   - Si este comando se recibe cuando la posición de impresión está en [Ancho del área de impresión + 1], la impresora ejecuta la impresión del búfer lleno de la línea actual y procesa la tabulación horizontal desde el inicio de la siguiente línea.
	//   - La configuración predeterminada de la posición de tabulación horizontal para el rollo de papel es la fuente A (12 × 24) cada 8 caracteres (9°, 17°, 25°, ... columna).
	//
	// Referencia:
	//   ESC D
	HT byte = 0x09 // Tabulación Horizontal

	// LF representa el comando Imprimir y Alimentar Línea.
	//
	// Nombre:
	//   Imprimir y alimentar línea
	//
	// Formato:
	//   ASCII: LF
	//   Hex: 0A
	//   Decimal: 10
	//
	// Descripción:
	//   Imprime los datos en el búfer de impresión y alimenta una línea según el espaciado de línea actual.
	//
	// Nota:
	//   Este comando establece la posición de impresión al inicio de la línea.
	//
	// Referencia:
	//   ESC 2, ESC 3
	LF byte = 0x0A // Imprimir y Alimentar Línea

	// FF representa el comando Imprimir y Regresar al Modo Estándar en Modo Página.
	//
	// Nombre:
	//   Imprimir y regresar al modo estándar en modo página
	//
	// Formato:
	//   ASCII: FF
	//   Hex: 0C
	//   Decimal: 12
	//
	// Descripción:
	//   Imprime los datos en el búfer de impresión de forma colectiva y regresa al modo estándar.
	//
	// Detalles:
	//   - Los datos del búfer se eliminan después de ser impresos.
	//   - El área de impresión configurada por ESC W se restablece a la configuración predeterminada.
	//   - La impresora no ejecuta el corte de papel.
	//   - Este comando establece la posición de impresión al inicio de la línea.
	//   - Este comando está habilitado únicamente en modo página.
	//
	// Referencia:
	//   ESC FF, ESC L, ESC S
	FF byte = 0x0C // Imprimir y Regresar al Modo Estándar en Modo Página

	// CAN representa el comando Cancelar datos de impresión en modo página.
	//
	// Nombre:
	//   Cancelar datos de impresión en modo página
	//
	// Formato:
	//   ASCII: CAN
	//   Hex: 18
	//   Decimal: 24
	//
	// Descripción:
	//   En modo página, elimina todos los datos de impresión en el área imprimible actual.
	//
	// Detalles:
	//   - Este comando está habilitado únicamente en modo página.
	//   - Si existen datos en el área de impresión previamente especificada que también están en el área de impresión actualmente especificada, estos se eliminan.
	//
	// Referencia:
	//   ESC L, ESC W
	CAN byte = 0x18 // Cancelar datos de impresión en modo página

	// ESC representa el carácter de control Escape (0x1B), utilizado como prefijo
	// para la mayoría de los comandos de control de la impresora.
	//
	// Nombre:
	//   Escape (prefijo de comando)
	//
	// Formato:
	//   ASCII: ESC
	//   Hex: 1B
	//   Decimal: 27
	//
	// Descripción:
	//   ESC es un carácter de escape que precede a muchos comandos ESC/POS.
	//   Indica que el siguiente byte (o secuencia) representa una instrucción
	//   de control para la impresora.
	//
	// Detalles:
	//   - La mayoría de los comandos de formato, espaciado, inicialización y alineación
	//     comienzan con este byte.
	//   - Es obligatorio para interpretar correctamente comandos como ESC @, ESC a n, etc.
	//   - No tiene efecto por sí solo; siempre debe ir seguido de un comando válido.
	//
	// Referencia:
	//   ESC
	ESC byte = 0x1B // Escape (prefijo de comando)

	SP  byte = 0x20 // Espacio (carácter de espacio en blanco)
	FS  byte = 0x1C
	GS  byte = 0x1D
	NUL byte = 0x00
)

var (
	// PrintDataPageMode representa el comando para imprimir datos en modo página.
	//
	// Nombre:
	//   Imprimir datos en modo página
	//
	// Formato:
	//   ASCII: ESC FF
	//   Hex: 1B 0C
	//   Decimal: 27 12
	//
	// Descripción:
	//   En modo página, imprime todos los datos almacenados en el búfer dentro del área de impresión de manera colectiva.
	//
	// Detalles:
	//   - Este comando está habilitado únicamente en modo página.
	//   - Después de imprimir, la impresora no borra los datos almacenados en el búfer, los valores configurados para ESC T y ESC W, ni la posición para almacenar datos de caracteres.
	//
	// Referencia:
	//   FF, ESC L, ESC S
	PrintDataPageMode = []byte{ESC, FF}
)

// SetRightSideCharSpacing representa el comando para configurar el espaciado a la derecha de los caracteres.
//
// Nombre:
//
//	Configurar espaciado a la derecha de los caracteres
//
// Formato:
//
//	ASCII: ESC SP n
//	Hex: 1B 20 n
//	Decimal: 27 32 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Configura el espaciado de caracteres en el lado derecho del carácter a [n unidades de movimiento horizontal o vertical].
//
// Detalles:
//   - El espaciado a la derecha de los caracteres en modo de doble ancho es el doble del valor normal. Cuando los caracteres se agrandan, el espaciado a la derecha es n veces el valor normal.
//   - Este comando no afecta la configuración de caracteres kanji.
//   - Este comando establece valores de manera independiente en cada modo (modo estándar y modo página).
//   - Las unidades de movimiento horizontal y vertical se especifican mediante el comando GS P. Cambiar la unidad de movimiento horizontal o vertical no afecta el espaciado actual del lado derecho.
//   - El comando GS P puede cambiar la unidad de movimiento horizontal (y vertical). Sin embargo, el valor no puede ser menor que la cantidad mínima de movimiento horizontal y debe estar en unidades pares de la cantidad mínima de movimiento horizontal.
//   - En modo estándar, se utiliza la unidad de movimiento horizontal.
//   - En modo página, la unidad de movimiento horizontal o vertical depende de la posición inicial del área imprimible configurada mediante ESC T:
//     1. Cuando la posición inicial se configura en la esquina superior izquierda o inferior derecha del área imprimible utilizando ESC T, se utiliza la unidad de movimiento horizontal (x).
//     2. Cuando la posición inicial se configura en la esquina superior derecha o inferior izquierda del área imprimible utilizando ESC T, se utiliza la unidad de movimiento vertical (y).
//   - El espaciado máximo del lado derecho es de 255/180 pulgadas. Cualquier configuración que exceda el máximo se convierte automáticamente al valor máximo.
//
// Valor por Defecto:
//
//	n = 0
//
// Referencia:
//
//	GS P
func SetRightSideCharSpacing(n byte) []byte {
	return []byte{ESC, SP, n}
}

// SelectPrintModes representa el comando para seleccionar modos de impresión.
//
// Nombre:
//
//	Seleccionar modos de impresión
//
// Formato:
//
//	ASCII: ESC ! n
//	Hex: 1B 21 n
//	Decimal: 27 33 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Selecciona modos de impresión utilizando n de la siguiente manera:
//
//	Bit Off/On | Hex | Decimal | Función
//	------------------------------------
//	0 Off      | 00  | 0       | Fuente de carácter A (12 × 24).
//	  On       | 01  | 1       | Fuente de carácter B (9 × 17).
//	1          | -   | -       | No definido.
//	2          | -   | -       | No definido.
//	3 Off      | 00  | 0       | Modo enfatizado no seleccionado.
//	  On       | 08  | 8       | Modo enfatizado seleccionado.
//	4 Off      | 00  | 0       | Modo de doble altura no seleccionado.
//	  On       | 10  | 16      | Modo de doble altura seleccionado.
//	5 Off      | 00  | 0       | Modo de doble ancho no seleccionado.
//	  On       | 20  | 32      | Modo de doble ancho seleccionado.
//	6          | -   | -       | No definido.
//	7 Off      | 00  | 0       | Modo subrayado no seleccionado.
//	  On       | 80  | 128     | Modo subrayado seleccionado.
//
// Detalles:
//   - Cuando se seleccionan ambos modos, doble altura y doble ancho, se imprimen caracteres de tamaño cuádruple.
//   - La impresora puede subrayar todos los caracteres, pero no puede subrayar espacios establecidos por HT o caracteres rotados 90° en sentido horario.
//   - El grosor del subrayado es el seleccionado por ESC , independientemente del tamaño del carácter.
//   - Cuando algunos caracteres en una línea tienen doble altura o más, todos los caracteres en la línea se alinean en la línea base.
//   - ESC E también puede activar o desactivar el modo enfatizado. Sin embargo, la configuración del último comando recibido es la efectiva.
//   - ESC —también puede activar o desactivar el modo subrayado. Sin embargo, la configuración del último comando recibido es la efectiva.
//   - GS ! también puede seleccionar el tamaño de los caracteres. Sin embargo, la configuración del último comando recibido es la efectiva.
//   - El modo enfatizado es efectivo para caracteres alfanuméricos y Kanji. Todos los modos de impresión, excepto el modo enfatizado, son efectivos solo para caracteres alfanuméricos.
//
// Valor por Defecto:
//
//	n = 0
//
// Referencia:
//
//	ESC -, ESC E, GS !
func SelectPrintModes(n byte) []byte {
	return []byte{ESC, '!', n}
}

// PrintRasterBitImage representa el comando para imprimir una imagen de bits en modo raster.
//
// Nombre:
//
//	Imprimir imagen de bits en modo raster
//
// Formato:
//
//	ASCII: GS v 0 m xL xH yL yH d1...dk
//	Hex: 1D 76 30 m xL xH yL yH d1...dk
//	Decimal: 29 118 48 m xL xH yL yH d1...dk
//
// Rango:
//   - m: 0 ≤ m ≤ 3, 48 ≤ m ≤ 51
//   - xL: 0 ≤ xL ≤ 255
//   - xH: 0 ≤ xH ≤ 255
//   - yL: 0 ≤ yL ≤ 255
//   - d: 0 ≤ d ≤ 255
//   - k = (xL + xH × 256) × (yL + yH × 256) (k ≥ 0)
//
// Descripción:
//
//	Selecciona el modo de imagen de bits raster. El valor de m selecciona el modo, como se indica a continuación:
//	- m = 0, 48: Modo normal (200 DPI vertical y horizontal).
//	- m = 1, 49: Modo de doble ancho (200 DPI vertical y 100 DPI horizontal).
//	- m = 2, 50: Modo de doble altura (100 DPI vertical y 200 DPI horizontal).
//	- m = 3, 51: Modo cuádruple (100 DPI vertical y horizontal).
//
// Detalles:
//   - En modo estándar, este comando es efectivo solo cuando no hay datos en el búfer de impresión.
//   - Este comando no tiene efecto en todos los modos de impresión (tamaño de caracteres, enfatizado, doble impacto, invertido, subrayado, impresión en blanco/negro, etc.) para imágenes de bits raster.
//   - Si el ancho del área de impresión configurado con GS L y GS W es menor que el ancho mínimo, el área de impresión se extiende al ancho mínimo solo en la línea en cuestión. El ancho mínimo es:
//   - 1 punto en modos normal (m=0, 48) y doble altura (m=2, 50).
//   - 2 puntos en modos doble ancho (m=1, 49) y cuádruple (m=3, 51).
//   - Los datos fuera del área de impresión se leen y se descartan punto por punto.
//   - La posición para imprimir caracteres posteriores en imágenes de bits raster se especifica mediante:
//   - HT (Tabulación Horizontal).
//   - ESC $ (Establecer posición de impresión absoluta).
//   - ESC \ (Establecer posición de impresión relativa).
//   - GS L (Establecer margen izquierdo).
//   - Si la posición para imprimir caracteres posteriores no es múltiplo de 8, la velocidad de impresión puede disminuir.
//   - La configuración de ESC a (Seleccionar justificación) también es efectiva para imágenes de bits raster.
//   - Cuando este comando se recibe durante la definición de macro, la impresora termina la definición de macro y comienza a ejecutar este comando. La definición de este comando debe ser borrada.
//   - El valor de d indica los datos de imagen de bits. Configurar un bit en 1 imprime un punto, mientras que configurarlo en 0 no imprime un punto.
func PrintRasterBitImage(m, xL, xH, yL, yH byte, d []byte) []byte {
	cmd := []byte{GS, 'v', '0', m, xL, xH, yL, yH}
	return append(cmd, d...)
}
