package update

import (
	service "github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

func levelSensorUpdate(sensor *models.LevelSensor, controller *profilux.Controller) error {
	var err error
	sensor.OperationMode, err = controller.GetLevelSensorMode(sensor.Index)
	if err != nil {
		return err
	}

	sensor.HasTwoInputs = sensor.OperationMode == types.LevelSensorAutoTopOffWith2Sensors ||
		sensor.OperationMode == types.LevelSensorWaterChangeAndAutoTopOff ||
		sensor.OperationMode == types.LevelSensorWaterChange ||
		sensor.OperationMode == types.LevelSensorMinMaxControl

	sensor.HasWaterChange = sensor.OperationMode == types.LevelSensorWaterChangeAndAutoTopOff ||
		sensor.OperationMode == types.LevelSensorWaterChange

	sensor.DisplayName, err = controller.GetLevelName(sensor.Index)
	if err != nil {
		return err
	}

	state, err := controller.GetLevelSensorState(sensor.Index)
	if err != nil {
		return err
	}
	sensor.AlarmState = state.Alarm
	sensor.WaterMode = state.WaterMode

	source1, err := controller.GetLevelSource1(sensor.Index)
	if err != nil {
		return err
	}
	sensorState, err := controller.GetLevelSensorCurrentState(source1)
	if err != nil {
		return err
	}
	sensor.Value = sensorState.Undelayed
	sensor.SensorIndex = source1

	if sensor.HasTwoInputs {
		source2, err := controller.GetLevelSource2(sensor.Index)
		if err != nil {
			return err
		}

		sensorState2, err := controller.GetLevelSensorCurrentState(source2)
		if err != nil {
			return err
		}
		sensor.SecondSensor = sensorState2.Undelayed
		sensor.SecondSensorIndex = source2
	}

	return nil
}

func levelSensorUpdateState(sensor *models.LevelSensor, controller *profilux.Controller) error {
	state, err := controller.GetLevelSensorState(sensor.Index)
	if err != nil {
		return err
	}

	sensor.AlarmState = state.Alarm
	sensor.WaterMode = state.WaterMode

	sensorState, err := controller.GetLevelSensorCurrentState(sensor.SensorIndex)
	if err != nil {
		return err
	}
	sensor.Value = sensorState.Undelayed

	if sensor.HasTwoInputs {
		sensorState2, err := controller.GetLevelSensorCurrentState(sensor.SecondSensorIndex)
		if err != nil {
			return err
		}

		sensor.SecondSensor = sensorState2.Undelayed
	}

	return nil
}

// LevelSensors update the state
func LevelSensors(profiluxController *profilux.Controller, repo service.Controller) error {
	count, err := profiluxController.GetLevelSenosrCount()
	if err != nil {
		return err
	}
	for i := 0; i < count; i++ {
		mode, err := profiluxController.GetLevelSensorMode(i)
		if err != nil {
			return err
		}

		sensor, err := repo.GetLevelSensor(models.GetID(models.LevelSensorType, i))
		if err != nil && err != service.ErrNotFound {
			return nil
		}

		found := err != service.ErrNotFound
		if mode != types.LevelSensorNotEnabled {
			if !found {
				sensor = models.NewLevelSensor(i)
			}

			err = levelSensorUpdate(&sensor, profiluxController)
			if err != nil {
				return err
			}

			err = repo.SetLevelSensor(sensor)
			if err != nil {
				return err
			}

		} else {
			if found {
				err = repo.DeleteLevelSensor(sensor.ID)
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}
