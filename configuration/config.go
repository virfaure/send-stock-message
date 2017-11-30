package configuration

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Routines  int
	Messages  int
	Skus      int
	Sources   int
	Endpoint  string
	Clients   []string
	Templates []string

	Queues struct {
		Reindex string
		Export  string
		Region  string
		Profile string
	}
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
