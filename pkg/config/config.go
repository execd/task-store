package config

import (
	"github.com/execd/task-store/pkg/model"
	"github.com/BurntSushi/toml"
)

type Parser interface {
	ParseConfig(configPath string) *model.Config
}

type ParserImpl struct{}

func NewParserImpl() *ParserImpl {
	return &ParserImpl{}
}

func (p *ParserImpl) ParseConfig(configPath string) *model.Config {
	config := new(model.Config)
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		panic(err.Error())
	}
	return config
}
