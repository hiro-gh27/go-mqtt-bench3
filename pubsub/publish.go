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
var maxIntarval float64
var qos byte
var publishPid string
var baseTopic string
var trial int
var count int
var publishDebug = true

func initPubOpts(opts PublishOptions) {
	baseTopic = opts.Topic
	count = opts.Count
	messageSize = opts.MessageSize
	publishPidStr := strconv.FormatInt(int64(os.Getpid()), 16)
	publishPid = fmt.Sprintf("%05s", publishPidStr)
	//publishPid = strconv.FormatInt(int64(os.Getpid()), 16)
	maxIntarval = opts.MaxInterval
	trial = opts.TrialNum
	qos = opts.Qos
	fmt.Printf("pid=%s\n", publishPid)
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
	pResults := make([]PublishResult, count)
	startTimeGaps := make([]time.Duration, count)
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
			//fmt.Printf("interval=%s\n", waitTime)
			firstFlag = false
		} else {
			message = getMessage(messageSize - len(clientID) - 35 - 2) //30 => nanoTimeStamp, 2=> "//"
			if maxIntarval > 0 {
				waitTime = time.Duration(maxIntarval * 1000000)
				waitTime = waitTime - startTimeGaps[index-1]
				if waitTime > 0 {
					time.Sleep(waitTime)
				}
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
		pResults[index] = pResult

		if publishDebug {
			fmt.Printf("clientID=%s, publishTime=%s\n", clientID, startTime)
		}

		// 1個前の実行時間から理想的な実行時間を求める, その理想的な時間と, 今回行われた時間の差分を
		// 次の待ち時間から減らすことで, 誤差が積み重なっていくのを避けるってゆう配慮...
		if index > 0 {
			idealStartTime := pResults[0].StartTime.Add(time.Duration(int(maxIntarval*1000000) * index))
			if startTime.Sub(idealStartTime) > 0 {
				startTimeGaps[index] = startTime.Sub(idealStartTime)
			}
		}
	}
	return pResults

}

func periodAsync() {

}

// AsyncPublish is
func AsyncPublish(opts PublishOptions) []PublishResult {
	initPubOpts(opts)
	var pResults []PublishResult
	pResultPacks := make([][]PublishResult, opts.ClientNum)
	wg := &sync.WaitGroup{}
	freeze := &sync.WaitGroup{}
	freeze.Add(1)
	for index := 0; index < len(opts.Clients); index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			id := opts.StartID + index
			pResultPacks[index] = aspub(id, opts.Clients[index], freeze)
		}(index)
	}

	// オプションで, 実行時間を指定している場合に同期して実行する.
	if opts.ExecuteTime.IsZero() {
		time.Sleep(3 * time.Second)
	} else {
		gapTimer := time.NewTimer(opts.ExecuteTime.Sub(time.Now()))
		<-gapTimer.C
	}
	//fmt.Printf("execute time=%s", time.Now())
	freeze.Done()

	wg.Wait()
	fmt.Printf("pResultPacks len=%d", len(pResultPacks))
	for _, packs := range pResultPacks {
		for _, p := range packs {
			pResults = append(pResults, p)
		}
	}
	/*
		sort.Sort(durationSort(tds))
		for _, t := range tds {
			fmt.Printf("td=%s\n", t)
		}
	*/

	//	sort.Sort(pResultSort(pResults))

	return pResults
}

// LoadPublish is
func LoadPublish() {

}
