package repo

import (
	"github.com/cjburchell/profiluxmqtt/models"

	"github.com/pkg/errors"
)

// ErrNotFound when a element is not found
var ErrNotFound = errors.New("not found")

// Controller repo interface
type Controller interface {
	GetInfo() (models.Info, error)
	GetDigitalInputs() ([]models.DigitalInput, error)
	GetDosingPumps() ([]models.DosingPump, error)
	GetLevelSensors() ([]models.LevelSensor, error)
	GetLights() ([]models.Light, error)
	GetLPorts() ([]models.LPort, error)
	GetProbes() ([]models.Probe, error)
	GetProgrammableLogics() ([]models.ProgrammableLogic, error)
	GetCurrentPumps() ([]models.CurrentPump, error)
	GetSPorts() ([]models.SPort, error)

	GetCurrentPump(id string) (models.CurrentPump, error)
	GetDigitalInput(id string) (models.DigitalInput, error)
	GetDosingPump(id string) (models.DosingPump, error)
	GetLevelSensor(id string) (models.LevelSensor, error)
	GetLight(id string) (models.Light, error)
	GetLPort(id string) (models.LPort, error)
	GetProbe(id string) (models.Probe, error)
	GetProgrammableLogic(id string) (models.ProgrammableLogic, error)
	GetSPort(id string) (models.SPort, error)

	SetInfo(info models.Info) error
	SetDigitalInput(input models.DigitalInput) error
	SetDosingPump(pump models.DosingPump) error
	SetLevelSensor(models.LevelSensor) error
	SetLight(models.Light) error
	SetLPort(models.LPort) error
	SetProbe(models.Probe) error
	SetProgrammableLogic(models.ProgrammableLogic) error
	SetCurrentPump(models.CurrentPump) error
	SetSPort(models.SPort) error

	DeleteDigitalInput(id string) error
	DeleteDosingPump(id string) error
	DeleteLevelSensor(id string) error
	DeleteLight(id string) error
	DeleteLPort(id string) error
	DeleteProbe(id string) error
	DeleteProgrammableLogic(id string) error
	DeleteCurrentPump(id string) error
	DeleteSPort(id string) error
}

type controller struct {
	data map[string]map[string]interface{}
}

// NewController create a new controller
func NewController() Controller {
	c := controller{}
	c.data = make(map[string]map[string]interface{})
	return c
}

const (
	infoKey              = "reefstatus:info"
	digitalInputKey      = "reefstatus:digitalinput"
	dosingPumpKey        = "reefstatus:dosingpump"
	levelSensorKey       = "reefstatus:levelSensor"
	lightsKey            = "reefstatus:light"
	lPortKey             = "reefstatus:lPort"
	probeKey             = "reefstatus:probe"
	programmableLogicKey = "reefstatus:programmableLogic"
	pumpKey              = "reefstatus:pump"
	sPortKey             = "reefstatus:sPort"
)

func (controller controller) getItem(key string, id string) (interface{}, error) {
	if _, ok := controller.data[key]; ok {
		if _, ok2 := controller.data[key][id]; ok2 {
			item := controller.data[key][id]
			return item, nil
		}
		return nil, ErrNotFound
	}

	return nil, ErrNotFound
}

func (controller controller) setItem(key string, id string, item interface{}) error {
	if _, ok := controller.data[key]; ok {
		controller.data[key][id] = item
	} else {
		controller.data[key] = make(map[string]interface{})
		controller.data[key][id] = item
	}

	return nil
}

func (controller controller) deleteItem(key string, id string) error {
	if _, ok := controller.data[key]; ok {
		if _, ok2 := controller.data[key][id]; ok2 {
			delete(controller.data, key)
			return nil
		}
		return ErrNotFound
	}

	return ErrNotFound
}

func (controller controller) GetInfo() (models.Info, error) {
	info, err := controller.getItem(infoKey, "info")
	if err == ErrNotFound {
		return models.NewInfo(), nil
	}

	return info.(models.Info), err
}

func (controller controller) GetDigitalInputs() ([]models.DigitalInput, error) {
	items := make([]models.DigitalInput, 0)
	if result, ok := controller.data[digitalInputKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.DigitalInput))
		}
	}
	return items, nil
}
func (controller controller) GetDosingPumps() ([]models.DosingPump, error) {
	items := make([]models.DosingPump, 0)
	if result, ok := controller.data[dosingPumpKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.DosingPump))
		}
	}
	return items, nil
}
func (controller controller) GetLevelSensors() ([]models.LevelSensor, error) {
	items := make([]models.LevelSensor, 0)
	if result, ok := controller.data[levelSensorKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.LevelSensor))
		}
	}
	return items, nil
}
func (controller controller) GetLights() ([]models.Light, error) {
	items := make([]models.Light, 0)
	if result, ok := controller.data[lightsKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.Light))
		}
	}
	return items, nil
}
func (controller controller) GetLPorts() ([]models.LPort, error) {
	items := make([]models.LPort, 0)
	if result, ok := controller.data[lPortKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.LPort))
		}
	}
	return items, nil
}
func (controller controller) GetProbes() ([]models.Probe, error) {
	items := make([]models.Probe, 0)
	if result, ok := controller.data[probeKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.Probe))
		}
	}
	return items, nil
}
func (controller controller) GetProgrammableLogics() ([]models.ProgrammableLogic, error) {
	items := make([]models.ProgrammableLogic, 0)
	if result, ok := controller.data[programmableLogicKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.ProgrammableLogic))
		}
	}
	return items, nil
}
func (controller controller) GetCurrentPumps() ([]models.CurrentPump, error) {
	items := make([]models.CurrentPump, 0)
	if result, ok := controller.data[pumpKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.CurrentPump))
		}
	}
	return items, nil
}
func (controller controller) GetSPorts() ([]models.SPort, error) {
	items := make([]models.SPort, 0)
	if result, ok := controller.data[sPortKey]; ok {
		for _, value := range result {
			items = append(items, value.(models.SPort))
		}
	}
	return items, nil
}

func (controller controller) GetDigitalInput(id string) (models.DigitalInput, error) {
	result, err := controller.getItem(digitalInputKey, id)
	if err != nil {
		return models.DigitalInput{}, err
	}
	return result.(models.DigitalInput), err
}
func (controller controller) GetCurrentPump(id string) (models.CurrentPump, error) {
	result, err := controller.getItem(pumpKey, id)
	if err != nil {
		return models.CurrentPump{}, err
	}
	return result.(models.CurrentPump), err
}

func (controller controller) GetDosingPump(id string) (models.DosingPump, error) {
	result, err := controller.getItem(dosingPumpKey, id)
	if err != nil {
		return models.DosingPump{}, err
	}
	return result.(models.DosingPump), err
}
func (controller controller) GetLevelSensor(id string) (models.LevelSensor, error) {
	result, err := controller.getItem(levelSensorKey, id)
	if err != nil {
		return models.LevelSensor{}, err
	}
	return result.(models.LevelSensor), err
}
func (controller controller) GetLight(id string) (models.Light, error) {
	result, err := controller.getItem(lightsKey, id)
	if err != nil {
		return models.Light{}, err
	}
	return result.(models.Light), err
}
func (controller controller) GetLPort(id string) (models.LPort, error) {
	result, err := controller.getItem(lPortKey, id)
	if err != nil {
		return models.LPort{}, err
	}
	return result.(models.LPort), err
}
func (controller controller) GetProbe(id string) (models.Probe, error) {
	result, err := controller.getItem(probeKey, id)
	if err != nil {
		return models.Probe{}, err
	}
	return result.(models.Probe), err
}
func (controller controller) GetProgrammableLogic(id string) (models.ProgrammableLogic, error) {
	result, err := controller.getItem(programmableLogicKey, id)
	if err != nil {
		return models.ProgrammableLogic{}, err
	}
	return result.(models.ProgrammableLogic), err
}
func (controller controller) GetSPort(id string) (models.SPort, error) {
	result, err := controller.getItem(sPortKey, id)
	if err != nil {
		return models.SPort{}, err
	}
	return result.(models.SPort), err
}

func (controller controller) SetInfo(info models.Info) error {
	return controller.setItem(infoKey, "info", info)
}
func (controller controller) SetDigitalInput(input models.DigitalInput) error {
	return controller.setItem(digitalInputKey, input.ID, input)
}
func (controller controller) SetDosingPump(pump models.DosingPump) error {
	return controller.setItem(dosingPumpKey, pump.ID, pump)
}
func (controller controller) SetLevelSensor(sensor models.LevelSensor) error {
	return controller.setItem(levelSensorKey, sensor.ID, sensor)
}
func (controller controller) SetLight(item models.Light) error {
	return controller.setItem(lightsKey, item.ID, item)
}
func (controller controller) SetLPort(item models.LPort) error {
	return controller.setItem(lPortKey, item.ID, item)
}
func (controller controller) SetProbe(item models.Probe) error {
	return controller.setItem(probeKey, item.ID, item)
}
func (controller controller) SetProgrammableLogic(item models.ProgrammableLogic) error {
	return controller.setItem(programmableLogicKey, item.ID, item)
}
func (controller controller) SetCurrentPump(item models.CurrentPump) error {
	return controller.setItem(pumpKey, item.ID, item)
}
func (controller controller) SetSPort(item models.SPort) error {
	return controller.setItem(sPortKey, item.ID, item)
}

func (controller controller) DeleteDigitalInput(id string) error {
	return controller.deleteItem(digitalInputKey, id)
}
func (controller controller) DeleteDosingPump(id string) error {
	return controller.deleteItem(dosingPumpKey, id)
}
func (controller controller) DeleteLevelSensor(id string) error {
	return controller.deleteItem(levelSensorKey, id)
}
func (controller controller) DeleteLight(id string) error {
	return controller.deleteItem(lightsKey, id)
}
func (controller controller) DeleteLPort(id string) error {
	return controller.deleteItem(lPortKey, id)
}
func (controller controller) DeleteProbe(id string) error {
	return controller.deleteItem(probeKey, id)
}
func (controller controller) DeleteProgrammableLogic(id string) error {
	return controller.deleteItem(programmableLogicKey, id)
}
func (controller controller) DeleteCurrentPump(id string) error {
	return controller.deleteItem(pumpKey, id)
}
func (controller controller) DeleteSPort(id string) error {
	return controller.deleteItem(sPortKey, id)
}
