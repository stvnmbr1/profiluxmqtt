package models

import (
	"fmt"
)

// BaseInfo base information for most models
type BaseInfo struct {
	ID          string
	Type        string
	DisplayName string
	Units       string
}

type IBaseInfo interface {
	GetID() string
	GetType() string
}

// GetID Generates an id
func GetID(typ string, index int) string {
	return fmt.Sprintf("%s%d", typ, 1+index)
}

func (i *BaseInfo) GetID() string {
	return i.ID
}

func (i *BaseInfo) GetType() string {
	return i.Type
}
