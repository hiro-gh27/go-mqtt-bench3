package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	pubsub "github.com/hiro-gh27/go-mqtt-bench2/pubsub"
)

var asyncFlag bool

func main() {
	rand.Seed(time.Now().UnixNano())
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	var clients []MQTT.Client
	opts := initOption()
	if opts.AsyncFlag {
		clients = pubsub.AsyncConnect(opts)
	} else {
		clients = pubsub.SyncConnect(opts)
	}
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
