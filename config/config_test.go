package config

import "testing"

func TestNewConfigInstance(t *testing.T) {
	err := NewConfigInstance()
	if err != nil {
		t.Error(err)
	}
	t.Logf("ConfigInstance=[%+v]", ConfigInstance)
}
