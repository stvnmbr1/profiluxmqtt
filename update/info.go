package update

import (
	"github.com/cjburchell/profiluxmqtt/profilux/types"
	"time"

        "github.com/cjburchell/profiluxmqtt/profilux/code"


	service "github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux"
	logger "github.com/cjburchell/uatu-go"
)

func info(controller *profilux.Controller, repo service.Controller) error {

        log := logger.Create(logger.Settings{
                MinLogLevel:  logger.DEBUG,
                ServiceName:  "Profilux MQTT",
                LogToConsole: true,
                UseHTTP:      false,
                UsePubSub:    false,
        })


	info, err := repo.GetInfo()
	if err != nil && err != service.ErrNotFound {
		return err
	}

	info.LastUpdate = time.Now()

        log.Debug("get software version")
	info.SoftwareVersion, err = controller.GetSoftwareVersion()
	if err != nil {
		return err
	}
        log.Debug("get model")

	info.Model, err = controller.GetModel()
	if err != nil {
		return err
	}

	info.SerialNumber, err = controller.GetSerialNumber()
	if err != nil {
		return err
	}

	info.SoftwareDate, err = controller.GetSoftwareDate()
	if err != nil {
		return err
	}

//KH Director

        info.KHDSoftwareVersion, err = controller.GetKHDSoftwareVersion(0)
        if err != nil {
                return err
        }

	info.KHDSerialNumber, err = controller.GetKHDSerialNumber(0)
        if err != nil {
                return err
        }


	info.KHDKHMeasurement, err = controller.GetKHDKHMeasurement(0)
        if err != nil {
                return err
        }

//TEMPERATURE

	info.Temperature, err = controller.GetTemperature(0)
        if err != nil {
                return err
        }

// Dosing pumps

	info.SA_PUMP1_NAME, err = controller.GetPumpName(code.SA_PUMP1_NAME)
       	if err != nil {
               	return err
       	}

        info.SA_PUMP1_REMAINING_ML, err = controller.GetPumpRemainingML(code.SA_PUMP1_REMAINING_ML)
        if err != nil {
                return err
        }

        info.SA_PUMP1_REMAINING_DAYS, err = controller.GetPumpRemainingDays(code.SA_PUMP1_REMAINING_DAYS)
        if err != nil {
                return err
        }

        info.SA_PUMP1_DAILY_DOSE, err = controller.GetPumpDailyDose(code.SA_PUMP1_DAILY_DOSE_ML)
        if err != nil {
                return err
        }

        info.SA_PUMP1_CONT_CAPACITY, err = controller.GetPumpContCapacity(code.SA_PUMP1_CONT_CAPACITY)
        if err != nil {
                return err
        }

        info.SA_PUMP2_NAME, err = controller.GetPumpName(code.SA_PUMP2_NAME)
        if err != nil {
                return err
        }

        info.SA_PUMP2_REMAINING_ML, err = controller.GetPumpRemainingML(code.SA_PUMP2_REMAINING_ML)
        if err != nil {
                return err
        }

        info.SA_PUMP2_REMAINING_DAYS, err = controller.GetPumpRemainingDays(code.SA_PUMP2_REMAINING_DAYS)
        if err != nil {
                return err
        }

	info.SA_PUMP2_DAILY_DOSE, err = controller.GetPumpDailyDose(code.SA_PUMP2_DAILY_DOSE_ML)
        if err != nil {
                return err
        }

        info.SA_PUMP2_CONT_CAPACITY, err = controller.GetPumpContCapacity(code.SA_PUMP2_CONT_CAPACITY)
        if err != nil {
                return err
        }

	info.SA_PUMP3_NAME, err = controller.GetPumpName(code.SA_PUMP3_NAME)
        if err != nil {
                return err
        }

        info.SA_PUMP3_REMAINING_ML, err = controller.GetPumpRemainingML(code.SA_PUMP3_REMAINING_ML)
        if err != nil {
                return err
        }

        info.SA_PUMP3_REMAINING_DAYS, err = controller.GetPumpRemainingDays(code.SA_PUMP3_REMAINING_DAYS)
        if err != nil {
                return err
        }

        info.SA_PUMP3_DAILY_DOSE, err = controller.GetPumpDailyDose(code.SA_PUMP3_DAILY_DOSE_ML)
        if err != nil {
                return err
        }

        info.SA_PUMP3_CONT_CAPACITY, err = controller.GetPumpContCapacity(code.SA_PUMP3_CONT_CAPACITY)
        if err != nil {
                return err
        }

        info.SA_PUMP4_NAME, err = controller.GetPumpName(code.SA_PUMP4_NAME)
        if err != nil {
                return err
        }

        info.SA_PUMP4_REMAINING_ML, err = controller.GetPumpRemainingML(code.SA_PUMP4_REMAINING_ML)
        if err != nil {
                return err
        }

        info.SA_PUMP4_REMAINING_DAYS, err = controller.GetPumpRemainingDays(code.SA_PUMP4_REMAINING_DAYS)
        if err != nil {
                return err
        }

        info.SA_PUMP4_DAILY_DOSE, err = controller.GetPumpDailyDose(code.SA_PUMP4_DAILY_DOSE_ML)
        if err != nil {
                return err
        }

        info.SA_PUMP4_CONT_CAPACITY, err = controller.GetPumpContCapacity(code.SA_PUMP4_CONT_CAPACITY)
        if err != nil {
                return err
        }



//	info.DeviceAddress, err = controller.GetDeviceAddress()
//	if err != nil {
//		return err
//	}
//
//	info.Latitude, err = controller.GetLatitude()
//	if err != nil {
//		return err
//	}
//
//	info.Longitude, err = controller.GetLongitude()
//	if err != nil {
//		return err
//	}

//	info.MoonPhase, err = controller.GetMoonPhase()
//	if err != nil {
//		return err
//	}
//	info.Alarm, err = controller.GetAlarm()
//	if err != nil {
//		return err
//	}
//
//	info.OperationMode, err = controller.GetOperationMode()
//	if err != nil {
//		return err
//	}
//
//	if info.OperationMode == types.OperationModeManualSockets {
//		controller.SetOperationMode(types.OperationModeManualSockets)
//	}
//
//	for i := 0; i < 4; i++ {
//		err = updateMaintenanceMode(controller, &info, i)
//		if err != nil {
//			return err
//		}
//	}
//
//	reminderCount, err := controller.GetReminderCount()
//	if err != nil {
//		return err
//	}
//
//	for i := 0; i < reminderCount; i++ {
//		err = updateReminder(controller, &info, i)
//		if err != nil {
//			return err
//		}
//	}
//
	log.Debug("info returned")
	return repo.SetInfo(info)
}

func updateMaintenanceMode(controller *profilux.Controller, info *models.Info, index int) error {
	var maintenance *models.Maintenance
	for idx, item := range info.Maintenance {
		if item.Index == index {
			maintenance = &info.Maintenance[idx]
			break
		}
	}

	add := false
	if maintenance == nil {
		maintenance = models.NewMaintenance(index)
		add = true
	}

	var err error
	maintenance.Duration, err = controller.GetMaintenanceDuration(maintenance.Index)
	if err != nil {
		return err
	}
	maintenance.DisplayName, err = controller.GetMaintenanceText(maintenance.Index)
	if err != nil {
		return err
	}
	maintenance.IsActive, err = controller.IsMaintenanceActive(maintenance.Index)
	if err != nil {
		return err
	}
	maintenance.TimeLeft, err = controller.GetMaintenanceTimeLeft(maintenance.Index)
	if err != nil {
		return err
	}

	if add {
		info.Maintenance = append(info.Maintenance, *maintenance)
	}

	return nil
}

func updateReminder(controller *profilux.Controller, info *models.Info, index int) error {
	var reminder *models.Reminder
	var reminderIndex = 0
	for i, item := range info.Reminders {
		if item.Index == index {
			reminder = &info.Reminders[i]
			reminderIndex = i
			break
		}
	}

	isActive, err := controller.IsReminderActive(index)
	if err != nil {
		return err
	}

	if !isActive {
		if reminder != nil {
			info.Reminders = append(info.Reminders[:reminderIndex], info.Reminders[reminderIndex+1:]...)
		}
		return nil
	}

	add := false
	if reminder == nil {
		reminder = models.NewReminder(index)
		add = true
	}

	err = reminderUpdate(reminder, controller)
	if err != nil {
		return err
	}

	if add {
		info.Reminders = append(info.Reminders, *reminder)
	}

	return nil
}

// InfoState update
func InfoState(controller *profilux.Controller, repo service.Controller) error {
	info, err := repo.GetInfo()
	if err != nil {
		return err
	}

	info.LastUpdate = time.Now()

	info.KHDKHMeasurement, err = controller.GetKHDKHMeasurement(0)
        if err != nil {
                return err
        }

        info.Temperature, err = controller.GetTemperature(0)
        if err != nil {
                return err
        }


//	info.Alarm, err = controller.GetAlarm()
//	if err != nil {
//		return err
//	}
//	info.OperationMode, err = controller.GetOperationMode()
//	if err != nil {
//		return err
//	}
//
//	if info.OperationMode == types.OperationModeManualSockets {
//		controller.SetOperationMode(types.OperationModeManualSockets)
//	}
//
//	info.MoonPhase, err = controller.GetMoonPhase()
//	if err != nil {
//		return err
//	}
//
//	for index := range info.Maintenance {
//		info.Maintenance[index].IsActive, err = controller.IsMaintenanceActive(info.Maintenance[index].Index)
//		if err != nil {
//			return err
//		}
//
//		info.Maintenance[index].TimeLeft, err = controller.GetMaintenanceTimeLeft(info.Maintenance[index].Index)
//		if err != nil {
//			return err
//		}
//	}
//
//	for _, item := range info.Reminders {
//		err = reminderUpdate(&item, controller)
//		if err != nil {
//			return err
//		}
//	}
//
	return repo.SetInfo(info)
}

func reminderUpdate(reminder *models.Reminder, controller *profilux.Controller) error {
	err := reminderUpdateState(reminder, controller)
	if err != nil {
		return err
	}

	reminder.Text, err = controller.GetReminderText(reminder.Index)
	if err != nil {
		return err
	}

	reminder.Period, err = controller.GetReminderPeriod(reminder.Index)
	if err != nil {
		return err
	}

	reminder.IsRepeating, err = controller.GetReminderIsRepeating(reminder.Index)
	return err
}

func reminderUpdateState(reminder *models.Reminder, controller *profilux.Controller) error {
	var err error
	reminder.Next, err = controller.GetReminderNext(reminder.Index)
	if err != nil {
		return err
	}

	if reminder.Next.Before(time.Now()) {
		reminder.IsOverdue = types.CurrentStateOn
	} else {
		reminder.IsOverdue = types.CurrentStateOff
	}

	return nil
}
