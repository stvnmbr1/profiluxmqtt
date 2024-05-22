package update

import (
	service "github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

func digitalInputs(controller *profilux.Controller, repo service.Controller) error {
	count, err := controller.GetDigitalInputCount()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		mode, err := controller.GetDigitalInputFunction(i)
		if err != nil {
			return err
		}

		sensor, err := repo.GetDigitalInput(models.GetID(models.DigitalInputType, i))
		if err != nil && err != service.ErrNotFound {
			return err
		}

		found := err != service.ErrNotFound
		if mode != types.DigitalInputFunctionNotUsed {
			if !found {
				sensor = models.NewDigitalInput(i)
			}

			sensor.Function, err = controller.GetDigitalInputFunction(sensor.Index)
			if err != nil {
				return err
			}

			sensor.Value, err = controller.GetDigitalInputState(sensor.Index)
			if err != nil {
				return err
			}

			err = repo.SetDigitalInput(sensor)
			if err != nil {
				return err
			}

		} else {
			if found {
				err = repo.DeleteDigitalInput(sensor.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
