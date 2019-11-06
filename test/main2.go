package main

import (
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"math/big"
	"os"
	"time"
)

var valueInput0 *big.Int
var valueInput1 *big.Int
var input0 string
var input1 string

func init() {
	flag.StringVar(&input0, "startValue", "0", "value input")
	flag.StringVar(&input1, "endValue", "0", "value input")
	valueInput0 = new(big.Int)
	valueInput1 = new(big.Int)
}

func main() {
	flag.Parse()

	if _, ok := valueInput0.SetString(input0, 10); !ok {
		os.Exit(-1)
	}
	if _, ok := valueInput1.SetString(input1, 10); !ok {
		os.Exit(-1)
	}

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

	inc := new(big.Int).SetInt64(1)
	for ; valueInput0.Cmp(valueInput1) < 0; valueInput0.Add(valueInput0, inc) {
		fmt.Printf("%v\n", valueInput0)
		if token := c.Publish("primes", 0, false, fmt.Sprintf("%d", valueInput0)); token.Wait() && token.Error() != nil {
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
