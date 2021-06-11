package homeassistant

import (
	"encoding/json"
	"fmt"

	"gitlab.com/dmnbars/noolite/internal/noolite"

	"go.uber.org/zap"
)

type Switch struct {
	config        SwitchConfig
	logger        *zap.SugaredLogger
	mqttClient    MqttClient
	nooliteClient NooliteClient
}

type SwitchConfig struct {
	Channel int    `json:"channel"`
	ID      string `json:"id"`
	Name    string `json:"name"`
}

func (p SwitchConfig) commandTopic() string {
	return fmt.Sprintf("%s/switch/%s/command", prefix, p.ID)
}

func (p SwitchConfig) configTopic() string {
	return fmt.Sprintf("%s/switch/%s/config", prefix, p.ID)
}

func (p SwitchConfig) stateTopic() string {
	return fmt.Sprintf("%s/switch/%s/state", prefix, p.ID)
}

func NewSwitch(
	config SwitchConfig,
	log *zap.SugaredLogger,
	mqttClient MqttClient,
	nooliteClient NooliteClient,
) (Switch, error) {
	sw := Switch{config: config, logger: log, mqttClient: mqttClient, nooliteClient: nooliteClient}

	sw.nooliteClient.AddHandler(sw.config.Channel, sw.responseHandler)

	if err := sw.sendConfig(); err != nil {
		return Switch{}, err
	}

	if err := mqttClient.AddHandler(config.commandTopic(), sw.commandHandler); err != nil {
		return Switch{}, err
	}

	return sw, nil
}

func (s *Switch) responseHandler(response noolite.Response) {
	if response.GetCommand() == noolite.CmdSendState {
		if response.GetD2() == 0 {
			s.setStateOff()
			return
		}
		if response.GetD2() == 1 {
			s.setStateOn()
			return
		}
		return
	}

	if response.GetCommand() == noolite.CmdSwitch {
		s.nooliteClient.Send(noolite.NewCommand(
			noolite.ModeFTX,
			noolite.CommandCtrSend,
			s.config.Channel,
			noolite.CmdReadState,
		))
	}
}

func (s *Switch) commandHandler(payload string) error {
	switch payload {
	case payloadOn:
		s.turnOn()
		return nil

	case payloadOff:
		s.turnOff()
		return nil
	}

	return fmt.Errorf("unknown payload: %s", payload)
}

func (s *Switch) turnOn() {
	s.nooliteClient.Send(noolite.NewCommand(
		noolite.ModeFTX,
		noolite.CommandCtrSend,
		s.config.Channel,
		noolite.CmdOn,
	))
}

func (s *Switch) turnOff() {
	s.nooliteClient.Send(noolite.NewCommand(
		noolite.ModeFTX,
		noolite.CommandCtrSend,
		s.config.Channel,
		noolite.CmdOff,
	))
}

func (s *Switch) setStateOn() {
	s.setState(payloadOn)
}

func (s *Switch) setStateOff() {
	s.setState(payloadOff)
}

func (s *Switch) setState(payload string) {
	if err := s.mqttClient.Send(s.config.stateTopic(), payload); err != nil {
		s.logger.Errorw(
			"can't change state",
			"id", s.config.ID,
			"state", payload,
		)
	}
}

func (s *Switch) sendConfig() error {
	config := map[string]string{
		"name":          s.config.Name,
		"unique_id":     s.config.ID,
		"command_topic": s.config.commandTopic(),
		"state_topic":   s.config.stateTopic(),
	}

	payload, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return s.mqttClient.Send(s.config.configTopic(), string(payload))
}
