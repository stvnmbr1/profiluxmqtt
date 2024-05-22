package associations

import (
	"github.com/cjburchell/profiluxmqtt/data/repo"
	"github.com/cjburchell/profiluxmqtt/models"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

func getAssociatedModeItemID(mode types.PortMode, repo repo.Controller) (string, string) {
	item := GetAssociatedModeItem(mode, repo)
	if item != nil {
		return item.GetID(), item.GetType()
	}

	return "", ""
}

func GetAssociatedModeItem(mode types.PortMode, repo repo.Controller) models.IBaseInfo {
	index := mode.Port
	if mode.IsProbe {
		probe, err := repo.GetProbe(models.GetID(models.ProbeType, index-1))
		if err == nil {
			return &probe
		}
	}

	switch mode.DeviceMode {
	case types.DeviceModeLights:
		light, err := repo.GetLight(models.GetID(models.LightType, index))
		if err == nil {
			return &light
		}
	case types.DeviceModeTimer:
		timer, err := repo.GetDosingPump(models.GetID(models.DosingPumpType, index))
		if err == nil {
			return &timer
		}
	case types.DeviceModeWater:
		level, err := repo.GetLevelSensor(models.GetID(models.LevelSensorType, index))
		if err == nil {
			return &level
		}

	case types.DeviceModeDrainWater:
		level, err := repo.GetLevelSensor(models.GetID(models.LevelSensorType, index))
		if err == nil {
			return &level
		}

	case types.DeviceModeWaterChange:
		level, err := repo.GetLevelSensor(models.GetID(models.LevelSensorType, index))
		if err == nil {
			return &level
		}

	case types.DeviceModeCurrentPump:
		pump, err := repo.GetCurrentPump(models.GetID(models.CurrentPumpType, index))
		if err == nil {
			return &pump
		}

	case types.DeviceModeProgrammableLogic:
		logic, err := repo.GetProgrammableLogic(models.GetID(models.ProgrammableLogicType, index))
		if err == nil {
			return &logic
		}
	}

	return nil
}

func GetLogicInputs(item models.ProgrammableLogic, repo repo.Controller) []types.PortMode {
	leafs := make([]types.PortMode, 0)
	input1 := GetAssociatedModeItem(item.Input1, repo)
	if input1 != nil {
		if input1.GetType() == models.ProgrammableLogicType {
			leafs = append(leafs, GetLogicInputs(input1.(models.ProgrammableLogic), repo)...)
		} else {
			leafs = append(leafs, item.Input1)
		}
	}

	input2 := GetAssociatedModeItem(item.Input2, repo)
	if input2 != nil {
		if input2.GetType() == models.ProgrammableLogicType {
			leafs = append(leafs, GetLogicInputs(input2.(models.ProgrammableLogic), repo)...)
		} else {
			leafs = append(leafs, item.Input2)
		}
	}

	return leafs
}

// Update the associations
func Update(repo repo.Controller) {

	logicItems, _ := repo.GetProgrammableLogics()
	for _, logic := range logicItems {
		logic.Input1.Id, logic.Input1.Type = getAssociatedModeItemID(logic.Input1, repo)
		logic.Input2.Id, logic.Input2.Type = getAssociatedModeItemID(logic.Input2, repo)
		repo.SetProgrammableLogic(logic)
	}

	sPorts, _ := repo.GetSPorts()
	for _, port := range sPorts {
		port.Mode.Id, port.Mode.Type = getAssociatedModeItemID(port.Mode, repo)
		repo.SetSPort(port)
	}

	lPorts, _ := repo.GetLPorts()
	for _, port := range lPorts {
		port.Mode.Id, port.Mode.Type = getAssociatedModeItemID(port.Mode, repo)
		repo.SetLPort(port)
	}
}
