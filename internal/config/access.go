package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/bovinxx/auth-service/internal/models"
	"github.com/spf13/viper"
)

type AccessConfig interface {
	AccessRule() map[string]AccessRuleConfig
}

type AccessRuleConfig struct {
	Roles  []models.Role `mapstructure:"roles"`
	Public bool          `mapstructure:"public"`
}

type config struct {
	AccessRules map[string]AccessRuleConfig
}

func NewAccessConfig() (AccessConfig, error) {
	v := viper.New()
	v.SetConfigFile("./configs/access.yaml")
	v.SetConfigType("yaml")

	v.SetTypeByDefaultValue(true)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read access config: %w", err)
	}

	accessRules := v.GetStringMap("access_rules")

	cfg := config{
		AccessRules: make(map[string]AccessRuleConfig),
	}

	for path, rule := range accessRules {
		ruleMap, ok := rule.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid rule format for path %s", path)
		}

		rolesInterface, ok := ruleMap["roles"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid roles format for path %s", path)
		}

		var roles []models.Role
		for _, roleInterface := range rolesInterface {
			roleStr, ok := roleInterface.(string)
			if !ok {
				return nil, fmt.Errorf("invalid role type for path %s", path)
			}

			role, err := stringToRole(roleStr)
			if err != nil {
				return nil, fmt.Errorf("invalid role in path %s: %w", path, err)
			}
			roles = append(roles, role)
		}

		public, ok := ruleMap["public"].(bool)
		if !ok {
			return nil, fmt.Errorf("invalid public format for path %s", path)
		}

		cfg.AccessRules[path] = AccessRuleConfig{
			Roles:  roles,
			Public: public,
		}
	}

	log.Printf("Parsed config: %+v", cfg.AccessRules)
	return &cfg, nil
}

func stringToRole(roleStr string) (models.Role, error) {
	switch roleStr {
	case "admin":
		return models.RoleAdmin, nil
	case "user":
		return models.RoleUser, nil
	default:
		return 0, fmt.Errorf("unknown role: %s", roleStr)
	}
}

func (cfg *config) AccessRule() map[string]AccessRuleConfig {
	return cfg.AccessRules
}
