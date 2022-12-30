package config

// Config is the configuration for the program
type Config struct {
	UrlFilePath string
}

// NewConfig returns a new Config
func NewConfig() Config {
	return Config{
		UrlFilePath: "urls.json",
	}
}
