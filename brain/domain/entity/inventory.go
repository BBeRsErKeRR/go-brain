package entity

import "fmt"

type InventoryRequest struct {
	InventoryPaths string
	Chd            string
	IsJson         bool
}

func (ir *InventoryRequest) String() string {
	return fmt.Sprintf("InventoryPaths: %s, Chd: %s, IsJson %v", ir.InventoryPaths, ir.Chd, ir.IsJson)
}
