package main

import (
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

var valueInput int64

func init() {
	flag.Int64Var(&valueInput, "valueInput", 0, "value input")
}

func main() {
	flag.Parse()
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID(fmt.Sprintf("go-test-prod-%d", time.Now().UnixNano()))

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer c.Disconnect(50)

	for i := int64(0); i < valueInput; i++ {
		if token := c.Publish("primes", 0, false, fmt.Sprintf("%d", i)); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}
	for i := int64(0); i < 10; i++ {
		if token := c.Publish("primes", 0, false, fmt.Sprintf("%d", -1)); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}

	/**
	  //Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	  //from the server after sending each message
	  for i := 0; i < 5; i++ {
	    text := fmt.Sprintf("this is msg #%d!", i)
	    token := c.Publish("go-mqtt/sample", 0, false, text)
	    token.Wait()
	  }

	  time.Sleep(3 * time.Second)

	  //unsubscribe from /go-mqtt/sample
	  if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
	    fmt.Println(token.Error())
	    os.Exit(1)
	  }

	  /**/

	time.Sleep(1 * time.Second)
}
