package main

import (
	"flag"
	"fmt"
	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/snub/senter"
	"os"
	"os/signal"
)

var configFileName string

var defaultHandler mqtt.MessageHandler = func(client *mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	message := msg.Payload()
	logger.Printf("defaultHandler, topic: %s, message: %s\n", topic, message)
	logger.Println("does nothing")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	const (
		defaultFileName = "senter-config.json"
		usage           = "Senter configuration file with database and MQTT settings"
	)
	flag.StringVar(&configFileName, "c", defaultFileName, usage)
	flag.StringVar(&configFileName, "config", defaultFileName, usage)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	logger.Println("... sensor controller daemon ...")
	logger.Printf("using configuration: %s\n", configFileName)

	config, err := LoadConfig(configFileName)
	if err != nil {
		os.Exit(1)
	}

	err = senter.InitDatabase(config.Database)
	if err != nil {
		logger.Printf("unable initialize database: %s\n", err)
	}
	defer senter.CloseDatabase()
	senter.EnableDatabaseLogger()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Mqtt.Broker)
	opts.SetClientID("senter-ctrl")
	opts.SetDefaultPublishHandler(defaultHandler)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Panicf("unable to connect MQTT broker: %s\n", token.Error())
	} else {
		logger.Printf("connected to MQTT broker: %s\n", config.Mqtt.Broker)
	}

	for topic, topicHandler := range topicsAndHandlers {
		if token := c.Subscribe(topic, byte(0), topicHandler); token.Wait() && token.Error() != nil {
			logger.Panicf("unable to subscribe to topic: %s - %s\n", topic, token.Error())
		} else {
			logger.Printf("subscribed to topic: %s\n", topic)
		}
	}

	death := make(chan os.Signal, 1)
	signal.Notify(death, os.Interrupt, os.Kill)

Terminates:
	for {
		select {
		case <-death:
			logger.Println("signal received...")
			break Terminates
		}
	}

	logger.Println("stopping...")

	for topic, _ := range topicsAndHandlers {
		if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
			logger.Panicf("unable to unsubscribe topic: %s - %s\n", topic, token.Error())
		} else {
			logger.Printf("unsubscribed topic: %s\n", topic)
		}
	}

	c.Disconnect(250)

	logger.Println("done.")
}
