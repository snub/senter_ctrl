package main

import (
	"flag"
	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/snub/senter"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
)

const (
	//startupTopic string = "/controller/+/startup"

	//sensorCountTopic     string = "/controller/+/sensors/count"
	//sensorDiscoveryTopic string = "/controller/+/sensors/discovery"
	//sensorTempTopic      string = "/controller/+/sensor/+/temp"

	controllerTopic string = "/controller/#"
)

var (
	regexStartup   *regexp.Regexp = regexp.MustCompile("^/controller/(.+)/startup$")
	reqexDiscovery *regexp.Regexp = regexp.MustCompile("^/controller/(.+)/sensors/discovery$")
	reqexTemp      *regexp.Regexp = regexp.MustCompile("^/controller/(.+)/sensor/(.+)/temp$")
)

// TODO error handling
var handler mqtt.MessageHandler = func(client *mqtt.MqttClient, msg mqtt.Message) {
	topic := msg.Topic()
	message := msg.Payload()
	logger.Printf("topic: %s, message: %s\n", topic, message)

	// check for startup topic
	match := regexStartup.FindStringSubmatch(topic)
	logger.Printf("startup topic match: %v\n", match)
	if len(match) == 2 {
		macAddress := match[1]
		timestamp, err := strconv.Atoi(string(message))
		if err != nil {
			logger.Printf("unable to convert timestamp: %s\n", err)
		}
		logger.Printf("macAddress: %s, timestamp: %d\n", macAddress, timestamp)
		controller := senter.LoadControllerByMacAddress(macAddress)
		controller.SetLastStartup(int64(timestamp))
		controller.Save()
		logger.Printf("controller: %v\n", controller)
	}

	// check for discovery topic
	match = reqexDiscovery.FindStringSubmatch(topic)
	logger.Printf("discovery topic match: %v\n", match)
	if len(match) == 2 {
		controllerMacAddress := match[1]
		logger.Printf("controllerMacAddress: %s\n", controllerMacAddress)

		deviceAddress := string(message)
		sensor := senter.LoadSensorByDeviceAddress(deviceAddress)
		sensor.Save()
		logger.Printf("sensor: %v\n", sensor)
	}

	// check for temperature topic
	// TODO check is sensor is stored in database
	match = reqexTemp.FindStringSubmatch(topic)
	logger.Printf("temperature topic match: %v\n", match)
	if len(match) == 3 {
		controllerMacAddress := match[1]
		sensorDeviceAddress := match[2]
		logger.Printf("controllerMacAddress: %s, sensorDeviceAddress: %s\n", controllerMacAddress, sensorDeviceAddress)
		splitMsg := strings.Split(string(message), ",")
		if len(splitMsg) == 2 {
			timestamp, err := strconv.Atoi(splitMsg[0])
			if err != nil {
				logger.Printf("unable to convert timestamp: %s\n", err)
			}
			value, err := strconv.ParseFloat(splitMsg[1], 32)
			if err != nil {
				logger.Printf("unable to parse temperature: %s\n", err)
			}
			logger.Printf("timestamp: %d, value: %f\n", timestamp, value)

			sensor := senter.LoadSensorByDeviceAddress(sensorDeviceAddress)
			temperature := senter.NewTemperature(sensor, int64(timestamp), float32(value))
			temperature.Save()
			logger.Printf("temperature: %v\n", temperature)
		}
	}
}

var configFileName string

func init() {
	const (
		defaultFileName = "senter-config.json"
		usage           = "Senter configuration file with database and MQTT settings"
	)
	flag.StringVar(&configFileName, "c", defaultFileName, usage)
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	logger.Println("... sensor temperature storer ...")
	logger.Printf("using configuration: %s\n", configFileName)

	config, err := LoadConfig(configFileName)
	logger.Printf("loaded config: %+v\n", config)
	if err != nil {
		os.Exit(1)
	}

	err = senter.InitDatabase(config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.Database)
	if err != nil {
		logger.Printf("unable initialize database: %s\n", err)
	}
	defer senter.CloseDatabase()
	senter.EnableDatabaseLogger()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Mqtt.Broker)
	opts.SetDefaultPublishHandler(handler)

	c := mqtt.NewClient(opts)
	_, err = c.Start()
	if err != nil {
		panic(err)
	}

	filter, err := mqtt.NewTopicFilter(controllerTopic, 0)
	if err != nil {
		panic(err)
	}

	if receipt, err := c.StartSubscription(nil, filter); err != nil {
		logger.Println(err)
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
			logger.Println("Signal received...")
			break Terminates
		}
	}

	logger.Println("Stopping...")

	if receipt, err := c.EndSubscription(controllerTopic); err != nil {
		logger.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	c.Disconnect(250)
}
