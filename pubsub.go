/*
Package pubsub contains utility functions for working with local broker mqtt.
*/
package pubsub

import (
	"fmt"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Timestamp float64     `json:"timestamp"`
}

//this type store the status of connection
type PubSub struct {
	Conn MQTT.Client
	Err  chan error
}

var subscriptions map[string]chan []byte

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	// fmt.Printf("TOPIC: %s\n", msg.Topic())
	// fmt.Printf("MSG: %s\n", msg.Payload())

	vch, ok := subscriptions[string(msg.Topic())]
	if ok {
		select {
		case vch <- msg.Payload():
		default:
		}
	}
}

var onConnection MQTT.OnConnectHandler = func(c MQTT.Client) {
	log.Println("OnConnection MQTT")
	for k := range subscriptions {
		t := c.Subscribe(k, 0, nil)
		if t.WaitTimeout(3*time.Second) && t.Error() != nil {
			log.Printf("error subzcription: %s", t.Error())
		}
	}
}

//NewConnection return the PubSub object. nameClient is the client name in local broker.
func NewConnection(nameClient string) (*PubSub, error) {

	p := &PubSub{}

	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID(nameClient)
	opts.SetDefaultPublishHandler(f)
	opts.SetOnConnectHandler(onConnection)
	opts.SetAutoReconnect(true)
	p.Conn = MQTT.NewClient(opts)
	token := p.Conn.Connect()
	ok := token.WaitTimeout(30 * time.Second)
	switch {
	case !ok:
		return nil, fmt.Errorf("Timeout Error at the beginning of the connection")
	case token.Error() != nil:
		return nil, token.Error()
	}

	p.Err = make(chan error)

	return p, nil
}

//New return the PubSub object. Without start connection, nameClient is the client name in local broker.
func New(nameClient string) *PubSub {

	p := &PubSub{}
	subscriptions = make(map[string]chan []byte)

	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID(nameClient)
	opts.SetDefaultPublishHandler(f)
	opts.SetOnConnectHandler(onConnection)
	opts.SetAutoReconnect(true)
	p.Conn = MQTT.NewClient(opts)
	return p
}

//Start connection
func (p *PubSub) Start() error {
	token := p.Conn.Connect()
	ok := token.WaitTimeout(30 * time.Second)
	switch {
	case !ok:
		return fmt.Errorf("Timeout Error at the beginning of the connection")
	case token.Error() != nil:
		return token.Error()
	}

	p.Err = make(chan error)

	return nil
}

//AddSubscription add new subcription before Connection
func (p *PubSub) AddSubscription(topic string, ch chan []byte) {
	subscriptions[topic] = ch
}

//Publish should be executed with a Go routine. The channel obtains the content that must be sent to the local broker.
func (p *PubSub) Publish(topic string, ch <-chan string) {

	if p.Conn == nil {
		p.Err <- fmt.Errorf("Nil Connection Error, execute NewConnection()")
		return
	}

	for msg := range ch {
		if msg == "EOF" {
			log.Println("FINISH publish function")
			return
		}
		token := p.Conn.Publish(topic, 0, false, msg)
		if ok := token.WaitTimeout(10 * time.Second); !ok {
			p.Err <- fmt.Errorf("timeout Error in publish")
		}
		// log.Printf("TOPIC: %s; message: %s\n", topic, msg)
	}
}

//Disconnect and close error channel
func (p *PubSub) Disconnect() {
	close(p.Err)
	p.Conn.Disconnect(250)
}
