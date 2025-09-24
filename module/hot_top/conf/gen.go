package conf

import (
	"os"
	"path/filepath"
)

func GenConf(filename string) error {
	if globalConfig == nil {
		return nil
	}
	if filename == "" {
		filename = "config.sample.toml"
	}

	if wd, err := os.Getwd(); err != nil {
		return err
	} else {
		vip.SetConfigType("toml")
		return vip.WriteConfigAs(filepath.Join(wd, filename))
	}

}
