package update

import (
	service "github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

func probeUpdate(probe *models.Probe, controller *profilux.Controller) error {
	var err error
	probe.SensorType, err = controller.GetSensorType(probe.Index)
	if err != nil {
		return err
	}

	probe.SensorMode, err = controller.GetSensorMode(probe.Index)
	if err != nil {
		return err
	}

	probe.Format, err = controller.GetSensorFormat(probe.Index)
	if err != nil {
		return err
	}

	probe.Units = probe.GetUnits()
	probe.DisplayName, err = controller.GetProbeName(probe.Index)
	if err != nil {
		return err
	}

	nominalValue, err := controller.GetSensorNominalValue(probe.Index, probe.GetValueMultiplier())
	if err != nil {
		return err
	}
	probe.SetNominalValue(nominalValue)

	alarmDeviation, err := controller.GetSensorAlarmDeviation(probe.Index, probe.GetValueMultiplier())
	if err != nil {
		return err
	}
	probe.SetAlarmDeviation(alarmDeviation)

	probe.AlarmEnable, err = controller.GetSensorAlarmEnable(probe.Index)
	if err != nil {
		return err
	}

	return probeUpdateState(probe, controller)
}

func probeUpdateState(probe *models.Probe, controller *profilux.Controller) error {
	var err error
	probe.AlarmState, err = controller.GetSensorAlarm(probe.Index)
	if err != nil {
		return err
	}

	value, err := controller.GetSensorValue(probe.Index, probe.GetValueMultiplier())
	if err != nil {
		return err
	}
	probe.SetValue(value)

	probe.OperationHours, err = controller.GetProbeOperationHours(probe.Index)
	return err
}

func probes(profiluxController *profilux.Controller, repo service.Controller) error {

	count, err := profiluxController.GetSensorCount()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		sensorType, err := profiluxController.GetSensorType(i)
		if err != nil {
			return err
		}

		var mode types.SensorMode
		active := false
		if sensorType != types.SensorTypeFree && sensorType != types.SensorTypeNone {
			mode, err = profiluxController.GetSensorMode(i)
			if err != nil {
				return err
			}

			active, err = profiluxController.GetSensorActive(i)
			if err != nil {
				return err
			}
		}

		probe, err := repo.GetProbe(models.GetID(models.ProbeType, i))
		if err != nil && err != service.ErrNotFound {
			return err
		}

		found := err != service.ErrNotFound

		if active && mode == types.SensorModeNormal {
			if !found {
				probe = models.NewProbe(i)
			}

			err = probeUpdate(&probe, profiluxController)
			if err != nil {
				return err
			}

			err = repo.SetProbe(probe)
			if err != nil {
				return err
			}

		} else {
			if found {
				err = repo.DeleteProbe(probe.ID)
				if err != nil {
					return err
				}
			}
		}

	}

	return err
}
