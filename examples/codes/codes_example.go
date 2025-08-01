package main

import (
	"fmt"
	"log"
	"pos-daemon.adcon.dev/pkg/posprinter"
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/encoding"
	"pos-daemon.adcon.dev/pkg/posprinter/profile"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol/escpos"
)

func main() {
	// Configuración de impresoras para probar
	printers := []struct {
		Name     string
		CharSets []int // Charsets reportados por el fabricante
	}{
		{
			Name:     "80mm EC-PM-80250 x",
			CharSets: []int{16},
		},
		{
			Name:     "58mm GOOJPRT PT-210",
			CharSets: []int{8, 9, 10, 17, 18, 20, 21},
		},
		{
			Name:     "58mm GP-58N x",
			CharSets: []int{},
		},
		// Agregar tu tercera impresora aquí
	}

	// Texto de prueba con caracteres especiales en español
	testTexts := []string{
		"Acentos: áéíóú ÁÉÍÓÚ",
		"Eñe: ñ Ñ",
		"Diéresis: ü Ü",
		"Símbolos: ¡¿ € $ ¢ £ ¥",
		"Especiales: ª º",
	}

	// Probar cada impresora
	for _, printer := range printers {
		fmt.Printf("\n=== Probando %s ===\n", printer.Name)

		// Conectar a la impresora
		conn, err := connector.NewWindowsPrintConnector(printer.Name)
		if err != nil {
			log.Printf("Error conectando a %s: %v", printer.Name, err)
			continue
		}
		defer func(conn *connector.WindowsPrintConnector) {
			err := conn.Close()
			if err != nil {
				log.Printf("Error al cerrar el conector de %s: %v", printer.Name, err)
			}
		}(conn)

		// Crear perfil personalizado
		prof := profile.CreateProfile80mm()
		prof.CharacterSets = printer.CharSets
		prof.Model = printer.Name

		// Crear protocolo e impresora
		proto := escpos.NewESCPOSProtocol()
		p, err := posprinter.NewGenericPrinter(proto, conn, prof)
		if err != nil {
			log.Printf("Error creando impresora: %v", err)
			continue
		}
		defer func(p *posprinter.GenericPrinter) {
			err := p.Close()
			if err != nil {
				log.Printf("Error al cerrar la impresora %s: %v", printer.Name, err)
			}
		}(p)

		// Imprimir encabezado
		if err := p.SetJustification(command.AlignCenter); err != nil {
			log.Printf("Error estableciendo alineación centrada: %v", err)
			continue
		}

		if err := p.SetEmphasis(true); err != nil {
			log.Printf("Error activando negrita: %v", err)
			continue
		}
		err = p.TextLn(fmt.Sprintf("TEST CODIFICACIÓN - %s", printer.Name))
		if err != nil {
			log.Printf("Error imprimiendo encabezado: %v", err)
			continue
		}

		if err := p.SetEmphasis(false); err != nil {
			log.Printf("Error desactivando negrita: %v", err)
			continue
		}

		if err := p.Feed(1); err != nil {
			log.Printf("Error alimentando papel: %v", err)
			continue
		}

		// Probar cada charset soportado
		for _, charset := range printer.CharSets {
			// Verificar que el charset esté en nuestro Registry
			if _, exists := encoding.Registry[charset]; !exists {
				continue
			}

			err := p.SetJustification(command.AlignLeft)
			if err != nil {
				log.Printf("Error estableciendo alineación izquierda: %v", err)
				continue
			}

			if err := p.SetEmphasis(true); err != nil {
				log.Printf("Error activando negrita: %v", err)
				continue
			}
			err = p.TextLn(fmt.Sprintf("=== Charset %d (%s) ===",
				charset, encoding.Registry[charset].Name))
			if err != nil {
				log.Printf("Error imprimiendo encabezado de charset: %v", err)
				continue
			}
			if err := p.SetEmphasis(false); err != nil {
				log.Printf("Error desactivando negrita: %v", err)
				continue
			}

			// Cambiar al charset
			if err := p.SetCharacterSet(charset); err != nil {
				err := p.TextLn(fmt.Sprintf("Error: %v", err))
				if err != nil {
					log.Printf("Error imprimiendo mensaje de error: %v", err)
					continue
				}
				continue
			}

			// Imprimir textos de prueba
			for _, text := range testTexts {
				if err := p.TextLn(text); err != nil {
					err := p.TextLn(fmt.Sprintf("Error imprimiendo: %v", err))
					if err != nil {
						log.Printf("Error imprimiendo texto: %v", err)
						continue
					}
				}
			}

			if err := p.Feed(1); err != nil {
				log.Printf("Error alimentando papel: %v", err)
				continue
			}
		}

		// Cortar

		if err := p.Feed(2); err != nil {
			log.Printf("Error alimentando papel: %v", err)
			continue
		}

		if err := p.Cut(command.CutFeed, 2); err != nil {
			log.Printf("Error cortando papel: %v", err)
			continue
		}
	}
}
