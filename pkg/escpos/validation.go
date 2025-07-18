package escpos

import (
	"fmt"
	cons "pos-daemon.adcon.dev/pkg/escpos/constants"
	"regexp"
	"strings"
)

// --- Funciones de Ayuda para Validación ---
// Estas funciones validan los argumentos y devuelven un error si son inválidos.

// validateBoolean es en gran parte redundante en Go debido al tipado estático.

func validateFloat(test float64, min, max float64, source, argument string) error {
	if test < min || test > max {
		return fmt.Errorf("el argumento '%s' (%f) dado a %s debe estar en el rango %f a %f", argument, test, source, min, max)
	}
	return nil
}

func validateInteger(test, min, max int, source, argument string) error {
	return validateIntegerMulti(test, [][]int{{min, max}}, source, argument)
}

func ValidateBarcodeTextPosition(pos cons.BarcodeTextPos) error {
	if pos < cons.TextNone || pos > cons.TextBoth {
		return fmt.Errorf("posición de texto inválida: %d", pos)
	}
	return nil
}

func ValidateBarcodeType(barcode cons.BarcodeType) error {
	switch barcode {
	case cons.UpcA, cons.UpcE, cons.Jan13, cons.Jan8,
		cons.Code39, cons.Itf, cons.Codabar,
		cons.Code93, cons.Code128:
		return nil
	default:
		return fmt.Errorf("tipo de código de barras inválido: %d", barcode)
	}
}

func ValidateJustifyMode(mode cons.Justify) error {
	switch mode {
	case cons.Left, cons.Right, cons.Center:
		return nil
	default:
		return fmt.Errorf("tipo de justificación inválida: %d", mode)
	}
}

func ValidateFont(font cons.Font) error {
	switch font {
	case cons.A, cons.B:
		return nil
	default:
		return fmt.Errorf("tipo de fuente inválida: %d", font)
	}
}

func ValidateUnderline(under cons.UnderlineMode) error {
	switch under {
	case cons.NoUnderline, cons.Single, cons.Double:
		return nil
	default:
		return fmt.Errorf("tipo de fuente inválida: %d", under)
	}
}

func validateIntegerMulti(test int, ranges [][]int, source, argument string) error {
	match := false
	for _, r := range ranges {
		if len(r) != 2 {
			// Esto indica un error interno en cómo se llama a esta función de validación
			return fmt.Errorf("error interno: validateIntegerMulti recibió un rango inválido %v", r)
		}
		if test >= r[0] && test <= r[1] {
			match = true
			break
		}
	}

	if !match {
		// Construir el mensaje de rango similar a PHP
		rangeStrs := make([]string, len(ranges))
		for i, r := range ranges {
			rangeStrs[i] = fmt.Sprintf("%d-%d", r[0], r[1])
		}
		rangeStr := strings.Join(rangeStrs, ", ")
		if len(ranges) > 1 {
			// Reemplazar la última coma con " o " si hay más de un rango
			lastCommaIndex := strings.LastIndex(rangeStr, ", ")
			if lastCommaIndex != -1 {
				rangeStr = rangeStr[:lastCommaIndex+2] + "o " + rangeStr[lastCommaIndex+2:]
			}
		}

		return fmt.Errorf("el argumento '%s' (%d) dado a %s debe estar en el rango %s", argument, test, source, rangeStr)
	}
	return nil
}

// validateString es en gran parte redundante en Go debido al tipado estático.
// El chequeo de PHP sobre objetos con __toString no aplica directamente en Go.
func validateString(test string, source, argument string) error {
	// En Go, el tipado estático ya asegura que es una cadena si el argumento es string.
	// La función se mantiene por completitud del port, pero siempre devuelve nil.
	return nil
}

// Cache para expresiones regulares compiladas
var regexCache = make(map[string]*regexp.Regexp)

func validateStringRegex(test string, regexPattern string, argument string) error {
	// Compilar la regex si no está en caché
	re, ok := regexCache[regexPattern]
	if !ok {
		var err error
		re, err = regexp.Compile(regexPattern)
		if err != nil {
			// Error interno: la regex proporcionada no es válida
			return fmt.Errorf("error interno: falló al compilar la regex '%s': %w", regexPattern, err)
		}
		regexCache[regexPattern] = re
	}

	if !re.MatchString(test) {
		// El mensaje de error de PHP incluía el nombre de la función fuente,
		// pero aquí el argumento 'argument' ya describe qué valor es.
		return fmt.Errorf("el argumento '%s' ('%s') es inválido. Debe coincidir con la regex '%s'", argument, test, regexPattern)
	}
	return nil
}
