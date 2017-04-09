package pubsub

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type sortResults []ConnectResult

func (x sortResults) Len() int { return len(x) }
func (x sortResults) Less(i, j int) bool {
	itime := x[i].StartTime
	jtime := x[j].StartTime
	dtime := jtime.Sub(itime)
	return dtime > 0
}
func (x sortResults) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

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

// RandomMessage is
func RandomMessage(strlen int) string {
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
	return string(message)

}

// DumpConnectResults is
func DumpConnectResults(cResults []ConnectResult) {
	sort.Sort(sortResults(cResults))
	for _, r := range cResults {
		fmt.Printf("ID=%s, sTime=%s, eTime=%s, Durtime=%s \n",
			r.ClientID, r.StartTime, r.EndTime, r.DurTime)
	}

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
