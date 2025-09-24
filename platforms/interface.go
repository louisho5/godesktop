package platforms

// AppConfig holds the configuration for creating an app
type AppConfig struct {
	Name     string
	Width    int
	Height   int
	IconPath string
	URL      string
}

// Platform defines the interface for different platform builders
type Platform interface {
	CreateApp(config AppConfig) error
	GetFileExtension() string
	GetPlatformName() string
}
