package config

import (
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	cfg, err := NewConfig()
	if err != nil {
		t.Error(err)
	}

	if cfg == nil {
		t.Error(err)
	}
}
