package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
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
	nanoAverageRTT := float64(totalRTT.Nanoseconds()) / subscribeNumber
	microAverageRTT := nanoAverageRTT / 1000
	millAverageRTT := microAverageRTT / 1000

	fmt.Printf("subscribeNum=%f, totalRTT=%s\n", subscribeNumber, totalRTT)
	fmt.Printf("RTT: nano=>%fns, micro=>%fμs, nano=>%fms\n", nanoAverageRTT, microAverageRTT, millAverageRTT)
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
	startID := flag.Int("id", 0, "startID -> N_cliensID")
	hostNum := flag.Float64("h", 1, "number of subscribe hosts")

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

	var connectedClients []MQTT.Client
	if *startID > 0 {
		connectedClients = pubsub.SpecificConnect(*broker, *clients, *startID)
	} else {
		connectedClients = pubsub.NomalConnect(*broker, *clients)
	}

	var options pubsub.SubscribeOptions
	options.Qos = byte(*qos)
	options.Topic = *topic
	options.HostNum = *hostNum
	options.ClientNum = len(connectedClients)
	options.Clients = connectedClients
	options.StartID = *startID

	return options
}
