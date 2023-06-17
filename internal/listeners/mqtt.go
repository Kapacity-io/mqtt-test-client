package listeners

import (
	"fmt"
	"sync"

	"github.com/kapacity-io/mqtt-test-client/internal/mqtt"
)

type MQTTListener struct {
	mqtt *mqtt.MQTTService

	wg      sync.WaitGroup
	closeCh chan struct{}
}

func NewMQTTListener(mqtt *mqtt.MQTTService) *MQTTListener {

	return &MQTTListener{
		mqtt:    mqtt,
		closeCh: make(chan struct{}),
	}

}

func (cli *MQTTListener) Listen() {

	cli.wg.Add(1)
	go cli.startListeningMessages()

}

func (cli *MQTTListener) Stop() {

	close(cli.closeCh)
	cli.wg.Wait()

}

func (cli *MQTTListener) startListeningMessages() {

	defer cli.wg.Done()

	for {

		select {

		case msg, ok := <-cli.mqtt.Get():
			if !ok {
				fmt.Println("mqtt topic listener channel closed")
				return
			}

			fmt.Printf("\n--------------------------\n"+
				"received message from topic \"%s\":\n%s"+
				"\n------------------------------\n\n",
				msg.Topic, string(msg.Payload))

		case <-cli.closeCh:

			fmt.Println("stopping mqtt topic listener")
			return

		}
	}

}
