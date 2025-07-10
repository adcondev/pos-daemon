package config

type Wrapper struct {
	Data Config `json:"data"`
}

type Config struct {
	Printer  string
	DebugLog bool
}
