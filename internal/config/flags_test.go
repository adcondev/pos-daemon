package config_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"pos-daemon.adcon.dev/internal/config"
)

func ExampleParseFlags() {
	// Simular flags de línea de comandos
	os.Args = []string{"programa", "-printer", "EC-PM-80250", "-debug"}
	
	// Reiniciar flag para el ejemplo
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	
	cfg := config.ParseFlags()
	
	// Verificar los valores
	if cfg.Printer == "EC-PM-80250" && cfg.DebugLog {
		fmt.Println("Configuración cargada correctamente")
	}
	
	// Output: Configuración cargada correctamente
}

func TestParseFlags(t *testing.T) {
	// Guardar estado original
	originalArgs := os.Args
	originalCommandLine := flag.CommandLine
	
	// Restaurar estado al final
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalCommandLine
	}()
	
	// Casos de prueba
	tests := []struct {
		name          string
		args          []string
		expectedPrinter string
		expectedDebug   bool
	}{
		{
			name:          "valores por defecto",
			args:          []string{"programa"},
			expectedPrinter: "",
			expectedDebug:   false,
		},
		{
			name:          "solo printer",
			args:          []string{"programa", "-printer", "EC-PM-80250"},
			expectedPrinter: "EC-PM-80250",
			expectedDebug:   false,
		},
		{
			name:          "solo debug",
			args:          []string{"programa", "-debug"},
			expectedPrinter: "",
			expectedDebug:   true,
		},
		{
			name:          "ambos parámetros",
			args:          []string{"programa", "-printer", "TestPrinter", "-debug"},
			expectedPrinter: "TestPrinter",
			expectedDebug:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar args para esta prueba
			os.Args = tt.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			
			cfg := config.ParseFlags()
			
			if cfg.Printer != tt.expectedPrinter {
				t.Errorf("Printer = %v, esperado %v", cfg.Printer, tt.expectedPrinter)
			}
			
			if cfg.DebugLog != tt.expectedDebug {
				t.Errorf("DebugLog = %v, esperado %v", cfg.DebugLog, tt.expectedDebug)
			}
		})
	}
}