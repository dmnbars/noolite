package config

import (
	"encoding/json"
	"reflect"

	"github.com/caarlos0/env/v6"
	"gitlab.com/dmnbars/noolite/internal/homeassistant"
)

type Config struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`

	SerialPort string `env:"SERIAL_PORT"`

	MqttHost     string `env:"MQTT_HOST"`
	MqttPort     string `env:"MQTT_PORT"`
	MqttClientID string `env:"MQTT_CLIENT_ID"`
	MqttUsername string `env:"MQTT_USERNAME,unset"`
	MqttPassword string `env:"MQTT_PASSWORD,unset"`

	PowerOutlets []homeassistant.PowerOutletConfig `env:"POWER_OUTLETS" envDefault:"[]"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.ParseWithFuncs(&cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf([]homeassistant.PowerOutletConfig{}): powerOutletsParser,
	}); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

func powerOutletsParser(value string) (interface{}, error) {
	var powerOutlets []homeassistant.PowerOutletConfig
	if err := json.Unmarshal([]byte(value), &powerOutlets); err != nil {
		return nil, err
	}

	return powerOutlets, nil
}
