package interfaceadapters

import (
	"github.com/bberserkerr/go-brain/brain/domain/entity"
	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate/inventory"
)

type AnsibleInventoryI interface {
	LoadInventory() error
	GetInventoryData() map[string]interface{}
	GetHeadGroups() []string
	GetGroup(groupName string) *inventory.AnsibleGroup
	GetAllGroups() []*inventory.AnsibleGroup
	GetGroupsNames() []string
}

type AnsibleInventoryBuilderI interface {
	Build(inventoryRequest *entity.InventoryRequest) (*inventory.AnsibleInventory, error)
}
