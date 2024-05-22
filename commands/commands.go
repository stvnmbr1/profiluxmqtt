package commands

import (
	"github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
	"github.com/cjburchell/profiluxmqtt/update"
	logger "github.com/cjburchell/uatu-go"
)

func FeedPause(enabled bool, repo repo.Controller, log logger.ILog, config profilux.Settings) {
	profiluxController, err := profilux.NewController(config)
	if err != nil {
		log.Errorf(err, "unable to connect")
		return
	}

	defer profiluxController.Disconnect()

	err = profiluxController.FeedPause(0, enabled)
	if err != nil {
		log.Errorf(err, "Unable to send feed pause")
		return
	}

	err = update.InfoState(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable update state")
		return
	}
}

func ResetReminder(index int, repo repo.Controller, log logger.ILog, config profilux.Settings) {

	var reminder *models.Reminder

	info, err := repo.GetInfo()
	if err != nil {
		log.Errorf(err, "Unable to get info")
		return
	}

	for _, item := range info.Reminders {
		if item.Index == index {
			reminder = &item
			break
		}
	}

	if reminder == nil {
		log.Warnf("unable to find reminder")
		return
	}

	profiluxController, err := profilux.NewController(config)
	if err != nil {
		log.Errorf(err, "unable to connect")
		return
	}

	defer profiluxController.Disconnect()

	if reminder.IsRepeating {
		err = profiluxController.ResetReminder(index, reminder.Period)
		if err != nil {
			log.Errorf(err, "Unable Reset Reminder")
			return
		}

	} else {
		err = profiluxController.ClearReminder(index)
		if err != nil {
			log.Errorf(err, "Unable Clear Reminder")
			return
		}
	}

	err = update.InfoState(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable update state")
		return
	}
}
func Maintenance(index int, enable bool, repo repo.Controller, log logger.ILog, config profilux.Settings) {
	profiluxController, err := profilux.NewController(config)
	if err != nil {
		log.Errorf(err, "unable to connect")
		return
	}

	defer profiluxController.Disconnect()

	var maintenance *models.Maintenance
	info, err := repo.GetInfo()
	if err != nil {
		log.Errorf(err, "Unable to get info")
		return
	}

	for _, item := range info.Maintenance {
		if item.Index == index {
			maintenance = &item
			break
		}
	}

	if maintenance == nil {
		log.Warnf("unable to find reminder")
		return
	}

	err = profiluxController.Maintenance(enable, index)
	if err != nil {
		log.Errorf(err, "Unable to set Maintenance")
		return
	}

	err = update.InfoState(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable update state")
		return
	}
}
func ClearLevelAlarm(id string, repo repo.Controller, log logger.ILog, config profilux.Settings) {

	var sensor *models.LevelSensor
	items, err := repo.GetLevelSensors()
	if err != nil {
		log.Warnf("unable to find level sensor %s", id)
		return
	}

	for _, level := range items {
		if level.ID == id {
			sensor = &level
			break
		}
	}

	if sensor == nil {
		log.Warnf("unable to find level sensor %s", id)
		return
	}

	profiluxController, err := profilux.NewController(config)
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	err = profiluxController.ClearLevelAlarm(sensor.Index)
	if err != nil {
		log.Errorf(err, "Unable Clear Level Alarm")
		return
	}

	err = update.LevelSensors(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable to update Level Sensors")
		return
	}

	err = update.InfoState(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable update state")
		return
	}
}
func WaterChange(id string, repo repo.Controller, log logger.ILog, config profilux.Settings) {
	var sensor *models.LevelSensor

	sensors, err := repo.GetLevelSensors()
	if err != nil {
		log.Errorf(err, "Unable to get level sensors")
		return
	}

	for _, level := range sensors {
		if level.ID == id {
			sensor = &level
			break
		}
	}

	if sensor == nil {
		log.Warnf("unable to find level snsoro %s", id)
		return
	}

	profiluxController, err := profilux.NewController(config)
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	err = profiluxController.WaterChange(sensor.Index)
	if err != nil {
		log.Errorf(err, "Unable to do Water Change")
		return
	}

	err = update.LevelSensors(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable to update Level Sensors")
		return
	}

	err = update.InfoState(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable update state")
		return
	}
}

func ManualSockets(enable bool, repo repo.Controller, log logger.ILog, config profilux.Settings) {
	profiluxController, err := profilux.NewController(config)
	if err != nil {
		log.Errorf(err, "unable to connect")
		return
	}
	defer profiluxController.Disconnect()

	mode := types.OperationModeNormal
	if enable {
		mode = types.OperationModeManualSockets
	}

	err = profiluxController.SetOperationMode(mode)
	if err != nil {
		log.Errorf(err, "Unable to send feed pause")
		return
	}

	err = update.InfoState(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable update state")
		return
	}
}

func SetSocketState(id string, enable bool, repo repo.Controller, log logger.ILog, config profilux.Settings) {
	profiluxController, err := profilux.NewController(config)
	if err != nil {
		log.Errorf(err, "unable to connect")
		return
	}

	defer profiluxController.Disconnect()

	port, err := repo.GetSPort(id)
	if err != nil {
		log.Errorf(err, "Unable find socket with id of %s", id)
		return
	}

	err = profiluxController.SetSPortValue(port.PortNumber, enable)
	if err != nil {
		log.Errorf(err, "Unable to set socket id")
		return
	}

	err = update.SPorts(profiluxController, repo)
	if err != nil {
		log.Errorf(err, "Unable update ports")
		return
	}
}
