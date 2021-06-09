package mqtt

import (
	"crypto/tls"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.SugaredLogger
	client mqtt.Client
}

func NewClient(
	logger *zap.SugaredLogger,
	host string,
	port string,
	clientID string,
	user string,
	password string,
) (*Client, error) {
	server := fmt.Sprintf("tcp://%s:%s", host, port)

	connOpts := mqtt.NewClientOptions().
		AddBroker(server).
		SetUsername(user).
		SetPassword(password).
		SetClientID(clientID).
		SetCleanSession(true)

	connOpts.SetTLSConfig(
		&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}, // nolint
	)

	client := mqtt.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	logger.Debugw("connected to mqtt", "server", server)

	return &Client{logger: logger, client: client}, nil
}

func (c *Client) Send(topic string, payload string) error {
	if token := c.client.Publish(
		topic,
		0,
		true,
		payload,
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	c.logger.Debugw("message has been send", "topic", topic, "payload", payload)

	return nil
}

func (c *Client) AddHandler(topic string, handler func(payload string) error) error {
	if token := c.client.Subscribe(
		topic,
		0,
		func(client mqtt.Client, message mqtt.Message) {
			payload := string(message.Payload())

			c.logger.Debugw(
				"handle message",
				"fromTopic", message.Topic(),
				"payload", payload,
				"topic", topic,
			)

			if message.Topic() != topic {
				return
			}

			if err := handler(payload); err != nil {
				c.logger.Errorw(
					"can't handle payload",
					"payload", payload,
					"topic", topic,
					"error", err,
				)
			}
		},
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
