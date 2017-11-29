package app

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Routines int
	Messages int
	Skus int
	Sources int
	Queue string
	Endpoint string
	Clients []string
	Templates []string
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
