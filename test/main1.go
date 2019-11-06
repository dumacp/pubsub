package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var numbers chan *big.Int
var lenChannels int

func init() {
	numbers = make(chan *big.Int)
	flag.IntVar(&lenChannels, "lenChannels", 1, "channel count to implement")
	runtime.GOMAXPROCS(8)

}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	//fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())

	value, ok := new(big.Int).SetString(string(msg.Payload()), 10)
	fmt.Printf("input:: %v\n", value)

	if !ok {
		log.Println("Error in parse input")
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

func IsPrimeSqrt(value *big.Int) bool {
	fmt.Printf("input IsPrimeSqrt: %v\n", value)
	valueF := new(big.Float).SetInt(value)
	i := new(big.Float).SetInt64(2)
	inc := new(big.Float).SetInt64(1)
	zs := new(big.Float)
	for ; i.Cmp(zs.Sqrt(valueF)) <= 0; i.Add(i, inc) {
		//fmt.Printf("Z Sqrt: %v\n", zs)
		ii, _ := i.Int(new(big.Int))
		m := new(big.Int)
		z := new(big.Int)
		if _, m = z.DivMod(value, ii, m); m.Cmp(new(big.Int).SetInt64(0)) == 0 {
			return false
		}

	}
	return value.Cmp(new(big.Int).SetInt64(1)) > 0
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
	flag.Parse()
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
	if token := c.Subscribe("$share/TEST/primes", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	if token := c.Subscribe("control", 0, nil); token.Wait() && token.Error() != nil {
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
	chs := make([]chan *big.Int, lenChannels)
	selects := make([]reflect.SelectCase, lenChannels)
	results := make([]*big.Int, 0)

	for i := 0; i < lenChannels; i++ {
		chs[i] = make(chan *big.Int)
		go func(j int) {
			for v := range chs[j] {
				if IsPrimeSqrt(v) {
					results = append(results, v)
				}
			}
		}(i)
	}

	flagT := false
	var t1 time.Time
	for v := range numbers {
		if !flagT {
			t1 = time.Now()
			flagT = true
		}
		if v.Cmp(new(big.Int)) < 0 {
			break
		}

		for i := 0; i < lenChannels; i++ {
			selects[i] = reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(chs[i]),
				Send: reflect.ValueOf(v),
			}
		}
		reflect.Select(selects)

	}
	fmt.Printf("%v ", len(results))
	fmt.Printf("\n\n%s\n", time.Now().Sub(t1))
	for i := 0; i < lenChannels; i++ {
		close(chs[i])
	}
	/**
	SieveOfEratosthenes(100)
	/**/
}
