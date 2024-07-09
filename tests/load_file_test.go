package test

import (
	"hugo_backend/config"
	"testing"
)

func TestLoadFromFile_FromEmbedded(t *testing.T) {
	err := config.LoadFromFile()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Loaded config from embedded file:%v", config.GlobalConfigInstance.UserName)
}
