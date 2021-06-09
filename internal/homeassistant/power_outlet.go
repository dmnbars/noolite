package homeassistant

import (
	"encoding/json"
	"fmt"

	"gitlab.com/dmnbars/noolite/internal/noolite"

	"go.uber.org/zap"
)

type PowerOutlet struct {
	config        PowerOutletConfig
	logger        *zap.SugaredLogger
	mqttClient    MqttClient
	nooliteClient NooliteClient
}

type PowerOutletConfig struct {
	Channel int    `json:"channel"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Prefix  string `json:"prefix"`
}

func (p PowerOutletConfig) prefix() string {
	if p.Prefix == "" {
		return defaultPrefix
	}

	return p.Prefix
}

func (p PowerOutletConfig) commandTopic() string {
	return fmt.Sprintf("%s/switch/%s/command", p.prefix(), p.Id)
}

func (p PowerOutletConfig) configTopic() string {
	return fmt.Sprintf("%s/switch/%s/config", p.prefix(), p.Id)
}

func (p PowerOutletConfig) stateTopic() string {
	return fmt.Sprintf("%s/switch/%s/state", p.prefix(), p.Id)
}

func NewPowerOutlet(
	config PowerOutletConfig,
	log *zap.SugaredLogger,
	mqttClient MqttClient,
	nooliteClient NooliteClient,
) (PowerOutlet, error) {
	p := PowerOutlet{config: config, logger: log, mqttClient: mqttClient, nooliteClient: nooliteClient}

	p.nooliteClient.AddHandler(p.config.Channel, p.responseHandler)

	if err := p.sendConfig(); err != nil {
		return PowerOutlet{}, err
	}

	if err := mqttClient.AddHandler(config.commandTopic(), p.commandHandler); err != nil {
		return PowerOutlet{}, err
	}

	return p, nil
}

func (p *PowerOutlet) responseHandler(response noolite.Response) {
	if response.GetCommand() == noolite.CmdSendState {
		if response.GetD2() == 0 {
			p.setStateOff()
			return
		}
		if response.GetD2() == 1 {
			p.setStateOn()
			return
		}
		return
	}

	if response.GetCommand() == noolite.CmdSwitch {
		p.nooliteClient.Send(noolite.NewCommand(
			noolite.ModeFTX,
			noolite.CommandCtrSend,
			p.config.Channel,
			noolite.CmdReadState,
		))
	}
}

func (p *PowerOutlet) commandHandler(payload string) error {
	switch payload {
	case payloadOn:
		p.turnOn()
		return nil

	case payloadOff:
		p.turnOff()
		return nil
	}

	return fmt.Errorf("unknown payload: %s", payload)
}

func (p *PowerOutlet) turnOn() {
	p.nooliteClient.Send(noolite.NewCommand(
		noolite.ModeFTX,
		noolite.CommandCtrSend,
		p.config.Channel,
		noolite.CmdOn,
	))
}

func (p *PowerOutlet) turnOff() {
	p.nooliteClient.Send(noolite.NewCommand(
		noolite.ModeFTX,
		noolite.CommandCtrSend,
		p.config.Channel,
		noolite.CmdOff,
	))
}

func (p *PowerOutlet) setStateOn() {
	p.setState(payloadOn)
}

func (p *PowerOutlet) setStateOff() {
	p.setState(payloadOff)
}

func (p *PowerOutlet) setState(payload string) {
	if err := p.mqttClient.Send(p.config.stateTopic(), payload); err != nil {
		p.logger.Errorw(
			"can't change state",
			"id", p.config.Id,
			"state", payload,
		)
	}
}

func (p *PowerOutlet) sendConfig() error {
	config := map[string]string{
		"name":          p.config.Name,
		"unique_id":     p.config.Id,
		"command_topic": p.config.commandTopic(),
		"state_topic":   p.config.stateTopic(),
	}

	payload, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return p.mqttClient.Send(p.config.configTopic(), string(payload))
}
