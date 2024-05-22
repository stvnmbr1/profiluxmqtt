package types

import "fmt"

type OperationMode string

const (
	OperationModeNormal             = "Normal"
	OperationModeDiagnostics        = "Diagnostics"
	OperationModeLightTest          = "LightTest"
	OperationModeMaintenance        = "Maintenance"
	OperationModeManualSockets      = "ManualSockets"
	OperationModeManualIllumination = "ManualIllumination"
	OperationModeCanTransparent     = "CanTransparent"
)

var operationModeMap = map[int]string{
	0: OperationModeNormal,
	1: OperationModeDiagnostics,
	3: OperationModeLightTest,
	4: OperationModeMaintenance,
	5: OperationModeManualSockets,
	6: OperationModeManualIllumination,
	7: OperationModeCanTransparent,
}
var operationModeIdMap = map[string]int{
	OperationModeNormal:             0,
	OperationModeDiagnostics:        1,
	OperationModeLightTest:          3,
	OperationModeMaintenance:        4,
	OperationModeManualSockets:      5,
	OperationModeManualIllumination: 6,
	OperationModeCanTransparent:     7,
}

func GetOperationMode(id int) string {
	if val, ok := operationModeMap[id]; ok {
		return val
	}

	return fmt.Sprintf("Unknown Operation Mode (%d???)", id)
}

func GetOperationIndex(id string) int {
	if val, ok := operationModeIdMap[id]; ok {
		return val
	}
	return 0
}
