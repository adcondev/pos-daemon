package internal

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BoolFlex Necesitado ya que Go no entiende 0 y 1 como true/false
type BoolFlex bool

func (b *BoolFlex) UnmarshalJSON(data []byte) error {
	// Intenta como bool
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		*b = BoolFlex(boolVal)
		return nil
	}
	// Intenta como string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		if strVal == "1" {
			*b = true
		} else if strVal == "0" {
			*b = false
		} else if strVal == "true" {
			*b = true
		} else if strVal == "false" {
			*b = false
		} else if strVal == "" {
			*b = false
		} else {
			return fmt.Errorf("valor no soportado para BoolFlex: %s", strVal)
		}
		return nil
	}
	return fmt.Errorf("no se pudo deserializar BoolFlex: %s", string(data))
}

// IntFlex Necesitado para manejar enteros con cadenas vacías como 0
type IntFlex int

func (i *IntFlex) UnmarshalJSON(data []byte) error {
	// Intenta como entero
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*i = IntFlex(intVal)
		return nil
	}
	// Intenta como string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		if strVal == "" {
			*i = 0
		} else {
			parsedInt, err := strconv.Atoi(strVal)
			if err != nil {
				return fmt.Errorf("valor no soportado para IntFlex: %s", strVal)
			}
			*i = IntFlex(parsedInt)
		}
		return nil
	}
	return fmt.Errorf("no se pudo deserializar IntFlex: %s", string(data))
}
