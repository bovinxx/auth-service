package config

import (
	"errors"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type AccessConfig interface {
	AccessRule() map[string]AccessRuleConfig
}

type AccessRuleConfig struct {
	Role   []string `mapstructure:"roles"`
	Public bool     `mapstructure:"public"`
}

type config struct {
	AccessRules map[string]AccessRuleConfig `mapstructure:"access_rules"`
}

func NewAccessConfig() (AccessConfig, error) {
	v := viper.New()
	v.SetConfigFile("./configs/access.yaml")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.New("failed to read access config")
	}

	accessRules := v.GetStringMap("access_rules")

	var cfg config
	cfg.AccessRules = make(map[string]AccessRuleConfig)

	err := mapstructure.Decode(accessRules, &cfg.AccessRules)
	if err != nil {
		return nil, errors.New("failed to parse access config")
	}

	return &cfg, nil
}

func (cfg *config) AccessRule() map[string]AccessRuleConfig {
	return cfg.AccessRules
}
