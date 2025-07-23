package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BoolFlex permite deserializar valores booleanos que vienen en diferentes formatos
// como true/false, 0/1, "0"/"1", "true"/"false" o cadenas vacías
type BoolFlex bool

// UnmarshalJSON implementa la interfaz json.Unmarshaler para BoolFlex
func (b *BoolFlex) UnmarshalJSON(data []byte) error {
	// Maneja el caso de null
	if string(data) == "null" {
		*b = false
		return nil
	}

	// Intenta como bool
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		*b = BoolFlex(boolVal)
		return nil
	}

	// Intenta como string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		switch strVal {
		case "1", "true":
			*b = true
		case "0", "false", "":
			*b = false
		default:
			return fmt.Errorf("valor no soportado para BoolFlex: %s", strVal)
		}
		return nil
	}

	return fmt.Errorf("no se pudo deserializar BoolFlex: %s", string(data))
}

// IntFlex permite deserializar valores enteros que pueden venir como números
// o como cadenas, incluso vacías (que se convierten a 0)
type IntFlex int

// UnmarshalJSON implementa la interfaz json.Unmarshaler para IntFlex
func (i *IntFlex) UnmarshalJSON(data []byte) error {
	// Maneja el caso de null
	if string(data) == "null" {
		*i = 0
		return nil
	}
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
			return nil
		}

		parsedInt, err := strconv.Atoi(strVal)
		if err != nil {
			return fmt.Errorf("valor no soportado para IntFlex: %s", strVal)
		}
		*i = IntFlex(parsedInt)
		return nil
	}

	return fmt.Errorf("no se pudo deserializar IntFlex: %s", string(data))
}
