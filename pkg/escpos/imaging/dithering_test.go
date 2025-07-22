package imaging

import (
	"image"
	"image/color"
	"testing"
)

// TestResizeImage verifica que la función mantiene la proporción correcta
// y transfiere adecuadamente los colores de la imagen original
func TestResizeImage(t *testing.T) {
	// Creamos una imagen de prueba de 4x2 con colores diferentes
	src := image.NewRGBA(image.Rect(0, 0, 4, 2))

	// Llenamos con patrones de colores reconocibles
	src.Set(0, 0, color.RGBA{R: 255, A: 255})                 // Rojo
	src.Set(1, 0, color.RGBA{G: 255, A: 255})                 // Verde
	src.Set(2, 0, color.RGBA{B: 255, A: 255})                 // Azul
	src.Set(3, 0, color.RGBA{R: 255, G: 255, A: 255})         // Amarillo
	src.Set(0, 1, color.RGBA{R: 255, B: 255, A: 255})         // Magenta
	src.Set(1, 1, color.RGBA{G: 255, B: 255, A: 255})         // Cian
	src.Set(2, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // Blanco
	src.Set(3, 1, color.RGBA{A: 255})                         // Negro

	// Redimensionamos a la mitad del ancho (2px)
	out := ResizeImage(src, 2)
	bounds := out.Bounds()

	// Verificamos que la proporción se mantenga (4:2 -> 2:1)
	if bounds.Dx() != 2 || bounds.Dy() != 1 {
		t.Fatalf("proporción incorrecta: got %dx%d, want 2x1", bounds.Dx(), bounds.Dy())
	}

	// Verificamos que los colores se conserven (mediante muestreo de vecinos más cercanos)
	// En la reducción 4x2 -> 2x1, deberíamos obtener los colores de las posiciones 0,0 y 2,0
	rojo := color.RGBAModel.Convert(out.At(0, 0)).(color.RGBA)
	if rojo.R != 255 || rojo.G != 0 || rojo.B != 0 {
		t.Errorf("primer píxel incorrecto: got %v, want rojo", rojo)
	}

	azul := color.RGBAModel.Convert(out.At(1, 0)).(color.RGBA)
	if azul.R != 0 || azul.G != 0 || azul.B != 255 {
		t.Errorf("segundo píxel incorrecto: got %v, want azul", azul)
	}
}

// TestProcessImageWithDithering_FloydStein verifica el procesamiento de imágenes
// con el algoritmo de dithering Floyd-Steinberg
func TestProcessImageWithDithering_FloydStein(t *testing.T) {
	// Creamos una imagen en escala de grises con un gradiente simple
	src := image.NewGray(image.Rect(0, 0, 3, 3))

	// Llenamos con un gradiente de 0 a 240 en pasos de 30
	valor := uint8(0)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			src.SetGray(x, y, color.Gray{Y: valor})
			valor += 30
		}
	}

	// Procesamos la imagen con dithering Floyd-Steinberg
	outImg, err := ProcessImageWithDithering(src, FloydStein, 3)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	// Verificamos que el resultado sea una imagen Gray
	out, ok := outImg.(*image.Gray)
	if !ok {
		t.Fatalf("resultado no es *image.Gray")
	}

	// Verificamos que todos los píxeles sean binarios (0 o 255)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			v := out.GrayAt(x, y).Y
			if v != 0 && v != 255 {
				t.Errorf("pixel (%d,%d)=%d no binario", x, y, v)
			}
		}
	}

	// Verificamos que el patrón resultante tenga una distribución de error
	// Nota: No podemos predecir exactamente el resultado, pero comprobamos
	// que haya una mezcla de blancos y negros proporcional al gradiente original
	blancos := 0
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if out.GrayAt(x, y).Y == 255 {
				blancos++
			}
		}
	}

	// Esperamos aproximadamente 5 píxeles blancos de los 9 totales
	// ya que el gradiente va de 0 a 240 con una media cercana a 120
	if blancos < 3 || blancos > 6 {
		t.Errorf("distribución de píxeles blancos (%d de 9) fuera del rango esperado", blancos)
	}
}

// TestProcessImageWithDithering_Ordered verifica el procesamiento de imágenes
// con dithering ordenado (matriz 4x4)
func TestProcessImageWithDithering_Ordered(t *testing.T) {
	// Creamos una imagen gris uniforme para ver el patrón de dithering
	src := image.NewGray(image.Rect(0, 0, 4, 4))

	// Establecemos un valor uniforme cercano al umbral (128)
	valorUniforme := uint8(120) // Justo por debajo del umbral
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			src.SetGray(x, y, color.Gray{Y: valorUniforme})
		}
	}

	// Procesamos la imagen con dithering ordenado
	outImg, err := ProcessImageWithDithering(src, Ordered, 4)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	out, ok := outImg.(*image.Gray)
	if !ok {
		t.Fatalf("resultado no es *image.Gray")
	}

	// Verificamos que todos los píxeles sean binarios
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			v := out.GrayAt(x, y).Y
			if v != 0 && v != 255 {
				t.Errorf("pixel (%d,%d)=%d no binario", x, y, v)
			}
		}
	}

	// El valor 120 está justo por debajo del umbral, por lo que el patrón
	// resultante debería tener más píxeles negros que blancos según la matriz
	blancos := 0
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if out.GrayAt(x, y).Y == 255 {
				blancos++
			}
		}
	}

	// Para el valor 120, esperamos aproximadamente 5-6 píxeles blancos de 16
	if blancos > 8 {
		t.Errorf("demasiados píxeles blancos (%d de 16) para valor uniforme 120", blancos)
	}
}

// TestResizeImageMaintainsContent verifica que el redimensionado preserva
// el contenido visual reconocible de la imagen original
func TestResizeImageMaintainsContent(t *testing.T) {
	// Creamos una imagen con un patrón simple y reconocible (una cruz)
	src := image.NewGray(image.Rect(0, 0, 5, 5))

	// Inicializamos a 0 (negro)
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			src.SetGray(x, y, color.Gray{Y: 0})
		}
	}

	// Dibujamos una cruz blanca en el centro
	for x := 0; x < 5; x++ {
		src.SetGray(x, 2, color.Gray{Y: 255}) // Línea horizontal
	}
	for y := 0; y < 5; y++ {
		src.SetGray(2, y, color.Gray{Y: 255}) // Línea vertical
	}

	// Redimensionamos a 3x3
	out := ResizeImage(src, 3)
	grayOut, ok := out.(*image.RGBA) // ResizeImage devuelve un RGBA
	if !ok {
		t.Fatalf("resultado no es *image.RGBA")
	}

	// Verificamos que la cruz siga visible en el centro
	// La cruz original tenía píxeles blancos en:
	// - La fila central (y=2)
	// - La columna central (x=2)

	// En la imagen 3x3, estos serían:
	// - La fila central (y=1)
	// - La columna central (x=1)

	// Verificamos el píxel central (debe ser blanco)
	centralPixel := grayOut.At(1, 1)
	r, g, b, _ := centralPixel.RGBA()
	if r == 0 && g == 0 && b == 0 {
		t.Errorf("píxel central debería ser blanco, pero es negro")
	}

	// La cruz debería tener al menos los píxeles horizontales y verticales alrededor del centro
	horizontalPixel := grayOut.At(0, 1)
	r, g, b, _ = horizontalPixel.RGBA()
	if r == 0 && g == 0 && b == 0 {
		t.Errorf("píxel horizontal (0,1) debería ser blanco, pero es negro")
	}

	verticalPixel := grayOut.At(1, 0)
	r, g, b, _ = verticalPixel.RGBA()
	if r == 0 && g == 0 && b == 0 {
		t.Errorf("píxel vertical (1,0) debería ser blanco, pero es negro")
	}
}
