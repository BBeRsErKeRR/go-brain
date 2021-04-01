package entity

import (
	"fmt"
)

type DIffEntity struct {
	name string
	obj  map[string]any
}

func NewDIffEntity(name string, obj map[string]any) *DIffEntity {
	return &DIffEntity{
		name: name,
		obj:  obj,
	}
}

func (d *DIffEntity) GetObj() map[string]any {
	return d.obj
}

func (d *DIffEntity) GetName() string {
	return d.name
}

func (d *DIffEntity) String() string {
	return fmt.Sprintf("DIffEntity(name=%s, obj=%v)", d.name, d.obj)
}

type DiffRequest struct {
	Groups []string
}

type DiffEntitySelector struct {
	StrategyType string                 `json:"strategy_type"`
	Parameters   map[string]interface{} `json:"parameters"`
}
