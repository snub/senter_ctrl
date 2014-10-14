package main

import (
	"fmt"
	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"os"
	"os/signal"
	"regexp"
    "strconv"
)

const (
	startupTopic string = "/controller/+/startup"

	sensorCountTopic     string = "/controller/+/sensors/count"
	sensorDiscoveryTopic string = "/controller/+/sensors/discovery"
	sensorTempTopic      string = "/controller/+/sensor/+/temp"

	controllerTopic string = "/controller/#"
)

var regexStartup *regexp.Regexp = regexp.MustCompile("^/controller/(.+)/startup$")

var handler mqtt.MessageHandler = func(client *mqtt.MqttClient, msg mqtt.Message) {
	topic := msg.Topic()
	message := msg.Payload()
	fmt.Printf("TOPIC: %s\n", topic)
	fmt.Printf("MSG: %s\n", message)
	match := regexStartup.FindStringSubmatch(topic)
	fmt.Printf("match: %v\n", match)
	if len(match) > 1 {
		mac := match[1]
        timestamp, err := strconv.Atoi(string(message))
        if err != nil {
            fmt.Println(err)
        }
		fmt.Printf("mac: %s\n", mac)
		if !existsSensorController(mac) {
			_ = saveSensorController(mac, timestamp)
		}
	}
}

func init() {
	err := connectDatabase("localhost", "3306", "root", "qwerty123", "sensor_center")
    if err != nil {
        fmt.Println(err)
    }
}

func main() {
	fmt.Printf("Sensor Controller Temperature Storer\n")

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetDefaultPublishHandler(handler)

	c := mqtt.NewClient(opts)
	_, err := c.Start()
	if err != nil {
		panic(err)
	}

	filter, err := mqtt.NewTopicFilter(controllerTopic, 0)
	if err != nil {
		panic(err)
	}

	if receipt, err := c.StartSubscription(nil, filter); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	death := make(chan os.Signal, 1)
	signal.Notify(death, os.Interrupt, os.Kill)

Terminates:
	for {
		select {
		case <-death:
			fmt.Println("Signal received...")
			break Terminates
		}
	}

	fmt.Println("Stopping...")

	//unsubscribe from /go-mqtt/sample
	if receipt, err := c.EndSubscription(controllerTopic); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	c.Disconnect(250)
}
