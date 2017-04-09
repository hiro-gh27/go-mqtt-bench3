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
	DurTime   time.Duration // =[endtime - starttime]
	Topic     string        // = MessageID(basetopic/clientID/tiralNum)
	ClientID  string        //
	Message   string
}

// SubscribeResult is
type SubscribeResult struct {
	timestamp time.Time //
	MessageID string    //
}
