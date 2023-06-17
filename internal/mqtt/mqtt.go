package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageFromBroker struct {
	Topic   string
	Payload []byte
}

type MQTTService struct {
	queue     chan MessageFromBroker
	clients   map[string]*mqtt.Client
	closeCh   chan struct{}
	wg        *sync.WaitGroup
	closeOnce sync.Once

	mu sync.Mutex

	clientID string
}

func NewMQTTService() *MQTTService {

	ms := MQTTService{
		queue:   make(chan MessageFromBroker, 10),
		clients: make(map[string]*mqtt.Client),
		closeCh: make(chan struct{}),
		wg:      &sync.WaitGroup{},

		clientID: "",
	}

	return &ms
}

func (m *MQTTService) Put(msg MessageFromBroker) {
	select {
	case m.queue <- msg:
	case <-m.closeCh:
		log.Println("MQTTMessageQueue is closed. Unable to put a job.")
	}
}

func (m *MQTTService) Get() <-chan MessageFromBroker {
	return m.queue
}

func (m *MQTTService) Close() {
	m.closeOnce.Do(func() {
		close(m.closeCh)
		m.wg.Wait()
		close(m.queue)
	})

	for clientid, c := range m.clients {
		fmt.Printf("disconnecting client_id: %s\n", clientid)
		(*c).Disconnect(1500)
		fmt.Printf("disconnected client_id: %s\n", clientid)
	}
}

func (m *MQTTService) RegisterMqttClient(
	clientID string,
	brokerURL string,
	topicsToSub []string,
	deviceCert string,
	deviceCertKey string,
	rootCA string,
) error {

	tlsConfig, err := createTLSConfigWithCertKeyPair(deviceCert, deviceCertKey, rootCA)
	if err != nil {
		return err
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetTLSConfig(tlsConfig)
	opts.SetCleanSession(true)

	// set reconnect callback
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		fmt.Println("client connected, (re)subscribing to topics")

		wgSubsDone := sync.WaitGroup{}
		detailsSubsDone := make(chan string, len(topicsToSub))

		for _, topic := range topicsToSub {

			m.wg.Add(1)
			wgSubsDone.Add(1)
			go m.subscribe(c, topic, &wgSubsDone, detailsSubsDone)
		}

		wgSubsDone.Wait()
		close(detailsSubsDone)

		// Read all messages from the channel
		fmt.Printf("Done with subscribing. Result(s):\n")
		for detail := range detailsSubsDone {
			fmt.Printf("\t--> %s\n", detail)
		}

	})

	mqtt.ERROR = log.New(os.Stdout, "[mqtt_driver][ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[mqtt_driver][CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[mqtt_driver][WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[mqtt_driver][DEBUG] ", 0)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	m.mu.Lock()
	m.clients[clientID] = &client
	m.mu.Unlock()

	m.clientID = clientID

	return nil
}

func (m *MQTTService) subscribe(
	client mqtt.Client,
	topic string,
	wgSubDone *sync.WaitGroup,
	subDetailsChan chan<- string,
) {

	defer m.wg.Done()

	qos := byte(0)

	receiverFunc := func(c mqtt.Client, msg mqtt.Message) {
		message := MessageFromBroker{
			Topic:   msg.Topic(),
			Payload: msg.Payload(),
		}
		m.Put(message)
	}

	token := client.Subscribe(topic, qos, receiverFunc)
	token.Wait()

	if token.Error() != nil {
		errmsg := fmt.Sprintf("error on subscribing to topic %s: %v", topic, token.Error())
		subDetailsChan <- errmsg
		wgSubDone.Done()
		return
	}
	okmsg := fmt.Sprintf("succesfully subscribed to topic %s", topic)
	subDetailsChan <- okmsg
	wgSubDone.Done()

	<-m.closeCh
	fmt.Printf("unsubscribing from topic: \"%s\"\n", topic)
	client.Unsubscribe(topic)
	fmt.Printf("succesfully unsubscribed from topic: \"%s\"\n", topic)
}

func (m *MQTTService) UserInputPubHandler(topic string, message string) {

	go func() {

		m.mu.Lock()
		mqttclient := m.clients[m.clientID]
		m.mu.Unlock()

		err := m.publish(mqttclient, topic, message)
		if err != nil {
			fmt.Printf("errored on publishing message: \"%s\" to topic \"%s\". Reason: %v,",
				message, topic, err,
			)
		}

	}()
}

func (m *MQTTService) publish(
	client *mqtt.Client,
	topic string,
	payload any,
) error {

	token := (*client).Publish(topic, 0, false, payload)

	_ = token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	return nil

}

func createTLSConfigWithCertKeyPair(certFile, keyFile, caFile string) (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return tlsConfig, nil
}
