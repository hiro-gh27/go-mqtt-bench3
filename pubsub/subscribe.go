package pubsub

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var subscribeQos byte
var subscribeBaseTopic string

func subscribe(client MQTT.Client, id int, ch chan SubscribeResult) {
	//topic := fmt.Sprintf(subscribeBaseTopic+"%s"+"/"+"#", id)
	topic := fmt.Sprintf(subscribeBaseTopic+"%d", id)
	var handller MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		subscribeTime := time.Now()

		var sResult SubscribeResult
		sResult.SubscribeTime = subscribeTime
		sResult.ClientID = id
		sResult.Topic = msg.Topic()
		sResult.MessageID = string(msg.Payload()[:25])
		ch <- sResult
		//messageID := msg.Payload()[:25]
		//fmt.Printf("msg payload size= %d", len(msg.Payload()))
		//fmt.Printf("ClientID=%d, topic=%s, messageID=%d\n", id, msg.Topic(), msg.MessageID())
		//fmt.Printf("messageID=%s", string(messageID))
	}
	token := client.Subscribe(topic, subscribeQos, handller)
	if token.Wait() && token.Error() != nil {
		fmt.Printf("Subscribe Error: %s\n", token.Error())
	}
}

// Subscribe is
func Subscribe(opts SubscribeOptions) []SubscribeResult {
	var sResults []SubscribeResult
	sResultChan := make(chan SubscribeResult)
	//wg := &sync.WaitGroup{}
	//wg.Add(1)
	subscribeQos = opts.Qos
	subscribeBaseTopic = opts.Topic
	for id := 0; id < opts.ClientNum; id++ {
		client := opts.Clients[id]
		subscribe(client, id, sResultChan)
	}

	signalchan := make(chan os.Signal, 1)
	signal.Notify(signalchan, os.Interrupt)

	go func() {
		fmt.Println("exit: ctrl + c")
		for {

			//sResult :=
			/*
				fmt.Printf("subTime=%s, topic=%s, ClienID=%d, MessageID=%s",
					sResult.SubscribeTime, sResult.Topic, sResult.ClientID, sResult.MessageID)
			*/
			sResults = append(sResults, <-sResultChan)
			//fmt.Println("subscribe and add informastion")
		}
	}()

	<-signalchan
	fmt.Println("get signal!!")
	//runtime.Goexit()
	close(sResultChan)
	return sResults
}
