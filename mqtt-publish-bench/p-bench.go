package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"time"

	"math"

	MQTT "github.com/eclipse/paho.mqtt.golang"
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
		pResults = pubsub.AsyncPublish(opts)
	} else {
		fmt.Println("--- SyncMode ---")
		pResults = pubsub.SyncPublish(opts)
	}

	fmt.Println("--- Result ---")
	publishNum := float64(len(pResults))
	var DEBUG bool
	var allTimeStamps []time.Time
	for _, p := range pResults {
		allTimeStamps = append(allTimeStamps, p.StartTime, p.WaitStartTime, p.EndTime)
		if DEBUG {
			fmt.Printf("clientID=%s, start=%s, end=%s, Durtime=%s\n",
				p.ClientID, p.StartTime, p.EndTime, p.DurTime)
			fmt.Printf("waitstart=%s, wait=%s, total=%s\n",
				p.WaitStartTime, p.WaitDuration, p.TotalDuration)
		}
	}

	// Publish thoughtput
	fmt.Printf("%d*%d=%dpublish\n", opts.ClientNum, opts.Count, len(pResults))
	sort.Sort(pubsub.TimeSort(allTimeStamps))
	allTimeDuration := allTimeStamps[len(allTimeStamps)-1].Sub(allTimeStamps[0])
	allMillTimeDuration := float64(allTimeDuration.Nanoseconds()) / math.Pow10(6)
	publishThoughtput := publishNum / allMillTimeDuration
	fmt.Printf("totalDuration=%fms, publishThoughtput=%fpub/ms\n",
		allMillTimeDuration, publishThoughtput)

	// brokerが配送処理中にdisconnectすると, 余計な負荷がかかると思う.
	// なので, ctrl+cを入力されるまで待つ
	signalchan := make(chan os.Signal, 1)
	signal.Notify(signalchan, os.Interrupt)
	<-signalchan

	pubsub.SyncDisconnect(opts.Clients)
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
	asyncmode := flag.Bool("async", true, "ture mean asyncmode")
	trial := flag.Int("trial", 1, "trial is number of how many loops are")
	exp2Num := flag.Float64("exp2", 0, "publish/ms")
	convertTime := flag.Int("time", -1, "when program start")
	startID := flag.Int("id", 0, "Number of start clientID")
	//pubPerMillSecond := flag.Float64("pub/ms", 10, "publish/ms")
	flag.Parse()

	if len(os.Args) < 1 {
		fmt.Println("### Error ###")
		flag.Usage()
		os.Exit(0)
	}

	if broker == nil || *broker == "" || *broker == "tcp://{host}:{port}" {
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

	var connectedClients []MQTT.Client
	if *startID > 0 {
		connectedClients = pubsub.SpecificConnect(*broker, *clients, *startID)
	} else {
		connectedClients = pubsub.NomalConnect(*broker, *clients)
	}
	// 安定させる
	time.Sleep(3 * time.Second)

	var maxInterval float64
	thoughtput := math.Exp2(*exp2Num)
	maxInterval = float64(*clients) / thoughtput

	var options pubsub.PublishOptions
	options.Qos = byte(*qos)
	options.Retain = *retain
	options.Topic = *topic
	options.MessageSize = *size
	options.ClientNum = len(connectedClients)
	options.Count = *count
	options.MaxInterval = maxInterval
	options.AsyncFlag = *asyncmode
	options.Clients = connectedClients
	options.TrialNum = *trial
	options.ExecuteTime = executeTime
	options.StartID = *startID

	fmt.Println("--- exec info ---")
	fmt.Printf("qos=%b\n", options.Qos)
	fmt.Printf("clientNum=%d\n", options.ClientNum)
	fmt.Printf("count=%d\n", options.Count)
	fmt.Printf("messsageSize=%dbyte\n", options.MessageSize)
	fmt.Printf("thoughtput=%5fpub/ms\n", thoughtput)
	fmt.Printf("interval=%fms\n", options.MaxInterval)

	return options
}
