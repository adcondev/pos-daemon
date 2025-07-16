package local_config

type LocalConfig struct {
	Data LocalConfigData `json:"data"`
}

type LocalConfigData struct {
	Printer  string
	DebugLog bool
}
