package main

import (
	"fmt"
	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/snub/senter"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	cmdTopic string = "/controller/%s/cmd"

	cmdIp       string = "ip"
	cmdInterval string = "interval"
	cmdNtp      string = "ntp"
)

var (
	topicsAndHandlers map[string]mqtt.MessageHandler = map[string]mqtt.MessageHandler{
		"/controller/+/hello":             helloHandler,
		"/controller/+/ip":                setupHandler,
		"/controller/+/interval":          setupHandler,
		"/controller/+/ntp":               setupHandler,
		"/controller/+/sensors/discovery": discoveryHandler,
		"/controller/+/sensor/+/tmp":      tmpHandler,
	}

	regexHello     *regexp.Regexp = regexp.MustCompile("^/controller/(.+)/hello$")
	regexSetup     *regexp.Regexp = regexp.MustCompile("^/controller/(.+)/(.+)$")
	reqexDiscovery *regexp.Regexp = regexp.MustCompile("^/controller/(.+)/sensors/discovery$")
	reqexTmp       *regexp.Regexp = regexp.MustCompile("^/controller/(.+)/sensor/(.+)/tmp$")

	getCmds map[string]int = map[string]int{
		cmdIp:       1,
		cmdInterval: 3,
		cmdNtp:      5,
	}

	setCmds map[string]int = map[string]int{
		cmdIp:       2,
		cmdInterval: 4,
		cmdNtp:      6,
	}

	clockSkew int64 = 30
)

func helloHandler(client *mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	message := msg.Payload()
	logger.Printf("helloHandler, topic: %s, message: %s\n", topic, message)

	var macAddress string = ""

	match := regexHello.FindStringSubmatch(topic)
	if len(match) == 2 {
		logger.Println("processing hello topic...")
		macAddress = match[1]
		timestamp, err := strconv.Atoi(string(message))
		if err != nil {
			logger.Printf("unable to convert timestamp: %s\n", err)
		}
		logger.Printf("mac address: %s, timestamp: %d\n", macAddress, timestamp)
		controller := senter.LoadControllerByMacAddress(macAddress)
		controller.SetLastStartup(int64(timestamp))
		controller.Save()
		logger.Printf("saved controller: %v\n", controller)
	}

	if macAddress != "" {
		pubTopic := fmt.Sprintf(cmdTopic, macAddress)
		for name, value := range getCmds {
			logger.Printf("publishing command: %s (%d) to %s\n", name, value, pubTopic)
			client.Publish(pubTopic, byte(0), false, []byte(strconv.Itoa(value)))
		}
	}
}

// TODO: check what happens if some object is not found
func setupHandler(client *mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	message := msg.Payload()
	logger.Printf("setupHandler, topic: %s, message: %s\n", topic, message)

	match := regexSetup.FindStringSubmatch(topic)
	logger.Printf("match: %v\n", match)
	if len(match) == 3 {
		logger.Println("processing setup topic...")
		controllerMacAddress := match[1]
		cmd := match[2]

		logger.Printf("mac address: %s, command: %s\n", controllerMacAddress, cmd)

		controller := senter.LoadControllerByMacAddress(controllerMacAddress)
		controllerConfig := senter.LoadControllerConfigByControllerId(controller.Id)

		logger.Printf("loaded controller %s config: %v\n", controllerMacAddress, controllerConfig)

		_, ok := setCmds[cmd]
		if ok {
			pubTopic := fmt.Sprintf(cmdTopic, controllerMacAddress)
			switch cmd {
			case cmdIp:
				if !controllerConfig.IpAddress.Valid {
					controllerConfig.IpAddress.String = string(message)
					controllerConfig.IpAddress.Valid = true
					controllerConfig.Save()
					return
				}
				if controllerConfig.IpAddress.String != string(message) {
					cmdValue := fmt.Sprintf("%d,%s", setCmds[cmd], controllerConfig.IpAddress.String)
					logger.Printf("publishing setup command: %s with args %s to %s\n", cmd, cmdValue, pubTopic)
					client.Publish(pubTopic, byte(0), false, []byte(cmdValue))
				}
			case cmdInterval:
				currentInterval, err := strconv.Atoi(string(message))
				if err == nil {
					if !controllerConfig.UpdateInterval.Valid {
						controllerConfig.UpdateInterval.Int64 = int64(currentInterval)
						controllerConfig.UpdateInterval.Valid = true
						controllerConfig.Save()
						return
					}
					if controllerConfig.UpdateInterval.Int64 != int64(currentInterval) {
						cmdValue := fmt.Sprintf("%d,%d", setCmds[cmd], controllerConfig.UpdateInterval.Int64)
						logger.Printf("publishing setup command: %s with args %s to %s\n", cmd, cmdValue, pubTopic)
						client.Publish(pubTopic, byte(0), false, []byte(cmdValue))
					}
				}
			case cmdNtp:
				if !controllerConfig.NtpIpAddress.Valid {
					controllerConfig.NtpIpAddress.String = string(message)
					controllerConfig.NtpIpAddress.Valid = true
					controllerConfig.Save()
					return
				}
				if controllerConfig.NtpIpAddress.String != string(message) {
					cmdValue := fmt.Sprintf("%d,%s", setCmds[cmd], controllerConfig.NtpIpAddress.String)
					logger.Printf("publishing setup command: %s with args %s to %s\n", cmd, cmdValue, pubTopic)
					client.Publish(pubTopic, byte(0), false, []byte(cmdValue))
				}
			default:
				logger.Printf("unhandled command: %s", cmd)
			}
		}

	}
}

// TODO: what happens if sensor is connected to new controller
func discoveryHandler(client *mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	message := msg.Payload()
	logger.Printf("discoveryHandler, topic: %s, message: %s\n", topic, message)

	match := reqexDiscovery.FindStringSubmatch(topic)
	if len(match) == 2 {
		logger.Println("processing discovery topic...")
		controllerMacAddress := match[1]
		deviceAddress := string(message)
		logger.Printf("controller mac address: %s, sensor device address: %s", controllerMacAddress, deviceAddress)

		controller := senter.LoadControllerByMacAddress(controllerMacAddress)
		if controller.New() {
			logger.Printf("unsaved contoller detected, ignoring discovered sensor")
			return
		}

		sensor := senter.LoadSensorByDeviceAddress(deviceAddress)
		if sensor.New() {
			logger.Printf("discovered new sensor with device address: %s\n", deviceAddress)
			logger.Println("saving newly discovered sensor")
			sensor.SetControllerId(controller.Id)
			sensor.Create()
			logger.Printf("created sensor: %v\n", sensor)
		} else {
			logger.Printf("existing sensor: %v\n", sensor)
		}
	}
}

func tmpHandler(client *mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	message := msg.Payload()
	logger.Printf("tmpHandler, topic: %s, message: %s\n", topic, message)

	match := reqexTmp.FindStringSubmatch(topic)
	//logger.Printf("temperature topic match: %v\n", match)
	if len(match) == 3 {
		logger.Println("processing temperature topic...")
		controllerMacAddress := match[1]
		sensorDeviceAddress := match[2]
		splitMsg := strings.Split(string(message), ",")
		if len(splitMsg) == 2 {
			timestamp, err := strconv.ParseInt(splitMsg[0], 10, 0)
			if err != nil {
				logger.Printf("unable to convert timestamp: %s\n", err)
			}

			timeDiff := time.Now().UTC().Unix() - timestamp
			if timeDiff < 0 {
				timeDiff = -timeDiff
			}
			if timeDiff > clockSkew {
				logger.Printf("controller and center time difference is more than %d seconds\n", clockSkew)
				logger.Printf("sensor reading timestamp: %s\n", time.Unix(timestamp, 0).UTC())
				logger.Println("discarding...")
				return
			}

			value, err := strconv.ParseFloat(splitMsg[1], 32)
			if err != nil {
				logger.Printf("unable to parse temperature: %s\n", err)
			}
			logger.Printf("controller mac address: %s, sensor device address: %s, timestamp: %d, value: %f\n", controllerMacAddress, sensorDeviceAddress, timestamp, value)

			controller := senter.LoadControllerByMacAddress(controllerMacAddress)
			if controller.New() {
				logger.Printf("unsaved contoller detected, ignoring sensor temperature")
				return
			}

			sensor := senter.LoadSensorByDeviceAddress(sensorDeviceAddress)
			if sensor.New() {
				logger.Printf("detected unsaved sensor with device address: %s\n", sensorDeviceAddress)
				logger.Println("saving unsaved sensor")
				sensor.SetControllerId(controller.Id)
				sensor.Create()
				logger.Printf("created sensor: %v\n", sensor)
			}
			temperature := senter.NewTemperature(sensor, int64(timestamp), float32(value))
			temperature.Create()
			logger.Printf("created temperature: %v\n", temperature)
		}
	}

}
