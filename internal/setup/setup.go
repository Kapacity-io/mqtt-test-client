package setup

import (
	"errors"
	"flag"
	"strings"
)

// Configs holds the configuration values.
type Configs struct {
	BrokerAddr  string
	DeviceCert  string
	DeviceKey   string
	RootCA      string
	ClientID    string
	TopicsToSub []string
}

// String returns a string representation of the configs.
// would be used in debug printing if we are unsure
// that the Configs struct would be properly populated.
func (c Configs) String() string {
	var sb strings.Builder
	sb.WriteString("mqtt-broker-addr: ")
	sb.WriteString(c.BrokerAddr)
	sb.WriteString("\nmqtt-broker-device-cert: ")
	sb.WriteString(c.DeviceCert)
	sb.WriteString("\nmqtt-broker-device-key: ")
	sb.WriteString(c.DeviceKey)
	sb.WriteString("\nmqtt-broker-root-cert: ")
	sb.WriteString(c.RootCA)
	sb.WriteString("\nmqtt-broker-client-id: ")
	sb.WriteString(c.ClientID)
	sb.WriteString("\nmqtt-topics-to-sub: ")
	sb.WriteString("\n")
	for _, t := range c.TopicsToSub {
		sb.WriteString("\t")
		sb.WriteString(t)
		sb.WriteString("\n")
	}

	return sb.String()
}

// Validate checks the validity of the configs.
func (c Configs) Validate() error {
	if c.BrokerAddr == "" {
		return errors.New("mqtt-broker-addr cannot be empty")
	}
	if c.DeviceCert == "" {
		return errors.New("mqtt-broker-device-cert cannot be empty")
	}
	if c.DeviceKey == "" {
		return errors.New("mqtt-broker-device-key cannot be empty")
	}
	return nil
}

func ParseConfigsFromFlags() Configs {
	// Define the flags.
	mqttBrokerAddr := flag.String("mqtt-broker-addr", "", "URL of the MQTT broker, for example: \"tcps://foo-ats.iot.eu-north-1.amazonaws.com:8883/mqtt\"")
	mqttRootCa := flag.String("mqtt-broker-root-ca", "", "root-CA, for example \"AmazonRootCA1.pem\"")
	mqttBrokerDeviceCert := flag.String("mqtt-broker-device-cert", "", "Filepath to the cert file")
	mqttBrokerDeviceKey := flag.String("mqtt-broker-device-key", "", "Filepath to the key for the cert file")
	mqttBrokerClientID := flag.String("mqtt-broker-client-id", "mqtt-test-client", "client ID used for connecting to the broker")
	mqttTopicsToSubscribe := flag.String("mqtt-topics-to-sub", "", "Topic to subscribe, if multiple add topics separated by a comma")

	flag.Parse()

	topics := strings.Split(*mqttTopicsToSubscribe, ",")

	// Populate the config struct.
	config := Configs{
		BrokerAddr:  *mqttBrokerAddr,
		DeviceCert:  *mqttBrokerDeviceCert,
		DeviceKey:   *mqttBrokerDeviceKey,
		RootCA:      *mqttRootCa,
		ClientID:    *mqttBrokerClientID,
		TopicsToSub: topics,
	}

	return config
}
