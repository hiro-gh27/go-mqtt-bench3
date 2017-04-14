package pubsub

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Base is mean topic
const Base = "go-mqtt-bench/"

// ExecOptions is
type ExecOptions struct {
	Broker      string // Broker URI
	Qos         byte   // QoS(0|1|2)
	Retain      bool   // Retain
	Debug       bool   //デバック
	Topic       string // Topicのルート
	Method      string // 実行メソッド
	ClientNum   int    // クライアントの同時実行数
	Count       int    // 1クライアント当たりのメッセージ数
	MessageSize int    // 1メッセージのサイズ(byte)
	SleepTime   int    // 実行前の待機時間(ms)
	MaxInterval int    // メッセージ毎の実行間隔時間(ms)
	Test        bool   //テスト
	TrialNum    int    //試行回数
	SynBacklog  int    //net.ipv4.tcp_max_syn_backlog =
	AsyncFlag   bool   //ture mean asyncmode
}

// ConnectOptions is
type ConnectOptions struct {
	Broker      string  // Broker URI
	AsyncFlag   bool    //ture mean asyncmode
	ClientNum   int     // クライアントの同時実行数
	MaxInterval float64 // メッセージ毎の実行間隔時間(ms)
}

// PublishOptions is
type PublishOptions struct {
	Qos         byte          // QoS(0|1|2)
	Retain      bool          // Retain
	Topic       string        // Topicのルート
	ClientNum   int           // クライアントの同時実行数
	Count       int           // 1クライアント当たりのメッセージ数
	MessageSize int           // 1メッセージのサイズ(byte)
	MaxInterval float64       // メッセージ毎の実行間隔時間(ms)
	AsyncFlag   bool          //ture mean asyncmode
	Clients     []MQTT.Client //クライアントをスライスで確保!!
	TrialNum    int
}

// SubscribeOptions is
type SubscribeOptions struct {
	Qos       byte          // QoS(0|1|2)
	Topic     string        // Topicのルート
	ClientNum int           // クライアントの同時実行数
	Clients   []MQTT.Client //クライアントをスライスで確保!!
}

// RttOption is
type RttOption struct {
	Qos         byte
	Retain      bool   // Retain
	Topic       string // Topicのルート
	MessageSize int    // 1メッセージのサイズ(byte)
	MaxInterval int    // メッセージ毎の実行間隔時間(ms)
}

// LoadOptions is
type LoadOptions struct {
}
