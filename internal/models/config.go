package models

// LocalConfig es la estructura raíz para la configuración local
type Config struct {
	Data ConfigData `json:"data"`
}

// ConfigData contiene la configuración local de la aplicación
type ConfigData struct {
	// Configuración general
	Printer  string `json:"printer"`   // Nombre de la impresora a utilizar
	DebugLog bool   `json:"debug_log"` // Habilitar logs de depuración

	// Configuración de puerto serial
	SerialBaudRate int    `json:"serial_baud_rate"` // Velocidad en baudios
	SerialDataBits int    `json:"serial_data_bits"` // Bits de datos (típicamente 8)
	SerialStopBits int    `json:"serial_stop_bits"` // Bits de parada (típicamente 1)
	SerialParity   string `json:"serial_parity"`    // Paridad (none, odd, even)
}
