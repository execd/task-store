package config

import (
	"github.com/BurntSushi/toml"
	"github.com/execd/task-store/pkg/model"
)

// Parser : config parser
type Parser interface {
	ParseConfig(configPath string) *model.Config
}

// ParserImpl : implementation of a config parser
type ParserImpl struct{}

// NewParserImpl : build a new config parser
func NewParserImpl() *ParserImpl {
	return &ParserImpl{}
}

// ParseConfig : parse the given config file, panics if parsing fails
func (p *ParserImpl) ParseConfig(configPath string) *model.Config {
	config := new(model.Config)
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		panic(err.Error())
	}
	return config
}
