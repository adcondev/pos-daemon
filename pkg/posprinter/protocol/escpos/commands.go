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

	// DLE representa el byte de "Data Link Escape" en ESC/POS.
	//
	// Nombre:
	//   Data Link Escape (DLE)
	//
	// Valor:
	//   Hex: 0x10
	//   Decimal: 16
	//
	// Descripción:
	//   DLE se utiliza en protocolos de comunicación, incluyendo ESC/POS, para señalar el inicio
	//   de una secuencia de comandos o para distinguir entre datos y comandos en la transmisión.
	DLE byte = 0x10 // Data Link Escape (DLE)

	// EOT representa el byte de "End Of Transmission" en ESC/POS.
	//
	// Nombre:
	//   End Of Transmission (EOT)
	//
	// Valor:
	//   Hex: 0x04
	//   Decimal: 4
	//
	// Descripción:
	//   EOT se utiliza para indicar el final de una transmisión. En el contexto de ESC/POS, se emplea
	//   en comandos de transmisión en tiempo real para marcar el final de la comunicación o
	//   para solicitar información de estado.
	EOT byte = 0x04 // Fin de transmisión (End of Transmission)

	// ENQ representa el byte de "Enquiry" en ESC/POS.
	//
	// Nombre:
	//   Enquire (ENQ)
	//
	// Valor:
	//   Hex: 0x05
	//   Decimal: 5
	//
	// Descripción:
	//   ENQ se utiliza para solicitar una respuesta o confirmación del dispositivo. En sistemas ESC/POS,
	//   es comúnmente empleado en comandos de solicitud en tiempo real para recuperar información o
	//   reiniciar estados.
	ENQ byte = 0x05 // Solicitud de información (Enquiry)

	// DC4 representa el byte "Device Control 4" en ESC/POS.
	//
	// Nombre:
	//   Device Control 4 (DC4)
	//
	// Valor:
	//   Hex: 0x14
	//   Decimal: 20
	//
	// Descripción:
	//   DC4 es un byte de control utilizado en la comunicación de datos para funciones específicas.
	//   En el contexto de ESC/POS, DC4 forma parte de comandos que generan pulsos en tiempo real,
	//   tales como el comando DLE DC4 n m t.
	DC4 byte = 0x14 // Comando de control de dispositivo 4 (Device Control 4)

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

// TransmitRealTimeStatus representa el comando para transmitir el estado de la impresora en tiempo real.
//
// Nombre:
//
//	Transmisión de estado en tiempo real
//
// Formato:
//
//	ASCII: DLE EOT n
//	Hex:   10 04 n
//	Decimal: 16 4 n
//
// Rango:
//
//	1 ≤ n ≤ 4
//
// Descripción:
//
//	Transmite en tiempo real el estado seleccionado de la impresora según el valor de n, de la siguiente manera:
//	  n = 1: Transmitir estado de la impresora
//	  n = 2: Transmitir estado fuera de línea
//	  n = 3: Transmitir estado de error
//	  n = 4: Transmitir estado del sensor del papel continuo
//
// Detalles:
//   - El estado se transmite cada vez que se recibe la secuencia de datos <10>H<04>H<n>, con 1 ≤ n ≤ 4.
//     Ejemplo: En "ESC * m nL nH d1...dk", se tiene: d1 = <10>H, d2 = <04>H, d3 = <01>H.
//   - Este comando no debe ser utilizado dentro de la secuencia de datos de otro comando que consista en 2 o más bytes.
//     Ejemplo: Si se intenta transmitir "ESC 3 n" a la impresora, pero DTR (o DSR en la computadora host) cambia a MARK antes de que se transmita n,
//     y luego DLE EOT 3 interrumpe antes de que se reciba n, el código <10>H de DLE EOT 3 se procesa como el código para "ESC 3 <10>H".
//   - Este comando es efectivo incluso si la impresora no ha sido seleccionada mediante ESC = (seleccionar dispositivo periférico).
//   - La impresora transmite el estado actual, donde cada estado se representa con un dato de un byte.
//   - La transmisión se realiza sin confirmar si el host es capaz de recibir los datos.
//   - La impresora ejecuta el comando en cuanto lo recibe, incluso si está fuera de línea, el búfer de recepción está lleno o hay un estado de error
//     en modelos de interfaz serial.
//   - En modelos de interfaz paralela, este comando no se ejecuta cuando la impresora está ocupada, pero sí se ejecuta si está fuera de línea
//     o hay error cuando el DIP switch 2-1 está activado.
//   - Cuando se habilita Auto Status Back (ASB) mediante el comando GS a, se debe distinguir entre el estado transmitido por DLE EOT y el estado ASB.
//
// Parámetros de estado según el valor de n:
//
//	n = 1: Estado de la impresora
//	  Bit   | Off/On | Hex  | Decimal | Función
//	  ----- | ------ | ---- | ------- | -------------------------------------------------------------
//	   0    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//	        |   1    | 02   |    2    | No utilizado. Fijado en On.
//	   2    |   0    | 00   |    0    | Señal del cajón (pin 3 del conector) en nivel LOW.
//	        |   1    | 04   |    4    | Señal del cajón (pin 3 del conector) en nivel LOW.
//	   3    |   0    | 00   |    0    | En línea.
//	        |   1    | 08   |    8    | Fuera de línea.
//	   4    |   1    | 10   |   16    | No utilizado. Fijado en On.
//	   5,6  |   -    |  -   |    -    | Indefinido.
//	   7    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//
//	n = 2: Estado fuera de línea
//	  Bit   | Off/On | Hex  | Decimal | Función
//	  ----- | ------ | ---- | ------- | -------------------------------------------------------------
//	   0    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//	   1    |   1    | 02   |    2    | No utilizado. Fijado en On.
//	   2    |   0    | 00   |    0    | La tapa está cerrada.
//	        |   1    | 04   |    4    | La tapa está abierta.
//	   3    |   0    | 00   |    0    | Papel no se alimenta mediante el botón FEED.
//	        |   1    | 08   |    8    | Papel se alimenta mediante el botón FEED.
//	   4    |   1    | 10   |   16    | No utilizado. Fijado en On.
//	   5    |   0    | 00   |    0    | No hay tope final de papel.
//	        |   1    | 20   |   32    | Se está deteniendo la impresión.
//	   6    |   0    | 00   |    0    | No hay error.
//	        |   1    | 40   |   64    | Se produce un error.
//	   7    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//
//	n = 3: Estado de error
//	  Bit   | Off/On | Hex  | Decimal | Función
//	  ----- | ------ | ---- | ------- | -------------------------------------------------------------
//	   0    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//	   1    |   1    | 02   |    2    | No utilizado. Fijado en On.
//	   2    |   -    |  -   |    -    | Indefinido.
//	   3    |   0    | 00   |    0    | No hay error en la autocortadora.
//	        |   1    | 08   |    8    | Error en la autocortadora.
//	   4    |   1    | 10   |   16    | No utilizado. Fijado en On.
//	   5    |   0    | 00   |    0    | No hay error irrecuperable.
//	        |   1    | 20   |   32    | Se produce un error irrecuperable.
//	   6    |   0    | 00   |    0    | No hay error auto-recuperable.
//	        |   1    | 40   |   64    | Se produce un error auto-recuperable.
//	   7    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//
//	n = 4: Estado del sensor de papel continuo
//	  Bit   | Off/On | Hex  | Decimal | Función
//	  ----- | ------ | ---- | ------- | -------------------------------------------------------------
//	   0    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//	   0    |   1    | 02   |    2    | No utilizado. Fijado en On.
//	   2,3  |   0    | 00   |    0    | Sensor de fin de papel: papel adecuado.
//	        |   1    |   -  |    -    | Fin de papel detectado (sensor activado).
//	   5,6  |   0C   | 10   |   16    | Se detecta fin de papel mediante el sensor.
//	   4    |   0    | 00   |    0    | No utilizado. Fijado en On.
//	        |   0    | 00   |    0    | Sensor de papel: papel presente.
//	   1    |   1    | 60   |   96    | Sensor de fin de papel: fin de papel detectado.
//	   7    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//
// Referencia:
//
//	DLE ENQ, GS a, GS r
func TransmitRealTimeStatus(n byte) []byte {
	cmd := []byte{DLE, EOT, n}
	return cmd
}

// RequestPrinterRecovery representa el comando para realizar una petición en tiempo real a la impresora.
//
// Nombre:
//
//	Petición en tiempo real a la impresora
//
// Formato:
//
//	ASCII: DLE ENQ n
//	Hex:   10 05 n
//	Decimal: 16 5 n
//
// Rango:
//
//	1 ≤ n ≤ 2
//
// Descripción:
//
//	Responde a una solicitud de la computadora host. El valor de n especifica la solicitud de la siguiente manera:
//	  n = 1: Recuperar de un error y reiniciar la impresión desde la línea en la que ocurrió el error.
//	  n = 2: Recuperar de un error después de borrar los búferes de recepción e impresión.
//
// Detalles:
//   - Este comando es efectivo únicamente cuando ocurre un error en la autocortadora.
//   - La impresora comienza a procesar los datos al recibir este comando.
//   - Se ejecuta incluso si la impresora está fuera de línea, el búfer de recepción está lleno o se presenta un estado de error en modelos de interfaz serial.
//   - En modelos de interfaz paralela, el comando no se ejecuta cuando la impresora está ocupada; sin embargo, se ejecuta cuando la impresora está fuera de línea o hay error, si el DIP switch 2-1 está activado.
//   - El estado también se transmite cada vez que se recibe la secuencia de datos <10>H <05>H <n> (1 ≤ n ≤ 2).
//     Ejemplo: En el comando "ESC * m nL nH dk", se tiene: d1 = <10>H, d2 = <05>H, d3 = <01>H.
//   - Este comando no debe incluirse dentro de otra secuencia de comandos que consista en dos o más bytes.
//   - DLE ENQ 2 permite que la impresora se recupere de un error luego de borrar los datos en los búferes de recepción e impresión, conservando las configuraciones (por ejemplo, las establecidas por ESC !, ESC 3, etc.) vigentes al ocurrir el error.
//   - La impresora se puede inicializar completamente utilizando este comando junto con ESC @.
//   - Cuando la impresora se deshabilita con ESC = (seleccionar dispositivo periférico), las funciones de recuperación de error (DLE ENQ 1 y DLE ENQ 2) quedan habilitadas y las demás funciones se desactivan.
//
// Referencia:
//
//	DLE EOT
func RequestPrinterRecovery(n byte) []byte {
	cmd := []byte{DLE, ENQ, n}
	return cmd
}

// GenerateRealTimePulse representa el comando DLE DC4 n m t para generar un pulso en tiempo real.
//
// Nombre:
//
//	Generar pulso en tiempo real
//
// Formato:
//
//	ASCII: DLE DC4 n m t
//	Hex:   10 14 n m t
//	Decimal: 16 20 n m t
//
// Rango:
//   - n: n = 1
//   - m: m ∈ {0, 1}
//     m = 0: Pin 2 del conector de expulsión del cajón.
//     m = 1: Pin 5 del conector de expulsión del cajón.
//   - t: 1 ≤ t ≤ 8
//
// Descripción:
//
//	Genera el pulso especificado por t en el pin de conector indicado por m.
//	El tiempo de encendido del pulso es t × 100 ms y el tiempo de apagado es también t × 100 ms.
//
// Detalles:
//   - Si la impresora se encuentra en estado de error al procesar este comando, éste se ignora.
//   - Si el pulso se está enviando al pin especificado mientras se ejecuta ESC p o DEL DC4, el comando se ignora.
//   - La impresora ejecuta este comando en cuanto lo recibe.
//   - Con un modelo de interfaz serial, este comando se ejecuta incluso si la impresora está fuera de línea, el búfer de recepción está lleno o hay un estado de error.
//   - Con un modelo de interfaz paralela, el comando no se ejecuta cuando la impresora está ocupada; sin embargo, se ejecuta cuando está fuera de línea o hay error si el DIP switch 2-1 está activado.
//   - Si los datos de impresión incluyen cadenas de caracteres idénticas a este comando, la impresora realizará la misma operación especificada por este comando. Se debe tener en cuenta este comportamiento.
//   - Este comando no debe usarse dentro de la secuencia de datos de otro comando que consista en 2 o más bytes.
//   - Es efectivo incluso cuando la impresora está deshabilitada mediante ESC = (selección de dispositivo periférico).
//
// Referencia:
//
//	ESC p
func GenerateRealTimePulse(n, m, t byte) []byte {
	cmd := []byte{DLE, DC4, n, m, t}
	return cmd
}

// SetRightSideCharacterSpacing representa el comando para configurar el espaciado a la derecha de los caracteres.
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
func SetRightSideCharacterSpacing(n byte) []byte {
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

// SetAbsolutePrintPosition representa el comando para establecer la posición de impresión absoluta.
//
// Nombre:
//
//	Establecer posición de impresión absoluta
//
// Formato:
//
//	ASCII: ESC $ nL nH
//	Hex:   1B 24 nL nH
//	Decimal: 27 36 nL nH
//
// Rango:
//
//	0 ≤ nL ≤ 255
//	0 ≤ nH ≤ 255
//
// Descripción:
//
//	Establece la distancia desde el comienzo de la línea hasta la posición en la que se imprimirán los caracteres posteriores.
//
// Detalles:
//   - La distancia se calcula como: (nL + nH × 256) × (unidad de movimiento vertical u horizontal) en pulgadas.
//   - Las configuraciones fuera del área imprimible especificada se ignoran.
//   - La unidad de movimiento horizontal y vertical se especifica mediante el comando GS P. Este comando puede modificar dichas unidades,
//     aunque el valor no puede ser menor que la cantidad mínima de movimiento horizontal y debe ser un múltiplo par de dicha cantidad mínima.
//   - En modo estándar se utiliza la unidad de movimiento horizontal (x).
//   - En modo página, la unidad de movimiento puede ser horizontal o vertical, dependiendo de la posición inicial del área imprimible configurada con ESC T:
//     1. Si la posición inicial se establece en la esquina superior izquierda o inferior derecha, se usa la unidad horizontal (x).
//     2. Si la posición inicial se establece en la esquina superior derecha o inferior izquierda, se usa la unidad vertical (y).
//
// Referencia:
//
//	ESC \, GS $, GS P
func SetAbsolutePrintPosition(nL, nH byte) []byte {
	return []byte{ESC, '$', nL, nH}
}

// SelectOrCancelUserDefinedCharset representa el comando para seleccionar o cancelar el conjunto de caracteres definido por el usuario.
//
// Nombre:
//
//	Seleccionar o cancelar conjunto de caracteres definido por el usuario
//
// Formato:
//
//	ASCII: ESC % n
//	Hex:   1B 25 n
//	Decimal: 27 37 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Selecciona o cancela el conjunto de caracteres definido por el usuario.
//
// Detalles:
//   - Cuando el bit menos significativo (LSB) de n es 0, se cancela el conjunto de caracteres definido por el usuario.
//   - Cuando el bit menos significativo de n es 1, se selecciona el conjunto de caracteres definido por el usuario.
//   - Al cancelar el conjunto definido por el usuario, se selecciona automáticamente el conjunto de caracteres interno.
//
// Valor por Defecto:
//
//	n = 0
//
// Referencia:
//
//	ESC &, ESC ?
func SelectOrCancelUserDefinedCharset(n byte) []byte {
	return []byte{ESC, '%', n}
}

// DefineUserDefinedCharacters representa el comando para definir caracteres personalizados.
//
// Nombre:
//
//	Definir caracteres personalizados
//
// Formato:
//
//	ASCII: ESC & y c1 c2 [x1 d1...d(y × x1)] ... [xk d1...d(y × xk)]
//	Hex:   1B 26 y c1 c2 [x1 d1...d(y×x1)] ... [xk d1...d(y×xk)]
//
// Rango:
//   - y: Especifica el número de bytes en la dirección vertical (por ejemplo, y = 3).
//   - c1, c2: Rango de códigos de caracteres a definir; 32 (0x20) ≤ c1 ≤ c2 ≤ 126 (0x7E).
//   - x: Número de puntos en la dirección horizontal.
//     Para la Fuente A (12×24): 0 ≤ x ≤ 12.
//     Para la Fuente B (9×17): 0 ≤ x ≤ 9.
//   - d (datos de puntos): 0 ≤ d ≤ 255.
//   - La cantidad de datos para cada carácter es (y × x) bytes.
//
// Descripción:
//
//	Define caracteres personalizados utilizando un conjunto de datos que especifican el patrón de puntos de cada
//	carácter. Se pueden definir múltiples caracteres consecutivos asignados a códigos de carácter desde c1 hasta c2.
//	Para definir un único carácter, se utiliza c1 = c2.
//
// Detalles:
//   - y especifica el número de bytes en la dirección vertical.
//   - c1 indica el código del primer carácter a definir y c2, el código final.
//   - x define el número de puntos (dots) en la dirección horizontal para cada carácter.
//   - El patrón de puntos se configura de izquierda a derecha; los puntos restantes a la derecha se dejan en blanco.
//   - Cada carácter se define con (y × x) bytes, donde cada bit determina si se imprime (bit = 1) o no imprime (bit = 0) un punto.
//   - Es posible definir diferentes patrones de caracteres para cada fuente. Para seleccionar la fuente deseada, utilice el comando ESC !
//   - No es posible definir un carácter personalizado y una imagen de bits descargada simultáneamente. La ejecución de este comando
//     borrará cualquier imagen de bits descargada.
//   - La definición de caracteres personalizados se borra cuando se ejecuta cualquiera de los siguientes comandos:
//     ① ESC @
//     ② ESC ?
//     ③ FS q
//     ④ GS *
//     ⑤ Al reiniciar la impresora o apagar la alimentación.
//   - Para caracteres definidos en la Fuente B (9×17), únicamente es efectivo el bit más significativo del tercer byte en la dirección vertical.
//   - Al cancelar la definición personalizada mediante otros comandos (por ejemplo, ESC %, ESC ?), se selecciona el conjunto de caracteres interno.
//
// Referencia:
//
//	ESC %, ESC ?
func DefineUserDefinedCharacters(y, c1, c2 byte, data ...[]byte) []byte {
	cmd := []byte{ESC, '&', y, c1, c2}
	for _, d := range data {
		cmd = append(cmd, d...)
	}
	return cmd
}

// SelectBitImageMode representa el comando para seleccionar el modo de imagen de bits.
//
// Nombre:
//
//	Seleccionar modo de imagen de bits
//
// Formato:
//
//	ASCII: ESC * m nL nH d1...dk
//	Hex:   1B 2A m nL nH d1...dk
//	Decimal: 27 42 m nL nH d1...dk
//
// Rango:
//   - m: Puede ser 0, 1, 32 o 33.
//   - nL: 0 ≤ nL ≤ 255
//   - nH: 0 ≤ nH ≤ 3
//   - d (datos): 0 ≤ d ≤ 255
//
// Descripción:
//
//	Selecciona un modo de imagen de bits utilizando el parámetro m para determinar
//	el número de puntos especificados por nL y nH, de la siguiente manera:
//
//	Modos según m:
//	  * m = 0: Modo de 8 puntos, densidad simple (8-dot single-density).
//	    - Dirección vertical: 8 puntos.
//	    - Dirección horizontal: El número de puntos es nL + nH × 256.
//	    - Densidad: 67 DPI.
//	  * m = 1: Modo de 8 puntos, doble densidad (8-dot double-density).
//	    - Dirección vertical: 8 puntos.
//	    - Dirección horizontal: El número de puntos es nL + nH × 256.
//	    - Densidad: 67 DPI y 100 DPI (dependiente del contexto).
//	  * m = 32: Modo de 24 puntos, densidad simple (24-dot single-density).
//	    - Dirección vertical: 24 puntos.
//	    - Dirección horizontal: El número de puntos es nL + nH × 256.
//	    - Densidad: 200 DPI.
//	  * m = 33: Modo de 24 puntos, doble densidad (24-dot double-density).
//	    - Dirección vertical: 24 puntos.
//	    - Dirección horizontal: El número de puntos es (nL + nH × 256) × 3.
//	    - Densidad: 200 DPI.
//
// Detalles:
//   - Si el valor de m no se encuentra dentro de los rangos especificados, nL y los datos
//     siguientes se procesarán como datos normales.
//   - Los parámetros nL y nH indican el número de puntos de la imagen de bits en la dirección horizontal,
//     calculándose como nL + nH × 256.
//   - Si los datos de la imagen de bits exceden el número de puntos que se pueden imprimir en una línea,
//     los datos excedentes serán ignorados.
//   - d representa los datos de la imagen de bits; se debe establecer el bit correspondiente a 1 para imprimir un punto,
//     o a 0 para no imprimirlo.
//   - Si el ancho del área de impresión establecido por GS L y GS W es menor que el requerido por los datos enviados
//     con el comando ESC *, se realizará lo siguiente en la línea en cuestión (sin exceder el área máxima de impresión):
//     ① Se extiende hacia la derecha el ancho del área de impresión para acomodar los datos.
//     ② Si el paso ① no proporciona el ancho suficiente, se reduce el margen izquierdo para acomodar los datos.
//   - Después de imprimir la imagen de bits, la impresora regresa al modo de procesamiento de datos normal.
//   - Este comando no se ve afectado por los modos de impresión (enfatizado, doble impacto, subrayado,
//     tamaño de caracteres o impresión en blanco/negro), excepto en el modo de impresión invertida (upside-down).
//
// Referencia:
//
//	GS L, GS W, ESC \, GS P
func SelectBitImageMode(m, nL, nH byte, data []byte) []byte {
	cmd := []byte{ESC, '*', m, nL, nH}
	return append(cmd, data...)

}

// TurnUnderlineMode representa el comando ESC - n para activar o desactivar el modo subrayado.
//
// Nombre:
//
//	Activar/Desactivar modo subrayado
//
// Formato:
//
//	ASCII: ESC - n
//	Hex:   1B 2D n
//	Decimal: 27 45 n
//
// Rango:
//
//	n puede tomar uno de los siguientes valores:
//	  • 0, 48: Desactiva el modo subrayado.
//	  • 1, 49: Activa el modo subrayado con grosor de 1 punto.
//	  • 2, 50: Activa el modo subrayado con grosor de 2 puntos.
//	Valores válidos: 0 ≤ n ≤ 2, 48 ≤ n ≤ 50
//
// Descripción:
//
//	Activa o desactiva el modo subrayado en la impresora basado en el valor de n recibido.
//
// Detalles:
//   - La impresora puede subrayar todos los caracteres (incluyendo el espaciado derecho), pero no puede subrayar
//     los espacios establecidos por HT.
//   - No es posible subrayar caracteres rotados 90° en sentido horario ni caracteres invertidos de blanco y negro.
//   - Cuando el modo subrayado se desactiva (n = 0 o 48), los datos posteriores no se subrayan y el grosor del subrayado
//     configurado previamente se mantiene sin cambios (el grosor por defecto es de 1 punto).
//   - Cambiar el tamaño de los caracteres no afecta el grosor de subrayado actual.
//   - El modo subrayado también puede activarse o desactivarse mediante ESC !. Sin embargo, solo se considera efectivo
//     el último comando recibido.
//   - Este comando no afecta la configuración de los caracteres Kanji.
//
// Valor por defecto:
//
//	n = 0
//
// Referencia:
//
//	ESC !
func TurnUnderlineMode(n byte) []byte {
	cmd := []byte{ESC, '-', n}
	return cmd
}

// SelectDefaultLineSpacing representa el comando ESC 2 para seleccionar el espaciado de línea por defecto.
//
// Nombre:
//
//	Seleccionar espaciado de línea por defecto
//
// Formato:
//
//	ASCII: ESC 2
//	Hex:   1B 32
//	Decimal: 27 50
//
// Descripción:
//
//	Selecciona un espaciado de línea de 1/6 de pulgada (aproximadamente 4,23 mm).
//
// Detalles:
//   - El espaciado de línea se puede configurar de manera independiente en el modo estándar y en el modo página.
//
// Referencia:
//
//	ESC 3
func SelectDefaultLineSpacing() []byte {
	cmd := []byte{ESC, '2'}
	return cmd
}

// SetLineSpacing representa el comando ESC 3 n para configurar el espaciado de línea.
//
// Nombre:
//
//	Configurar espaciado de línea
//
// Formato:
//
//	ASCII: ESC 3 n
//	Hex:   1B 33 n
//	Decimal: 27 51 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Establece el espaciado de línea a [n × (unidad de movimiento vertical u horizontal)] pulgadas.
//
// Detalles:
//   - El espaciado de línea se puede configurar de manera independiente en el modo estándar y en el modo página.
//   - La unidad de movimiento horizontal y vertical se especifica mediante el comando GS P. Cambiar la unidad horizontal o vertical
//     no afecta el espaciado de línea actual.
//   - En el modo estándar se utiliza la unidad de movimiento vertical (y).
//   - En el modo página, el comando funciona de la siguiente manera según la posición inicial del área imprimible definida con ESC T:
//     ① Si la posición inicial se establece en la esquina superior izquierda o inferior derecha, se usa la unidad vertical (y).
//     ② Si la posición inicial se establece en la esquina superior derecha o inferior izquierda, se usa la unidad horizontal (x).
//   - La cantidad máxima de alimentación de papel es de 1016 mm (40 pulgadas). Aunque se configure un valor mayor,
//     la impresora alimenta el papel únicamente hasta 1016 mm (40 pulgadas).
//
// Valor por Defecto:
//
//	Espaciado de línea equivalente a aproximadamente 4,23 mm (1/6 de pulgada).
//
// Referencia:
//
//	ESC 2, GS P
func SetLineSpacing(n byte) []byte {
	cmd := []byte{ESC, '3', n}
	return cmd
}

// SetPeripheralDevice representa el comando ESC = n para seleccionar el dispositivo periférico.
//
// Nombre:
//
//	Seleccionar dispositivo periférico
//
// Formato:
//
//	ASCII: ESC = n
//	Hex:   1B 3D n
//	Decimal: 27 61 n
//
// Rango:
//
//	1 ≤ n ≤ 255
//
// Descripción:
//
//	Selecciona el dispositivo al que la computadora host envía datos, utilizando el valor de n de la siguiente manera:
//	  - Si el bit correspondiente está apagado (n = 0x00), la impresora se desactiva (Printer disabled).
//	  - Si el bit correspondiente está encendido (n = 0x01), la impresora se activa (Printer enabled).
//	  - Los valores de n de 1 a 7 están indefinidos.
//
// Detalles:
//   - Cuando la impresora está deshabilitada, ignora todos los datos entrantes, excepto los comandos de recuperación de errores
//     (DLE EOT, DLE ENQ, DLE DC4), hasta que se active nuevamente utilizando este comando.
//
// Valor por Defecto:
//
//	n = 1
func SetPeripheralDevice(n byte) []byte {
	cmd := []byte{ESC, '=', n}
	return cmd
}

// CancelUserDefinedCharacters representa el comando ESC ? n para cancelar caracteres definidos por el usuario.
//
// Nombre:
//
//	Cancelar caracteres definidos por el usuario
//
// Formato:
//
//	ASCII: ESC ? n
//	Hex:   1B 3F n
//	Decimal: 27 63 n
//
// Rango:
//
//	32 ≤ n ≤ 126
//
// Descripción:
//
//	Cancela los caracteres definidos por el usuario. Después de cancelar, se imprime el patrón correspondiente
//	al carácter interno.
//
// Detalles:
//   - El comando elimina el patrón definido para el código de carácter especificado por n en la fuente seleccionada mediante ESC !.
//   - Si un carácter definido por el usuario no ha sido previamente definido para el código especificado, la impresora ignora este comando.
//
// Referencia:
//
//	ESC &, ESC %
func CancelUserDefinedCharacters(n byte) []byte {
	cmd := []byte{ESC, '?', n}
	return cmd
}

// InitializePrinter representa el comando ESC @ para inicializar la impresora.
//
// Nombre:
//
//	Inicializar impresora
//
// Formato:
//
//	ASCII: ESC @
//	Hex:   1B 40
//	Decimal: 27 64
//
// Descripción:
//
//	Borra los datos en el búfer de impresión y restablece el modo de la impresora al que estaba vigente al encenderla.
//
// Detalles:
//   - La configuración de los DIP switches no se vuelve a verificar.
//   - Los datos en el búfer de recepción no se borran.
//   - La definición de macros no se limpia.
//   - Los datos de la imagen de bits en la memoria NV no se borran.
//   - Los datos de la memoria NV del usuario no se borran.
func InitializePrinter() []byte {
	cmd := []byte{ESC, '@'}
	return cmd
}

// SetHorizontalTabPositions representa el comando ESC D n1...nk NUL para establecer posiciones de tabulación horizontales.
//
// Nombre:
//
//	Establecer posiciones de tabulación horizontales
//
// Formato:
//
//	ASCII: ESC D n1...nk NUL
//	Hex:   1B 44 n1...nk 00
//	Decimal: 27 68 n1...nk 0
//
// Rango:
//   - n: 1 ≤ n ≤ 255 (cada n especifica un número de columna desde el inicio de la línea)
//   - k: 0 ≤ k ≤ 32 (k indica la cantidad total de posiciones de tabulación que se pueden establecer)
//
// Descripción:
//
//	Configura las posiciones de tabulación horizontales. Cada valor n representa la columna en la que se establecerá una posición de tabulación,
//	contada desde el inicio de la línea. Al enviar un código NUL (0) al final, se indica el fin de la secuencia de tabulaciones.
//	Si se envía ESC D NUL, se cancelan todas las posiciones de tabulación horizontales previamente definidas.
//
// Detalles:
//   - Las posiciones de tabulación se almacenan como el producto de [ancho de carácter × n], medido desde el inicio de la línea.
//     El ancho de carácter incluye el espaciado a la derecha, y los caracteres de doble ancho se configuran con el doble del ancho normal.
//   - Este comando borra las configuraciones de tabulación horizontales anteriores.
//   - Al establecer n = 8, la posición de impresión se mueve a la columna 9 mediante el envío del carácter HT.
//   - Se pueden establecer hasta 32 posiciones de tabulación (k = 32). Cualquier dato que exceda de 32 valores se procesa como datos normales.
//   - Se deben enviar los valores de [n] en orden ascendente y finalizar con un código NUL (0).
//   - Si un valor de [n] es menor o igual que el valor anterior (n[k] ≤ n[k-1]), se concluye la configuración y los datos siguientes se interpretan como datos normales.
//   - Las posiciones de tabulación previamente configuradas no cambian, incluso si el ancho del carácter se modifica posteriormente.
//   - El ancho de carácter se memoriza de forma independiente para el modo estándar y el modo página.
//
// Valor por defecto:
//
//	Las posiciones de tabulación por defecto se establecen a intervalos de 8 caracteres (por ejemplo, columnas 9, 17, 25, ...)
//	para la Fuente A (12x24).
//
// Referencia:
//
//	HT
func SetHorizontalTabPositions(n []byte) []byte {
	cmd := []byte{ESC, 'D'}
	cmd = append(cmd, n...)
	cmd = append(cmd, NUL) // NUL terminator
	return cmd
}

// TurnEmphasizedMode representa el comando ESC E n para activar o desactivar el modo enfatizado.
//
// Nombre:
//
//	Activar/Desactivar modo enfatizado
//
// Formato:
//
//	ASCII: ESC E n
//	Hex:   1B 45 n
//	Decimal: 27 69 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Activa o desactiva el modo enfatizado basado en el valor del bit menos significativo (LSB) de n:
//	  - Si el LSB de n es 0, el modo enfatizado se desactiva.
//	  - Si el LSB de n es 1, el modo enfatizado se activa.
//
// Detalles:
//   - Solo el bit menos significativo (LSB) de n es utilizado para determinar el estado del modo enfatizado.
//   - Este comando y ESC ! activan o desactivan el modo enfatizado de la misma manera. Es importante tener cuidado
//     al utilizar ambos comandos simultáneamente, ya que el último comando recibido prevalecerá.
//
// Valor por defecto:
//
//	n = 0
//
// Referencia:
//
//	ESC !
func TurnEmphasizedMode(n byte) []byte {
	cmd := []byte{ESC, 'E', n}
	return cmd
}

// PrintAndFeedPaper representa el comando ESC J n para imprimir y alimentar el papel.
//
// Nombre:
//
//	Imprimir y alimentar papel
//
// Formato:
//
//	ASCII: ESC J n
//	Hex:   1B 4A n
//	Decimal: 27 74 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Imprime los datos en el búfer de impresión y alimenta el papel en una cantidad equivalente a
//	[n × unidad de movimiento vertical u horizontal] pulgadas.
//
// Detalles:
//   - Después de completar la impresión, este comando establece la posición inicial de impresión al comienzo de la línea.
//   - La cantidad de alimentación de papel configurada por este comando no afecta los valores establecidos por ESC 2 o ESC 3.
//   - La unidad de movimiento horizontal y vertical se especifica mediante el comando GS P.
//   - El comando GS P puede cambiar las unidades de movimiento vertical y horizontal. Sin embargo, el valor no puede ser menor
//     que la cantidad mínima de movimiento vertical y debe ser un múltiplo par de dicha cantidad mínima.
//   - En modo estándar, la impresora utiliza la unidad de movimiento vertical (y).
//   - En modo página, el comando funciona de la siguiente manera según la posición inicial del área imprimible definida con ESC T:
//     ① Si la posición inicial se establece en la esquina superior izquierda o inferior derecha, se usa la unidad vertical (y).
//     ② Si la posición inicial se establece en la esquina superior derecha o inferior izquierda, se usa la unidad horizontal (x).
//   - La cantidad máxima de alimentación de papel es de 1016 mm (40 pulgadas). Si el valor configurado excede este límite,
//     se ajustará automáticamente al máximo permitido.
//
// Referencia:
//
//	GS P
func PrintAndFeedPaper(n byte) []byte {
	cmd := []byte{ESC, 'J', n}
	return cmd
}

// SelectPageMode cambia la impresora del modo estándar al modo página.
//
// Formato:
//
//	ASCII: ESC L
//	Hex:   1B 4C
//	Decimal: 27 76
//
// Descripción:
//
//	Activa el modo página en la impresora. Este comando es válido solo cuando se procesa al comienzo de una línea en modo estándar
//	y no tiene efecto si ya se encuentra en modo página.
//
// Detalles:
//   - Después de completar la impresión utilizando el comando FF o ESC S, la impresora vuelve al modo estándar.
//   - Este comando establece la posición donde se almacenan los datos en el búfer según la posición especificada por ESC T dentro del área de impresión definida por ESC W.
//   - Los siguientes comandos se configuran con valores para el modo página, donde se pueden establecer valores de forma independiente en modo estándar y modo página:
//     ① Configurar el espaciado de caracteres hacia la derecha: ESC SP, FS S
//     ② Seleccionar el espaciado de línea predeterminado: ESC 2, ESC 3
//   - En modo página, solo es posible configurar valores para los siguientes comandos; estos comandos no se ejecutan:
//     ① Activar/desactivar el modo de rotación 90° en sentido horario: ESC V
//     ② Seleccionar justificación: ESC a
//     ③ Activar/desactivar el modo de impresión invertida: ESC {
//     ④ Configurar margen izquierdo: GS L
//     ⑤ Configurar el ancho del área imprimible: GS W
//   - El siguiente comando se ignora en modo página:
//     ① Ejecutar impresión de prueba: GS ( A
//   - Los siguientes comandos no están disponibles en modo página:
//     ① Imprimir imagen NV: FS p
//     ② Definir imagen NV: FS q
//     ③ Escribir en la memoria NV del usuario: FS g 1
//     ④ Imprimir imagen de bits rasterizada: GS v 0
//   - La impresora vuelve al modo estándar cuando se enciende, se reinicia o se utiliza el comando ESC @.
//
// Referencia:
//
//	FF, CAN, ESC FF, ESC S, ESC T, ESC W, GS $, GS \
func SelectPageMode() []byte {
	cmd := []byte{ESC, 'L'}
	return cmd
}

// SelectCharacterFont permite configurar la fuente de los caracteres en la impresora.
//
// Formato:
//
//	ASCII: ESC M n
//	Hex:   1B 4D n
//	Decimal: 27 77 n
//
// Rango:
//
//	n = 0, 1, 48, 49
//
// Descripción:
//
//	Selecciona la fuente de caracteres que será utilizada por la impresora:
//	  - n = 0 o 48: Se selecciona la fuente de caracteres A (12 × 24).
//	  - n = 1 o 49: Se selecciona la fuente de caracteres B (9 × 17).
//
// Detalles:
//   - Este comando determina el tamaño de fuente activo para la impresión de texto.
func SelectCharacterFont(n byte) []byte {
	cmd := []byte{ESC, 'M', n}
	return cmd
}

// SelectInternationalCharacterSet configura el conjunto de caracteres internacionales de la impresora.
//
// Formato:
//
//	ASCII: ESC R n
//	Hex:   1B R n
//	Decimal: 27 R n
//
// Rango:
//
//	0 ≤ n ≤ 15
//
// Descripción:
//
//	Selecciona un conjunto de caracteres internacionales basado en el valor de n, según la siguiente tabla:
//	  n   Conjunto de caracteres
//	  0   Estados Unidos
//	  1   Francia
//	  2   Alemania
//	  3   Reino Unido
//	  4   Dinamarca
//	  5   Suecia
//	  6   Italia
//	  7   España
//	  8   Japón
//	  9   Noruega
//	  10  Dinamarca
//	  11  España
//	  12  Latino
//	  13  Chino
//	  14  Corea
//	  15  Eslovenia/Croacia
//
// Predeterminado:
//
//	Para el modelo Simplificado Chino: n = 15; Para modelos distintos del Simplificado Chino: n = 0
//
// Nota:
//
//	Los conjuntos de caracteres para Eslovenia/Croacia y China solo son compatibles con el modelo Simplificado Chino.
func SelectInternationalCharacterSet(n byte) []byte {
	cmd := []byte{ESC, 'R', n}
	return cmd
}

// SelectStandardMode cambia la impresora del modo página al modo estándar.
//
// Formato:
//
//	ASCII: ESC S
//	Hex:   1B 53
//	Decimal: 27 83
//
// Descripción:
//
//	Activa el modo estándar en la impresora. Este comando es válido solo en el modo página.
//
// Detalles:
//   - Los datos almacenados en el búfer en modo página se eliminan al cambiar al modo estándar.
//   - La posición de impresión se establece al comienzo de la línea.
//   - El área de impresión configurada por ESC W se inicializa.
//   - Los siguientes comandos se configuran con valores para el modo estándar, donde los valores pueden establecerse de forma independiente para modo estándar y modo página:
//     ① Configurar el espaciado de caracteres hacia la derecha: ESC SP, FS S
//     ② Seleccionar el espaciado de línea predeterminado: ESC 2, ESC 3
//   - Los siguientes comandos están habilitados solo para configurar en modo estándar:
//     ① Configurar el área de impresión en modo página: ESC W
//     ② Seleccionar la dirección de impresión en modo página: ESC T
//   - Los siguientes comandos se ignoran en modo estándar:
//     ① Configurar posición de impresión vertical absoluta en modo página: GS $
//     ② Configurar posición de impresión vertical relativa en modo página: GS \
//   - El modo estándar se selecciona automáticamente cuando se enciende la impresora, se reinicia o se utiliza el comando ESC @.
//
// Referencia:
//
//	FF, ESC FF, ESC L
func SelectStandardMode() []byte {
	cmd := []byte{ESC, 'S'}
	return cmd
}
