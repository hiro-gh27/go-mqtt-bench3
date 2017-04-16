package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"sort"

	pubsub "github.com/hiro-gh27/go-mqtt-bench3/pubsub"
)

const base = "go-mqtt-bench/"

func main() {
	// init random seed
	rand.Seed(time.Now().UnixNano())

	// use max CPU's
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	// execute publish benchmark
	opts := initOption()
	var pResults []pubsub.PublishResult
	if opts.AsyncFlag {
		fmt.Println("--- AsyncMode ---")
		pResults = pubsub.AsyncPublish(opts)
	} else {
		fmt.Println("--- SyncMode ---")
		pResults = pubsub.SyncPublish(opts)
	}

	var total []time.Time
	var DEBUG bool
	for _, p := range pResults {
		total = append(total, p.StartTime)
		/*
			fmt.Printf("clientID=%s, start=%s, end=%s, Durtime=%s\n",
				p.ClientID, p.StartTime, p.EndTime, p.DurTime)
		*/
		if DEBUG {
			fmt.Printf("waitstart=%s, wait=%s, total=%s\n",
				p.WaitStartTime, p.WaitDuration, p.TotalDuration)
		}

	}
	sort.Sort(pubsub.TimeSort(total))
	totalTime := total[len(total)-1].Sub(total[0])
	millThroughput := float64(totalTime.Nanoseconds()) / float64(len(pResults)) / 1000000
	publishPerMillsecond := float64(1) / millThroughput
	fmt.Printf("\ntotal count = %d, total=%s, nanoTime=%d, throughput=%fpub/ms\n",
		len(pResults), totalTime, totalTime.Nanoseconds(), publishPerMillsecond /*totalTime.Nanoseconds()/int64(len(pResults))*/)
	pubsub.SyncDisconnect(opts.Clients)

	// export elasticseaech
}

func initOption() pubsub.PublishOptions {
	// for connect
	broker := flag.String("broker", "tcp://{host}:{port}", "URI of MQTT broker (required)")
	clients := flag.Int("clients", 10, "Number of clients")

	// for publish
	qos := flag.Int("qos", 0, "MQTT QoS(0|1|2)")
	retain := flag.Bool("retain", false, "MQTT Retain")
	topic := flag.String("topic", base, "Base topic")
	count := flag.Int("count", 10, "Number of loops per client")
	size := flag.Int("size", 1024, "Message size per publish (byte)")
	asyncmode := flag.Bool("async", false, "ture mean asyncmode")
	trial := flag.Int("trial", 1, "trial is number of how many loops are")
	pubPerMillSecond := flag.Float64("pub/ms", 10, "publish/ms")
	convertTime := flag.Int("time", -1, "when program start")
	flag.Parse()

	if len(os.Args) < 1 {
		fmt.Println("### Error ###")
		flag.Usage()
		os.Exit(0)
	}

	if broker == nil || *broker == "" || *broker == "tcp://{host}:{port}" {
		fmt.Println("Use Default Broker= tcp://10.0.0.4:1883")
		*broker = "tcp://10.0.0.4:1883"
	}

	var executeTime time.Time
	if *convertTime != -1 {
		executeTime = pubsub.GetExecuteTime(*convertTime)
		fmt.Printf("time=%s\n", executeTime)
		if executeTime.IsZero() {
			os.Exit(0)
		}
	}

	connectedClients := pubsub.NomalConnect(*broker, *clients)

	var options pubsub.PublishOptions
	options.Qos = byte(*qos)
	options.Retain = *retain
	options.Topic = *topic
	options.MessageSize = *size
	options.ClientNum = len(connectedClients)
	options.Count = *count
	options.MaxInterval = float64(*clients) / *pubPerMillSecond
	options.AsyncFlag = *asyncmode
	options.Clients = connectedClients
	options.TrialNum = *trial
	options.ExecuteTime = executeTime

	fmt.Printf("\n max Interval=%f \n", options.MaxInterval)

	return options
}
