package config

import (
	"github.com/antihax/optional"
	"github.com/spf13/viper"
)

const (
	OutputDefault = ""
	OutputTable   = "table"
	OutputJson    = "json"
	OutputYaml    = "yaml"
)

func GetOptionalString(flag string) optional.String {
	val := viper.GetString(flag)
	if val == "" {
		return optional.EmptyString()
	}

	return optional.NewString(val)
}
