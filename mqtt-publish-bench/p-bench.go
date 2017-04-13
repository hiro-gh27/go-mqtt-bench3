package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

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

	for _, p := range pResults {
		/*
			fmt.Printf("clientID=%s, start=%s, end=%s, Durtime=%s\n",
				p.ClientID, p.StartTime, p.EndTime, p.DurTime)
		*/
		fmt.Printf("lead=%s, wait=%s, total=%s", p.LeadDuration, p.WaitDuration, p.TotalDuration)
	}
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
	size := flag.Int("size", 100, "Message size per publish (byte)")
	intervalTime := flag.Int("interval", 0, "Interval time per message (ms)")
	asyncmode := flag.Bool("async", false, "ture mean asyncmode")
	trial := flag.Int("trial", 1, "trial is number of how many loops are")

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
	// make clients
	connectedClients := pubsub.NomalConnect(*broker, *clients)

	var options pubsub.PublishOptions
	options.Qos = byte(*qos)
	options.Retain = *retain
	options.Topic = *topic
	options.MessageSize = *size
	options.ClientNum = len(connectedClients)
	options.Count = *count
	options.MaxInterval = *intervalTime
	options.AsyncFlag = *asyncmode
	options.Clients = connectedClients
	options.TrialNum = *trial

	return options
}
