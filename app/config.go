package app

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Messages int
	Skus int
	Sources int
	Endpoint string
}

func Load(file string) (cfg Config, err error) {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return
	}

	return
}
