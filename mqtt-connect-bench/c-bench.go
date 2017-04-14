package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

	"sort"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	pubsub "github.com/hiro-gh27/go-mqtt-bench3/pubsub"
)

var asyncFlag bool

func main() {
	rand.Seed(time.Now().UnixNano())
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	var clients []MQTT.Client
	var cResults []pubsub.ConnectResult
	opts := initOption()
	pid := strconv.FormatInt(int64(os.Getpid()), 16)
	fmt.Println(pid)
	if opts.AsyncFlag {
		fmt.Println("--- AsyncMode ---")
		cResults, clients = pubsub.AsyncConnect(opts)
	} else {
		fmt.Println("--- SyncMode ---")
		cResults, clients = pubsub.SyncConnect(opts)
	}

	/**
	* keep mode
	 */

	sort.Sort(pubsub.SortResults(cResults))

	var totalRTT time.Duration
	var leadTotal time.Duration
	var starts []time.Time
	//var wStarts []time.Time
	for _, r := range cResults {
		/*
			fmt.Printf("startTime=%s, waitStartTime=%s, endTime=%s\n",
				r.StartTime, r.WaitStartTime, r.EndTime)
		*/
		/*

			fmt.Printf("Lead=%s, wait=%s, total=%s\n",
				r.LeadDuration, r.WaitDuration, r.TotalDuration)
		*/

		/*
			fmt.Printf("WaitStartTime=%s, RTT=%s\n", r.WaitStartTime, r.WaitDuration)
		*/
		totalRTT = totalRTT + r.WaitDuration
		leadTotal = leadTotal + r.LeadDuration
		starts = append(starts, r.StartTime)
		starts = append(starts, r.WaitStartTime)
		starts = append(starts, r.EndTime)
	}
	/*
		for _, t := range starts {
			fmt.Printf("start=%s\n", t)
		}
	*/
	/*
		sort.Sort(pubsub.TimeSort(starts))
		fmt.Printf("\n sort now \n")
		for _, t := range starts {
			fmt.Printf("start=%s\n", t)
		}
	*/
	/*
		d := starts[len(starts)-1].Sub(starts[0])
		fmt.Printf("\n dur=%s スループット=%s\n", d, d/5000)
	*/
	//num := len(clients)
	//fmt.Printf("Lead: total=%s RTT: total=%s, ave=%s\n", leadTotal, totalRTT, totalRTT/5000)

	/*
		同期バージョン, 非同期バージョンの戻り値をcResultにして
		ここの位置で, Exportできるといいのではと考えている.
	*/
	pubsub.SyncDisconnect(clients)

}

func initOption() pubsub.ConnectOptions {

	broker := flag.String("broker", "tcp://{host}:{port}", "URI of MQTT broker (required)")
	clients := flag.Int("clients", 10, "Number of clients")
	maxinterval := flag.Int("interval", 0, "Interval time per message (ms)")
	asyncmode := flag.Bool("async", false, "ture mean asyncmode")
	flag.Parse()

	if len(os.Args) < 1 {
		fmt.Println("call here")
		flag.Usage()
		os.Exit(0)
	}

	if broker == nil || *broker == "" || *broker == "tcp://{host}:{port}" {
		fmt.Println("Use Default Broker= tcp://10.0.0.4:1883")
		*broker = "tcp://10.0.0.4:1883"
	}

	options := pubsub.ConnectOptions{}
	options.Broker = *broker
	options.ClientNum = *clients
	options.MaxInterval = *maxinterval
	options.AsyncFlag = *asyncmode
	return options
}
