package main

import (
	"flag"

	pubsub "github.com/hiro-gh27/go-mqtt-bench3/pubsub"
)

func main() {
}

func initOption() pubsub.RttOption {
	qos := flag.Int("qos", 0, "MQTT QoS(0|1|2)")
	retain := flag.Bool("retain", false, "MQTT Retain")
	topic := flag.String("topic", pubsub.Base, "Base topic")
	size := flag.Int("size", 100, "Message size per publish (byte)")
	intervalTime := flag.Int("interval", 0, "Interval time per message (ms)")

	//var options pubsub.RttOptio
	//return options
}
