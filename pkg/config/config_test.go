package config

import (
	"github.com/execd/task-store/pkg/model"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var context = GinkgoT()

var _ = Describe("config", func() {
	Describe("parsing a config file", func() {
		It("should parse a well defined config file", func() {
			// Arrange
			parser := NewParserImpl()
			expectedConfig := &model.Config{
				Manager: model.ManagerInfo{
					ExecutionQueueSize: 10,
					TaskQueueSize:      10,
				},
			}

			// Act
			config := parser.ParseConfig("config.toml")

			// Assert
			assert.Equal(context, expectedConfig, config)
		})

		It("should panic if it fails to build config", func() {
			// Arrange
			parser := NewParserImpl()

			// Assert
			assert.Panics(context, func() { parser.ParseConfig("nada") })
		})
	})
})
