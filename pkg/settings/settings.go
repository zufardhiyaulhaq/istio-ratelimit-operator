package settings

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	RateLimitServiceImage string `required:"true" envconfig:"RATE_LIMIT_SERVICE_IMAGE" default:"envoyproxy/ratelimit:4c088856"`
}

func NewSettings() (Settings, error) {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
