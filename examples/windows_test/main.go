// examples/windows_test/main.go
package main

import (
	"github.com/skip2/go-qrcode"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	conn "pos-daemon.adcon.dev/pkg/posprinter/connector"
	imaging2 "pos-daemon.adcon.dev/pkg/posprinter/imaging"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
	"time"
)

func main() {
	// Configuración básica del logger
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Iniciando prueba de impresora...")

	// --- CONFIGURACIÓN DE LA IMPRESORA ---
	// Define el nombre EXACTO de la impresora instalada en Windows
	// Editar este valor según tu configuración
	printerName := "58mm GOOJPRT PT-210" // Cambia esto al nombre de tu impresora
	// printerName := "80mm EC-PM-80250"

	log.Printf("Intentando conectar a la impresora: %s", printerName)

	// Crear el conector para la impresora
	connector, err := conn.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Error al crear el conector para '%s': %v", printerName, err)
	}

	defer func() {
		log.Println("Cerrando el conector de la impresora.")
		if closeErr := connector.Close(); closeErr != nil {
			log.Printf("Error al cerrar el conector: %v", closeErr)
		}
	}()
	// Después de crear el conector y abrir la conexión
	log.Println("Conector Bluetooth creado exitosamente.")

	// Añadir un retraso para estabilizar la conexión antes de inicializar
	log.Println("Esperando para estabilizar conexión Bluetooth...")
	time.Sleep(2 * time.Second)

	// Crear instancia de ESCPrinter
	printer, err := escpos.NewPrinter(connector, nil)
	if err != nil {
		log.Fatalf("Error al crear e inicializar la impresora: %v", err)
	}
	log.Println("Impresora inicializada correctamente.")

	// --- PRUEBAS DE IMPRESIÓN ---

	// Encabezado
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetEmphasis(true); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("*** PRUEBA DE IMPRESORA ***\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if err = printer.Text(timestamp + "\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetEmphasis(false); err != nil {
		log.Printf("Error: %v", err)
	}

	// --- PRUEBA 1: ALINEACIONES ---
	if err = printer.Text("=== PRUEBA DE ALINEACIONES ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Izquierda
	if err = printer.SetJustification(escpos.Left); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Texto alineado a la IZQUIERDA\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Centro
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Texto alineado al CENTRO\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Derecha
	if err = printer.SetJustification(escpos.Right); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Texto alineado a la DERECHA\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// --- PRUEBA 2: ESTILOS DE TEXTO ---
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("=== PRUEBA DE ESTILOS ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Left); err != nil {
		log.Printf("Error: %v", err)
	}

	// Normal
	if err = printer.Text("Texto normal\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Énfasis (negrita)
	if err = printer.SetEmphasis(true); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Texto con énfasis (negrita)\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetEmphasis(false); err != nil {
		log.Printf("Error: %v", err)
	}

	// Subrayado
	if err = printer.SetUnderline(1); err != nil { // 1 = subrayado simple
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Texto con subrayado simple\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetUnderline(2); err != nil { // 2 = subrayado doble
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Texto con subrayado doble\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetUnderline(0); err != nil { // 0 = sin subrayado
		log.Printf("Error: %v", err)
	}

	// --- PRUEBA ADICIONAL: VOCALES ACENTUADAS ---
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("=== PRUEBA DE ACENTOS ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Left); err != nil {
		log.Printf("Error: %v", err)
	}

	// Minúsculas con acentos
	if err = printer.Text("Vocales minúsculas: á é í ó ú ü ñ\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Mayúsculas con acentos
	if err = printer.Text("Vocales mayúsculas: Á É Í Ó Ú Ü Ñ\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Ejemplo de oración con acentos
	if err = printer.Text("Oración: El niño está jugando con el pingüino.\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// --- PRUEBA 3: TAMAÑOS DE FUENTE ---
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("=== PRUEBA DE TAMAÑOS ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Tamaño normal
	if err = printer.SetTextSize(1, 1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Tamaño Normal (1,1)\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Ancho doble
	if err = printer.SetTextSize(2, 1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Ancho Doble (2,1)\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Alto doble
	if err = printer.SetTextSize(1, 2); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Alto Doble (1,2)\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Doble ancho y alto
	if err = printer.SetTextSize(2, 2); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Doble (2,2)\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Cambiar a fuente B y mostrar los mismos tamaños
	if err = printer.SetFont(1); err != nil { // Fuente B
		log.Printf("Error: %v", err)
	}

	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Fuente B:\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Tamaño normal con fuente B
	if err = printer.SetTextSize(1, 1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Tamaño Normal (1,1)\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Ancho doble con fuente B
	if err = printer.SetTextSize(2, 1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Ancho Doble (2,1)\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Alto doble con fuente B
	if err = printer.SetTextSize(1, 2); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Alto Doble (1,2)\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Doble ancho y alto con fuente B
	if err = printer.SetTextSize(2, 2); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Doble (2,2)\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Restablecer tamaño normal y fuente A
	if err = printer.SetTextSize(1, 1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetFont(0); err != nil {
		log.Printf("Error: %v", err)
	}

	// --- PRUEBA 4: FUENTES ---
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("=== PRUEBA DE FUENTES ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Fuente A (predeterminada)
	if err = printer.SetFont(0); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Fuente A (0) - Estándar\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Fuente B (más pequeña/condensada)
	if err = printer.SetFont(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Fuente B (1) - Condensada\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Restaurar fuente A
	if err = printer.SetFont(0); err != nil {
		log.Printf("Error: %v", err)
	}

	// --- PRUEBA 5: ANCHO DE LÍNEA ---
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("=== PRUEBA DE ANCHO DE LÍNEA ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Left); err != nil {
		log.Printf("Error: %v", err)
	}

	// Fuente A - caracteres por línea
	if err = printer.SetFont(0); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Fuente A - 75 caracteres:\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("123456789012345678901234567890123456789012345678901234567890123456789012345\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Fuente B - caracteres por línea
	if err = printer.SetFont(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("Fuente B - 75 caracteres:\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("123456789012345678901234567890123456789012345678901234567890123456789012345\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Restaurar fuente A
	if err = printer.SetFont(0); err != nil {
		log.Printf("Error: %v", err)
	}

	// --- PRUEBA 6: CÓDIGO DE BARRAS ---
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("=== PRUEBA DE CÓDIGO DE BARRAS ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Configurar el código de barras
	if err = printer.SetBarcodeHeight(80); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetBarcodeWidth(3); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetBarcodeTextPosition(escpos.TextBelow); err != nil {
		log.Printf("Error: %v", err)
	}

	// UPC-A
	if err = printer.Text("UPC-A:\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Barcode("012345678901", escpos.UpcA); err != nil {
		log.Printf("Error: %v", err)
	}

	// CODE39
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("CODE39:\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Barcode("CODE39TEST", escpos.Code39); err != nil {
		log.Printf("Error: %v", err)
	}

	// --- PRUEBA 7: CÓDIGO QR ---
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("=== PRUEBA DE CÓDIGO QR ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	qrText := "https://github.com/AdConDev/pos-daemon"
	if err = printer.Text("QR: " + qrText + "\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Generar código QR con diferentes niveles de corrección
	// Nivel Alto (30% recuperación de datos)
	qr, err := qrcode.New(qrText, qrcode.High)
	if err != nil {
		log.Fatalf("Error generando QR: %v", err)
	}

	// Obtener imagen del QR y convertirla
	qrImage := qr.Image(256)
	escposQR := escpos.NewEscposImage(qrImage, 128)

	// Imprimir QR
	if err = printer.BitImage(escposQR, imaging2.ImgDefault); err != nil {
		log.Printf("Error al imprimir QR: %v", err)
	}

	// --- PRUEBA 8: IMAGEN ---
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.SetJustification(escpos.Center); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Text("=== PRUEBA DE IMAGEN ===\n\n"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Cargar imagen - intenta primero JPEG, luego PNG
	logoPath := "./img/perro.jpeg"
	if _, err = os.Stat(logoPath); os.IsNotExist(err) {
		logoPath = "./img/perro.png"
		if _, err = os.Stat(logoPath); os.IsNotExist(err) {
			logoPath = "./img/logo.jpeg"
			if _, err = os.Stat(logoPath); os.IsNotExist(err) {
				logoPath = "./img/logo.png"
			}
		}
	}

	// Abrir el archivo de imagen
	logoFile, err := os.Open(logoPath)
	if err != nil {
		log.Printf("Error abriendo imagen (%s): %v", logoPath, err)
	} else {
		defer func(logoFile *os.File) {
			err = logoFile.Close()
			if err != nil {
				log.Printf("Error al cerrar el archivo de imagen: %v", err)
			}
		}(logoFile)

		// Decodificar imagen
		imgLogo, format, err := image.Decode(logoFile)
		if err != nil {
			log.Printf("Error decodificando imagen: %v", err)
		} else {
			log.Printf("Imagen cargada: %s (formato %s)", logoPath, format)

			// Imprimir imagen con dithering para mejor calidad
			if err := printer.ImageWithDithering(imgLogo, imaging2.ImgDefault, imaging2.FloydStein, imaging2.DefaultPrintSize); err != nil {
				log.Printf("Error al imprimir imagen: %v", err)
			}
		}
	}

	// Finalizar impresión
	if err = printer.Feed(4); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Cut(escpos.CUT_FULL, 0); err != nil {
		log.Printf("Error: %v", err)
	}

	log.Println("Prueba de impresión completada.")
}
