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
var pid string
var baseTopic string
var trial int
var count int

func initPubOpts(opts PublishOptions2) {
	baseTopic = opts.Topic
	count = opts.Count
	messageSize = opts.MessageSize
	pid = strconv.FormatInt(int64(os.Getpid()), 16)
	maxIntarval = maxInterval
	trial = opts.TrialNum
	qos = opts.Qos
}

// "sync publish"
func spub(id int, clinet MQTT.Client, trialNum int) PublishResult {
	var pResult PublishResult
	message := RandomMessage(messageSize)
	clientID := fmt.Sprintf("%s-%d", pid, id)
	topic := fmt.Sprintf(baseTopic+"%s"+"/"+"%d", clientID, trialNum)
	startTime := time.Now()
	token := clinet.Publish(topic, qos, false, message)
	token.Wait()
	endTime := time.Now()

	pResult.StartTime = startTime
	pResult.EndTime = endTime
	pResult.DurTime = endTime.Sub(startTime)
	pResult.Topic = topic
	pResult.ClientID = clientID
	fmt.Printf("### dtime=%s, clientID=%s, topic=%s ###\n",
		pResult.DurTime, pResult.ClientID, pResult.Topic)
	return pResult
}

// SyncPublish is
func SyncPublish(opts PublishOptions2) {
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

	message := RandomMessage(messageSize)
	clientID := fmt.Sprintf("%s-%d", pid, id)
	topic := fmt.Sprintf(baseTopic+"%s"+"/"+"%d", clientID, 0)
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
	pResults = append(pResults, vals)

	for index := 1; index < count; index++ {
		message = RandomMessage(messageSize)
		topic = fmt.Sprintf(baseTopic+"%s"+"/"+"%d", clientID, index)
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
		pResults = append(pResults, vals)
	}
	return pResults
}

// AsyncPublish is
func AsyncPublish(opts PublishOptions2) {
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
