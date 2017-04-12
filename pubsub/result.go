package pubsub

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// ConnectResult is
type ConnectResult struct {
	//Time Stamps
	StartTime     time.Time // when trial for connect start
	WaitStartTime time.Time // when publish end, and wait start
	EndTime       time.Time // when trial for connect end
	//Durations
	LeadDuration  time.Duration
	WaitDuration  time.Duration
	TotalDuration time.Duration
	//General　info
	Client   MQTT.Client //client
	ClientID string      // identification times
}

// PublishResult is
type PublishResult struct {
	StartTime time.Time     // when trial for connect start
	EndTime   time.Time     // when trial for connect end
	DurTime   time.Duration // = endtime - starttime
	Topic     string        // = basetopic/id
	ClientID  string        // = pid + id(index)
	MessageID string        // = time.Stamp()
	Message   string        // RandomMessage
}

// SubscribeResult is
type SubscribeResult struct {
	SubscribeTime time.Time //
	Topic         string    //
	ClientID      string    // 大して意味ない気がする, チェック用になるかな
	MessageID     string    //

	PublisherID string // 作ろうと思ったけど, 毎回長さが変わる可能性あるからやめた.
	PublishTime time.Time
}
