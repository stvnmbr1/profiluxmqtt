package update

import (
	data "github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux"
)

func currentPumps(controller *profilux.Controller, repo data.Controller) error {
	for i := 0; i < controller.GetCurrentPumpCount(); i++ {
		pump, err := repo.GetCurrentPump(models.GetID(models.CurrentPumpType, i))
		if err != nil && err != data.ErrNotFound {
			return err
		}

		found := err != data.ErrNotFound
		isAssigned, err := controller.GetIsCurrentPumpAssigned(i)

		if isAssigned {
			if !found {
				pump = models.NewCurrentPump(i)
			}

			pump.Value, err = controller.GetCurrentPumpValue(pump.Index)
			if err != nil {
				return err
			}

			err = repo.SetCurrentPump(pump)
			if err != nil {
				return err
			}

		} else {
			if found {
				err = repo.DeleteCurrentPump(pump.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
