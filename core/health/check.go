package health

import "github.com/jrapoport/gothic/config"

// Status is the health check status string
const Status = "bela lugosi's dead"

// Health check response
type Health struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

// Check returns a health check response
func Check(c *config.Config) Health {
	return Health{
		c.Name,
		c.Version(),
		Status,
	}
}
