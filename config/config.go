// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"github.com/kussj/jolokiabeat/common"
)

type Config struct {
	Jolokiabeat JolokiabeatConfig
}

type JolokiabeatConfig struct {
	Period	string `yaml:"period"`
	Url     string `yaml:"url"`	
	Queries []common.QueryConfig `yaml:"queries"`
}
