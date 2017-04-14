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
	rand.Seed(time.Now().UnixNano())

	// use max CPU's
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	opts := initOption()

	sRestults := pubsub.Subscribe(opts)
	var totalRTT time.Duration
	for _, s := range sRestults {
		RTT := s.SubscribeTime.Sub(s.PublishTime)
		fmt.Printf("ClientID=%s, pubTime=%s, subTime=%s, RTT=%s\n",
			s.ClientID, s.PublishTime, s.SubscribeTime, RTT)
		totalRTT += RTT
	}

	subscribeNumber := float64(len(sRestults))
	averageRTT := float64(totalRTT.Nanoseconds()) / subscribeNumber

	fmt.Printf("subscribeNum=%f, totalRTT=%s, ave=%f\n", subscribeNumber, totalRTT, averageRTT)

	pubsub.SyncDisconnect(opts.Clients)
	/*
		export
	*/
}

func initOption() pubsub.SubscribeOptions {
	// for connect
	broker := flag.String("broker", "tcp://{host}:{port}", "URI of MQTT broker (required)")
	clients := flag.Int("clients", 10, "Number of clients")
	// for subscribe
	qos := flag.Int("qos", 0, "MQTT QoS(0|1|2)")
	topic := flag.String("topic", base, "Base topic")

	flag.Parse()
	if len(os.Args) < 0 {
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

	var options pubsub.SubscribeOptions
	options.Qos = byte(*qos)
	options.Topic = *topic
	options.ClientNum = len(connectedClients)
	options.Clients = connectedClients

	return options
}
