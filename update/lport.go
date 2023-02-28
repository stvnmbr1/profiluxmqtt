package update

import (
	service "github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

func lPortUpdate(port *models.LPort, controller *profilux.Controller) error {
	var err error
	port.Mode, err = controller.GetLPortFunction(port.PortNumber)
	if err != nil {
		return err
	}

	port.Value, err = controller.GetLPortValue(port.PortNumber)
	return err
}

func lPortUpdateState(port *models.LPort, controller *profilux.Controller) error {
	var err error
	port.Value, err = controller.GetLPortValue(port.PortNumber)
	return err
}

func lPorts(profiluxController *profilux.Controller, repo service.Controller) error {
	count, err := profiluxController.GetLPortCount()
	if err != nil {
		return err
	}

	for portNumber := 0; portNumber < count; portNumber++ {
		port, err := repo.GetLPort(models.GetID(models.LPortType, portNumber))
		if err != nil && err != service.ErrNotFound {
			return err
		}

		found := err != service.ErrNotFound

		mode, err := profiluxController.GetLPortFunction(portNumber)
		if err != nil {
			return err
		}

		if mode.DeviceMode != types.DeviceModeAlwaysOff {
			if !found {
				port = models.NewLPort(portNumber)
			}

			err = lPortUpdate(&port, profiluxController)
			if err != nil {
				return err
			}

			err = repo.SetLPort(port)
			if err != nil {
				return err
			}

		} else {
			if found {
				err = repo.DeleteLPort(port.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
