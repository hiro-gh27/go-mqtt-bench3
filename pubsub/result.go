package pubsub

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// ConnectResult is
type ConnectResult struct {
	StartTime time.Time     // when trial for connect start
	EndTime   time.Time     // when trial for connect end
	DurTime   time.Duration // =[endtime - starttime]
	Client    MQTT.Client   //client
	ClientID  string        // identification times
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
	ClientID      int       // 大して意味ない気がする, チェック用になるかな
	MessageID     string    //
}
