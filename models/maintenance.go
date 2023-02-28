package models

import (
	"fmt"
	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

// Maintenance model
type Maintenance struct {
	DisplayName string
	Index       int
	IsActive    types.CurrentState
	Duration    int
	TimeLeft    int
}

// NewMaintenance creates object
func NewMaintenance(index int) *Maintenance {
	var maintenance Maintenance
	maintenance.Index = index
	maintenance.DisplayName = fmt.Sprintf("Maintenance%d", index+1)
	return &maintenance
}
