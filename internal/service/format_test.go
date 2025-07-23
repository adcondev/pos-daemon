package service

import "testing"

func TestCountChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single character", "a", 1},
		{"Multiple characters", "hello", 5},
		{"Unicode characters", "こんにちは", 5},
		{"Emoji characters", "😊👍", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountChars(tt.input)
			if result != tt.expected {
				t.Errorf("CountChars(%q) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		width    int
		padChar  rune
		expected string
	}{
		{"String más corta que width", "abc", 5, ' ', "  abc"},
		{"String igual a width", "abc", 3, ' ', "abc"},
		{"String más larga que width", "abcde", 3, ' ', "abcde"},
		{"Padding con carácter especial", "abc", 5, '*', "**abc"},
		{"String vacía", "", 3, '-', "---"},
		{"Width cero", "abc", 0, ' ', "abc"},
		{"Caracteres Unicode", "你好", 4, '~', "~~你好"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadLeft(tt.s, tt.width, tt.padChar)
			if result != tt.expected {
				t.Errorf("PadLeft(%q, %d, %q) = %q; want %q",
					tt.s, tt.width, tt.padChar, result, tt.expected)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		width    int
		padChar  rune
		expected string
	}{
		{"String más corta que width", "abc", 5, ' ', "abc  "},
		{"String igual a width", "abc", 3, ' ', "abc"},
		{"String más larga que width", "abcde", 3, ' ', "abcde"},
		{"Padding con carácter especial", "abc", 5, '*', "abc**"},
		{"String vacía", "", 3, '-', "---"},
		{"Width cero", "abc", 0, ' ', "abc"},
		{"Caracteres Unicode", "你好", 4, '~', "你好~~"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadRight(tt.s, tt.width, tt.padChar)
			if result != tt.expected {
				t.Errorf("PadRight(%q, %d, %q) = %q; want %q",
					tt.s, tt.width, tt.padChar, result, tt.expected)
			}
		})
	}
}

func TestPadCenter(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		width    int
		padChar  rune
		expected string
	}{
		{"String más corta que width (padding par)", "abc", 7, ' ', "  abc  "},
		{"String más corta que width (padding impar)", "abc", 6, ' ', " abc  "},
		{"String igual a width", "abc", 3, ' ', "abc"},
		{"String más larga que width", "abcde", 3, ' ', "abcde"},
		{"Padding con carácter especial", "abc", 7, '*', "**abc**"},
		{"String vacía", "", 4, '-', "----"},
		{"Width cero", "abc", 0, ' ', "abc"},
		{"Caracteres Unicode", "你好", 6, '~', "~~你好~~"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadCenter(tt.s, tt.width, tt.padChar)
			if result != tt.expected {
				t.Errorf("PadCenter(%q, %d, %q) = %q; want %q",
					tt.s, tt.width, tt.padChar, result, tt.expected)
			}
		})
	}
}

func TestSubstr(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		length   int
		expected string
	}{
		{"Extraer menos caracteres que la longitud", "abcde", 3, "abc"},
		{"Extraer todos los caracteres", "abc", 3, "abc"},
		{"Longitud mayor que string", "abc", 5, "abc"},
		{"String vacía", "", 3, ""},
		{"Longitud cero", "abc", 0, ""},
		{"Caracteres Unicode", "你好世界", 2, "你好"},
		{"Extraer un carácter", "abcde", 1, "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Substr(tt.s, tt.length)
			if result != tt.expected {
				t.Errorf("Substr(%q, %d) = %q; want %q",
					tt.s, tt.length, result, tt.expected)
			}
		})
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name     string
		f        float64
		decimals int
		expected string
	}{
		{"Dos decimales", 123.456, 2, "123.46"},
		{"Sin decimales", 123.456, 0, "123"},
		{"Más decimales que los presentes", 123.4, 3, "123.400"},
		{"Número negativo", -123.456, 2, "-123.46"},
		{"Cero", 0.0, 2, "0.00"},
		{"Valor pequeño", 0.00123, 5, "0.00123"},
		{"Redondeo hacia arriba", 0.12345, 4, "0.1235"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFloat(tt.f, tt.decimals)
			if result != tt.expected {
				t.Errorf("FormatFloat(%f, %d) = %q; want %q",
					tt.f, tt.decimals, result, tt.expected)
			}
		})
	}
}

func TestSplitString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		length   int
		expected []string
	}{
		{"Dividir en partes iguales", "abcdef", 2, []string{"ab", "cd", "ef"}},
		{"Última parte más corta", "abcde", 2, []string{"ab", "cd", "e"}},
		{"Longitud igual al string", "abcde", 5, []string{"abcde"}},
		{"Longitud mayor que string", "abc", 5, []string{"abc"}},
		{"String vacía", "", 2, []string{}},
		{"Longitud cero", "abc", 0, nil},
		{"Longitud negativa", "abc", -1, nil},
		{"Caracteres Unicode", "你好世界", 2, []string{"你好", "世界"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitString(tt.s, tt.length)

			if tt.expected == nil && result != nil {
				t.Errorf("SplitString(%q, %d) = %v; want nil",
					tt.s, tt.length, result)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("SplitString(%q, %d) = %v; want %v",
					tt.s, tt.length, result, tt.expected)
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("SplitString(%q, %d)[%d] = %q; want %q",
						tt.s, tt.length, i, result[i], tt.expected[i])
				}
			}
		})
	}
}
