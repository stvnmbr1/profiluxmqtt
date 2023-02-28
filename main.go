package main

import (
	"fmt"
	"github.com/cjburchell/profiluxmqtt/commands"
	"github.com/cjburchell/profiluxmqtt/update"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	appSettings "github.com/cjburchell/profiluxmqtt/settings"
	"github.com/cjburchell/tools-go/env"

	"github.com/cjburchell/profiluxmqtt/data/repo"

	"github.com/cjburchell/settings-go"
	logger "github.com/cjburchell/uatu-go"
	"github.com/eclipse/paho.mqtt.golang"

	profiluxmqtt "github.com/cjburchell/profiluxmqtt/mqtt"
)

func main() {
	log := logger.Create(logger.Settings{
		MinLogLevel:  logger.INFO,
		ServiceName:  "Profilux MQTT",
		LogToConsole: true,
		UseHTTP:      false,
		UsePubSub:    false,
	})

	log.Printf("Starting Up!")
	set := settings.Get(env.Get("ConfigFile", ""))
	appConfig := appSettings.Get(set)

	controllerRepo := repo.NewController()

	mqttOptions := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("mqtt://%s:%d", appConfig.MqttHost, appConfig.MqttPort)).SetClientID("profilux-mqtt")
	mqttOptions.SetOrderMatters(false)       // Allow out of order messages (use this option unless in order delivery is essential)
	mqttOptions.ConnectTimeout = time.Second // Minimal delays on connect
	mqttOptions.WriteTimeout = time.Second   // Minimal delays on writes
	mqttOptions.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	mqttOptions.PingTimeout = time.Second    // local broker so response should be quick
	mqttOptions.ConnectRetry = true
	mqttOptions.AutoReconnect = true

	mqttOptions.OnConnectionLost = func(cl mqtt.Client, err error) {
		log.Warn("connection lost")
	}
	mqttOptions.OnConnect = func(mqtt.Client) {
		log.Print("connection established")
	}
	mqttOptions.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		log.Print("attempting to reconnect")
	}

	mqttClient := mqtt.NewClient(mqttOptions)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Print("Connection is up")

	defer mqttClient.Disconnect(100)

	mqttClient.SubscribeMultiple(map[string]byte{
		"profiluxmqtt/+/Maintenance/+/command": 1}, func(_ mqtt.Client, message mqtt.Message) {
		tokens := strings.Split(message.Topic(), "/")
		index, _ := strconv.Atoi(tokens[3])
		log.Printf("Maintenance: %d", index)
		commands.Maintenance(index, string(message.Payload()) == "ON", controllerRepo, log, appConfig.Connection)
	})

	mqttClient.SubscribeMultiple(map[string]byte{
		"profiluxmqtt/+/Reminders/+/reset": 1}, func(_ mqtt.Client, message mqtt.Message) {
		tokens := strings.Split(message.Topic(), "/")
		index, _ := strconv.Atoi(tokens[3])
		log.Printf("Reset Reminders: %d", index)
		commands.ResetReminder(index, controllerRepo, log, appConfig.Connection)
	})

	mqttClient.SubscribeMultiple(map[string]byte{
		"profiluxmqtt/+/Level/+/clearalarm": 1}, func(_ mqtt.Client, message mqtt.Message) {
		tokens := strings.Split(message.Topic(), "/")
		log.Printf("Clear Alarm: %s", tokens[3])
		commands.ClearLevelAlarm(tokens[3], controllerRepo, log, appConfig.Connection)
	})

	mqttClient.SubscribeMultiple(map[string]byte{
		"profiluxmqtt/+/Level/+/startwaterchange": 1}, func(_ mqtt.Client, message mqtt.Message) {
		tokens := strings.Split(message.Topic(), "/")
		log.Printf("Water Change: %s", tokens[3])
		commands.WaterChange(tokens[3], controllerRepo, log, appConfig.Connection)
	})

	mqttClient.SubscribeMultiple(map[string]byte{
		"profiluxmqtt/+/Controller/feedpause": 1}, func(_ mqtt.Client, message mqtt.Message) {
		log.Printf("Start Feed Pause")
		commands.FeedPause(true, controllerRepo, log, appConfig.Connection)
	})

	var profiluxMqtt profiluxmqtt.ProfiluxMqtt

	log.Debugf("Getting Data from Controller")
	for {
		var err = update.All(controllerRepo, log, appConfig.Connection)
		if err == nil {
			break
		}

		log.Error(err, "Unable to do first update")
		log.Debugf("RefreshSettings Sleeping for %s", logRate.String())
		<-time.After(logRate)
		continue
	}

	profiluxMqtt.UpdateMQTT(controllerRepo, mqttClient, log, false)
	profiluxMqtt.UpdateHomeAssistant(controllerRepo, mqttClient, log, false)

	go RunUpdateConfig(controllerRepo, mqttClient, log, appConfig, &profiluxMqtt)
	RunUpdate(controllerRepo, mqttClient, log, appConfig, &profiluxMqtt)
}

const logRate = time.Second * 1
const logAllRate = time.Minute * 1

func RunUpdateConfig(controller repo.Controller, mqttClient mqtt.Client, log logger.ILog, config appSettings.Config, profiluxMqtt *profiluxmqtt.ProfiluxMqtt) {
	c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-c:
			profiluxMqtt.PublishMQTT(mqttClient, log, "status", "offline", true)
			log.Debug("Exit Application")
			return
		case <-time.After(logAllRate):
			log.Print("Updating Config")
			var err = update.All(controller, log, config.Connection)
			if err != nil {
				log.Errorf(err, "Unable to update")
			} else {
				profiluxMqtt.UpdateMQTT(controller, mqttClient, log, true)
				profiluxMqtt.UpdateHomeAssistant(controller, mqttClient, log, true)
			}
		}
	}
}

func RunUpdate(controller repo.Controller, mqttClient mqtt.Client, log logger.ILog, config appSettings.Config, profiluxMqtt *profiluxmqtt.ProfiluxMqtt) {
	c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-c:
			log.Debug("Exit Application")
			return
		case <-time.After(logRate):
			log.Print("Updating State")
			var err = update.State(controller, log, config.Connection)
			if err != nil {
				log.Errorf(err, "Unable to update state")
			} else {
				profiluxMqtt.UpdateMQTT(controller, mqttClient, log, false)
			}
		}
	}
}
