/*
Package pubsub contains utility functions for working with local broker mqtt.
*/
package pubsub

import (
	"log"
	"time"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)


//this type store the status of connection
type PubSub struct {
	Conn	MQTT.Client
	Err	chan error
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
  fmt.Printf("TOPIC: %s\n", msg.Topic())
  fmt.Printf("MSG: %s\n", msg.Payload())
}

//NewConnection return the PubSub object. nameClient is the client name in local broker.
func NewConnection(nameClient string) (*PubSub, error) {

	p := &PubSub{}

	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID(nameClient)
	opts.SetDefaultPublishHandler(f)
	p.Conn =  MQTT.NewClient(opts)
	token := p.Conn.Connect();
	ok := token.WaitTimeout(30 * time.Second)
	switch {
	case !ok:
		return nil, errors.New("Timeout Error at the beginning of the connection")
	case token.Error() != nil:
		return nil, token.Error()
	}

	p.Err = make(chan error)

	return p, nil
}

//Publish should be executed with a Go routine. The channel obtains the content that must be sent to the local broker.
func (p *PubSub) Publish(topic string, ch <-chan string) {

	if p.Conn == nil {
		p.Err <- errors.New("Nil Connection Error, execute NewConnection()")
		return
	}

        for msg:= range ch {
		if msg == "EOF" {
			log.Println("FINISH publish function")
			return
		}
		token := p.Conn.Publish(topic, 0, false, msg)
		if ok := token.WaitTimeout(10 * time.Second); !ok {
			p.Err <- errors.New("timeout Error in publish")
		}
		log.Printf("message: %s\n",msg)
	}
}

//Disconnect and close error channel
func (p *PubSub) Disconnect() {
	close(p.Err)
	p.Conn.Disconnect(250)
}
