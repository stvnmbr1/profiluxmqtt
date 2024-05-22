package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/cjburchell/profiluxmqtt/data/repo"
//	"github.com/cjburchell/profiluxmqtt/profilux/types"
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

func (profiMqtt *ProfiluxMqtt) UpdateMQTT(controllerRepo repo.Controller, mqttClient mqtt.Client, log logger.ILog, forceUpdate bool) {
	info, _ := controllerRepo.GetInfo()
	msg, _ := json.Marshal(info)

/////////// new code - KH Director


	khdControllerName := fmt.Sprintf("KH_Director_%v", info.KHDSerialNumber)

        profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/KHMeasurement", fmt.Sprintf("%v",info.KHDKHMeasurement), forceUpdate)
//        profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/LastMeasurement", fmt.Sprintf("%v",info.KHDLastMeasurement), forceUpdate)
//        profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/MeasurementPerDay", fmt.Sprintf("%v",info.KHDKHMeasurement), forceUpdate)
//        profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/SetValue", fmt.Sprintf("%v",info.KHDKHMeasurement), forceUpdate)
//        profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/UpperLimit", fmt.Sprintf("%v",info.KHDKHMeasurement), forceUpdate)
//        profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/LowerLimit", fmt.Sprintf("%v",info.KHDKHMeasurement), forceUpdate)
//	profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/PH", fmt.Sprintf("%v",info.KHDKHMeasurement), forceUpdate)
//	profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/ADC1", fmt.Sprintf("%v",info.KHDADC1), forceUpdate)
//        profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+"/ADC2", fmt.Sprintf("%v",info.KHDADC2), forceUpdate)



///////////// old code

	controllerName := fmt.Sprintf("%s_%d%s", sanitize(string(info.Model)), info.SerialNumber, suffix)
	profiMqtt.PublishMQTT(mqttClient, log, "status", "online", forceUpdate)
	profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/Controller/data", string(msg), forceUpdate)

///////////

////////// new code - Standalone doser

//        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/IPAddress", fmt.Sprintf("%v",info.SA_IP_ADDRESS), forceUpdate)
//        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/WIFISSID", fmt.Sprintf("%v",info.SA_WIFI_SSID), forceUpdate)
//        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/SoftwareDate", fmt.Sprintf("%v",info.SA_SOFTWAREDATE), forceUpdate)

        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/Temperature", fmt.Sprintf("%v",info.Temperature), forceUpdate)

//////////////

        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_1_name",fmt.Sprintf("%v",info.SA_PUMP1_NAME), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_1_remaining_ml", fmt.Sprintf("%v",info.SA_PUMP1_REMAINING_ML), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_1_remaining_days", fmt.Sprintf("%v",info.SA_PUMP1_REMAINING_DAYS), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_1_daily_dose", fmt.Sprintf("%v",info.SA_PUMP1_DAILY_DOSE), forceUpdate)

        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_2_name",fmt.Sprintf("%v",info.SA_PUMP2_NAME), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_2_remaining_ml", fmt.Sprintf("%v",info.SA_PUMP2_REMAINING_ML), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_2_daily_dose", fmt.Sprintf("%v",info.SA_PUMP2_DAILY_DOSE), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_2_remaining_days", fmt.Sprintf("%v",info.SA_PUMP2_REMAINING_DAYS), forceUpdate)

        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_3_name",fmt.Sprintf("%v",info.SA_PUMP3_NAME), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_3_remaining_ml", fmt.Sprintf("%v",info.SA_PUMP3_REMAINING_ML), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_3_daily_dose", fmt.Sprintf("%v",info.SA_PUMP3_DAILY_DOSE), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_3_remaining_days", fmt.Sprintf("%v",info.SA_PUMP3_REMAINING_DAYS), forceUpdate)

        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_4_name",fmt.Sprintf("%v",info.SA_PUMP4_NAME), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_4_remaining_ml", fmt.Sprintf("%v",info.SA_PUMP4_REMAINING_ML), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_4_daily_dose", fmt.Sprintf("%v",info.SA_PUMP4_DAILY_DOSE), forceUpdate)
        profiMqtt.PublishMQTT(mqttClient, log, controllerName+"/pump_4_remaining_days", fmt.Sprintf("%v",info.SA_PUMP4_REMAINING_DAYS), forceUpdate)


//	for i:=1; i<5; i++{
//
//		profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+fmt.Sprintf("/pump_%v_name",i), fmt.Sprintf("%v",fmt.Sprintf("info.SA_PUMP%v_NAME",i)), forceUpdate)
  //              profiMqtt.PublishMQTT(mqttClient, log, khdControllerName+fmt.Sprintf("/pump_%v_remaining_ml",i), fmt.Sprintf("%v",fmt.Sprintf("info.SA_PUMP%v_REMAINING_ML",i)), forceUpdate)
//
//	}


////////




}
