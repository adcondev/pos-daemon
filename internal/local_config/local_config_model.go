package local_config

type LocalConfig struct {
	Data LocalConfigData `json:"data"`
}

// Add these fields to your LocalConfigData struct
type LocalConfigData struct {
	// existing fields
	Printer  string `json:"printer"`
	DebugLog bool   `json:"debug_log"`

	// Add these new fields for serial configuration
	SerialBaudRate int    `json:"serial_baud_rate,omitempty"`
	SerialDataBits int    `json:"serial_data_bits,omitempty"`
	SerialStopBits int    `json:"serial_stop_bits,omitempty"`
	SerialParity   string `json:"serial_parity,omitempty"`
}
