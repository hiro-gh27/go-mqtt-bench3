package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"

	"sort"

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
	var allTimeStamp []time.Time
	for _, s := range sRestults {
		allTimeStamp = append(allTimeStamp, s.PublishTime, s.SubscribeTime)
		RTT := s.SubscribeTime.Sub(s.PublishTime)
		fmt.Printf("ClientID=%s, pubTime=%s, subTime=%s, RTT=%s\n",
			s.ClientID, s.PublishTime, s.SubscribeTime, RTT)
		totalRTT += RTT
	}

	sort.Sort(pubsub.TimeSort(allTimeStamp))
	for _, t := range allTimeStamp {
		fmt.Println(t)
	}

	pubsubDuration := allTimeStamp[len(allTimeStamp)-1].Sub(allTimeStamp[0])
	fmt.Printf("finalSubTime=%s, firstPubTime=%s, duration=%s\n",
		allTimeStamp[len(allTimeStamp)-1], allTimeStamp[0], pubsubDuration)

	nanoPubsubDuration := float64(pubsubDuration.Nanoseconds())
	totalClientsNum := float64(len(sRestults)) * opts.HostNum
	//pubsubThoughtput := nanoPubsubDuration / totalClientsNum
	fmt.Printf("totalClientNum=%f, nanoPubsubDuration=%f, pow10(6)=%f\n",
		totalClientsNum, nanoPubsubDuration, float64(math.Pow10(6)))
	pubsubThoughtput := totalClientsNum / (nanoPubsubDuration / float64(math.Pow10(6)))
	fmt.Printf("thoughtput=%fps/ms\n", pubsubThoughtput)
	/*
		if nanoPubsubDuration < math.Pow10(3) {
			pubsubThoughtput := totalClientsNum / nanoPubsubDuration
			fmt.Printf("thoughtput=%fps/ns\n", pubsubThoughtput)
		} else if nanoPubsubDuration < math.Pow10(6) {
			pubsubThoughtput := totalClientsNum / (nanoPubsubDuration / float64(math.Pow10(3)))
			fmt.Printf("thoughtput=%fps/μs\n", pubsubThoughtput)
		} else if nanoPubsubDuration < math.Pow10(9) {
			fmt.Printf("totalClientNum=%f, nanoPubsubDuration=%f, pow10(6)=%f\n",
				totalClientsNum, nanoPubsubDuration, float64(math.Pow10(6)))
			pubsubThoughtput := totalClientsNum / (nanoPubsubDuration / float64(math.Pow10(6)))
			fmt.Printf("thoughtput=%fps/ms\n", pubsubThoughtput)
		} else {
			fmt.Printf("totalClientNum=%f, nanoPubsubDuration=%f, pow10(9)=%f\n",
				totalClientsNum, nanoPubsubDuration, float64(math.Pow10(9)))
			pubsubThoughtput := totalClientsNum / (nanoPubsubDuration / float64(math.Pow10(9)))
			fmt.Printf("thoughtput=%fps/s\n", pubsubThoughtput)
		}
	*/

	subscribeNumber := float64(len(sRestults))
	nanoAverageRTT := float64(totalRTT.Nanoseconds()) / subscribeNumber
	microAverageRTT := nanoAverageRTT / 1000
	millAverageRTT := microAverageRTT / 1000

	fmt.Printf("subscribeNum=%f, totalRTT=%s\n", subscribeNumber, totalRTT)
	fmt.Printf("RTT: nano=>%fns, micro=>%fμs, mill=>%fms\n", nanoAverageRTT, microAverageRTT, millAverageRTT)

	// ブローカーがフリーズすると切断できない, 強制終了も可能にしておく.
	go func() {
		signalchan := make(chan os.Signal, 1)
		signal.Notify(signalchan, os.Interrupt)
		<-signalchan
		os.Exit(0)
	}()

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
