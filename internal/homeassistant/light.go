package homeassistant

import (
	"encoding/json"
	"fmt"

	"gitlab.com/dmnbars/noolite/internal/noolite"

	"go.uber.org/zap"
)

type Light struct {
	config        LightConfig
	logger        *zap.SugaredLogger
	mqttClient    MqttClient
	nooliteClient NooliteClient
}

type LightConfig struct {
	Channel int    `json:"channel"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Prefix  string `json:"prefix"`
}

func (l LightConfig) commandTopic() string {
	return fmt.Sprintf("%s/light/%s/command", prefix, l.ID)
}

func (l LightConfig) configTopic() string {
	return fmt.Sprintf("%s/light/%s/config", prefix, l.ID)
}

func (l LightConfig) stateTopic() string {
	return fmt.Sprintf("%s/light/%s/state", prefix, l.ID)
}

func NewLight(
	config LightConfig,
	log *zap.SugaredLogger,
	mqttClient MqttClient,
	nooliteClient NooliteClient,
) (Light, error) {
	light := Light{config: config, logger: log, mqttClient: mqttClient, nooliteClient: nooliteClient}

	light.nooliteClient.AddHandler(light.config.Channel, light.responseHandler)

	if err := light.sendConfig(); err != nil {
		return Light{}, err
	}

	if err := mqttClient.AddHandler(config.commandTopic(), light.commandHandler); err != nil {
		return Light{}, err
	}

	return light, nil
}

func (l *Light) responseHandler(response noolite.Response) {
	command := response.GetCommand()
	if command == noolite.CmdSendState {
		if response.GetD2() == 0 {
			l.setStateOff()
			return
		}
		if response.GetD2() == 1 {
			l.setStateOn()
			return
		}
		return
	}

	if command == noolite.CmdSwitch || command == noolite.CmdOn || command == noolite.CmdOff {
		l.nooliteClient.Send(noolite.NewCommand(
			noolite.ModeFTX,
			noolite.CommandCtrSend,
			l.config.Channel,
			noolite.CmdReadState,
		))
	}
}

type lightCommandPayload struct {
	State string `json:"state"`
}

func (l *Light) commandHandler(rawPayload string) error {
	var payload lightCommandPayload
	if err := json.Unmarshal([]byte(rawPayload), &payload); err != nil {
		return fmt.Errorf("can't unmarshal raw payload: %s", rawPayload)
	}

	switch payload.State {
	case payloadOn:
		l.turnOn()
		return nil

	case payloadOff:
		l.turnOff()
		return nil
	}

	return fmt.Errorf("unknown payload: %s", payload)
}

func (l *Light) turnOn() {
	l.nooliteClient.Send(noolite.NewCommand(
		noolite.ModeFTX,
		noolite.CommandCtrSend,
		l.config.Channel,
		noolite.CmdOn,
	))
}

func (l *Light) turnOff() {
	l.nooliteClient.Send(noolite.NewCommand(
		noolite.ModeFTX,
		noolite.CommandCtrSend,
		l.config.Channel,
		noolite.CmdOff,
	))
}

func (l *Light) setStateOn() {
	l.setState(payloadOn)
}

func (l *Light) setStateOff() {
	l.setState(payloadOff)
}

func (l *Light) setState(state string) {
	payload, err := json.Marshal(lightCommandPayload{
		State: state,
	})
	if err != nil {
		l.logger.Errorw(
			"can't marshal payload state",
			"id", l.config.ID,
			"state", state,
			"error", err,
		)
	}

	if err := l.mqttClient.Send(l.config.stateTopic(), string(payload)); err != nil {
		l.logger.Errorw(
			"can't change state",
			"id", l.config.ID,
			"state", payload,
			"error", err,
		)
	}
}

func (l *Light) sendConfig() error {
	config := map[string]string{
		"name":          l.config.Name,
		"unique_id":     l.config.ID,
		"command_topic": l.config.commandTopic(),
		"state_topic":   l.config.stateTopic(),
		"schema":        "json",
	}

	payload, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return l.mqttClient.Send(l.config.configTopic(), string(payload))
}
