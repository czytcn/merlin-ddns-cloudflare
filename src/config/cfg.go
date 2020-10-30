package config

import "github.com/BurntSushi/toml"

var (
	Obj Cfg
)

type Cfg struct {
	Api ApiSetting `toml:"api"`
	App AppSetting `toml:"app"`
}

func LoadCfg(cfgFile string) error {
	_, err := toml.DecodeFile(cfgFile, &Obj)
	if err != nil {
		return err
	}
	return nil
}
