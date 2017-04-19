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
	subscribePidStr := strconv.FormatInt(int64(os.Getpid()), 16)
	subscribePid = fmt.Sprintf("%05s", subscribePidStr)
	subscribeTimeZone, _ = time.LoadLocation("Asia/Tokyo")
}

func subscribe(client MQTT.Client, id int, ch chan SubscribeResult) {
	sid := fmt.Sprintf("%05d", id)
	topic := fmt.Sprintf(subscribeBaseTopic+"%s", sid)
	var handller MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		subscribeTime := time.Now()
		payload := msg.Payload()

		var sResult SubscribeResult
		sResult.SubscribeTime = subscribeTime
		sResult.ClientID = fmt.Sprintf("%s-%s", subscribePid, sid)
		sResult.Topic = msg.Topic()
		sResult.PublisherID = string(payload[:11])
		sResult.MessageID = string(payload[12:47])
		sResult.PublishTime, _ = time.Parse(RFC3339NanoForMQTT, sResult.MessageID)
		//fmt.Printf("publisherID=%s, get time=%s\n", sResult.PublisherID, sResult.SubscribeTime)
		ch <- sResult
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
	for index := 0; index < opts.ClientNum; index++ {
		client := opts.Clients[index]
		id := index + opts.StartID
		subscribe(client, id, sResultChan)
	}

	go func() {
		fmt.Printf("subscribePid=%s ", subscribePid)
		for {
			sResults = append(sResults, <-sResultChan)
		}
	}()

	signalchan := make(chan os.Signal, 1)
	signal.Notify(signalchan, os.Interrupt)
	<-signalchan
	close(sResultChan)
	fmt.Println()
	return sResults
}
