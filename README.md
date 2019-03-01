# pubsub


## Example:

```
package main

import (
	"log"
	"time"
	"fmt"
        "github.com/dumacp/pubsub"
)


func main() {
	pub, err := pubsub.NewConnection("go-test")
        if err != nil {
        	log.Fatal(err)
	}
        defer pub.Disconnect()
        msgChan := make(chan string)
        go pub.Publish("EVENTS/test", msgChan)

	t1 := time.NewTicker(10 * time.Second):
	defer t1.Stop()
	count := 0
	for {
		select {
		case <-t1.C:
			timeStamp := float64(time.Now().UnixNano())/1000000000
			count = count + 1
			msg := fmt.Sprintf("{\"timeStamp\": %f, \"value\": %v, \"type\": \"TEST\"}",timeStamp, count)
			msgChan <- msg
		case err := <-pub.Err:
			log.Printf("ERROR: %v\n", err)
		}
	}
}
```
