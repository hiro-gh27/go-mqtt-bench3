package pubsub

import (
	"fmt"

	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var subscribeQos byte
var subscribeBaseTopic string

func subscribe(client MQTT.Client, id int) {
	//topic := fmt.Sprintf(subscribeBaseTopic+"%s"+"/"+"#", id)
	topic := fmt.Sprintf(subscribeBaseTopic+"%d", id)
	var handller MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("ClientID=%d, topic=%s, messageID=%d\n", id, msg.Topic(), msg.MessageID())
	}
	token := client.Subscribe(topic, subscribeQos, handller)
	if token.Wait() && token.Error() != nil {
		fmt.Printf("Subscribe Error: %s\n", token.Error())
	}
}

// Subscribe is
func Subscribe(opts SubscribeOptions) []SubscribeResult {
	var sResults []SubscribeResult
	//ch := make(chan SubscribeResult)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	subscribeQos = opts.Qos
	subscribeBaseTopic = opts.Topic
	for id := 0; id < opts.ClientNum; id++ {
		client := opts.Clients[id]
		subscribe(client, id)
	}
	wg.Wait()
	/*
		go func() {
			for {
				close(ch)
			}
		}()
		for {
			sResults = append(sResults, <-ch)
		}
	*/
	return sResults
}
