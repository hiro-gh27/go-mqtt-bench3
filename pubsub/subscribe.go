package pubsub

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var subscribeQos byte
var subscribeBaseTopic string
var subscribePid string
var subscribeTimeZone *time.Location

func initSubOpts(opts SubscribeOptions) {
	subscribeQos = opts.Qos
	subscribeBaseTopic = opts.Topic
	subscribePid = strconv.FormatInt(int64(os.Getpid()), 16)
	subscribeTimeZone, _ = time.LoadLocation("Asia/Tokyo")
}

func subscribe(client MQTT.Client, id int, ch chan SubscribeResult) {
	//topic := fmt.Sprintf(subscribeBaseTopic+"%s"+"/"+"#", id)
	testID := 99
	sid := fmt.Sprintf("%05d", testID)
	topic := fmt.Sprintf(subscribeBaseTopic+"%s", sid)
	sid = fmt.Sprintf("%05d", id)
	var handller MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		subscribeTime := time.Now()
		payload := msg.Payload()

		var sResult SubscribeResult
		sResult.SubscribeTime = subscribeTime
		sResult.ClientID = fmt.Sprintf("%s-%s", subscribePid, sid)
		sResult.Topic = msg.Topic()
		sResult.PublisherID = string(payload[:10])
		sResult.MessageID = string(payload[11:46])
		sResult.PublishTime, _ = time.Parse(RFC3339NanoForMQTT, sResult.MessageID)
		//sResult.PublishTime, _ = time.ParseInLocation(MsgStampLayout, sResult.MessageID, subscribeTimeZone)
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
	initSubOpts(opts)

	var sResults []SubscribeResult
	sResultChan := make(chan SubscribeResult)
	for id := 0; id < opts.ClientNum; id++ {
		client := opts.Clients[id]
		subscribe(client, id, sResultChan)
	}

	go func() {
		fmt.Printf("subscribePid=%s, exit->ctrl + c\n", subscribePid)
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

	signalchan := make(chan os.Signal, 1)
	signal.Notify(signalchan, os.Interrupt)
	<-signalchan
	fmt.Println("get signal!!")
	//runtime.Goexit()
	close(sResultChan)
	return sResults
}
