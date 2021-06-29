package noolite

import (
	"github.com/tarm/serial"
	"go.uber.org/zap"
)

type Client struct {
	logger           *zap.SugaredLogger
	port             *serial.Port
	responses        chan Response
	commands         chan Command
	responseHandlers map[int]func(response Response)
}

func NewClient(l *zap.SugaredLogger, name string) (*Client, error) {
	config := &serial.Config{Name: name, Baud: 9600}
	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	client := &Client{
		logger:           l,
		port:             port,
		responses:        make(chan Response, 100),
		responseHandlers: map[int]func(response Response){},
		commands:         make(chan Command, 100),
	}

	go client.handleResponses()
	go client.handleCommands()
	go client.read()

	return client, nil
}

func (c *Client) AddHandler(channel int, handler func(response Response)) {
	c.responseHandlers[channel] = handler
}

func (c *Client) Send(command Command) {
	c.commands <- command
}

func (c *Client) BindRX(channel int) {
	c.Send(NewCommand(
		ModeRx,
		CommandCtrBindOn,
		channel,
		0,
	))
}

func (c *Client) BindFTX(channel int) {
	c.Send(NewCommand(
		ModeFTX,
		CommandCtrSend,
		channel,
		CmdBind,
	))
}

func (c *Client) handleResponses() {
	for response := range c.responses {
		c.logger.Infow("response", "response", response.String())

		if !response.IsSuccess() {
			c.logger.Errorw("response isn't success", "response", response.String())
			continue
		}

		channel := response.GetChannel()
		handler, ok := c.responseHandlers[channel]
		if !ok {
			c.logger.Infow(
				"no handlers for channel",
				"channel", channel,
			)
			continue
		}

		handler(response)
	}
}

func (c *Client) handleCommands() {
	for command := range c.commands {
		c.logger.Infow("command", "command", command.String())

		n, err := c.port.Write(command)
		if err != nil {
			c.logger.Errorw("can't write to serial port", "error", err)
		}

		c.logger.Debugw("write bytes", "number", n)
	}
}

func (c *Client) read() {
	c.logger.Debug("read bytes start")
	defer c.logger.Debug("read bytes stop")

	resp := make([]byte, 0)
	for {
		buf := make([]byte, 1)
		_, err := c.port.Read(buf)
		if err != nil {
			c.logger.Error(err)
			continue
		}

		if buf[0] == RespSt {
			resp = make([]byte, 0)
		}

		resp = append(resp, buf[0])

		if buf[0] == RespSp {
			response, err := NewResponse(resp)
			if err != nil {
				c.logger.Errorw(
					"can't create response",
					"error", err,
					"resp", string(resp),
				)
				continue
			}

			c.responses <- response
		}
	}
}

func calcCrc(data []byte) byte {
	var sum byte = 0
	for _, b := range data[:15] {
		sum += b
	}

	return sum & 0xFF
}
