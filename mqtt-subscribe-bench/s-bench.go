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

	fmt.Println("--- Result ---")
	var RTTs []time.Duration
	var totalRTT time.Duration
	var allTimeStamp []time.Time
	for _, s := range sRestults {
		allTimeStamp = append(allTimeStamp, s.PublishTime, s.SubscribeTime)
		RTT := s.SubscribeTime.Sub(s.PublishTime)
		RTTs = append(RTTs, RTT)
		totalRTT = totalRTT + RTT
	}

	// thoughtput
	sort.Sort(pubsub.TimeSort(allTimeStamp))
	pubsubDuration := allTimeStamp[len(allTimeStamp)-1].Sub(allTimeStamp[0])
	pubsubMillDuration := float64(pubsubDuration.Nanoseconds()) / math.Pow10(6)
	totalPubSubNum := float64(len(sRestults)) * opts.HostNum
	fmt.Printf("subscribeNum=%f, duration=%fms\n", totalPubSubNum, pubsubMillDuration)
	pubsubThoughtput := totalPubSubNum / pubsubMillDuration
	fmt.Printf("Thoughtput: %fps/ms\n", pubsubThoughtput)

	// RTT
	subscribeNumber := float64(len(sRestults))
	nanoAverageRTT := float64(totalRTT.Nanoseconds()) / subscribeNumber
	millAverageRTT := nanoAverageRTT / math.Pow10(6)
	fmt.Printf("subscribeNum=%f, totalRTT=%fns\n", subscribeNumber, nanoAverageRTT)
	fmt.Printf("RTT: %fms\n", millAverageRTT)

	// SD
	var nanoDispersions float64
	for _, RTT := range RTTs {
		nanoDispersion := (nanoAverageRTT - float64(RTT.Nanoseconds())) * (nanoAverageRTT - float64(RTT.Nanoseconds()))
		nanoDispersions = nanoDispersions + nanoDispersion
		fmt.Printf("RTT=%s, nanoDispersion=%fns\n", RTT, nanoDispersion)
	}
	dispersion := (nanoDispersions / float64(len(RTTs))) / math.Pow10(12)
	sd := math.Sqrt(dispersion)

	fmt.Printf("dispersion: =%f\nmssd: =%fms\n", dispersion, sd)

	// ブローカーがフリーズすると切断できない, 強制終了も可能にしておく.
	go func() {
		signalchan := make(chan os.Signal, 1)
		signal.Notify(signalchan, os.Interrupt)
		<-signalchan
		os.Exit(0)
	}()

	pubsub.SyncDisconnect(opts.Clients)
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

	fmt.Println("--- exec info ---")
	fmt.Printf("clientsNum=%d\n", options.ClientNum)
	fmt.Printf("qos=%b\n", options.Qos)

	return options
}
