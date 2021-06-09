package homeassistant

import "gitlab.com/dmnbars/noolite/internal/noolite"

type MqttClient interface {
	Send(topic string, payload string) error
	AddHandler(topic string, handler func(payload string) error) error
}

type NooliteClient interface {
	Send(command noolite.Command)
	AddHandler(channel int, handler func(response noolite.Response))
}
