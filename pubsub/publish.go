package pubsub

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var messageSize int
var maxIntarval int
var qos byte
var publishPid string
var baseTopic string
var trial int
var count int

func initPubOpts(opts PublishOptions) {
	baseTopic = opts.Topic
	count = opts.Count
	messageSize = opts.MessageSize
	publishPid = strconv.FormatInt(int64(os.Getpid()), 16)
	maxIntarval = maxInterval
	trial = opts.TrialNum
	qos = opts.Qos
}

// "sync publish"
func spub(id int, clinet MQTT.Client, trialNum int) PublishResult {
	var pResult PublishResult
	clientID := fmt.Sprintf("%s-%d", publishPid, id)
	//messageID := fmt.Sprintf("%s-%d-%d", pid, id, trialNum)
	//fmt.Printf("messageID=%s byte=%d", messageID, len(messageID))

	//messageID := time.Now().Format(time.StampNano)
	//fmt.Printf("messageID=%s, len=%d", messageID, len(messageID))
	//messageID, message := getMessageAndID(messageSize)
	//message = fmt.Sprintf("%s/%s", messageID, message)
	message := getMessage(messageSize)
	fmt.Printf("message=%s, size=%d", message, len(message))
	//topic := fmt.Sprintf(baseTopic+"%s"+"/"+"%d", clientID, trialNum)
	topic := fmt.Sprintf(baseTopic+"%d", id)

	startTime := time.Now()
	messageID := startTime.Format(time.StampNano)
	message = messageID + message
	//message = append(MessageID..., message...)
	//startTime.MarshalText()
	token := clinet.Publish(topic, qos, false, message)
	token.Wait()
	endTime := time.Now()

	fmt.Printf("message=%s, size=%d", message, len(message))
	pResult.StartTime = startTime
	pResult.EndTime = endTime
	pResult.DurTime = endTime.Sub(startTime)
	pResult.Topic = topic
	pResult.ClientID = clientID
	pResult.MessageID = messageID
	fmt.Printf("### dtime=%s, clientID=%s, topic=%s ###\n",
		pResult.DurTime, pResult.ClientID, pResult.Topic)
	return pResult
}

// SyncPublish is
func SyncPublish(opts PublishOptions) {
	initPubOpts(opts)
	var pResults []PublishResult
	for index := 0; index < opts.Count; index++ {
		for id := 0; id < len(opts.Clients); id++ {
			pr := spub(id, opts.Clients[id], index)
			pResults = append(pResults, pr)
		}
	}
}

// "async publish""
func aspub(id int, client MQTT.Client, freeze *sync.WaitGroup) []PublishResult {
	var pResults []PublishResult
	var waitTime time.Duration

	messageID, message := getMessageAndID(messageSize)
	clientID := fmt.Sprintf("%s-%d", publishPid, id)
	//topic := fmt.Sprintf(baseTopic+"%s"+"/"+"%d", clientID, 0)
	topic := fmt.Sprintf(baseTopic+"%d", id)
	if maxIntarval > 0 {
		waitTime = RandomInterval(maxIntarval)
	}
	freeze.Wait()
	if waitTime > 0 {
		time.Sleep(waitTime)
	}
	starttime := time.Now()
	token := client.Publish(topic, qos, false, message)
	token.Wait()
	endTime := time.Now()

	var vals PublishResult
	vals.StartTime = starttime
	vals.EndTime = endTime
	vals.DurTime = endTime.Sub(starttime)
	vals.Topic = topic
	vals.ClientID = clientID
	vals.MessageID = messageID
	pResults = append(pResults, vals)

	for index := 1; index < count; index++ {
		messageID, message = getMessageAndID(messageSize)
		//topic = fmt.Sprintf(baseTopic+"%s"+"/"+"%d", clientID, index)
		topic := fmt.Sprintf(baseTopic+"%d", id)
		if maxIntarval > 0 {
			waitTime = RandomInterval(maxIntarval)
			time.Sleep(waitTime)
		}
		starttime := time.Now()
		token := client.Publish(topic, qos, false, message)
		token.Wait()
		endTime := time.Now()

		vals = PublishResult{}
		vals.StartTime = starttime
		vals.EndTime = endTime
		vals.DurTime = endTime.Sub(starttime)
		vals.Topic = topic
		vals.ClientID = clientID
		vals.MessageID = messageID
		pResults = append(pResults, vals)
	}
	return pResults
}

// AsyncPublish is
func AsyncPublish(opts PublishOptions) {
	initPubOpts(opts)
	var pResults []PublishResult
	wg := &sync.WaitGroup{}
	syncStart := &sync.WaitGroup{}
	syncStart.Add(1)
	for id := 0; id < len(opts.Clients); id++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			re := aspub(id, opts.Clients[id], syncStart)
			pResults = append(pResults, re...)
		}(id)
	}
	time.Sleep(3 * time.Second)
	syncStart.Done()
	wg.Wait()
	for _, vals := range pResults {
		fmt.Printf("### dtime=%s, clientID=%s, topic=%s ###\n",
			vals.DurTime, vals.ClientID, vals.Topic)
	}
}

// LoadPublish is
func LoadPublish() {

}
