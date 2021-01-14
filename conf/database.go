package conf

// DatabaseConfig holds all the database related configuration.
type DatabaseConfig struct {
	Driver      string `json:"driver" default:"mysql"`
	URL         string `json:"url" required:"true"`
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	AutoMigrate bool   `json:"auto_migrate"`
}
