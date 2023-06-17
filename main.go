package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kapacity-io/mqtt-test-client/internal/listeners"
	"github.com/kapacity-io/mqtt-test-client/internal/mqtt"
	"github.com/kapacity-io/mqtt-test-client/internal/setup"
)

var (
	// These variables are set during build time with -ldflags
	Version string
	Build   string
)

func main() {

	fmt.Printf("\n--- PROGRAM STARTS ---\n"+
		"\tVersion: %s\n"+
		"\tBuild: %s\n",
		Version, Build)

	cfg := setup.ParseConfigsFromFlags()
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Invalid configs: %v", err)
		os.Exit(1)
	}
	fmt.Printf("options:\n%s\n\n", cfg.String())

	mqtt := mqtt.NewMQTTService()

	err := mqtt.RegisterMqttClient(
		cfg.ClientID,
		cfg.BrokerAddr,
		cfg.TopicsToSub,
		cfg.DeviceCert,
		cfg.DeviceKey,
		cfg.RootCA,
	)

	if err != nil {
		fmt.Printf("Failed to initialize client to MQTT broker: %v\n", err)
		os.Exit(1)
	}

	mqttlistener := listeners.NewMQTTListener(mqtt)
	mqttlistener.Listen()

	tcplistener, err := listeners.NewTCPListener(
		"localhost:9596",
		func(topic string, message string) {

			mqtt.UserInputPubHandler(topic, message)

		},
	)
	if err != nil {
		fmt.Println("Error setting up TCP server:", err)
		os.Exit(1)
	}

	tcplistener.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-c
	fmt.Printf("stop signal caught -> stopping services\n")

	tcplistener.Stop()
	mqttlistener.Stop()
	mqtt.Close()

	fmt.Printf("\n\n\t--- PROGRAM EXISTS ---\n\n")
	os.Exit(0)
}
