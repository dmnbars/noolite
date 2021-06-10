package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/dmnbars/noolite/internal/noolite"

	"gitlab.com/dmnbars/noolite/internal/config"
	"gitlab.com/dmnbars/noolite/internal/homeassistant"
	"gitlab.com/dmnbars/noolite/internal/mqtt"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logger, err := newLogger(cfg.Level)
	if err != nil {
		panic(err)
	}
	logger.Debugw("config", "config", cfg)

	nooliteClient, err := noolite.NewClient(logger.Named("noolite"), cfg.SerialPort)
	if err != nil {
		panic(err)
	}

	mqttClient, err := mqtt.NewClient(
		logger.Named("mqtt"),
		cfg.MqttHost,
		cfg.MqttPort,
		cfg.MqttClientID,
		cfg.MqttUsername,
		cfg.MqttPassword,
	)
	if err != nil {
		panic(err)
	}

	for _, powerOutletConfig := range cfg.PowerOutlets {
		_, err := homeassistant.NewPowerOutlet(
			powerOutletConfig,
			logger.Named(powerOutletConfig.ID),
			mqttClient,
			nooliteClient,
		)
		if err != nil {
			panic(err)
		}
	}

	for _, lightConfig := range cfg.Lights {
		_, err := homeassistant.NewLight(
			lightConfig,
			logger.Named(lightConfig.ID),
			mqttClient,
			nooliteClient,
		)
		if err != nil {
			panic(err)
		}
	}

	sig := <-waitExitSignal()
	logger.Infow("stopping by signal", "signal", sig.String())
}

func newLogger(level string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	atom := zap.NewAtomicLevel()
	err := atom.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	cfg.Level = atom

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func waitExitSignal() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return sigs
}
