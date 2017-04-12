package pubsub

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

//SortResults is
type SortResults []ConnectResult

// TimeSort is
type TimeSort []time.Time

// sortinterface
func (x SortResults) Len() int { return len(x) }
func (x SortResults) Less(i, j int) bool {
	itime := x[i].WaitStartTime
	jtime := x[j].WaitStartTime
	dtime := jtime.Sub(itime)
	return dtime > 0
}
func (x SortResults) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// time sort
func (x TimeSort) Len() int { return len(x) }
func (x TimeSort) Less(i, j int) bool {
	itime := x[i]
	jtime := x[j]
	dtime := jtime.Sub(itime)
	return dtime > 0
}
func (x TimeSort) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// MsgTsLayout means [Message_TimeStamp_Layout]
const (
	MsgStampLayout     = time.StampNano + " 2006"
	RFC3339NanoForMQTT = "2006-01-02T15:04:05.000000000Z07:00"
)

//use randomMessage
const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// RandomInterval is return duration time for asyncMode
func RandomInterval(max int) time.Duration {
	var td time.Duration
	if maxInterval > 0 {
		interval := rand.Intn(maxInterval)
		td = time.Duration(interval) * time.Millisecond
	}
	return td
}

func getMessageAndID(strlen int) (string, string) {
	nanoStamp := time.Now().Format(time.StampNano)
	strlen = strlen - len(nanoStamp)

	message := make([]byte, strlen)
	cache, remain := rand.Int63(), letterIdxMax
	for i := strlen - 1; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		idx := int(cache & letterIdxMask)
		if idx < len(letters) {
			message[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	str := fmt.Sprintf("%s%s", nanoStamp, string(message))
	return nanoStamp, str
}

func getMessage(strlen int) string {
	//	strlen = strlen - 25 // 25= len(time.nanostamp)
	if strlen < 0 {
		strlen = 1
	}

	message := make([]byte, strlen)
	cache, remain := rand.Int63(), letterIdxMax
	for i := strlen - 1; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		idx := int(cache & letterIdxMask)
		if idx < len(letters) {
			message[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	//str := fmt.Sprintf("%s%s", nanoStamp, string(message))
	return string(message)
}

func getMsgStamp() string {
	return time.Now().Format(MsgStampLayout)
}

// DumpConnectResults is
func DumpConnectResults(cResults []ConnectResult) {
	sort.Sort(SortResults(cResults))
	/*
		for _, r := range cResults {
			fmt.Printf("ID=%s, sTime=%s, eTime=%s, Durtime=%s \n",
				r.ClientID, r.StartTime, r.EndTime, r.DurTime)
		}
	*/

	fastTime := cResults[0].StartTime
	slowTime := cResults[len(cResults)-1].StartTime
	durtime := slowTime.Sub(fastTime)
	clientNum := int64(len(cResults))
	nanoTime := durtime.Nanoseconds()                 //nano秒に変換
	perClient := nanoTime / clientNum                 //1connectにかかったnano秒
	throuput := float64(1000000 / float64(perClient)) //1ms=1000000. 1ms/コネクション時間

	fmt.Printf("#### dtime= %s, clientNum=%d, duration=%dns, %dns/clinet, %f client/ms #### \n",
		durtime, clientNum, nanoTime, perClient, throuput)
}
