package config

import (
	"encoding/json"
	"fmt"
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

	Switches      []homeassistant.SwitchConfig       `env:"SWITCHES" envDefault:"[]"`
	Lights        []homeassistant.LightConfig        `env:"LIGHTS" envDefault:"[]"`
	BinarySensors []homeassistant.BinarySensorConfig `env:"BINARY_SENSORS" envDefault:"[]"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.ParseWithFuncs(&cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf([]homeassistant.SwitchConfig{}):       switchesParser,
		reflect.TypeOf([]homeassistant.LightConfig{}):        lightsParser,
		reflect.TypeOf([]homeassistant.BinarySensorConfig{}): binarySensorsParser,
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

	// TODO: make normal validation
	for _, item := range items {
		if !item.IsValid() {
			return nil, fmt.Errorf("switch config (#%s) isn't valid", item.ID)
		}
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

func binarySensorsParser(value string) (interface{}, error) {
	var items []homeassistant.BinarySensorConfig
	if err := json.Unmarshal([]byte(value), &items); err != nil {
		return nil, err
	}

	return items, nil
}
