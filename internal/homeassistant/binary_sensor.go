package homeassistant

import (
	"encoding/json"
	"fmt"

	"gitlab.com/dmnbars/noolite/internal/noolite"

	"go.uber.org/zap"
)

type BinarySensor struct {
	config        BinarySensorConfig
	logger        *zap.SugaredLogger
	mqttClient    MqttClient
	nooliteClient NooliteClient
}

type BinarySensorConfig struct {
	Channel int    `json:"channel"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
}

func (p BinarySensorConfig) configTopic() string {
	return fmt.Sprintf("%s/binary_sensor/%s/config", prefix, p.ID)
}

func (p BinarySensorConfig) stateTopic() string {
	return fmt.Sprintf("%s/binary_sensor/%s/state", prefix, p.ID)
}

func NewBinarySensor(
	config BinarySensorConfig,
	log *zap.SugaredLogger,
	mqttClient MqttClient,
	nooliteClient NooliteClient,
) (BinarySensor, error) {
	sw := BinarySensor{config: config, logger: log, mqttClient: mqttClient, nooliteClient: nooliteClient}

	sw.nooliteClient.AddHandler(sw.config.Channel, sw.responseHandler)

	if err := sw.sendConfig(); err != nil {
		return BinarySensor{}, err
	}

	return sw, nil
}

func (s *BinarySensor) responseHandler(response noolite.Response) {
	switch response.GetCommand() {
	case noolite.CmdOn:
		s.setStateOn()
		return

	case noolite.CmdOff:
		s.setStateOff()
		return
	}

	s.logger.Errorw(
		"unknown command",
		"id", s.config.ID,
		"cmd", response.GetCommand(),
		"response", response.String(),
	)
}

func (s *BinarySensor) setStateOn() {
	s.setState(payloadOn)
}

func (s *BinarySensor) setStateOff() {
	s.setState(payloadOff)
}

func (s *BinarySensor) setState(payload string) {
	if err := s.mqttClient.Send(s.config.stateTopic(), payload); err != nil {
		s.logger.Errorw(
			"can't change state",
			"id", s.config.ID,
			"state", payload,
		)
	}
}

func (s *BinarySensor) sendConfig() error {
	config := map[string]string{
		"name":        s.config.Name,
		"unique_id":   s.config.ID,
		"state_topic": s.config.stateTopic(),
		"icon":        s.config.Icon,
	}

	payload, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return s.mqttClient.Send(s.config.configTopic(), string(payload))
}
