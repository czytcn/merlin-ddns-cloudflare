package config

import "testing"

func TestLoadCfg(t *testing.T) {
	err := LoadCfg("../config.toml")
	if err != nil {
		t.Errorf("err:%v", err)
	}
	t.Log(O)
}
