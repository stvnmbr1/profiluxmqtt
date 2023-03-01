package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
	logger "github.com/cjburchell/uatu-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

func sanitize(text string) string {

	newText := strings.Replace(text, " ", "_", -1)
	newText = strings.Replace(newText, "/", "_", -1)
	newText = strings.Replace(newText, ".", "_", -1)
	newText = strings.Replace(newText, "&", "_", -1)
	return newText
}

type ProfiluxMqtt struct {
	data map[string]string
}

func (profiMqtt *ProfiluxMqtt) PublishMQTTOld(mqttClient mqtt.Client, log logger.ILog, topic string) {
	fullTopic := fmt.Sprintf("profiluxmqtt/%s", topic)
	if profiMqtt.data == nil {
		return
	} else {
		_, ok := profiMqtt.data[fullTopic]
		if !ok {
			return
		}
	}

	t := mqttClient.Publish(fullTopic, 1, false, profiMqtt.data[fullTopic])
	// Handle the token in a go routine so this loop keeps sending messages regardless of delivery status
	go func() {
		_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
		if t.Error() != nil {
			log.Warnf("ERROR PUBLISHING profiluxmqtt/%s", fullTopic)
		}
	}()
}

func (profiMqtt *ProfiluxMqtt) PublishMQTT(mqttClient mqtt.Client, log logger.ILog, topic string, payload string, forceUpdate bool) {
	fullTopic := fmt.Sprintf("profiluxmqtt/%s", topic)
	if profiMqtt.data == nil {
		profiMqtt.data = make(map[string]string)
	} else {
		if profiMqtt.data[fullTopic] == payload && !forceUpdate {
			return
		}
	}
	profiMqtt.data[fullTopic] = payload

	t := mqttClient.Publish(fullTopic, 1, false, payload)
	// Handle the token in a go routine so this loop keeps sending messages regardless of delivery status
	go func() {
		_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
		if t.Error() != nil {
			log.Warnf("ERROR PUBLISHING profiluxmqtt/%s", fullTopic)
		}
	}()
}

func (profiMqtt *ProfiluxMqtt) publishHA(mqttClient mqtt.Client, log logger.ILog, platform string, device string, topic string, payload []byte, forceUpdate bool) {
	fullTopic := fmt.Sprintf("homeassistant/%s/%s/%s/config", platform, device, topic)
	if profiMqtt.data == nil {
		profiMqtt.data = make(map[string]string)
	} else {
		if profiMqtt.data[fullTopic] == string(payload) && !forceUpdate {
			return
		}
	}
	profiMqtt.data[fullTopic] = string(payload)

	t := mqttClient.Publish(fullTopic, 1, false, payload)
	// Handle the token in a go routine so this loop keeps sending messages regardless of delivery status
	go func() {
		_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
		if t.Error() != nil {
			log.Warnf("ERROR PUBLISHING %s", fullTopic)
		}
	}()
}

type Device struct {
	Identifiers  string `json:"identifiers"`
	Name         string `json:"name"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
	Version      string `json:"hw_version"`
}

type HaBaseConfig struct {
	Device              Device `json:"device"`
	Name                string `json:"name"`
	UniqueId            string `json:"unique_id"`
	AvailabilityTopic   string `json:"availability_topic,omitempty"`
	DeviceClass         string `json:"device_class,omitempty"`
	PayloadAvailable    string `json:"payload_available"`
	PayloadNotAvailable string `json:"payload_not_available"`
	Icon                string `json:"icon, omitempty"`
}

type HaStateConfig struct {
	HaBaseConfig
	StateTopic        string `json:"state_topic"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
}

type HaButtonConfig struct {
	HaBaseConfig
	CommandTopic string `json:"command_topic"`
}

type HaSwitchConfig struct {
	HaBaseConfig
	StateTopic   string `json:"state_topic"`
	CommandTopic string `json:"command_topic"`
}

func (profiMqtt *ProfiluxMqtt) UpdateHomeAssistant(controllerRepo repo.Controller, mqttClient mqtt.Client, log logger.ILog, forceUpdate bool) {
	info, _ := controllerRepo.GetInfo()
	controllerName := fmt.Sprintf("%s_%d", sanitize(string(info.Model)), info.DeviceAddress)
	device := Device{
		Identifiers:  controllerName,
		Version:      fmt.Sprintf("%.2f", info.SoftwareVersion),
		Name:         string(info.Model),
		Model:        string(info.Model),
		Manufacturer: "GHL",
	}

	config := HaStateConfig{
		HaBaseConfig: HaBaseConfig{
			Device:              device,
			Name:                fmt.Sprintf("Alarm"),
			UniqueId:            strings.ToLower(fmt.Sprintf("%s_alarm", controllerName)),
			AvailabilityTopic:   "profiluxmqtt/status",
			PayloadAvailable:    "online",
			PayloadNotAvailable: "offline",
			DeviceClass:         "problem",
		},
		StateTopic: fmt.Sprintf("profiluxmqtt/%s/Controller/alarm", controllerName),
	}

	msg, _ := json.Marshal(config)
	profiMqtt.publishHA(mqttClient, log, "binary_sensor", controllerName, "Alarm", msg, forceUpdate)

	modeConfig := HaStateConfig{
		HaBaseConfig: HaBaseConfig{
			Device:              device,
			Name:                fmt.Sprintf("Mode"),
			UniqueId:            strings.ToLower(fmt.Sprintf("%s_alarm", controllerName)),
			AvailabilityTopic:   "profiluxmqtt/status",
			PayloadAvailable:    "online",
			PayloadNotAvailable: "offline",
		},
		StateTopic: fmt.Sprintf("profiluxmqtt/%s/Controller/mode", controllerName),
	}

	modeMsg, _ := json.Marshal(modeConfig)
	profiMqtt.publishHA(mqttClient, log, "sensor", controllerName, "Mode", modeMsg, forceUpdate)

	feedButtonConfig := HaButtonConfig{
		HaBaseConfig: HaBaseConfig{
			Device:              device,
			Name:                "Feed Pause",
			UniqueId:            fmt.Sprintf("%s_feedpause_button", controllerName),
			AvailabilityTopic:   "profiluxmqtt/status",
			PayloadAvailable:    "online",
			PayloadNotAvailable: "offline",
		},
		CommandTopic: fmt.Sprintf("profiluxmqtt/%s/Controller/feedpause", controllerName),
	}

	feedButtonMsg, _ := json.Marshal(feedButtonConfig)
	profiMqtt.publishHA(mqttClient, log, "button", controllerName, "FeedPause", feedButtonMsg, forceUpdate)

	for _, p := range info.Maintenance {
		name := fmt.Sprintf("Maintenance%d", p.Index)

		buttonConfig := HaSwitchConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                p.DisplayName,
				UniqueId:            fmt.Sprintf("%s_%s_button", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
			},
			CommandTopic: fmt.Sprintf("profiluxmqtt/%s/Maintenance/%d/command", controllerName, p.Index),
			StateTopic:   fmt.Sprintf("profiluxmqtt/%s/Maintenance/%d/state", controllerName, p.Index),
		}
		msg, _ := json.Marshal(buttonConfig)
		profiMqtt.publishHA(mqttClient, log, "switch", controllerName, name, msg, forceUpdate)
	}

	for _, p := range info.Reminders {
		name := fmt.Sprintf("Reminder%d", p.Index)

		buttonConfig := HaButtonConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("Reset %s", p.Text),
				UniqueId:            fmt.Sprintf("%s_%s_button", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
			},
			CommandTopic: fmt.Sprintf("profiluxmqtt/%s/Reminders/%d/reset", controllerName, p.Index),
		}

		msg, _ := json.Marshal(buttonConfig)
		profiMqtt.publishHA(mqttClient, log, "button", controllerName, name, msg, forceUpdate)

		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("%s Overdue", p.Text),
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         "problem",
			},
			StateTopic: fmt.Sprintf("profiluxmqtt/%s/Reminders/%d/state", controllerName, p.Index),
		}
		msg, _ = json.Marshal(stateConfig)
		profiMqtt.publishHA(mqttClient, log, "binary_sensor", controllerName, name, msg, forceUpdate)
	}

	probes, _ := controllerRepo.GetProbes()
	for _, p := range probes {
		name := p.ID

		deviceClass := ""
		switch p.SensorType {
		case types.SensorTypeTemperature:
			deviceClass = "temperature"
			break
		case types.SensorTypeHumidity:
			deviceClass = "temperature"
			break
		case types.SensorTypeAirTemperature:
			deviceClass = "temperature"
			break
		case types.SensorTypeVoltage:
			deviceClass = "voltage"
			break
		default:
			deviceClass = ""
		}

		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                p.DisplayName,
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         deviceClass,
			},
			StateTopic:        fmt.Sprintf("profiluxmqtt/%s/Probes/%s/convertedvalue", controllerName, name),
			UnitOfMeasurement: p.Units,
		}
		msg, _ := json.Marshal(stateConfig)
		profiMqtt.publishHA(mqttClient, log, "sensor", controllerName, name, msg, forceUpdate)
	}

	levelSensors, _ := controllerRepo.GetLevelSensors()
	for _, p := range levelSensors {
		name := p.ID
		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                p.DisplayName,
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         "moisture",
			},
			StateTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/state", controllerName, name),
		}
		msgLevel, _ := json.Marshal(stateConfig)
		profiMqtt.publishHA(mqttClient, log, "binary_sensor", controllerName, name+"_State", msgLevel, forceUpdate)

		alarmConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("%s Alarm", p.DisplayName),
				UniqueId:            fmt.Sprintf("%s_%s_alarm", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         "problem",
			},
			StateTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/alarm", controllerName, name),
		}
		msgAlarm, _ := json.Marshal(alarmConfig)
		profiMqtt.publishHA(mqttClient, log, "binary_sensor", controllerName, name+"_Alarm", msgAlarm, forceUpdate)

		clearAlarmConfig := HaButtonConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("%s Clear Alarm", p.DisplayName),
				UniqueId:            fmt.Sprintf("%s_%s_clear_button", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
			},
			CommandTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/clearalarm", controllerName, name),
		}
		msgClearAlarm, _ := json.Marshal(clearAlarmConfig)
		profiMqtt.publishHA(mqttClient, log, "button", controllerName, name, msgClearAlarm, forceUpdate)

		if p.HasTwoInputs {
			stateConfig2 := HaStateConfig{
				HaBaseConfig: HaBaseConfig{
					Device:              device,
					Name:                fmt.Sprintf("%s Sensor 2", p.DisplayName),
					UniqueId:            fmt.Sprintf("%s_%s_state2", controllerName, name),
					AvailabilityTopic:   "profiluxmqtt/status",
					PayloadAvailable:    "online",
					PayloadNotAvailable: "offline",
					DeviceClass:         "moisture",
				},
				StateTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/state2", controllerName, name),
			}
			config2Msg, _ := json.Marshal(stateConfig2)
			profiMqtt.publishHA(mqttClient, log, "binary_sensor", controllerName, name+"_State2", config2Msg, forceUpdate)
		}

		if p.HasWaterChange {
			waterChangeConfig := HaStateConfig{
				HaBaseConfig: HaBaseConfig{
					Device:              device,
					Name:                fmt.Sprintf("%s Water Change", p.DisplayName),
					UniqueId:            fmt.Sprintf("%s_%s_water_change", controllerName, name),
					AvailabilityTopic:   "profiluxmqtt/status",
					PayloadAvailable:    "online",
					PayloadNotAvailable: "offline",
					DeviceClass:         "moisture",
				},
				StateTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/waterchange", controllerName, name),
			}
			waterChangeMsg, _ := json.Marshal(waterChangeConfig)
			profiMqtt.publishHA(mqttClient, log, "sensor", controllerName, name+"_WaterChange", waterChangeMsg, forceUpdate)

			startWaterChange := HaButtonConfig{
				HaBaseConfig: HaBaseConfig{
					Device:              device,
					Name:                fmt.Sprintf("%s Start Water Change", p.DisplayName),
					UniqueId:            fmt.Sprintf("%s_%s_water_change_button", controllerName, name),
					AvailabilityTopic:   "profiluxmqtt/status",
					PayloadAvailable:    "online",
					PayloadNotAvailable: "offline",
				},
				CommandTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/startwaterchange", controllerName, name),
			}
			msgStartWaterChange, _ := json.Marshal(startWaterChange)
			profiMqtt.publishHA(mqttClient, log, "button", controllerName, name+"_StartWaterChange", msgStartWaterChange, forceUpdate)
		}
	}

	sockets, _ := controllerRepo.GetSPorts()
	for _, p := range sockets {
		name := p.ID
		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                p.DisplayName,
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         "power",
			},
			StateTopic: fmt.Sprintf("profiluxmqtt/%s/SPorts/%s/state", controllerName, name),
		}
		msgLevel, _ := json.Marshal(stateConfig)
		profiMqtt.publishHA(mqttClient, log, "binary_sensor", controllerName, name, msgLevel, forceUpdate)
	}

	lightPorts, _ := controllerRepo.GetLPorts()
	for _, p := range lightPorts {
		name := p.ID
		if p.DisplayName == "" {
			p.DisplayName = p.ID
		}

		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                p.DisplayName,
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
			},
			StateTopic:        fmt.Sprintf("profiluxmqtt/%s/LPorts/%s/state", controllerName, name),
			UnitOfMeasurement: "%",
		}
		msgLevel, _ := json.Marshal(stateConfig)
		profiMqtt.publishHA(mqttClient, log, "sensor", controllerName, name, msgLevel, forceUpdate)
	}
}

func (profiMqtt *ProfiluxMqtt) UpdateMQTT(controllerRepo repo.Controller, mqttClient mqtt.Client, log logger.ILog, forceUpdate bool) {
	info, _ := controllerRepo.GetInfo()
	msg, _ := json.Marshal(info)
	controllerName := sanitize(string(info.Model)) + "_" + fmt.Sprintf("%d", info.DeviceAddress)
	profiMqtt.PublishMQTT(mqttClient, log, "status", "online", forceUpdate)
	profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/Controller/data", string(msg), forceUpdate)
	profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/Controller/alarm", string(info.Alarm), forceUpdate)
	profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/Controller/mode", string(info.OperationMode), forceUpdate)

	for _, p := range info.Maintenance {
		path := fmt.Sprintf("%s/Maintenance/%d", controllerName, p.Index)
		data, _ := json.Marshal(p)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/data", string(data), forceUpdate)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/state", strings.ToUpper(string(p.IsActive)), forceUpdate)
	}

	for _, p := range info.Reminders {
		path := fmt.Sprintf("%s/Reminders/%d", controllerName, p.Index)
		data, _ := json.Marshal(p)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/data", string(data), forceUpdate)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/state", string(p.IsOverdue), forceUpdate)
	}

	probes, _ := controllerRepo.GetProbes()
	for _, p := range probes {
		path := fmt.Sprintf("%s/Probes/%s", controllerName, p.ID)
		if p.SensorType == types.SensorTypeAirTemperature {
			if p.Value > 35 || p.Value < -30 {
				if forceUpdate {
					profiMqtt.PublishMQTTOld(mqttClient, log, path+"/data")
					profiMqtt.PublishMQTTOld(mqttClient, log, path+"/state")
					profiMqtt.PublishMQTTOld(mqttClient, log, path+"/convertedvalue")
				}
				continue
			}
		}
		data, _ := json.Marshal(p)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/data", string(data), forceUpdate)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/state", fmt.Sprintf("%.2f", p.Value), forceUpdate)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/convertedvalue", fmt.Sprintf("%.2f", p.ConvertedValue), forceUpdate)
	}

	levelSensors, _ := controllerRepo.GetLevelSensors()
	for _, p := range levelSensors {
		path := fmt.Sprintf("%s/Level/%s", controllerName, p.ID)
		data, _ := json.Marshal(p)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/data", string(data), forceUpdate)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/state", string(p.Value), forceUpdate)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/alarm", string(p.AlarmState), forceUpdate)
		if p.HasTwoInputs {
			profiMqtt.PublishMQTT(mqttClient, log, path+"/state2", string(p.SecondSensor), forceUpdate)
		}

		if p.HasWaterChange {
			profiMqtt.PublishMQTT(mqttClient, log, path+"/waterchange", string(p.WaterMode), forceUpdate)
		}
	}

	sockets, _ := controllerRepo.GetSPorts()
	for _, p := range sockets {
		path := fmt.Sprintf("%s/SPorts/%s", controllerName, p.ID)
		data, _ := json.Marshal(p)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/data", string(data), forceUpdate)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/state", string(p.Value), forceUpdate)
	}

	lightPorts, _ := controllerRepo.GetLPorts()
	for _, p := range lightPorts {
		path := fmt.Sprintf("%s/LPorts/%s", controllerName, p.ID)
		data, _ := json.Marshal(p)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/data", string(data), forceUpdate)
		profiMqtt.PublishMQTT(mqttClient, log, path+"/state", fmt.Sprintf("%.2f", p.Value), forceUpdate)
	}
}
