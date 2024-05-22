package models

import (
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

// DosingPump creates a new item
type DosingPump struct {
	BaseInfo
	Channel  int
	Rate     int
	PerDay   int
// added
	MaxFlow			float64
	DailyDose		float64
	Name			string
	RemainingML		float64
	RemainingDays		float64
	ContainerCapacity	int
	ContainerMinimum	int
//
	Settings types.TimerSettings
}

// DosingPumpType name
const DosingPumpType = "Dosing"

// NewDosingPump creates a new pump
func NewDosingPump(index int) DosingPump {
	var pump DosingPump
	pump.Channel = index
	pump.Type = DosingPumpType
	pump.Units = "ml/day"
	pump.ID = GetID(DosingPumpType, index)
	return pump
}
