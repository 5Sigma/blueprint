package blueprint

import (
	"os"

	"github.com/spf13/viper"
)

// TemplateData holds configurational data read from
// codex.yaml and provides helper methods for the templates
type TemplateData struct{ viper *viper.Viper }

// NewData creates a new data instance for the template
func NewData(viper *viper.Viper) *TemplateData {
	return &TemplateData{viper: viper}
}

// Get returns a string value at a specific key in the
// configuration
func (data *TemplateData) Get(key string) string {
	return data.viper.GetString(key)
}

// GetStrings returns an array of strings
func (data *TemplateData) GetStrings(key string) []string {
	return data.viper.GetStringSlice(key)
}

// Env returns the value of an environment variable
func (data *TemplateData) Env(key string) string {
	return os.Getenv(key)
}
