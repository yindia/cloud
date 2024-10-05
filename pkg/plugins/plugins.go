package plugins

import (
	"fmt"
	"task/pkg/plugins/email"
	"task/pkg/plugins/query"
)

// Plugin interface defines the Run method for plugins
type Plugin interface {
	Run(parameters map[string]string) error
}

// NewPlugin returns a Plugin interface based on the provided type
func NewPlugin(pluginType string) (Plugin, error) {
	switch pluginType {
	case email.PLUGIN_NAME:
		return &email.Email{}, nil // Assuming Email is a struct that implements Plugin
	case query.PLUGIN_NAME:
		return &query.Query{}, nil // Assuming Query is a struct that implements Plugin
	default:
		return nil, fmt.Errorf("unknown plugin type: %s", pluginType)
	}
}
