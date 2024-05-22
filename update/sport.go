package update

import (
	service "github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

func sPortUpdate(port *models.SPort, controller *profilux.Controller) error {
	var err error
	port.Mode, err = controller.GetSPortFunction(port.PortNumber)
	if err != nil {
		return err
	}

	port.Value, err = controller.GetSPortValue(port.PortNumber)
	if err != nil {
		return err
	}

	port.IsActive = port.Value == types.CurrentStateOn
	port.DisplayName, err = controller.GetSPortName(port.PortNumber)
	return err
}

func sPortUpdateState(port *models.SPort, controller *profilux.Controller) error {
	var err error
	port.Value, err = controller.GetSPortValue(port.PortNumber)
	if err != nil {
		return err
	}

	port.IsActive = port.Value == types.CurrentStateOn
	return nil
}

func SPorts(profiluxController *profilux.Controller, repo service.Controller) error {
	count, err := profiluxController.GetSPortCount()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		port, err := repo.GetSPort(models.GetID(models.SPortType, i))
		if err != nil && err != service.ErrNotFound {
			return err
		}

		found := err != service.ErrNotFound
		mode, err := profiluxController.GetSPortFunction(i)
		if err != nil {
			return err
		}

		if mode.DeviceMode != types.DeviceModeAlwaysOff {
			if !found {
				port = models.NewSPort(i)
			}

			err = sPortUpdate(&port, profiluxController)
			if err != nil {
				return err
			}

			err = repo.SetSPort(port)
			if err != nil {
				return err
			}

		} else {
			if found {
				err = repo.DeleteSPort(port.ID)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
