package pubsub

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var maxInterval int
var average time.Time

// MQTT.clinet=nilに対してdisconnect要求するとpanicに陥るので, 中身があるかどうかをチェックする必要がある.
func iscompleat(results []ConnectResult) ([]MQTT.Client, bool) {
	var clietns []MQTT.Client
	haserr := false
	for _, r := range results {
		if r.Client != nil {
			clietns = append(clietns, r.Client)
		} else {
			haserr = true
		}
	}
	return clietns, haserr
}

/**
 * #### SyncMode ####
 */
func syncconnect(id int, broker string) ConnectResult {
	var cRresult ConnectResult
	prosessID := strconv.FormatInt(int64(os.Getpid()), 16)
	clientID := fmt.Sprintf("%s-%d", prosessID, id)
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	client := MQTT.NewClient(opts)

	startTime := time.Now()
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Printf("Connected error: %s\n", token.Error())
		client = nil
	}
	endTime := time.Now()

	cRresult.StartTime = startTime
	cRresult.EndTime = endTime
	cRresult.DurTime = endTime.Sub(startTime)
	cRresult.Client = client
	cRresult.ClientID = clientID
	return cRresult
}

// SyncConnect is
func SyncConnect(execOpts ConnectOptions) []MQTT.Client {
	var cResults []ConnectResult
	broker := execOpts.Broker
	for id := 0; id < execOpts.ClientNum; id++ {
		r := syncconnect(id, broker)
		cResults = append(cResults, r)
	}
	clients, haserr := iscompleat(cResults)
	if haserr {
		SyncDisconnect(clients)
		os.Exit(0)
	}
	time.Sleep(3000 * time.Millisecond)
	DumpConnectResults(cResults)
	return clients
}

/**
 * ### AsyncMode ###
 */
func asynconnect(id int, broker string, freeze *sync.WaitGroup) ConnectResult {
	var cRresult ConnectResult
	var waitTime time.Duration

	prosessID := strconv.FormatInt(int64(os.Getpid()), 16)
	clientID := fmt.Sprintf("%s-%d", prosessID, id)
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	client := MQTT.NewClient(opts)

	if maxInterval > 0 {
		waitTime = RandomInterval(maxInterval)
	}
	freeze.Wait()
	if waitTime > 0 {
		time.Sleep(waitTime)
	}

	startTime := time.Now()
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Printf("Connected error: %s\n", token.Error())
		client = nil
	}
	endTime := time.Now()

	cRresult.StartTime = startTime
	cRresult.EndTime = endTime
	cRresult.DurTime = endTime.Sub(startTime)
	cRresult.Client = client
	cRresult.ClientID = clientID
	return cRresult
}

// AsyncConnect is
func AsyncConnect(execOpts ConnectOptions) []MQTT.Client {
	var cResults []ConnectResult
	maxInterval = execOpts.MaxInterval
	wg := &sync.WaitGroup{}
	freeze := &sync.WaitGroup{}
	freeze.Add(1)
	broker := execOpts.Broker
	for id := 0; id < execOpts.ClientNum; id++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			r := asynconnect(id, broker, freeze)
			cResults = append(cResults, r)
		}(id)
	}
	time.Sleep(3 * time.Second)
	freeze.Done()
	wg.Wait()

	clients, haserr := iscompleat(cResults)
	if haserr {
		fmt.Println("### Error!! ###")
		SyncDisconnect(clients)
		os.Exit(0)
	}
	DumpConnectResults(cResults)

	return clients
}

// SyncDisconnect is
func SyncDisconnect(clinets []MQTT.Client) {
	for _, c := range clinets {
		c.Disconnect(250)
	}
}

// AsyncDisconnect is
func AsyncDisconnect(clients []MQTT.Client) {
	wg := sync.WaitGroup{}
	for _, c := range clients {
		wg.Add(1)
		go func(c MQTT.Client) {
			c.Disconnect(250)
			wg.Done()
		}(c)
	}
	wg.Wait()
}

// LoadConnect is
func LoadConnect() {

}

// NomalConnect is don't get time stamps so use by pub or sub.
func NomalConnect(broker string, number int) []MQTT.Client {
	var clients []MQTT.Client
	for index := 0; index < number; index++ {
		prosessID := strconv.FormatInt(int64(os.Getpid()), 16)
		clientID := fmt.Sprintf("%s-%d", prosessID, index)
		opts := MQTT.NewClientOptions()
		opts.AddBroker(broker)
		opts.SetClientID(clientID)
		client := MQTT.NewClient(opts)

		// connect and wait token or error
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			fmt.Printf("Connected error: %s\n", token.Error())
			client = nil
		}
		clients = append(clients, client)
	}

	var containClient []MQTT.Client
	for _, c := range clients {
		if c != nil {
			containClient = append(containClient, c)
		}
	}
	if len(containClient) < len(clients) {
		println("### Error!! ###")
		SyncDisconnect(containClient)
	}

	return clients
}
