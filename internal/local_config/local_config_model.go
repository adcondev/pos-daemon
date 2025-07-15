package local_config

type Wrapper struct {
	Data LocalConfig `json:"data"`
}

type LocalConfig struct {
	Printer  string
	DebugLog bool
}
