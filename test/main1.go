package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var numbers chan int64

func init() {
	numbers = make(chan int64)

}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	//fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())
	value, err := strconv.ParseInt(string(msg.Payload()), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	select {
	case numbers <- value:
	case <-time.After(1 * time.Second):
	}
}

func IsPrime(value int64) bool {
	for i := int64(2); i <= int64(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func IsPrimeSqrt(value int64) bool {
	for i := int64(2); i <= int64(math.Floor(math.Sqrt(float64(value)))); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func SieveOfEratosthenes(value int) {
	f := make([]bool, value)
	for i := 2; i <= int(math.Sqrt(float64(value))); i++ {
		if f[i] == false {
			for j := i * i; j < value; j += i {
				f[j] = true
			}
		}
	}
	for i := 2; i < value; i++ {
		if f[i] == false {
			fmt.Printf("%v ", i)
		}
	}
	fmt.Println("")
}

func main() {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID(fmt.Sprintf("go-test-%d", time.Now().UnixNano()))
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer c.Disconnect(50)

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe("TEST/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	/**
	for i := 1; i <= 100; i++ {
		if IsPrime(i) {
			fmt.Printf("%v ", i)
		}
	}
	fmt.Println("")
	/**/
	chs := make([]chan int64, lenChannels)

	selects := make([]reflect.SelectCase, lenChannels)

	for i := 0; i < lenChannels; i++ {
		chs[i] = make(chan int64)
		selects[i] = reflect.SelectCase{
			Dir:	reflect.SelectRecv
			reflect.
	}
	go func() {
		for v := range chs[0] {
			if IsPrimeSqrt(v) {
				fmt.Printf("%v ", v)
			}
		}
	}()
	go func() {
		for v := range chs[1] {
			if IsPrimeSqrt(v) {
				fmt.Printf("%v ", v)
			}
		}
	}()

	flagT := false
	var t1 time.Time
	for i := range numbers {
		//fmt.Printf("%v ", i)
		if !flagT {
			t1 = time.Now()
		}
		if i < int64(0) {
			break
		}
		select {
		case chs[0] <- i:
		case chs[1] <- i:
		}
	}
	close(chs[0])
	close(chs[1])
	fmt.Println("")
	/**
	SieveOfEratosthenes(100)
	/**/
}
