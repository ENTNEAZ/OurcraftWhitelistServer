package Config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var config ConfigStruct

func GetConfig() ConfigStruct {
	return config
}

func LoadConfig(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	if config.HTMLPath[len(config.HTMLPath)-1] != '/' {
		config.HTMLPath += "/"
	}

	if config.MiraiAddr[len(config.MiraiAddr)-1] != '/' {
		config.MiraiAddr += "/"
	}

	return nil
}
