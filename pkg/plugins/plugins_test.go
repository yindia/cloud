package plugins

import (
	"task/pkg/plugins/email"
	"task/pkg/plugins/query"
	"testing"
)

func TestNewPlugin(t *testing.T) {
	tests := []struct {
		name       string
		pluginType string
		want       Plugin
		wantErr    bool
	}{
		{
			name:       "SEND_EMAIL plugin",
			pluginType: email.PLUGIN_NAME,
			want:       &email.Email{},
			wantErr:    false,
		},
		{
			name:       "RUN_QUERY plugin",
			pluginType: query.PLUGIN_NAME,
			want:       &query.Query{},
			wantErr:    false,
		},
		{
			name:       "Unknown plugin type",
			pluginType: "UNKNOWN",
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPlugin(tt.pluginType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlugin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Errorf("NewPlugin() returned nil, want %T", tt.want)
				return
			}
			if _, ok := got.(Plugin); !ok {
				t.Errorf("NewPlugin() returned %T, which does not implement Plugin interface", got)
			}
		})
	}
}
