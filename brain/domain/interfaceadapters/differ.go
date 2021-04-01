package interfaceadapters

import (
	"github.com/bberserkerr/go-brain/brain/domain/entity"
	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate"
	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate/inventory"
)

type DifferI interface {
	GetPatch(donorDiff, recipientDiff *entity.DIffEntity) (*entity.PatchEntity, error)
}

type DiffEntitySelectorStrategyI interface {
	GetDiffRequest(donorI, recipientI *inventory.AnsibleInventory) (*aggregate.DiffEntitySelectorResponse, error)
}

type DiffEntitySelectorStrategyBuilderI interface {
	Build(*entity.DiffEntitySelector) DiffEntitySelectorStrategyI
}
