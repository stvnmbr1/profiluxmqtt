package settings

import (
	"github.com/cjburchell/profiluxmqtt/profilux"
	"github.com/cjburchell/settings-go"
)

func Get(settings settings.ISettings) Config {
	return Config{
		Connection: newConnectionSettings(settings),
		MqttHost:   settings.Get("MQTT_HOST", "localhost"),
		MqttPort:   settings.GetInt("MQTT_PORT", 1883),
	}
}

type Config struct {
	Connection profilux.Settings
	MqttHost   string
	MqttPort   int
}

func newConnectionSettings(settings settings.ISettings) (connection profilux.Settings) {
	connection.Address = settings.Get("PROFILUX_HOST", "192.168.3.10")
	connection.Port = settings.GetInt("PROFILUX_PORT", 80)
	connection.Protocol = profilux.ProtocolHTTP
	connection.ControllerAddress = 1
	return
}
