package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
	logger "github.com/cjburchell/uatu-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

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
	Icon                string `json:"icon,omitempty"`
	IconTemplate        string `json:"icon_template,omitempty"`
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

func contains(s []types.PortMode, e types.DeviceMode) bool {
	for _, a := range s {
		if a.DeviceMode == e {
			return true
		}
	}
	return false
}

func GetSensorMode(mode types.PortMode, controllerRepo repo.Controller) (class string, icon string) {
	class = ""
	icon = ""
	switch mode.DeviceMode {
	case types.DeviceModeLights:
		class = "light"
	case types.DeviceModeTimer:
		icon = "mdi:clock"
	case types.DeviceModeDecrease:
		if mode.IsProbe {
			probe, err := controllerRepo.GetProbe(models.GetID(models.ProbeType, mode.Port-1))
			if err == nil {
				switch probe.SensorType {
				case types.SensorTypeTemperature:
					class = "cold"
				case types.SensorTypeAirTemperature:
					class = "cold"
				}
			}
		}
	case types.DeviceModeIncrease:
		if mode.IsProbe {
			probe, err := controllerRepo.GetProbe(models.GetID(models.ProbeType, mode.Port-1))
			if err == nil {
				switch probe.SensorType {
				case types.SensorTypeTemperature:
					class = "heat"
				case types.SensorTypeAirTemperature:
					class = "heat"
				}
			}
		}
	case types.DeviceModeSubstrate:
		if mode.IsProbe {
			probe, err := controllerRepo.GetProbe(models.GetID(models.ProbeType, mode.Port-1))
			if err == nil {
				switch probe.SensorType {
				case types.SensorTypeTemperature:
					class = "heat"
				case types.SensorTypeAirTemperature:
					class = "heat"
				}
			}
		}
	case types.DeviceModeWater:
	case types.DeviceModeCurrentPump:
		icon = "mdi:waves"
	case types.DeviceModeDrainWater:
		icon = "mdi:water-pump"
	case types.DeviceModeProgrammableLogic:
		icon = "mdi:gate-or"
	case types.DeviceModeVariableIllumination:
		class = "light"
	case types.DeviceModeTempPTC:
	case types.DeviceModeDigitalInput:
		icon = "mdi:numeric-10-box-multiple-outline"
	case types.DeviceModeMaintenance:
		icon = "mdi:wrench"
	case types.DeviceModeThunderStorm:
		icon = "mdi:weather-lightning"
	case types.DeviceModeWaterChange:
		icon = "mdi:water-sync"
	case types.DeviceModeFilter:
	case types.DeviceModeProbeAlarm:
		icon = "mdi:bell"
	case types.DeviceModeAlarm:
		icon = "mdi:bell"
	case types.DeviceModeThunder:
		icon = "mdi:weather-lightning"
	}

	return
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
			Name:                fmt.Sprintf("%s Alarm", string(info.Model)),
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
			Name:                fmt.Sprintf("%s Mode", string(info.Model)),
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
			Name:                fmt.Sprintf("%s Feed Pause", string(info.Model)),
			UniqueId:            fmt.Sprintf("%s_feedpause_button", controllerName),
			AvailabilityTopic:   "profiluxmqtt/status",
			PayloadAvailable:    "online",
			PayloadNotAvailable: "offline",
			Icon:                "mdi:shaker",
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
				Name:                fmt.Sprintf("%s Maintenance %s", string(info.Model), p.DisplayName),
				UniqueId:            fmt.Sprintf("%s_%s_button", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				Icon:                "mdi:wrench",
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
				Name:                fmt.Sprintf("%s Reminder Reset %s", string(info.Model), p.Text),
				UniqueId:            fmt.Sprintf("%s_%s_button", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				Icon:                "mdi:restart",
			},
			CommandTopic: fmt.Sprintf("profiluxmqtt/%s/Reminders/%d/reset", controllerName, p.Index),
		}

		msg, _ := json.Marshal(buttonConfig)
		profiMqtt.publishHA(mqttClient, log, "button", controllerName, name, msg, forceUpdate)

		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("%s Reminder Overdue %s", string(info.Model), p.Text),
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         "problem",
				Icon:                "mdi:reminder",
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
		icon := ""
		switch p.SensorType {
		case types.SensorTypeTemperature:
			deviceClass = "temperature"
		case types.SensorTypeAirTemperature:
			deviceClass = "temperature"
		case types.SensorTypeHumidity:
			deviceClass = "humidity"
		case types.SensorTypeVoltage:
			deviceClass = "voltage"
		case types.SensorTypeConductivity:
			icon = "mdi:alpha-c-circle-outline"
		case types.SensorTypeConductivityF:
			icon = "mdi:alpha-c-circle-outline"
		case types.SensorTypeRedox:
			icon = "mdi:thermometer-probe"
		case types.SensorTypeOxygen:
			icon = "mdi:gas-cylinder"
		case types.SensorTypePH:
			icon = "mdi:ph"
		}

		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("%s Probe %s", string(info.Model), p.DisplayName),
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         deviceClass,
				Icon:                icon,
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
				Name:                fmt.Sprintf("%s Water Level %s", string(info.Model), p.DisplayName),
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
			},
			StateTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/state", controllerName, name),
		}
		msgLevel, _ := json.Marshal(stateConfig)
		profiMqtt.publishHA(mqttClient, log, "binary_sensor", controllerName, name+"_State", msgLevel, forceUpdate)

		alarmConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("%s Water Level Alarm %s", string(info.Model), p.DisplayName),
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
				Name:                fmt.Sprintf("%s Water Level Clear %s Alarm", string(info.Model), p.DisplayName),
				UniqueId:            fmt.Sprintf("%s_%s_clear_button", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				Icon:                "mdi:restore",
			},
			CommandTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/clearalarm", controllerName, name),
		}
		msgClearAlarm, _ := json.Marshal(clearAlarmConfig)
		profiMqtt.publishHA(mqttClient, log, "button", controllerName, name, msgClearAlarm, forceUpdate)

		if p.HasTwoInputs {
			stateConfig2 := HaStateConfig{
				HaBaseConfig: HaBaseConfig{
					Device:              device,
					Name:                fmt.Sprintf("%s Water Level 2 %s", string(info.Model), p.DisplayName),
					UniqueId:            fmt.Sprintf("%s_%s_state2", controllerName, name),
					AvailabilityTopic:   "profiluxmqtt/status",
					PayloadAvailable:    "online",
					PayloadNotAvailable: "offline",
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
					Name:                fmt.Sprintf("%s Water Change State %s", string(info.Model), p.DisplayName),
					UniqueId:            fmt.Sprintf("%s_%s_water_change", controllerName, name),
					AvailabilityTopic:   "profiluxmqtt/status",
					PayloadAvailable:    "online",
					PayloadNotAvailable: "offline",
					Icon:                "mdi:water-sync",
				},
				StateTopic: fmt.Sprintf("profiluxmqtt/%s/Level/%s/waterchange", controllerName, name),
			}
			waterChangeMsg, _ := json.Marshal(waterChangeConfig)
			profiMqtt.publishHA(mqttClient, log, "sensor", controllerName, name+"_WaterChange", waterChangeMsg, forceUpdate)

			startWaterChange := HaButtonConfig{
				HaBaseConfig: HaBaseConfig{
					Device:              device,
					Name:                fmt.Sprintf("%s Start Water Change %s", string(info.Model), p.DisplayName),
					UniqueId:            fmt.Sprintf("%s_%s_water_change_button", controllerName, name),
					AvailabilityTopic:   "profiluxmqtt/status",
					PayloadAvailable:    "online",
					PayloadNotAvailable: "offline",
					Icon:                "mdi:water-sync",
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
		class, icon := GetSensorMode(p.Mode, controllerRepo)
		if class == "" {
			class = "power"
		}

		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("%s Socket %s", string(info.Model), p.DisplayName),
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         class,
				Icon:                icon,
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

		class, icon := GetSensorMode(p.Mode, controllerRepo)
		stateConfig := HaStateConfig{
			HaBaseConfig: HaBaseConfig{
				Device:              device,
				Name:                fmt.Sprintf("%s Variable Socket %s", string(info.Model), p.DisplayName),
				UniqueId:            fmt.Sprintf("%s_%s_state", controllerName, name),
				AvailabilityTopic:   "profiluxmqtt/status",
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				DeviceClass:         class,
				Icon:                icon,
			},
			StateTopic:        fmt.Sprintf("profiluxmqtt/%s/LPorts/%s/state", controllerName, name),
			UnitOfMeasurement: "%",
		}
		msgLevel, _ := json.Marshal(stateConfig)
		profiMqtt.publishHA(mqttClient, log, "sensor", controllerName, name, msgLevel, forceUpdate)
	}
}
