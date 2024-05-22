package models

import (
	"time"

	"github.com/cjburchell/profiluxmqtt/profilux/types"
)

// Info model
type Info struct {
	Maintenance     []Maintenance
	OperationMode   types.OperationMode
	Model           types.Model
	SoftwareDate    time.Time
	DeviceAddress   int
	Latitude        float64
	Longitude       float64
	MoonPhase       float64
	Alarm           types.CurrentState
	SoftwareVersion float64
	SerialNumber    int
	LastUpdate      time.Time
	Reminders       []Reminder
	KHDSerialNumber int
	KHDSoftwareVersion float64
	KHDKHMeasurement float64
	KHDSoftwareDate	time.Time
	Temperature	float64
	SA_PUMP1_NAME	string
	SA_PUMP1_REMAINING_ML int
        SA_PUMP1_REMAINING_DAYS float64
	SA_PUMP1_DAILY_DOSE	int
	SA_PUMP1_CONT_CAPACITY	int
        SA_PUMP2_NAME   string
        SA_PUMP2_REMAINING_ML int
        SA_PUMP2_REMAINING_DAYS float64
        SA_PUMP2_DAILY_DOSE     int
        SA_PUMP2_CONT_CAPACITY  int
        SA_PUMP3_NAME   string
        SA_PUMP3_REMAINING_ML int
        SA_PUMP3_REMAINING_DAYS float64
        SA_PUMP3_DAILY_DOSE     int
        SA_PUMP3_CONT_CAPACITY  int
        SA_PUMP4_NAME   string
        SA_PUMP4_REMAINING_ML int
        SA_PUMP4_REMAINING_DAYS float64
        SA_PUMP4_DAILY_DOSE     int
        SA_PUMP4_CONT_CAPACITY  int
}

// NewInfo creates new object
func NewInfo() Info {
	var info Info
	info.Maintenance = make([]Maintenance, 0)
	info.Reminders = make([]Reminder, 0)

	return info
}

// IsP3 checks to see if the controller is a p3
func (info Info) IsP3() bool {
	return info.Model == types.ProfiLuxIII || info.Model == types.ProfiLuxIIIEx
}
