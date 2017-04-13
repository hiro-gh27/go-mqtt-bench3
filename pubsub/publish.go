package pubsub

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"sort"

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
	sid := fmt.Sprintf("%05d", id)
	clientID := fmt.Sprintf("%s-%s", publishPid, sid)
	message := getMessage(messageSize - len(clientID) - 35 - 2) //30 => nanoTimeStamp, 2=> / /
	topic := fmt.Sprintf(baseTopic+"%s", sid)

	startTime := time.Now()
	messageID := startTime.Format(RFC3339NanoForMQTT)
	message = clientID + "/" + messageID + "/" + message
	token := clinet.Publish(topic, qos, false, message)
	waitStartTime := time.Now()
	token.Wait()
	endTime := time.Now()
	//time.Sleep(100 * time.Millisecond)

	pResult.StartTime = startTime
	pResult.WaitStartTime = waitStartTime
	pResult.EndTime = endTime
	pResult.LeadDuration = waitStartTime.Sub(startTime)
	pResult.WaitDuration = endTime.Sub(waitStartTime)
	pResult.TotalDuration = endTime.Sub(startTime)
	pResult.DurTime = endTime.Sub(startTime)
	pResult.Topic = topic
	pResult.ClientID = clientID
	pResult.MessageID = messageID

	/*
		fmt.Printf("### dtime=%s, clientID=%s, topic=%s ###\n",
			pResult.DurTime, pResult.ClientID, pResult.Topic)
	*/
	return pResult
}

// SyncPublish is
func SyncPublish(opts PublishOptions) []PublishResult {
	initPubOpts(opts)
	var pResults []PublishResult
	for index := 0; index < opts.Count; index++ {
		for id := 0; id < len(opts.Clients); id++ {
			pr := spub(id, opts.Clients[id], index)
			pResults = append(pResults, pr)
		}
	}
	return pResults
}

// "async publish""
func aspub(id int, client MQTT.Client, freeze *sync.WaitGroup) []PublishResult {
	var pResults []PublishResult
	var waitTime time.Duration
	firstFlag := true
	sid := fmt.Sprintf("%05d", id)
	clientID := fmt.Sprintf("%s-%s", publishPid, sid)
	message := getMessage(messageSize - len(clientID) - 35 - 2) //30 => nanoTimeStamp, 2=> / /
	topic := fmt.Sprintf(baseTopic+"%s", sid)

	for index := 0; index < count; index++ {
		if firstFlag {
			if maxIntarval > 0 {
				waitTime = RandomInterval(maxIntarval)
			}
			freeze.Wait()
			if waitTime > 0 {
				time.Sleep(waitTime)
			}
			firstFlag = false
		} else {
			message = getMessage(messageSize - len(clientID) - 35 - 2) //30 => nanoTimeStamp, 2=> "//"
			if maxIntarval > 0 {
				waitTime = RandomInterval(maxIntarval)
				time.Sleep(waitTime)
			}
		}

		startTime := time.Now()
		messageID := startTime.Format(RFC3339NanoForMQTT)
		message = clientID + "/" + messageID + "/" + message
		token := client.Publish(topic, qos, false, message)
		waitStartTime := time.Now()
		token.Wait()
		endTime := time.Now()

		var pResult PublishResult
		pResult.StartTime = startTime
		pResult.WaitStartTime = waitStartTime
		pResult.EndTime = endTime
		pResult.LeadDuration = waitStartTime.Sub(startTime)
		pResult.WaitDuration = endTime.Sub(waitStartTime)
		pResult.TotalDuration = endTime.Sub(startTime)
		pResult.DurTime = endTime.Sub(startTime)
		pResult.Topic = topic
		pResult.ClientID = clientID
		pResult.MessageID = messageID
		pResults = append(pResults, pResult)
	}
	return pResults
}

// AsyncPublish is
func AsyncPublish(opts PublishOptions) []PublishResult {
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

	sort.Sort(pResultSort(pResults))

	return pResults
}

// LoadPublish is
func LoadPublish() {

}
