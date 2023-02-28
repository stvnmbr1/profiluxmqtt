package types

type CurrentState string

const CurrentStateOff = "OFF"
const CurrentStateOn = "ON"

func GetCurrentState(value int) CurrentState {
	if value != 0 {
		return CurrentStateOn
	}

	return CurrentStateOff
}
