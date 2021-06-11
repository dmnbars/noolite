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

	Switches []homeassistant.SwitchConfig `env:"SWITCHES" envDefault:"[]"`
	Lights   []homeassistant.LightConfig  `env:"LIGHTS" envDefault:"[]"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.ParseWithFuncs(&cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf([]homeassistant.SwitchConfig{}): switchesParser,
		reflect.TypeOf([]homeassistant.LightConfig{}):  lightsParser,
	}); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

func switchesParser(value string) (interface{}, error) {
	var items []homeassistant.SwitchConfig
	if err := json.Unmarshal([]byte(value), &items); err != nil {
		return nil, err
	}

	return items, nil
}

func lightsParser(value string) (interface{}, error) {
	var items []homeassistant.LightConfig
	if err := json.Unmarshal([]byte(value), &items); err != nil {
		return nil, err
	}

	return items, nil
}
