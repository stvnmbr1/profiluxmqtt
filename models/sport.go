package models

import (
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

// SPort model
type SPort struct {
	BaseInfo
	PortNumber int
	Mode       types.PortMode
	Value      types.CurrentState
	IsActive   bool
}

// SPortType name
const SPortType = "SPort"

// NewSPort creates new object
func NewSPort(index int) SPort {
	var probe SPort
	probe.Type = SPortType
	probe.Units = "State"
	probe.PortNumber = index
	probe.ID = GetID(SPortType, index)
	return probe
}
