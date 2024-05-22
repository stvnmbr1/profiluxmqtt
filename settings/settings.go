package settings

import (
	"github.com/cjburchell/profiluxmqtt/profilux"
	"github.com/cjburchell/settings-go"
)

func Get(settings settings.ISettings) Config {
	return Config{
		Connection: newConnectionSettings(settings),
		MqttHost:   settings.Get("MQTT_HOST", "192.168.0.160"),
		MqttPort:   settings.GetInt("MQTT_PORT", 1883),
		ClientID:   settings.Get("CLIENT_ID", "profilux-mqtt"),
	}
}

type Config struct {
	Connection profilux.Settings
	MqttHost   string
	MqttPort   int
	ClientID   string
}

func newConnectionSettings(settings settings.ISettings) (connection profilux.Settings) {
	connection.Address = settings.Get("PROFILUX_HOST", "192.168.0.30")
	connection.Port = settings.GetInt("PROFILUX_PORT", 10001)
	connection.Protocol = profilux.ProtocolSocket
	connection.ControllerAddress = 0
	return
}
